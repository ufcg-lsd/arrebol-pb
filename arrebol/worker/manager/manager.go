package manager

import (
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
)

type Manager struct {
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Join(w worker.Worker) (string, error) {
	queueId := m.selectQueue(w)
	w.QueueId = queueId
	return queueId, nil
}

func (m *Manager) selectQueue(w worker.Worker) string {
	return "default"
}