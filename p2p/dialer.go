package p2p

import (
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	"runtime"
	"sync"
)

type dialQueue struct {
	heap     dialQueueImpl
	lock     sync.Mutex
	items    map[peer.ID]*dialJob
	updateCh chan struct{}
	closeCh  chan struct{}
}

type dialWorker struct {
	id      int
	jobs    chan *dialJob
	results chan string
}

type dialJob struct {
	index    int
	addr     peer.AddrInfo
	priority uint64
}

type dialWorkerPool struct {
	maxDialers   int
	minDialers   int
	workers      []*dialWorker
	workersState map[*dialWorker]int
	jobsQueue    *dialQueue
}

func NewWorkerPool() *dialWorkerPool {

	maxWorkers := runtime.NumCPU() - 1

	dw := &dialWorkerPool{
		maxDialers:   maxWorkers,
		minDialers:   2,
		workers:      make([]*dialWorker, maxWorkers),
		workersState: make(map[*dialWorker]int),
		jobsQueue:    new(dialQueue),
	}

	for i := 0; i < maxWorkers; i++ {
		dw.NewWorker(i + 1)
	}
	fmt.Printf("added %x workers \n", len(dw.workers))
	return dw
}

func (dw *dialWorkerPool) NewWorker(i int) {
	worker := &dialWorker{
		id:      i,
		jobs:    make(chan *dialJob),
		results: make(chan string),
	}
	dw.workers[worker.id-1] = worker
	dw.workersState[worker] = 0 // 0 = idle initial state

}

func (dw *dialWorkerPool) dispatchJob(jobs chan *dialJob, results chan string) {
	maxJobs := len(jobs) / len(dw.workers)
	var wg sync.WaitGroup
	defer wg.Done()
	for w := 1; w <= len(dw.workers); w++ {
		wg.Add(1)

		for i := 1; i <= maxJobs; i++ {
			for j := range jobs {
				dw.workers[w].jobs <- j
				dw.workers[w].work(jobs, results)
				<-results
			}
		}
	}
	close(jobs)

}

func (w *dialWorker) work(jobs chan *dialJob, results chan string) {

	for j := range jobs {
		results <- "done" + j.addr.String()
	}

}

func (s *Server) DialFunc(addr []peer.AddrInfo) {
	var wg sync.WaitGroup
	wp := NewWorkerPool()
	jobs := wp.jobsQueue
	jobs.items = make(map[peer.ID]*dialJob)

	for i, p := range addr {
		jobs.items[p.ID] = &dialJob{index: i, addr: addr[i], priority: uint64(i)}
		wp.jobsQueue.heap.Push(jobs.items[p.ID])
	}
	defer wg.Done()
	wg.Add(len(wp.workers))

}

type dialQueueImpl []*dialJob

func (d dialQueueImpl) Len() int {
	return len(d)
}

func (d dialQueueImpl) Less(i, j int) bool {
	return d[i].priority < d[j].priority
}

func (d dialQueueImpl) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
	d[i].index = i
	d[j].index = j
}

func (d *dialQueueImpl) Push(x interface{}) {
	n := len(*d)
	item := x.(*dialJob) //nolint:forcetypeassert
	item.index = n
	*d = append(*d, item)
}

func (d *dialQueueImpl) Pop() interface{} {
	old := *d
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*d = old[0 : n-1]

	return item
}
