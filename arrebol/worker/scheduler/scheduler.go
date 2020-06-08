package scheduler

import (
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/storage"
)

type Scheduler struct {
	storage *storage.Storage
}

func NewScheduler(storage *storage.Storage) *Scheduler {
	return &Scheduler{
		storage:storage,
	}}

func (s *Scheduler) Join(w worker.Worker) (uint, error) {
	queueId := s.selectQueue(w)
	w.QueueID = queueId
	queue, err := s.storage.RetrieveQueue(queueId)
	if err != nil {
		return 0, err
	}
	queue.Workers = append(queue.Workers, &w)
	err = s.storage.SaveQueue(queue)
	if err != nil {
		return 0, err
	}
	return queueId, nil
}

func (s *Scheduler) selectQueue(w worker.Worker) uint {
	return 1
}