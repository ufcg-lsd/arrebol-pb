package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

const (
	KeysPathKey = "KEYS_PATH"
)

// Sign generates a digital signature of the message passed in.
func Sign(prv *rsa.PrivateKey, message []byte) (signature []byte, err error) {
	hash := sha256.New()
	hash.Write(message)
	d := hash.Sum(nil)
	signature, err = rsa.SignPSS(rand.Reader, prv, crypto.SHA256, d, nil)
	return
}

// Verifies an RSA digital signature for the given public key.
func Verify(pub *rsa.PublicKey, message, signature []byte) (err error) {
	hash := sha256.New()
	hash.Write(message)
	d := hash.Sum(nil)
	return rsa.VerifyPSS(pub, crypto.SHA256, d, signature, nil)
}

func GetPublicKey(workerID string) (*rsa.PublicKey, error) {
	keyName := workerID + ".pub"
	decodedKey, err := decodeKey(keyName)
	if err != nil {
		return nil, err
	}

	rsaKey, err := x509.ParsePKCS1PublicKey(decodedKey.Bytes)
	if err != nil {
		return nil, errors.New("Error on parsing public key " + err.Error())
	}

	return rsaKey, nil
}

func decodeKey(keyName string) (*pem.Block, error) {
	keysPath := os.Getenv(KeysPathKey)
	keyContent, err := ioutil.ReadFile(keysPath + "/" + keyName)

	if err != nil {
		return nil, errors.New("The key [" + keyName + " ] was not found")
	}

	decodedKey, rest := pem.Decode(keyContent)

	if len(rest) > 0 {
		return decodedKey, errors.New("Error on decoding key; the rest is not empty.")
	}

	return decodedKey, nil
}