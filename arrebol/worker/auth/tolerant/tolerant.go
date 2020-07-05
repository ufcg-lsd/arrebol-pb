package tolerant

import (
	"github.com/google/logger"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/token"
)

type SimpleAuthorizer struct{}

func NewAuthorizer() *SimpleAuthorizer {
	return &SimpleAuthorizer{}
}

func (sa *SimpleAuthorizer) Authorize(token *token.Token) error {
	wID, err := token.GetWorkerId()
	if err != nil {
		logger.Errorf("Unable to retrieve workerId: %s", err.Error())
	}
	logger.Infof("WorkerID %s retrieved with success", wID)
	// It just this for a simple authorization?
	return err
}
