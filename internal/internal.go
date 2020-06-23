package internal

import (
	"reflect"
)

// Type get real type
func Type(value reflect.Type) reflect.Type {
	if value.Kind() == reflect.Ptr {
		value = Type(value.Elem())
	}

	return value
}

// Value get real value
func Value(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Ptr {
		value = Value(value.Elem())
	}

	return value
}

// SetValue in a ptr
func SetValue(target reflect.Value, value reflect.Value) {
	if target.Kind() == reflect.Ptr {
		target.Set(reflect.New(target.Type().Elem()))

		SetValue(target.Elem(), value)
	} else {
		target.Set(value)
	}
}
