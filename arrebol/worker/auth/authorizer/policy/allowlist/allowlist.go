package allowlist

import (
	"bufio"
	"github.com/google/logger"
	"os"
)

const (
	ListFilePath = "ALLOW_LIST_PATH"
)

type allowList struct {
	list []string
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
