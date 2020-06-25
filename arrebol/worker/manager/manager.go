package manager

import (
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/storage"
	"log"
	"strconv"
)

type Manager struct {
	storage *storage.Storage
}

func NewManager(storage *storage.Storage) *Manager {
	return &Manager{
		storage:storage,
	}
}

//Join selects a queue for Worker to work on and joins him to that queue.
//If a problem occurs at the join, an error is returned and the queue ID returned is 0 by default,
//but the only error indicator is if the err variable is not null.
func (m *Manager) Join(w worker.Worker) (uint, error) {
	queueId := m.selectQueue(w)
	w.QueueID = queueId
	queue, err := m.storage.RetrieveQueue(queueId)
	if err != nil {
		return 0, err
	}
	queue.Workers = append(queue.Workers, &w)
	err = m.storage.SaveQueue(queue)
	if err != nil {
		return 0, err
	}
	log.Println("Worker [" + w.ID + "] has been assigned to queue [" + strconv.Itoa(int(queueId)) + "]")
	return queueId, nil
}

func (m *Manager) selectQueue(w worker.Worker) uint {
	log.Println("Selecting a queue for worker [" + w.ID + "]")
	return 1
}