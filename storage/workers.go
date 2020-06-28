package storage

import (
	"errors"
	"fmt"
	"github.com/google/logger"
	uuid "github.com/satori/go.uuid"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
)

var (
	SaveWorkerErr = errors.New("unable to create workflow")
)

func (s *Storage) RetrieveWorkersByQueueID(queueID uint) ([]*worker.Worker, error) {
	var workers []*worker.Worker

	logger.Infof("Retrieving workers of queue %d", queueID)
	err := s.driver.Where("queue_id = ?", queueID).Find(&workers).Error

	return workers, err
}

func (s *Storage) SaveWorker(w worker.Worker ) (uuid.UUID, error) {
	tx := s.driver.Begin()

	savedWorker := &worker.Worker{}
	if tx.Create(&w).Error != nil {
		tx.Rollback()
		errMsg := fmt.Sprintf("Rollback done because %s", SaveWorkerErr.Error())
		logger.Errorf(errMsg)
		return uuid.Nil, SaveWorkerErr
	} else {
		tx.Last(savedWorker)
		tx.Commit()
		sucMsg := fmt.Sprintf("Worker %s saved successfully", savedWorker.ID.String())
		logger.Infof(sucMsg)
		return savedWorker.ID, nil
	}
}