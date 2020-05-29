package manager

import (
	"encoding/json"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/token"
)

type Manager struct {
	auth auth.Auth
}

func NewManager() *Manager {
	auth := auth.NewDefaultAuth()
	return &Manager{auth:auth}
}

func (m *Manager) Join(signature string, w worker.Worker) (token.Token, error) {
	data, err := json.Marshal(w)
	if err != nil {
		return token.Token(nil), err
	}
	if err := m.auth.VerifySignature(w.ID, data, []byte(signature)); err != nil {
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