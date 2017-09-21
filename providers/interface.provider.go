package providers

import (
	"github.com/xxxtonixxx/gocodi/tokens"
)

type Interface struct {
	Token     *tokens.Token
	Value     interface{}
	IsFactory bool
}

func (provider *Interface) GetValue() interface{} {
	if provider.IsFactory {
		// fmt.Println("Is factory")
	}

	return provider.Value
}

func (provider *Interface) GetToken() *tokens.Token {
	return provider.Token
}

func (provider *Interface) Match(token interface{}) bool {
	return provider.Token.Match(tokens.New(token))
}
