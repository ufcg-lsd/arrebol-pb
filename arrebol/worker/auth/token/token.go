package token

type Token interface {
	GetPayloadField(key string) (interface{}, error)
	SetPayloadField(key string, value interface{}) (Token, error)
	Expired() bool
	String() string
}



