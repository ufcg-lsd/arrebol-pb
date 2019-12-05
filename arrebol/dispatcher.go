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
		jobsAccepted: make(chan *storage.Job),
		supervisors:  make(map[uint]*Supervisor),
		db:           db,
	}
}

func (d *Dispatcher) HireSupervisor(queue *storage.Queue) *Supervisor {
	d.mux.Lock()
	defer d.mux.Unlock()  // must be called after the return

	log.Printf("Hiring new supervisor to the queue %d", queue.ID)

	super := NewSupervisor(queue)
	d.supervisors[queue.ID] = super

	return super
}

func (d *Dispatcher) Start() {
	log.Println("Arrebol Dispatcher start accept jobs")
	d.initDefaultSupervisor()

	for job := range d.jobsAccepted {
		// only receive jobs that belong to a queue
		super := d.supervisors[job.QueueID]
		super.Collect(job)
	}
}

func (d *Dispatcher) initDefaultSupervisor() {
	var q *storage.Queue

	q, err := d.db.GetDefaultQueue()

	if err != nil && q != nil {
		super := d.HireSupervisor(q)
		go super.Start()
	}
}

func (d *Dispatcher) AcceptJob(job *storage.Job) {
	log.Printf("Job %d accepted\n", job.ID)
	job.State = storage.JobQueued
	d.jobsAccepted <- job
}


