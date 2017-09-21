package providers

import "github.com/xxxtonixxx/gocodi/tokens"

type Provider interface {
	GetValue() interface{}
	GetToken() *tokens.Token
	Match(token interface{}) bool
}
