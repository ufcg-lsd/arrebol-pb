package service

import (
	"errors"
	"github.com/ufcg-lsd/arrebol-pb/storage"
	"log"
)

type QueuesManager struct {
	Storage *storage.Storage
	Queues  []*storage.Queue
	Schedulers map[uint]Scheduler
}

func NewQueuesManager(s *storage.Storage, j *JobsHandler) *QueuesManager {
	queues := loadQueues(s)
	schedulers := loadSchedulers(queues, s, j)
	return &QueuesManager{Storage: s, Queues: queues, Schedulers: schedulers}
}

func loadSchedulers(queues []*storage.Queue, s *storage.Storage, j *JobsHandler) map[uint]Scheduler {
	schedulers := map[uint]Scheduler{}
	for _, queue := range queues {
		scheduler := NewScheduler(queue.ID, queue.SchedulingPolicy, j, s)
		go scheduler.Start()
		schedulers[queue.ID] = scheduler
	}
	return schedulers
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

func (q *QueuesManager) AddQueue(queue *storage.Queue, j *JobsHandler) error {
	err := q.Storage.SaveQueue(queue)

	if err != nil {
		return err
	}

	q.Queues = append(q.Queues, queue)
	scheduler := NewScheduler(queue.ID, queue.SchedulingPolicy, j, q.Storage)
	go scheduler.Start()
	q.Schedulers[queue.ID] = scheduler

	return nil
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

func (q *QueuesManager) GetQueueScheduler(queueId uint) (*Scheduler, error) {
	queueScheduler, ok := q.Schedulers[queueId]

	if !ok {
		return nil, errors.New("Queue not found")
	}

	return &queueScheduler, nil
}

func (q *QueuesManager) AddJob(queueId uint, j *storage.Job) error {
	var queue *storage.Queue
	for _, curr := range q.Queues {
		if curr.ID == queueId {
			queue = curr
			break
		}
	}

	if queue == nil {
		return errors.New("Queue " + string(queueId) + " not found")
	}

	queue.Jobs = append(queue.Jobs, j)

	return q.Storage.SaveQueue(queue)
}


