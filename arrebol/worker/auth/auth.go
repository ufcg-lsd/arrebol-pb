package auth

import (
	"encoding/json"
	"errors"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/allowlist"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/token"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/key"
	"github.com/ufcg-lsd/arrebol-pb/crypto"
	"log"
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

func (auth *Authenticator) Authenticate(rawPublicKey string, signature []byte, worker *worker.Worker) (token.Token, error) {
	data, err := json.Marshal(worker)
	if err != nil {
		return "", err
	}
	publicKey, err := crypto.ParsePublicKeyFromPemStr(rawPublicKey)
	if err != nil {
		return "", err
	}

	err = crypto.Verify(publicKey, data, signature)
	if err != nil {
		return "", err
	}
	log.Println("The worker [" + worker.ID + "] has a valid signature")

	if contains := auth.allowlist.Contains(worker.ID); !contains {
		return  "", errors.New("The worker [" + worker.ID + "] is not in the allowlist")
	}
	log.Println("The worker [" + worker.ID + "] is allowlisted")

	log.Println("Saving the worker [" + worker.ID + "] public key")
	if err := key.SavePublicKey(worker.ID, rawPublicKey); err != nil {return "", err}
	return auth.newToken(worker)
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