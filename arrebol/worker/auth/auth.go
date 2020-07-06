package auth

import (
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/authenticator"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/authorizer"
)

type Auth struct {
	Authorizer    authorizer.Authorizer
	Authenticator authenticator.Authenticator
}

func NewAuth() *Auth {
	return &Auth{
		Authorizer:    NewAuthorizer(),
		Authenticator: NewAuthenticator(),
	}
}

func NewAuthorizer() authorizer.Authorizer {
	return authorizer.NewAuthorizer()
}

func NewAuthenticator() authenticator.Authenticator {
	return authenticator.NewAuthenticator()
}
