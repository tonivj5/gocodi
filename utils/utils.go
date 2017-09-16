package utils

import "reflect"

func IsString(value interface{}) bool {
	return reflect.TypeOf(value).Kind() == reflect.String
}

func IsFunc(fn interface{}) bool {
	return reflect.TypeOf(fn).Kind() == reflect.Func
}
