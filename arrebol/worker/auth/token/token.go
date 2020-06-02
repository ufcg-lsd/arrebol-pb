package token

type Token interface {
	GetPayload(key string) string
	Expired() bool
	String() string
}



