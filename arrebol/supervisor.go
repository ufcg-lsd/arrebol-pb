package arrebol

import (
	"github.com/emanueljoivo/arrebol/storage"
	"log"
	"os"
	"strconv"
)

type Supervisor struct {
	queue        *storage.Queue
	workers      chan *Worker
	pendingTasks chan *storage.Task
	scheduler    *Scheduler
}

func NewSupervisor(queue *storage.Queue) *Supervisor {
	return &Supervisor{
		queue:        queue,
		pendingTasks: make(chan *storage.Task, 1000),
		scheduler:    NewScheduler(Fifo),
	}
}

// should be specific by node
func (s *Supervisor) HireWorkerPool(driver Driver) {
	switch driver {
	case Raw:
		log.Println("just support system level execution with static pool of workers")
		pool, _ := strconv.Atoi(os.Getenv("STATIC_WORKER_POOL"))

		for i := 0; i < pool; i++ {
			s.workers <- NewWorker(Raw)
		}

	case Docker:
		log.Println("not supported yet")
	default:
		log.Println("no worker type")
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
	task *storage.Task
}

func (s *Supervisor) AllocPlan(task *storage.Task) *AllocationPlan {
	return &AllocationPlan{
		task: task,
	}
}
