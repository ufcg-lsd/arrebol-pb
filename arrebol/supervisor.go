package arrebol

import (
	"github.com/emanueljoivo/arrebol/storage"
	"log"
	"sync"
)

type Supervisor struct {
	queue            *storage.Queue
	scheduler        *Scheduler
	mux              sync.Mutex
}

func NewSupervisor(queue *storage.Queue) *Supervisor {
	return &Supervisor{
		queue:        queue,
		scheduler:    NewScheduler(Fifo),
	}
}

// Starts the supervisor protocol with a static default scheduler
func (s *Supervisor) Start() {
	log.Printf("Supervisor of queue < %d > started\n", s.queue.ID)

	s.pokeScheduler()
}

func (s *Supervisor) Collect(job *storage.Job) {
	log.Printf("Collecting tasks of the job %d", job.ID)
	tasks := &job.Tasks
	for _, task := range *tasks {
		s.scheduler.AddTask(task)
	}
}

func (s *Supervisor) pokeScheduler() {
	log.Println("Scheduler woke up")
	s.scheduler.Start()
}
