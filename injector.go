package gocodi

import (
	"github.com/xxxtonixxx/gocodi/injector"
	"github.com/xxxtonixxx/gocodi/provider"
)

var mainInjector = injector.New()

func Provide(provider *provider.Provider) error {
	return mainInjector.Provide(provider)
}

func Get(token interface{}) interface{} {
	return mainInjector.Get(token)
}

func New() *injector.Injector {
	return injector.NewWithParent(mainInjector)
}
