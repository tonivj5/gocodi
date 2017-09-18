package container

import (
	"fmt"
	"reflect"

	"github.com/xxxtonixxx/gocodi/utils"
)

type Container struct {
	deps []Dependency
}

func (c *Container) Provide(provider *Provider) error {
	tokenType := reflect.TypeOf(provider.Provide)
	tokenValue := reflect.ValueOf(provider.Provide)
	valueType := reflect.TypeOf(provider.Value)
	valueValue := reflect.ValueOf(provider.Value)

	if !tokenValue.IsValid() {
		return fmt.Errorf("You must set Provide field")
	}

	useName := utils.IsString(provider.Provide)
	isInterface := utils.IsPtrToInterface(provider.Provide)
	// useFunction := utils.IsString(provider.Provide)
	if isInterface &&
		(!valueValue.IsValid() || !valueType.Implements(tokenType.Elem())) {
		return fmt.Errorf("%s does not implements %s", valueType, tokenType)
	} else if !useName && !isInterface &&
		(valueValue.IsValid() && !valueType.AssignableTo(tokenType)) {
		return fmt.Errorf("%s is not assignable to %s", valueType, tokenType)
	}

	c.deps = append(c.deps, Dependency{
		token:       provider.Provide,
		value:       provider.Value,
		useName:     useName,
		isInterface: isInterface,
	})

	return nil
}

func (c *Container) Get(token interface{}) interface{} {
	useName := utils.IsString(token)
	var name string

	if useName {
		name = token.(string)
	}

	for _, dep := range c.deps {
		if !dep.useName && useName {
			continue
		}

		var matchType bool
		if dep.useName {
			matchType = dep.token.(string) == name
		} else if dep.isInterface {
			typeOfProvided := reflect.TypeOf(dep.token).Elem()
			matchType = typeOfProvided == reflect.TypeOf(token).Elem()
		} else {
			typeOfProvided := reflect.TypeOf(dep.token)
			matchType = typeOfProvided.AssignableTo(reflect.TypeOf(token))

		}

		if matchType {
			if !dep.isInterface && !reflect.ValueOf(dep.value).IsValid() {
				dep.value = reflect.New(reflect.TypeOf(dep.token).Elem()).Interface()
				c.initializeDep(reflect.TypeOf(dep.token).Elem(), reflect.ValueOf(dep.value).Elem())
			}

			return dep.value
		}
	}

	return nil
}

func (c *Container) initializeDep(t reflect.Type, v reflect.Value) interface{} {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)
		tags := ft.Tag.Get("gocodi")

		if !f.CanSet() || tags == "-" {
			continue
		}

		var dep interface{}
		switch ft.Type.Kind() {
		// case reflect.Map:
		// 	fallthrough
		// case reflect.Slice:
		// 	fallthrough
		// case reflect.Chan:
		// 	fallthrough
		case reflect.Interface:
			dep := reflect.ValueOf(c.Get(reflect.New(ft.Type).Interface()))
			f.Set(dep)
		case reflect.Struct:
			dep := reflect.ValueOf(c.Get(ft.Type))
			f.Set(dep)
		case reflect.Ptr:
			if tags != "" {
				dep = c.Get(tags)
			} else {
				subdep := reflect.New(ft.Type.Elem()).Interface()
				dep = c.Get(subdep)
			}
		default:
			if tags != "" {
				dep = c.Get(tags)
			}
		}

		if reflect.ValueOf(dep).IsValid() {
			f.Set(reflect.ValueOf(dep))
		}
	}

	return v.Interface()
}

func New() *Container {
	return &Container{
		deps: make([]Dependency, 0),
	}
}

type Dependency struct {
	token       interface{}
	value       interface{}
	useName     bool
	useFactory  bool
	isInterface bool
}

type Provider struct {
	Provide interface{}
	Value   interface{}
	Factory interface{}
}
