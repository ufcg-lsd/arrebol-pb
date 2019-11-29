package arrebol

import (
	"github.com/emanueljoivo/arrebol/storage"
	"log"
)

type Scheduler struct {
	tasks chan *storage.Task
	workers chan Worker
}

func (s *Scheduler) NewScheduler(tasks chan *storage.Task) *Scheduler {
	s.tasks = tasks
}

func (s *Scheduler) Schedule() {
	for {
		currTask := <-s.tasks
		log.Println(currTask)
	}
}

