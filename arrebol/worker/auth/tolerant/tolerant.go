package tolerant

import (
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/token"
)

type SimpleAuthenticator struct{}

func NewAuthenticator() *SimpleAuthenticator {
	return &SimpleAuthenticator{}
}

func (sa *SimpleAuthenticator) Authenticate(rawPublicKey string, signature []byte, worker *worker.Worker) (token.Token, error) {
	return nil, nil
}
func (sa *SimpleAuthenticator) Authorize(token *token.Token) error {
	return nil
}
