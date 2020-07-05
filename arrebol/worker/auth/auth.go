package auth

import (
	"github.com/google/logger"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/allowlist"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/token"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/tolerant"
	"os"
	"strconv"
)

const AllowListConfKey = "ALLOW_ALL"

type Authenticator interface {
	Authenticate(rawPublicKey string, signature []byte, worker *worker.Worker) (token.Token, error)
	Authorize(token *token.Token) error
}

func NewAuthenticator() Authenticator {
	allow, err := strconv.ParseBool(os.Getenv(AllowListConfKey))

	if err != nil {
		logger.Fatalf("Cannot understand the flag: %s", err.Error())
	}
	if allow {
		return allowlist.NewAuthenticator()
	}
	return tolerant.NewAuthenticator()
}
