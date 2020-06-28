package storage

import (
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"log"
)

func (s *Storage) RetrieveWorkersByQueueID(queueID uint) ([]*worker.Worker, error) {
	var workers []*worker.Worker

	log.Printf("Retrieving workers of queue %d", queueID)
	err := s.driver.Where("queue_id = ?", queueID).Find(&workers).Error

	return workers, err
}
