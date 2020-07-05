package tolerant

import (
	"github.com/google/logger"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/token"
)

type Authorizer struct{}

func NewAuthorizer() *Authorizer {
	return &Authorizer{}
}

func (sa *Authorizer) Authorize(token *token.Token) error {
	wID, err := token.GetWorkerId()
	if err != nil {
		logger.Errorf("Unable to retrieve workerId: %s\n", err.Error())
	}
	logger.Infof("WorkerID %s retrieved with success\n", wID)
	// It just this for a simple authorization?
	return err
}
