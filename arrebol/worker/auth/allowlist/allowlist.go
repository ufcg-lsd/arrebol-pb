package allowlist

import (
	"bufio"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/service/errors"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/token"
	"log"
	"os"
)

const (
	ListFilePath = "ALLOW_LIST_PATH"
)

type allowList struct {
	list []string
}

type Authorizer struct {
	AllowList allowList
}

func NewAuthorizer() *Authorizer {
	return &Authorizer{AllowList: newAllowList()}
}

func (auth *Authorizer) Authorize(token *token.Token) error {
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
