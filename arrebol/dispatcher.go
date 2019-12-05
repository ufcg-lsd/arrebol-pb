package arrebol

import (
	"github.com/emanueljoivo/arrebol/storage"
	"log"
	"sync"
)

type Dispatcher struct {
	jobsAccepted chan *storage.Job
	supervisors  map[uint]*Supervisor
	db           *storage.Storage
	mux          sync.Mutex
}

func NewDispatcher(db *storage.Storage) *Dispatcher {
	return &Dispatcher{
		jobsAccepted: make(chan *storage.Job, 100),
		supervisors:  make(map[uint]*Supervisor),
		db:           db,
	}
}

func (d *Dispatcher) HireSupervisor(queue *storage.Queue) {
	d.mux.Lock()
	log.Printf("Hiring new supervisor to the queue %d", queue.ID)
	d.supervisors[queue.ID] = NewSupervisor(queue)
	d.mux.Unlock()
}

func (d *Dispatcher) Start() {
	log.Println("Arrebol Dispatcher start accept jobs")
	d.initDefaultSupervisor()

	for {
		job := <- d.jobsAccepted
		// only receive jobs that belong to a queue
		super := d.supervisors[job.QueueID]
		super.Collect(job)
	}
}

func (d *Dispatcher) initDefaultSupervisor() {
	var q *storage.Queue

	q, err := d.db.GetDefaultQueue()

	if err != nil && q != nil {
		d.HireSupervisor(q)
		super := d.supervisors[q.ID]
		go super.Start()
	}
}

func (d *Dispatcher) AcceptJob(job *storage.Job) {
	d.mux.Lock()
	log.Printf("Job %d accepted\n", job.ID)
	job.State = storage.JobQueued
	d.jobsAccepted <- job
	d.mux.Unlock()
}


