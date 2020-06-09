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

func GetPublicKey(keyPath string) (*rsa.PublicKey, error) {
	decodedKey, err := decodeRSAKey(keyPath)
	if err != nil {
		return nil, err
	}

	rsaKey, err := x509.ParsePKCS1PublicKey(decodedKey.Bytes)
	if err != nil {
		return nil, errors.New("Unable to parse RSA public key: " + err.Error())
	}

	return rsaKey, nil
}

func GetPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
	privPem, err := decodeRSAKey(keyPath)
	if err != nil {
		return nil, err
	}
	rsaKey, err := x509.ParsePKCS1PrivateKey(privPem.Bytes)
	if err != nil {
		return nil, errors.New("Unable to parse RSA public key: " + err.Error())
	}

	return rsaKey, nil
}

func ParsePublicKeyFromPemStr(content string) (*rsa.PublicKey, error) {
	block, rest := pem.Decode([]byte(content))
	if len(rest) > 0 {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	if block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("RSA public key is of the wrong type")
	}

	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub, nil
}

func SavePublicKey(path string, key *rsa.PublicKey) error {
	bytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return err
	}
	_pem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: bytes,
	})
	return ioutil.WriteFile(path, _pem, 0644)
}


func decodeRSAKey(keyPath string) (*pem.Block, error) {
	keyContent, err := ioutil.ReadFile(keyPath)

	if err != nil {
		return nil, errors.New("The key [" + keyPath + " ] was not found")
	}

	keyPem, rest := pem.Decode(keyContent)

	if len(rest) > 0 {
		return keyPem, errors.New("Error on decoding key; the rest is not empty.")
	}

	if keyPem.Type != "RSA PRIVATE KEY" && keyPem.Type != "RSA PUBLIC KEY" {
		return keyPem, errors.New("RSA private key is of the wrong type")
	}

	return keyPem, nil
}