package structmap

import (
	"reflect"
)

type Struct struct {
	reflect.Value
}

func NewStruct(s interface{}) (*Struct, error) {
	value := reflect.ValueOf(s)

	// if pointer get element
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return nil, ErrNotIsStruct
	}

	return &Struct{value}, nil
}

// Fields of struct
func (s *Struct) Fields() []Field {
	t := s.Value.Type()

	fields := make([]Field, 0)

	for i := 0; i < t.NumField(); i++ {
		field := Field{t.Field(i), s.Value.Field(i)}

		// we can't access the value of unexported fields
		if !field.IsExported() {
			continue
		}

		fields = append(fields, field)
	}

	return fields
}
