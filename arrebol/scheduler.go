package arrebol

import (
	"github.com/emanueljoivo/arrebol/storage"
	"log"
)

// no preemptive
type Scheduler struct {
	tasks chan *storage.Task
	policy Policy
}

type Policy uint

const (
	Fifo Policy = iota
)

func (rs Policy) String() string {
	return [...]string{"Fifo"}[rs]
}

func NewScheduler(policy Policy) *Scheduler {
	return &Scheduler{
		policy: policy,
	}
}

func (s *Scheduler) Schedule() {
	for {
		switch s.policy {

		case Fifo:
			currTask := <-s.tasks
			log.Println(currTask)
		}
	}
}

