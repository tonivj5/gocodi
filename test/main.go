package main

import (
	"fmt"

	di "github.com/xxxtonixxx/gocodi"
	"github.com/xxxtonixxx/gocodi/provider"
)

type Salutor interface {
	Hi() string
}

type Test struct {
	IP      string `gocodi:"ip"`
	test    string
	DepTest *TestDep `gocodi:"test"`
	Dep     *TestDep
	Saluto  Salutor
}

func (h *Test) Hi() string {
	return "Hi DI!"
}

func (h *Test) getIP() string {
	return h.IP
}

type TestDep struct {
	x int64
}

func main() {
	// You must set Provide property. If an error happened it will be returned
	err := di.Provide(&provider.Provider{})
	if err != nil {
		fmt.Printf("An error happened: %v\n", err)
	}
	// When you Get &Test{} to injector, it will create a new instance and return you
	di.Provide(&provider.Provider{Provide: new(Test)})
	// You can use a name as token and return a value
	di.Provide(&provider.Provider{Provide: "hi", UseValue: "I am a provider value using a string as token"})
	di.Provide(
		&provider.Provider{Provide: new(TestDep)},
	)
	di.Provide(
		&provider.Provider{Provide: "test", UseValue: &TestDep{x: 50}},
	)

	// You can set a int/string/map/array and whatever value using a string token
	di.Provide(
		&provider.Provider{Provide: "ip", UseValue: "192.168.1.1"},
	)

	// You can use an interface as provider and an struct
	// which implements the interface as value
	var token *Salutor
	err = di.Provide(&provider.Provider{Provide: token, UseValue: &Test{IP: "My IP"}})
	if err != nil {
		fmt.Printf("An error happenedd: %v\n", err)
		return
	}

	// So you can get the dependency from injector.
	testDI := di.Get(&Test{}).(*Test)
	fmt.Println(
		"-->",
		di.Get("hi"),
		testDI.DepTest,
		testDI.Dep,
		testDI.Hi(),
		testDI.getIP(),
		di.Get("test"),
		testDI != testDI.Saluto,
		testDI.Saluto,
		testDI.Saluto.Hi(),
	)

	type MyStruct struct {
		fieldOne string
	}

	// You can provide a struct (no pointer)
	err = di.Provide(&provider.Provider{Provide: MyStruct{}, UseValue: MyStruct{fieldOne: "hi"}})
	if err != nil {
		fmt.Printf("An error happenedd: %v\n", err)
		return
	}

	testMyStruct := di.Get(MyStruct{}).(MyStruct)
	fmt.Println(
		"2---->",
		testMyStruct,
		testMyStruct.fieldOne,
	)
}
