package authenticator

import (
	"encoding/json"
	"github.com/google/logger"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/token"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/key"
	"github.com/ufcg-lsd/arrebol-pb/crypto"
)

type Authenticator interface {
	Authenticate(rawPublicKey string, signature []byte, worker *worker.Worker) (token.Token, error)
}

type DefaultAuthenticator struct{}

func NewAuthenticator() Authenticator {
	return &DefaultAuthenticator{}
}

func (da *DefaultAuthenticator) Authenticate(rawPublicKey string, signature []byte, worker *worker.Worker) (token.Token, error) {
	data, err := json.Marshal(worker)
	if err != nil {
		logger.Errorln(err.Error())
		return "", err
	}
	publicKey, err := crypto.ParsePublicKeyFromPemStr(rawPublicKey)
	if err != nil {
		logger.Errorln(err.Error())
		return "", err
	}

	err = crypto.Verify(publicKey, data, signature)
	if err != nil {
		logger.Errorln(err.Error())
		return "", err
	}
	if err := key.SavePublicKey(worker.ID.String(), rawPublicKey); err != nil {
		logger.Errorln(err.Error())
		return "", err
	}
	logger.Infof("Worker %s authenticated with success\n", worker.ID.String())
	return newToken(worker)
}

func newToken(worker *worker.Worker) (token.Token, error) {
	var t token.Token
	var err error

	t, err = token.NewToken(worker)
	if err != nil {
		logger.Errorln(err.Error())
		return "", err
	}
	logger.Infof("Token to worker %s created with success\n", worker.ID.String())
	return t, nil
}
