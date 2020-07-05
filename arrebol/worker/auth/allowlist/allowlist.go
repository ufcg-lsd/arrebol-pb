package allowlist

import (
	"bufio"
	"encoding/json"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/service/errors"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/token"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/key"
	"github.com/ufcg-lsd/arrebol-pb/crypto"
	"log"
	"os"
)

const (
	ListFilePath = "ALLOW_LIST_PATH"
)

type allowList struct {
	list []string
}

type Authenticator struct {
	AllowList allowList
}

func NewAuthenticator() *Authenticator {
	return &Authenticator{AllowList: newAllowList()}
}

func (auth *Authenticator) Authorize(token *token.Token) error {
	// TODO authorize token
	var (
		err      error
		workerId string
	)
	workerId, err = token.GetWorkerId()
	if err != nil {
		return errors.New("error getting queueId from token")
	}
	if contains := auth.AllowList.contains(workerId); !contains {
		return errors.New("The worker [" + workerId + "] is not in the allowlist")
	}
	return err
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
	if err := key.SavePublicKey(worker.ID.String(), rawPublicKey); err != nil {
		return "", err
	}

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

func newAllowList() allowList {
	list, err := loadSourceFile()
	if err != nil {
		log.Fatal(err)
	}
	return allowList{list: list}
}

func (l *allowList) contains(workerId string) bool {
	for _, current := range l.list {
		if current == workerId {
			return true
		}
	}
	return false
}

func loadSourceFile() ([]string, error) {
	file, err := os.Open(os.Getenv(ListFilePath))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ids []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ids = append(ids, scanner.Text())
	}
	return ids, scanner.Err()
}
