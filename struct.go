package structmap

import (
	"reflect"
)

type strct struct {
	reflect.Value
}

func newStruct(s interface{}) (*strct, error) {
	value := reflect.ValueOf(s)

	// if pointer get element
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return nil, ErrNotIsStruct
	}

	return &strct{value}, nil
}

// Fields of struct
func (s *strct) Fields() []field {
	t := s.Value.Type()

	fields := make([]field, 0)

	for i := 0; i < t.NumField(); i++ {
		field := field{t.Field(i), s.Value.Field(i)}

		// we can't access the value of unexported fields
		if !field.IsExported() {
			continue
		}

		fields = append(fields, field)
	}

	return fields
}
