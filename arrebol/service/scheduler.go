package service

import (
	"errors"
	"github.com/ufcg-lsd/arrebol-pb/storage"
	"log"
	"time"
)

type Scheduler struct {
	P       storage.Policy
	QueueID uint
	Tasks   []*storage.Task
	Jh      *JobsHandler
	S       *storage.Storage
}

const (
	FIFO storage.Policy = iota
)

func NewScheduler(queueId uint, p storage.Policy, j *JobsHandler, s *storage.Storage) Scheduler {
	return Scheduler{
		QueueID: queueId,
		P:       p,
		Tasks:   []*storage.Task{},
		Jh:      j,
		S:       s,
	}
}

func (s *Scheduler) Start() {
	 s.retrieveDispatchedTasks()
	go s.feedPendingTasks()
}

func (s *Scheduler) retrieveDispatchedTasks() []*storage.Task{
	return s.S.RetrieveTasksFromQueueByState(s.QueueID, storage.TaskDispatched)
}

func (s *Scheduler) feedPendingTasks() {
	log.Print("Feeding scheduler tasks")
	for {
		//remove the already scheduled tasks
		for i, task := range s.Tasks {
			if task.State != storage.TaskDispatched {
				s.Tasks = append(s.Tasks[:i], s.Tasks[i+1:]...)
			}
		}

		if len(s.Tasks) <= 5 {
			tasks := s.Jh.GetPendingTasks(s.QueueID)
			for _, task := range tasks {
				if task.State == storage.TaskPending {
					task.State = storage.TaskDispatched
					s.S.SaveTask(task)
					s.Tasks = append(s.Tasks, task)
				}
			}
		}
		log.Println("taks:")
		log.Println(s.Tasks)
		time.Sleep(10*time.Second)
	}
}

func (s *Scheduler) Schedule(worker *storage.Worker) (*storage.Task, error){
	if len(s.Tasks) == 0 {
		return nil, errors.New("No task available")
	}

	switch s.P {
	case FIFO:
		task := s.Tasks[0]
		s.Tasks = s.Tasks[1:]
		task.State = storage.TaskRunning
		s.S.SaveTask(task)
		return task, nil
	default:
		return nil, errors.New("Policy not implemented yet")
	}
}

