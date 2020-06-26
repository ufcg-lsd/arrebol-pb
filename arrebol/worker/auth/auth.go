package auth

import (
	"encoding/json"
	"errors"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/allowlist"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/token"
	"github.com/ufcg-lsd/arrebol-pb/crypto"
	"os"
)

const KeysPathKey string = "KEYS_PATH"

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

	if err := saveWorkerKey(worker.ID, rawPublicKey); err != nil {return "", err}

	ok, err := CheckSignature(data, signature, worker.ID)
	if !ok || err != nil {
		return "", err
	}

	return auth.newToken(worker)
}

func CheckSignature(data []byte, signature []byte, workerId string) (bool, error) {
	keyName := workerId + ".pub"
	keyPath := os.Getenv(KeysPathKey) + "/" + keyName
	publicKey, err := crypto.GetPublicKey(keyPath)

	if err != nil {
		return false, err
	}

	err = crypto.Verify(publicKey, data, signature)

	if err != nil {
		return false, err
	}

	return true, nil
}

func saveWorkerKey(workerId, content string) error {
	publicKey, err := crypto.ParsePublicKeyFromPemStr(content)
	if err != nil {return err}
	path := os.Getenv(KeysPathKey) + "/" + workerId + ".pub"
	return crypto.SavePublicKey(path, publicKey)
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
	// TODO authorize token
	var (
		err error
		workerId string
	)
	workerId, err = token.GetWorkerId()
	if err != nil {
		return errors.New("error getting queueId from token")
	}
	if contains := auth.allowlist.Contains(workerId); !contains {
		return errors.New("The worker [" + workerId + "] is not in the allowlist")
	}
	return err
}