package arrebol

import (
	"github.com/emanueljoivo/arrebol/storage"
	"log"
)

type Scheduler struct {
	tasks chan *storage.Task
	workers []*Worker
}

func NewScheduler(tasks chan *storage.Task, ) *Scheduler {


	return &Scheduler{
		tasks:   tasks,
		workers: nil,
	}
}

func (s *Scheduler) Schedule() {
	for {
		currTask := <-s.tasks
		log.Println(currTask)
	}
}

