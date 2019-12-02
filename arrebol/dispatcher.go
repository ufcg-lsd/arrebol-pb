package arrebol

import (
	"github.com/emanueljoivo/arrebol/storage"
	"log"
	"sync"
)

type Dispatcher struct {
	jobsAccepted chan *storage.Job
	supervisor map[uint]*Supervisor
	mux sync.Mutex
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		jobsAccepted: make(chan *storage.Job, 100),
		supervisor: make(map[uint]*Supervisor, 0),
	}
}

func (m *Dispatcher) HireSupervisor(queue *storage.Queue) {
	m.mux.Lock()
	log.Printf("Hiring new supervisor to the queue %d", queue.ID)
	m.supervisor[queue.ID] = NewSupervisor(queue)
	m.mux.Unlock()
}

func (m *Dispatcher) Start() {
	log.Println("Arrebol Dispatcher start accept jobs")
	for {
		job := <- m.jobsAccepted
		super := m.supervisor[job.QueueID]
		super.Collect(job)
	}
}

func (m *Dispatcher) AcceptJob(job *storage.Job) {
	m.mux.Lock()
	log.Printf("Job %d accepted\n", job.ID)
	job.State = storage.JobQueued
	m.jobsAccepted <- job
	m.mux.Unlock()
}


