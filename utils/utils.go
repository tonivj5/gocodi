package utils

import "reflect"

func IsString(value interface{}) bool {
	return reflect.TypeOf(value).Kind() == reflect.String
}

func IsPrimitive(value interface{}) bool {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Struct,
		reflect.Ptr,
		reflect.Interface:

		return false
	default:
		return true
	}
}

func IsFunc(fn interface{}) bool {
	return reflect.TypeOf(fn).Kind() == reflect.Func
}

func IsPtrToInterface(value interface{}) bool {
	typeOf := reflect.TypeOf(value)

	return typeOf.Kind() == reflect.Ptr && typeOf.Elem().Kind() == reflect.Interface
}

func IsPtrToStruct(value interface{}) bool {
	typeOf := reflect.TypeOf(value)

	return typeOf.Kind() == reflect.Ptr && typeOf.Elem().Kind() == reflect.Struct
}

func IsStruct(value interface{}) bool {
	typeOf := reflect.TypeOf(value)

	return typeOf.Kind() == reflect.Struct
}

func IsSetted(value interface{}) bool {
	return reflect.ValueOf(value).IsValid()
}

func IsSameType(typeOf interface{}, value interface{}) bool {
	return reflect.TypeOf(value).AssignableTo(reflect.TypeOf(typeOf))
}

func ImplementsInterface(interfaceMustImplemented interface{}, value interface{}) bool {
	typeOfInterface := reflect.TypeOf(interfaceMustImplemented).Elem()

	return reflect.TypeOf(value).Implements(typeOfInterface)
}

func FuncReturnsSameType(typeOf interface{}, fn interface{}) bool {
	fnType := reflect.TypeOf(fn)

	return fnType.NumOut() == 1 && fnType.Out(0) == reflect.TypeOf(typeOf)
}
