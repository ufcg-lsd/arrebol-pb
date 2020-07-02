package service

import (
	"github.com/ufcg-lsd/arrebol-pb/arrebol/service/errors"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/storage"
	"time"
)

type TaskScheduler struct {
	P Policy
	QueueID uint
	PendingTasks []*storage.Task
	Jh *JobsHandler
	S *storage.Storage
}

type Policy uint8

const (
	FIFO Policy = iota
)

func NewTaskScheduler(queueId uint, p Policy, j *JobsHandler, s *storage.Storage) TaskScheduler {
	return TaskScheduler{
		QueueID: queueId,
		P: p,
		PendingTasks: []*storage.Task{},
		Jh: j,
		S: s,
	}
}

func (s *TaskScheduler) Start() {
	go s.feedPendingTasks()
}

func (s *TaskScheduler) feedPendingTasks() {
	for {
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

func (s *TaskScheduler) Schedule(worker *worker.Worker) (*storage.Task, error){
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

