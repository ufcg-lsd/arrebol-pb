package manager

import (
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
)

type Manager struct {
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Join(w worker.Worker) (uint, error) {
	queueId := m.selectQueue(w)
	w.QueueID = queueId
	return queueId, nil
}

func (m *Manager) selectQueue(w worker.Worker) uint {
	return 1
}