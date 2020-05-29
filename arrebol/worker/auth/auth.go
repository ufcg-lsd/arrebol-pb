package auth

import (
	"github.com/ufcg-lsd/arrebol-pb/arrebol/service/errors"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/key"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/token"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/whitelist"
	"github.com/ufcg-lsd/arrebol-pb/crypto"
)

type Auth interface {
	VerifySignature(workerId string, message, signature []byte) error
	CreateToken(worker *worker.Worker) (*token.Token, error)
	VerifyToken(token *token.Token) error
}

type DefaultAuth struct {
	workerKeyReader key.Reader
	tokenGenerator  token.Generator
	whitelist       whitelist.WhiteList
}

func NewDefaultAuth() Auth {
	reader := key.NewLocalReader()
	gen := token.NewSimpleGenerator()
	whitelist := whitelist.NewFileWhiteList()
	return &DefaultAuth{
		workerKeyReader: reader,
		tokenGenerator:  gen,
		whitelist:       whitelist,
	}
}

func (auth *DefaultAuth) VerifySignature(workerId string, message, signature []byte) error {
	publicKey, err := auth.workerKeyReader.GetPublicKey(workerId)
	if err != nil {
		return err
	}
	err = crypto.Verify(publicKey, message, signature)
	if err != nil {
		return err
	}
	return nil
}

func (auth *DefaultAuth) CreateToken(worker *worker.Worker) (*token.Token, error) {
	contains, err := auth.whitelist.Contains(worker.ID)
	if err != nil {
		return nil, err
	}
	if !contains {
		return nil, errors.New("The worker does not have permission")
	}
	token, err := auth.tokenGenerator.NewToken(worker)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (auth *DefaultAuth) VerifyToken(token *token.Token) error {
	return nil
}