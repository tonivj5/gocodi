package providers

import (
	"reflect"

	"github.com/xxxtonixxx/gocodi/tokens"
)

type Struct struct {
	Token     *tokens.Token
	Value     interface{}
	IsCreated bool
	IsFactory bool
}

func (provider *Struct) GetValue() interface{} {
	if provider.IsCreated {
		return provider.Value
	}

	var token reflect.Type
	if provider.Token.GetType() == tokens.PtrStruct {
		token = reflect.TypeOf(provider.Token.Get()).Elem()
	} else {
		token = reflect.TypeOf(provider.Token.Get())
	}

	provider.Value = reflect.New(token).Interface()
	provider.IsCreated = true

	return provider.Value
}

func (provider *Struct) GetToken() *tokens.Token {
	return provider.Token
}

func (provider *Struct) Match(token interface{}) bool {
	return provider.Token.Match(tokens.New(token))
}
