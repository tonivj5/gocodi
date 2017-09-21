package providers

import (
	"github.com/xxxtonixxx/gocodi/tokens"
)

type Primitive struct {
	Token     *tokens.Token
	Value     interface{}
	IsFactory bool
}

func (provider *Primitive) GetValue() interface{} {
	if provider.IsFactory {
		// fmt.Println("Is factory")
	}

	return provider.Value
}

func (provider *Primitive) GetToken() *tokens.Token {
	return provider.Token
}

func (provider *Primitive) Match(token interface{}) bool {
	return provider.Token.Match(tokens.New(token))
}
