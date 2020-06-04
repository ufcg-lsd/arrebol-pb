package token

type Token interface {
	GetPayload(key string) (interface{}, error)
	SetPayload(key string, value interface{}) (Token, error)
	Expired() bool
	String() string
}



