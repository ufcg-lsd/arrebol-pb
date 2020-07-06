package auth

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/google/logger"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/policy/allowlist"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/policy/tolerant"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/token"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/key"
	"github.com/ufcg-lsd/arrebol-pb/crypto"
)

const AllowListConfKey = "ALLOW_ALL"

type Authorizer interface {
	Authorize(token *token.Token) error
}

type Authenticator interface {
	Authenticate(rawPublicKey string, signature []byte, worker *worker.Worker) (token.Token, error)
}

type Auth struct {
	Authorizer    Authorizer
	Authenticator Authenticator
}

func NewAuth() *Auth {
	return &Auth{
		Authorizer:    NewAuthorizer(),
		Authenticator: NewAuthenticator(),
	}
}

func NewAuthorizer() Authorizer {
	allow, err := strconv.ParseBool(os.Getenv(AllowListConfKey))

	if err != nil {
		logger.Fatalf("Cannot understand the flag: %s", err.Error())
	}
	if allow {
		return allowlist.NewAuthorizer()
	}
	return tolerant.NewAuthorizer()
}

func NewAuthenticator() Authenticator {
	return &DefaultAuthenticator{}
}

type DefaultAuthenticator struct{}

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
