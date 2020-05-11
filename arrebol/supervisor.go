package arrebol

import (
	"github.com/ufcg-lsd/arrebol/storage"
	"log"
	"sync"
	"time"
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
		task.State = storage.TaskPending
		s.scheduler.AddTask(task)
	}
	go s.jobStateMonitor(job.ID)
}

func (s *Supervisor) jobStateMonitor(jobId uint) {
	for {
		job, _ := storage.DB.RetrieveJobByQueue(jobId, s.queue.ID)
		js := s.getJobState(*job)
		if job.State != js {
			storage.DB.SetJobState(job.ID, js)
			log.Printf("Updated Job [%d] to state [%s]", jobId, js.String())
		}
		if job.State == storage.JobFinished || job.State == storage.JobFailed {
			break
		}
		time.Sleep(3 * time.Second)
	}
}

func (s *Supervisor) getJobState(job storage.Job) storage.JobState {
	var jobState storage.JobState
	if s.isAll([]storage.TaskState{storage.TaskFailed}, job.Tasks) {
		jobState = storage.JobFailed
	} else if s.isAll([]storage.TaskState{storage.TaskFailed, storage.TaskFinished}, job.Tasks) {
		jobState = storage.JobFinished
	} else if s.isAll([]storage.TaskState{storage.TaskPending}, job.Tasks) {
		jobState = storage.JobQueued
	} else {
		jobState = storage.JobRunning
	}
	return jobState
}

func (s *Supervisor) isAll(states []storage.TaskState, tasks []*storage.Task) bool {
	for _, t := range tasks {
		if !contains(t.State, states) {
			return false
		}
	}
	return true
}

func contains(e storage.TaskState, arr []storage.TaskState) bool {
	for _, a := range arr {
		if a == e {
			return true
		}
	}
	return false
}

func (s *Supervisor) pokeScheduler() {
	log.Println("Scheduler woke up")
	s.scheduler.Start()
}
