package auth

import (
	"encoding/json"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/service/errors"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/key"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/token"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/whitelist"
	"github.com/ufcg-lsd/arrebol-pb/crypto"
)

type Auth interface {
	Validate(signature []byte, worker *worker.Worker) error
	Authenticate( worker *worker.Worker) (*token.Token, error)
	Authorize(token *token.Token) error
}

type JWTAuth struct {
	workerKeyReader key.Reader
	whitelist       whitelist.WhiteList
}

func NewJWTAuth() *Auth {
	reader := key.NewLocalReader()
	whitelist := whitelist.NewFileWhiteList()
	var auth Auth = &JWTAuth{
		workerKeyReader: reader,
		whitelist: whitelist,
	}
	return &auth
}


func (auth *JWTAuth) Validate(signature []byte, worker *worker.Worker) error {
	data, err := json.Marshal(worker)
	if err != nil {
		return err
	}
	publicKey, err := auth.workerKeyReader.GetPublicKey(worker.ID)
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

func (auth *JWTAuth) Authenticate(worker *worker.Worker) (*token.Token, error) {
	var t *token.Token
	var err error

	if t, err = auth.createToken(worker); err != nil {
		return nil, err
	}
	return t, nil

}

func (auth *JWTAuth) Authorize(token *token.Token) error {
	panic("implement me")
}

func (auth *JWTAuth) createToken(worker *worker.Worker) (*token.Token, error) {
	var t token.Token
	var err error

	t, err = token.NewJWToken(worker)
	if err != nil {
		return nil, err
	}
	return &t, nil
}