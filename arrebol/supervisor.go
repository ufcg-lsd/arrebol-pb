package arrebol

import (
	"github.com/emanueljoivo/arrebol/storage"
	"log"
)

type Supervisor struct {
	queue 	*storage.Queue
	workers chan *Worker
	pendingTasks chan *storage.Task
	scheduler *Scheduler
}

type Worker struct {
	id			string
	nodeAddr	string
}

const WorkerPoolAmount = 2

func NewSupervisor(queue *storage.Queue) *Supervisor {
	return &Supervisor{
		queue: queue,
		pendingTasks: make(chan *storage.Task, 1000),
		scheduler: NewScheduler(Fifo),
	}
}

func (s *Supervisor) HireWorkerPool(node *storage.ResourceNode) {
	node.State = storage.Allocated
	for i := 0; i < WorkerPoolAmount; i++ {
	}
}

func (s *Supervisor) Collect(job *storage.Job) {
	log.Printf("Collecting tasks of the job %d", job.ID)
	tasks := &job.Tasks
	for _, task := range *tasks {
		s.pendingTasks <- task
	}
}

type AllocationPlan struct {
	task   *storage.Task
}

func (s *Supervisor) AllocPlan(task *storage.Task) *AllocationPlan {
	return &AllocationPlan{
		task:   task,
	}
}

