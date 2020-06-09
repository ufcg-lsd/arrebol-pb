package key

import (
	"crypto/rsa"
	"github.com/ufcg-lsd/arrebol-pb/crypto"
	"os"
)

const (
	KeysPath = "KEYS_PATH"
)

func SavePublicKey(workerPublicKey *rsa.PublicKey) error {
	return nil
}

func GetPublicKey(workerId string) (*rsa.PublicKey, error) {
	keyName := workerId + ".pub"
	keyPath := os.Getenv(KeysPath) + "/" + keyName
	return crypto.GetPublicKey(keyPath)
}