package whitelist

import (
	"bufio"
	"log"
	"os"
)

const (
	WhiteListPath = "WHITE_LIST_PATH"
)

type WhiteList interface {
	Contains(id string) (bool)
}

type FileWhiteList struct {
	sourceFile string
}

func NewFileWhiteList() WhiteList {
	return &FileWhiteList{sourceFile: os.Getenv(WhiteListPath)}
}

func (l *FileWhiteList) Contains(workerId string) (bool) {
	ids, err := l.loadSourceFile()
	if err != nil {
		log.Fatal(err)
	}
	if contains(ids, workerId) {
		return true
	}
	return false
}

func (l *FileWhiteList) loadSourceFile() ([]string, error) {
	file, err := os.Open(os.Getenv(WhiteListPath))
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

func contains(ss []string, s string) bool {
	for _, a := range ss {
		if a == s {
			return true
		}
	}
	return false
}


