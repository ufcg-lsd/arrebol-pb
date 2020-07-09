package service

import (
	"errors"
	"github.com/ufcg-lsd/arrebol-pb/storage"
	"time"
)

type Scheduler struct {
	P            storage.Policy
	QueueID      uint
	PendingTasks []*storage.Task
	Jh           *JobsHandler
	S            *storage.Storage
}

const (
	FIFO storage.Policy = iota
)

func NewScheduler(queueId uint, p storage.Policy, j *JobsHandler, s *storage.Storage) Scheduler {
	return Scheduler{
		QueueID: queueId,
		P: p,
		PendingTasks: []*storage.Task{},
		Jh: j,
		S: s,
	}
}

func (s *Scheduler) Start() {
	go s.feedPendingTasks()
}

func (s *Scheduler) feedPendingTasks() {
	for {
		//remove the already scheduled tasks
		for i, task := range s.PendingTasks {
			if task.State != storage.TaskPending {
				s.PendingTasks = append(s.PendingTasks[:i], s.PendingTasks[i+1:]...)
			}
		}

		if len(s.PendingTasks) <= 5 {
			tasks := s.Jh.GetPendingTasks(s.QueueID)
			for _, task := range tasks {
				task.State = storage.TaskDispatched
				s.S.SaveTask(task)
				s.PendingTasks = append(s.PendingTasks, task)
			}
		}

		time.Sleep(10*time.Second)
	}
}

func (s *Scheduler) Schedule(worker *storage.Worker) (*storage.Task, error){
	if len(s.PendingTasks) == 0 {
		return nil, errors.New("No task available")
	}

	switch s.P {
	case FIFO:
		task := s.PendingTasks[0]
		s.PendingTasks = s.PendingTasks[1:]
		task.State = storage.TaskRunning
		s.S.SaveTask(task)
		return task, nil
	default:
		return nil, errors.New("Policy not implemented yet")
	}
}

