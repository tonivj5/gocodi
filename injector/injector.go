package injector

import (
	"fmt"
	"reflect"

	"github.com/xxxtonixxx/gocodi/provider"
	"github.com/xxxtonixxx/gocodi/providers"
	"github.com/xxxtonixxx/gocodi/tokens"
	"github.com/xxxtonixxx/gocodi/utils"
)

type Injector struct {
	deps   []providers.Provider
	parent *Injector
}

// NewChildInjector() *Injector
func (injector *Injector) Provide(provider *provider.Provider) error {
	if provider.Provide == nil {
		return fmt.Errorf("You must set a kind of Provider and its use(Value/Factory)/aliasOf")
	}

	var newProvider providers.Provider
	token := tokens.New(provider.Provide)

	switch token.GetType() {
	case tokens.Interface:
		if provider.UseValue == nil && provider.UseFactory == nil {
			return fmt.Errorf("You must provide a value or factory when token is %s", reflect.TypeOf(provider.Provide))
		}

		if !utils.ImplementsInterface(provider.Provide, provider.UseValue) {
			return fmt.Errorf(
				"%s does not implements %s",
				reflect.TypeOf(provider.UseValue),
				reflect.TypeOf(provider.Provide),
			)
		}

		newProvider = &providers.Interface{
			Token: token,
			Value: provider.UseValue,
		}
	case tokens.PtrStruct, tokens.Struct:
		// if provider.UseValue != nil {
		// 	if !utils.IsPtrToStruct(provider.UseValue) && !utils.IsStruct(provider.UseValue) {
		// 		return fmt.Errorf("When provide an struct or pointer to one, you must return an struct of same type")
		// 	}
		// } else if provider.UseFactory != nil {

		// } else {
		// 	return fmt.Errorf(`
		// 		If you provide a struct or pointer to one,
		// 		you must set UseValue or UseFactory
		// 		`)
		// }
		if provider.UseValue != nil &&
			(!utils.IsPtrToStruct(provider.UseValue) || !utils.IsStruct(provider.UseValue)) {
			return fmt.Errorf("If you provide a struct, the value must be a struct")
		}

		newProvider = &providers.Struct{
			Token:     token,
			Value:     provider.UseValue,
			IsCreated: provider.UseValue != nil,
		}
	case tokens.String:
		if provider.UseValue == nil && provider.UseFactory == nil {
			return fmt.Errorf("You must provide a value or factory when token is %s", reflect.TypeOf(provider.Provide))
		}

		newProvider = &providers.Primitive{
			Token: token,
			Value: provider.UseValue,
		}
	}

	injector.deps = append(injector.deps, newProvider)

	return nil
}

func (injector *Injector) Get(token interface{}) interface{} {
	for _, dep := range injector.deps {
		if dep.Match(token) {
			value := dep.GetValue()

			tokenType := dep.GetToken().GetType()
			switch tokenType {
			case tokens.Struct:
				typeOf := reflect.TypeOf(dep.GetToken().Get())
				valueOf := reflect.ValueOf(value)

				injector.deepInit(typeOf, valueOf)
			case tokens.PtrStruct:
				typeOf := reflect.TypeOf(dep.GetToken().Get()).Elem()
				valueOf := reflect.ValueOf(value).Elem()

				injector.deepInit(typeOf, valueOf)
			}

			return value
		}
	}

	return nil
}

func (injector *Injector) deepInit(typeOf reflect.Type, valueOf reflect.Value) {
	for i := 0; i < valueOf.NumField(); i++ {
		fieldValue := valueOf.Field(i)
		fieldType := typeOf.Field(i)
		tags := fieldType.Tag.Get("gocodi")

		if !fieldValue.CanSet() || tags == "-" {
			continue
		}

		var dep interface{}
		switch fieldType.Type.Kind() {
		case reflect.Interface:
			token := reflect.New(fieldType.Type).Interface()
			dep = injector.Get(token)
		case reflect.Struct:
			dep = injector.Get(fieldType.Type)
		case reflect.Ptr:
			if tags == "" {
				token := reflect.New(fieldType.Type.Elem()).Interface()
				dep = injector.Get(token)
			} else {
				dep = injector.Get(tags)
			}
		default:
			if tags != "" {
				dep = injector.Get(tags)
			}
		}

		if reflect.ValueOf(dep).IsValid() {
			fieldValue.Set(reflect.ValueOf(dep))
		} else {
			panic(fmt.Errorf(
				"%s can not be provided as dependency when resolved %s",
				fieldType.Type,
				typeOf,
			))
		}
	}
}

func New() *Injector {
	return &Injector{
		deps: make([]providers.Provider, 0),
	}
}

func NewWithParent(parent *Injector) *Injector {
	return &Injector{
		parent: parent,
		deps:   make([]providers.Provider, 0),
	}
}
