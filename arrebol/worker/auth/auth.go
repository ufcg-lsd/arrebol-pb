package auth

import (
	"encoding/json"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/service/errors"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/token"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/allowlist"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/key"
	"github.com/ufcg-lsd/arrebol-pb/crypto"
)

type Authenticator struct {
	allowlist allowlist.AllowList
}

func NewAuth() *Authenticator {
	allowlist := allowlist.NewAllowList()
	var auth = Authenticator{
		allowlist: allowlist,
	}
	return &auth
}

func (auth *Authenticator) Authenticate(signature []byte, worker *worker.Worker) (token.Token, error) {
	data, err := json.Marshal(worker)
	var token token.Token
	if err != nil {
		return token, err
	}
	publicKey, err := key.GetPublicKey(worker.ID)
	if err != nil {
		return token, err
	}
	err = crypto.Verify(publicKey, data, signature)
	if err != nil {
		return token, err
	}
	if contains := auth.allowlist.Contains(worker.ID); contains {
		token, err = auth.newToken(worker)
		return token, err
	} else {
		return  token, errors.New("The worker [" + worker.ID + "] is not in the allowlist")
	}
}

func (auth *Authenticator) newToken(worker *worker.Worker) (token.Token, error) {
	var t token.Token
	var err error

	t, err = token.NewToken(worker)
	if err != nil {
		return "", err
	}
	return t, nil
}

func (auth *Authenticator) Authorize(token *token.Token) error {
	panic("implement me")
}