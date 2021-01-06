package cast

import (
	"fmt"
	"reflect"

	"github.com/amulets/structmap"
)

type (
	// ConverterFunc has logic on how to convert
	ConverterFunc func(source reflect.Type, value reflect.Value) (result interface{}, err error)

	// OptionFunc define options
	OptionFunc func(*Cast)

	// Cast is a cast
	Cast struct {
		convertToKind map[reflect.Kind]ConverterFunc
		convertToType map[reflect.Type]ConverterFunc
	}
)

// ToType cast value to field type value
func ToType(options ...OptionFunc) *Cast {
	c := new(Cast)

	c.convertToKind = map[reflect.Kind]ConverterFunc{
		reflect.String:    c.toString,
		reflect.Int:       c.toInt,
		reflect.Uint:      c.toUint,
		reflect.Bool:      c.toBool,
		reflect.Float32:   c.toFloat,
		reflect.Map:       c.toMap,
		reflect.Struct:    c.toStruct,
		reflect.Slice:     c.toList,
		reflect.Array:     c.toList,
		reflect.Interface: c.toInterface,
	}

	c.convertToType = map[reflect.Type]ConverterFunc{}

	for _, option := range options {
		option(c)
	}

	return c
}

// Do has a ToType logic
func (c *Cast) Do(field *structmap.FieldPart) error {
	//Field value cannot is a pointer
	value := reflect.ValueOf(field.Value)

	toValue, err := c.toType(field.Type, value)
	if err != nil {
		switch err {
		case ErrEmptyValue:
			fallthrough
		case ErrNoCoveredType:
			// Do not have a coverter to this type
			return nil
		case ErrNoConvertible:
			err = fmt.Errorf("expected type '%s', got non-convertible type '%s'", field.Type, value.Type())
			fallthrough
		default:
			err = fmt.Errorf("%s: %s", field.Name, err)
			return err
		}
	}

	field.Value = toValue.Interface()

	return nil
}
