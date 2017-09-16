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
	if !reflect.ValueOf(provider.Provide).IsValid() {
		return fmt.Errorf("You must set Provide field")
	}

	useName := utils.IsString(provider.Provide)
	// useFunction := utils.IsString(provider.Provide)
	if !useName && reflect.ValueOf(provider.Value).IsValid() {
		typeOfProvided := reflect.TypeOf(provider.Provide)
		typeOfValue := reflect.TypeOf(provider.Value)

		if !typeOfValue.AssignableTo(typeOfProvided) {
			return fmt.Errorf("%s is not assignable to %s", typeOfValue, typeOfProvided)
		}

	}

	c.deps = append(c.deps, Dependency{
		token:   provider.Provide,
		value:   provider.Value,
		useName: useName,
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
		var matchType bool
		if !dep.useName {
			typeOfProvided := reflect.TypeOf(dep.token)
			matchType = typeOfProvided.AssignableTo(reflect.TypeOf(token))
		} else if dep.token.(string) == name {
			matchType = true
		}

		if matchType {
			if !reflect.ValueOf(dep.value).IsValid() {
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

		if !f.CanSet() {
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
		case reflect.Struct:
			dep := reflect.ValueOf(c.Get(ft.Type))
			f.Set(dep)
		case reflect.Ptr:
			useDep := ft.Tag.Get("gocodi")
			if useDep == "" {
				subdep := reflect.New(ft.Type.Elem()).Interface()
				dep = c.Get(subdep)
			} else {
				dep = c.Get(useDep)
			}
		default:
			useDep := ft.Tag.Get("gocodi")

			if useDep != "" {
				dep = c.Get(useDep)
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
	token      interface{}
	value      interface{}
	useName    bool
	useFactory bool
}

type Provider struct {
	Provide interface{}
	Value   interface{}
	Factory interface{}
}
