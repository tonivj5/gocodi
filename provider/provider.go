package provider

type Provider struct {
	Provide    interface{}
	UseValue   interface{}
	UseFactory interface{}
	deps       []interface{}
	aliasOf    interface{}
}
