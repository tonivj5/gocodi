package main

import (
	"fmt"

	"github.com/xxxtonixxx/gocodi/container"
)

type Test struct {
	IP      string `gocodi:"ip"`
	test    string
	DepTest *TestDep `gocodi:"test"`
	Dep     *TestDep
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
	di := container.New()
	// You must set Provide property. If an error happened it will be returned
	err := di.Provide(&container.Provider{})
	if err != nil {
		fmt.Printf("An error happened: %v\n", err)
	}
	// When you Get &Test{} to injector, it will create a new instance and return you
	di.Provide(&container.Provider{Provide: new(Test)})
	// You can use a name as token provider and set a value
	di.Provide(&container.Provider{Provide: "hi", Value: &Test{test: "testing DI!"}})
	di.Provide(
		&container.Provider{Provide: new(TestDep)},
	)
	di.Provide(
		&container.Provider{Provide: "test", Value: &TestDep{x: 50}},
	)

	// You can set a int/string/map/array and whatever value using a string token
	di.Provide(
		&container.Provider{Provide: "ip", Value: "192.168.1.1"},
	)

	// So you can get the dependency from injector.
	testDI := di.Get(&Test{}).(*Test)
	fmt.Println("------", testDI.DepTest, testDI.Dep, testDI.getIP())
}
