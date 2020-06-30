package service

import (
	"github.com/ufcg-lsd/arrebol-pb/storage"
	"log"
)

type QueuesManager struct {
	Storage *storage.Storage
	Queues  []*storage.Queue
	Schedulers map[uint]TaskScheduler
}

func NewQueuesManager(s *storage.Storage) *QueuesManager {
	queues := loadQueues(s)
	schedulers := loadSchedulers(queues)
	return &QueuesManager{Storage: s, Queues: queues, Schedulers: schedulers}
}

func loadSchedulers(queues []*storage.Queue) map[int]TaskScheduler {
	var schedulers map[uint]TaskScheduler
	for _, queue := range queues {
		scheduler := NewTaskScheduler(queue.ID, queue.SchedulingPolicy)
		go scheduler.Start()
		schedulers[queue.ID] = scheduler
	}
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
