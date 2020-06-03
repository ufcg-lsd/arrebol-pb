package auth

import (
	"encoding/json"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/service/errors"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/token"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/whitelist"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/key"
	"github.com/ufcg-lsd/arrebol-pb/crypto"
)

type Authenticator interface {
	Authenticate(signature []byte, worker *worker.Worker) error
	Authorize(token *token.Token) error
	NewToken(worker *worker.Worker) (*token.Token, error)
}

type JWTAuthenticator struct {
	whitelist       whitelist.WhiteList
}

func NewJWTAuth() *Authenticator {
	whitelist := whitelist.NewFileWhiteList()
	var auth Authenticator = &JWTAuthenticator{
		whitelist: whitelist,
	}
	return &auth
}

func (auth *JWTAuthenticator) Authenticate(signature []byte, worker *worker.Worker) error {
	data, err := json.Marshal(worker)
	if err != nil {
		return err
	}
	publicKey, err := key.GetPublicKey(worker.ID)
	if err != nil {
		return err
	}
	err = crypto.Verify(publicKey, data, signature)
	if err != nil {
		return err
	}
	if contains := auth.whitelist.Contains(worker.ID); !contains {
		return errors.New("The worker [" + worker.ID + "] does not have permission")
	}
	return nil
}

func (auth *JWTAuthenticator) NewToken(worker *worker.Worker) (*token.Token, error) {
	var t token.Token
	var err error

	t, err = token.NewJWToken(worker)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (auth *JWTAuthenticator) Authorize(token *token.Token) error {
	panic("implement me")
}