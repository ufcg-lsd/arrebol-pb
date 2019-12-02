package arrebol

import (
	"github.com/emanueljoivo/arrebol/storage"
	"log"
)

type Scheduler struct {
	tasks chan *storage.Task
	workers []*Worker
	policy Policy
}

type Policy uint

const (
	Fifo Policy = iota
)

func (rs Policy) String() string {
	return [...]string{"Fifo"}[rs]
}

func NewScheduler(tasks chan *storage.Task, policy Policy) *Scheduler {

	return &Scheduler{
		tasks:   tasks,
		workers: nil,
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

