package key

import (
	"crypto/rsa"
	"github.com/ufcg-lsd/arrebol-pb/crypto"
	"os"
)

const (
	KeysPath = "KEYS_PATH"
)

type Storage interface {
	SavePublicKey(workerPublicKey *rsa.PublicKey) error
	GetPublicKey(workerId string) (*rsa.PublicKey, error)
}

type Reader interface {
	GetPublicKey(workerId string) (*rsa.PublicKey, error)
}

type LocalStorage struct {
	sourceDir string
}

func NewLocalStorage() Storage {
	return &LocalStorage{os.Getenv(KeysPath)}
}

func NewLocalReader() Reader {
	return &LocalStorage{os.Getenv(KeysPath)}
}

func (s *LocalStorage) SavePublicKey(workerPublicKey *rsa.PublicKey) error {
	return nil
}

func (s *LocalStorage) GetPublicKey(workerId string) (*rsa.PublicKey, error) {
	keyName := workerId + ".pub"
	keyPath := s.sourceDir + "/" + keyName
	return crypto.GetPublicKey(keyPath)
}