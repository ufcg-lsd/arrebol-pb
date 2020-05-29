package manager

import (
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/token"
)

type Authorization struct {
	Signature []byte
	Message []byte
}

type Manager struct {
	auth auth.Auth
}

func NewManager() *Manager {
	auth := auth.NewDefaultAuth()
	return &Manager{auth:auth}
}

func (m *Manager) Join(a Authorization, w worker.Worker) (token.Token, error) {
	if err := m.auth.VerifySignature(w.ID, a.Message, a.Signature); err != nil {
		return nil, err
	}

	queueId := m.selectQueue(w)
	w.QueueId = queueId

	token, err := m.auth.CreateToken(&w)
	if err != nil {
		return nil, err
	}
	return *token, nil
}

func (m *Manager) selectQueue(w worker.Worker) string {
	return "default"
}