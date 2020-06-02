package token

import (
	//...
	// import the jwt-go library
	"github.com/dgrijalva/jwt-go"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/crypto"
	"log"
	"os"
	"time"

	//...
)

const (
	ArrebolPrivKeyPath = "ARREBOL_PRIV_KEY_PATH"
	ArrebolPubKeyPath  = "ARREBOL_PUB_KEY_PATH"
)

type JWToken struct {
	raw string
	QueueId string
	WorkerId string
}

type Claims struct {
	QueueId string `json:"QueueId"`
	WorkerId string `json:"WorkerId"`
	jwt.StandardClaims
}

func NewJWToken(worker *worker.Worker) (*JWToken, error){
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		QueueId:        worker.QueueId,
		WorkerId:       worker.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	privateKey, err := crypto.GetPrivateKey(os.Getenv(ArrebolPrivKeyPath))
	if err != nil {
		return nil, err
	}
	tokenStr, err := token.SignedString(privateKey)
	if err != nil {
		return &JWToken{}, err
	}
	return &JWToken{
		raw:      tokenStr,
		QueueId:  worker.QueueId,
		WorkerId: worker.ID,
	}, nil
}

func Parse(tokenString string) (*JWToken, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return crypto.GetPublicKey(os.Getenv(ArrebolPubKeyPath))
	})
	if err != nil {
		return nil, err
	}
	log.Println(token)
	return nil, nil
}

func (t *JWToken) String() string {
	return t.raw
}

func (*JWToken) Expired() bool {
	return false
}

func (*JWToken) GetPayload(key string) string {
	return ""
}

