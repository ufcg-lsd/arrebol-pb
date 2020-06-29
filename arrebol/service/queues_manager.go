package service

import (
	"github.com/ufcg-lsd/arrebol-pb/storage"
	"log"
)

type QueuesManager struct {
	Storage *storage.Storage
	Queues  []*storage.Queue
}

func NewQueuesManager(s *storage.Storage) *QueuesManager {
	return &QueuesManager{Storage: s, Queues: loadQueues(s)}
}

func loadQueues(s *storage.Storage) []*storage.Queue {
	queues, err := s.RetrieveQueues()

	if err != nil {
		log.Println("Error on retrieving queues, returning an empty array instead. Error: " + err.Error())
		return []*storage.Queue{}
	}

	return queues
}

func (q *QueuesManager) GetQueues() []*storage.Queue {
	return q.Queues
}

func (q *QueuesManager) AddQueue() {

}

func (q *QueuesManager) RemoveQueue() {

}

func (q *QueuesManager) GetQueue(queueId uint) (*storage.Queue, error) {
	queue, err := q.Storage.RetrieveQueue(queueId)

	if err != nil {
		return nil, err
	}

	return queue, nil
}
