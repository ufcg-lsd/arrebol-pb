package key

import (
	"crypto/rsa"
	"github.com/ufcg-lsd/arrebol-pb/crypto"
	"os"
)

const (
	KeysPath = "KEYS_PATH"
)

func SavePublicKey(workerId, content string) error {
	publicKey, err := crypto.ParsePublicKeyFromPemStr(content)
	if err != nil {
		return err
	}
	path := os.Getenv(KeysPath) + "/" + workerId + ".pub"
	return crypto.SavePublicKey(path, publicKey)
}

func GetPublicKey(workerId string) (*rsa.PublicKey, error) {
	keyName := workerId + ".pub"
	keyPath := os.Getenv(KeysPath) + "/" + keyName
	return crypto.GetPublicKey(keyPath)
}
