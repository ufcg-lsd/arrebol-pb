package arrebol

import (
	"github.com/emanueljoivo/arrebol/storage"
	"log"
)

type Supervisor struct {
	queue 	*storage.Queue
	workers      chan *Worker
	pendingTasks chan *storage.Task
}

func NewSupervisor(queue *storage.Queue) *Supervisor {
	return &Supervisor{
		queue: queue,
		workers: make(chan *Worker, 100),
		pendingTasks: make(chan *storage.Task, 1000),
	}
}

type AllocationPlan struct {
	worker *Worker
	task   *storage.Task
}

func (s *Supervisor) AllocPlan(task *storage.Task, worker *Worker) *AllocationPlan {
	return &AllocationPlan{
		worker: worker,
		task:   task,
	}
}

func (s *Supervisor) Collect(job *storage.Job) {
	log.Printf("Collecting tasks of the job %d", job.ID)
	tasks := &job.Tasks
	for _, task := range *tasks {
		s.pendingTasks <- task
	}
}

