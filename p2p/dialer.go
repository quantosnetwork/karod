package p2p

import (
	"context"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/quantosnetwork/karod/workers"
	"log"
	"runtime"
	"strconv"
	"time"
)

var MaxWorkers = runtime.NumCPU() - 1

func (s *Server) RunDialerWorkPool() {
	wp := workers.New(MaxWorkers)
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go wp.GenerateFrom(s.getDialJobs())
	go wp.Run(ctx)

	for {
		select {
		case r, ok := <-wp.Results():
			if !ok {
				continue
			}
			_, err := strconv.ParseInt(string(r.Descriptor.ID), 10, 64)
			if err != nil {
				log.Fatalf("unexpected error: %v", err)
			}
		case <-wp.Done:
			return
		default:
		}
	}

}

func getMaxJobsPerWorker(jobCnt int) int {
	return jobCnt / MaxWorkers
}

func (s *Server) getDialJobs() []workers.Job {
	queue := s.dialQueue
	jobs := make([]workers.Job, len(queue.items))
	for i, _ := range jobs {
		id, _ := uuid.NewUUID()
		jobs[i] = workers.Job{
			Descriptor: setJobDescriptor(id.String()),
			ExecFn:     execFn,
			Args:       queue.items[i],
		}
	}
	return jobs
}

var execFn = func(ctx context.Context, args interface{}) (interface{}, error) {
	return nil, nil
}

func setJobDescriptor(jobID string) workers.JobDescriptor {
	return workers.JobDescriptor{
		ID:    workers.JobID(jobID),
		JType: workers.JobType("p2pDialer"),
		Metadata: workers.JobMetadata{
			"dialer":      "ok",
			"timeCreated": time.Now().UnixNano(),
		},
	}
}

type dialQueue struct {
	items []*peer.AddrInfo
}

type dialJobArgs struct {
	toDial *peer.AddrInfo
}
