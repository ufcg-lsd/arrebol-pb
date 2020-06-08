package scheduler

import (
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
)

type Scheduler struct {
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

func (m *Scheduler) Join(w worker.Worker) (uint, error) {
	queueId := m.selectQueue(w)
	w.QueueID = queueId
	return queueId, nil
}

func (m *Scheduler) selectQueue(w worker.Worker) uint {
	return 1
}