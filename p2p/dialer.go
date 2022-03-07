package p2p

import (
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	"sync"
)

const MaxDialers = 4

const MinDialers = 2

type dialQueue struct {
	heap     dialQueueImpl
	lock     sync.Mutex
	items    map[peer.ID]*dialJob
	updateCh chan struct{}
	closeCh  chan struct{}
}

type dialWorker struct {
	id      int
	jobs    chan dialJob
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

	dw := &dialWorkerPool{
		maxDialers:   MaxDialers,
		minDialers:   MinDialers,
		workers:      make([]*dialWorker, MaxDialers),
		workersState: make(map[*dialWorker]int),
		jobsQueue:    new(dialQueue),
	}

	for i := 0; i < MaxDialers; i++ {
		dw.NewWorker()
	}
	fmt.Printf("added %x workers \n", len(dw.workers))
	return dw
}

func (dw *dialWorkerPool) NewWorker() {
	worker := &dialWorker{
		id:      len(dw.workers),
		jobs:    make(chan dialJob),
		results: make(chan string),
	}
	dw.workers[worker.id-1] = worker
	dw.workersState[worker] = 0 // 0 = idle initial state

}

func (dw *dialWorkerPool) dispatchJob(jobs <-chan dialJob, results chan<- string) {
	
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
