package allowlist

import (
	"bufio"
	"fmt"
	"github.com/google/logger"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/service/errors"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/token"
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

	if token.IsValid() {
		workerId, err = token.GetWorkerId()
		if err != nil {
			msg := fmt.Sprintf("Error getting the workerID from token %v\n", token.String())
			logger.Errorf(msg)
			return errors.New(msg)
		}
		if contains := auth.AllowList.contains(workerId); !contains {
			msg := fmt.Sprintf("The worker [%s] is not in the allowlist\n", workerId)
			logger.Errorf(msg)
			return errors.New(msg)
		}
	}
	logger.Infof("Token [%s] authorized\n", token.String())
	return err
}

func newAllowList() allowList {
	list, err := loadSourceFile()
	if err != nil {
		logger.Fatal(err)
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
