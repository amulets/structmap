package cast

import (
	"errors"
	"reflect"

	"github.com/amulets/structmap/internal"
)

var (
	ErrEmptyValue    = errors.New("empty value")
	ErrNoCoveredType = errors.New("no covered type")
	ErrNoConvertible = errors.New("non-convertible type")
)

func (c *Cast) toType(from reflect.Type, value reflect.Value) (to reflect.Value, err error) {
	// set received value to return
	to = value

	value = internal.Value(value, true)

	if !value.IsValid() {
		// Do not has a value
		to = reflect.Zero(from)
		err = ErrEmptyValue
		return
	}

	// check if type is a ptr and get real type
	fromType := internal.Type(from)

	convert, ok := c.convertToKind[ToKind(fromType)]
	if !ok {
		err = ErrNoCoveredType
		return
	}

	var rawValue interface{}
	if rawValue, err = convert(fromType, value); err == nil {
		to = reflect.ValueOf(rawValue)

		if !to.IsValid() {
			// TODO: return an error or instance a zero value (to = reflect.Zero(fromType))
			// skipping ConvertibleTo check
			err = errors.New("you need to return an error or a valid result")
			return
		}

		// convert its type ex: int32 to int
		if to.Type().ConvertibleTo(fromType) {
			to = to.Convert(fromType)
		}

		if from.Kind() == reflect.Ptr {
			pv := reflect.New(from.Elem())
			internal.SetValue(pv.Elem(), to)

			to = pv
		}
	}

	return
}
