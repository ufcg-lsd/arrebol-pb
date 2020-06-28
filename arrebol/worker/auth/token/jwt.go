package token

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/crypto"
	"os"
	"time"
)

const (
	ArrebolPrivKeyPath = "ARREBOL_PRIV_KEY_PATH"
	ArrebolPubKeyPath  = "ARREBOL_PUB_KEY_PATH"
	ExpirationTime = 10 * time.Minute
)

type Token string

type Claims struct {
	QueueId uint `json:"QueueId"`
	WorkerId string `json:"WorkerId"`
	jwt.StandardClaims
}

func NewToken(worker *worker.Worker) (Token, error){
	expirationTime := time.Now().Add(ExpirationTime)
	claims := &Claims{
		QueueId:        worker.QueueID,
		WorkerId:       worker.ID.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	signedToken, err := signToken(token)
	if err != nil {
		return "", err
	}
	return Token(signedToken), nil
}

func signToken(token *jwt.Token) (string, error) {
	privateKey, err := crypto.GetPrivateKey(os.Getenv(ArrebolPrivKeyPath))
	if err != nil {
		return "", err
	}
	return token.SignedString(privateKey)
}

func Parse(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return crypto.GetPublicKey(os.Getenv(ArrebolPubKeyPath))
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (t Token) String() string {
	return string(t)
}

func (t Token) Expired() bool {
	_, err := Parse(t.String())
	v, _ := err.(*jwt.ValidationError)

	if v.Errors == jwt.ValidationErrorExpired {
		return true
	}
	return false
}

func (t Token) GetPayloadField(key string) (interface{}, error) {
	token, err := Parse(t.String())
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims[key], nil
	} else {
		return nil, errors.New("Error while get payload from jwtoken")
	}
}

func (t Token) GetQueueId() (uint, error) {
	_queueId, err := t.GetPayloadField("QueueId")
	if err != nil {
		return 0, err
	}
	if queueId, ok := _queueId.(uint); ok {
		return queueId, nil
	}
	return 0, errors.New("QueueId was not found")
}

func (t Token) GetWorkerId() (string, error) {
	_workerId, err := t.GetPayloadField("WorkerId")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", _workerId), nil
}

func (t Token) SetPayloadField(key string, value interface{}) (Token, error) {
	token, err := Parse(t.String())
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claims[key] = value
		token.Claims = claims
		tokenStr, err := signToken(token)
		return Token(tokenStr), err
	} else {
		panic("Error while set payload from token")
	}
}

func (t Token) IsValid() bool {
	_, err := Parse(t.String())
	if err != nil {
		return false
	}
	return true
}