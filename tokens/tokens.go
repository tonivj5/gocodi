package tokens

import (
	"reflect"

	"github.com/xxxtonixxx/gocodi/utils"
)

type TokenType int

const (
	String TokenType = iota
	Interface
	Struct
	PtrStruct
	Unkown
)

type Token struct {
	Value interface{}
}

func (t *Token) Get() interface{} {
	return t.Value
}

func (t *Token) GetType() TokenType {
	if utils.IsPtrToInterface(t.Value) {
		return Interface
	} else if utils.IsPtrToStruct(t.Value) {
		return PtrStruct
	} else if utils.IsString(t.Value) {
		return String
	} else if utils.IsStruct(t.Value) {
		return Struct
	}

	return Unkown
}

func (t *Token) Match(token *Token) bool {
	if token.GetType() != t.GetType() {
		return false
	}

	switch token.GetType() {
	case Struct, PtrStruct:
		return reflect.TypeOf(t.Value) == reflect.TypeOf(token.Get())
	case Interface:
		return reflect.TypeOf(t.Value).Elem() == reflect.TypeOf(token.Get()).Elem()
	case String:
		return t.Get().(string) == token.Get().(string)
	}

	return false
}

func New(token interface{}) *Token {
	return &Token{
		Value: token,
	}
}
