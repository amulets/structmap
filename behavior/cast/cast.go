package cast

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/amulets/structmap"
	"github.com/amulets/structmap/behavior"
	"github.com/amulets/structmap/internal"
)

var (
	errEmptyValue    = errors.New("empty value")
	errNoCoveredType = errors.New("no covered type")
	errNoConvertible = errors.New("non-convertible type")
)

type kindConvert = map[reflect.Kind]func(source reflect.Type, value reflect.Value) (result interface{}, err error)

// Covered types to convert
var convertTo kindConvert

func init() {
	convertTo = kindConvert{
		reflect.String:    toString,
		reflect.Int:       toInt,
		reflect.Uint:      toUint,
		reflect.Bool:      toBool,
		reflect.Float32:   toFloat,
		reflect.Map:       toMap,
		reflect.Struct:    toStruct,
		reflect.Slice:     toList,
		reflect.Array:     toList,
		reflect.Interface: toInterface,
	}
}

// ToType cast value to field type value
var ToType = behavior.New(func(field *structmap.FieldPart) error {
	//Field value cannot is a pointer
	value := reflect.ValueOf(field.Value)

	toValue, err := toType(field.Type, value)
	if err != nil {
		switch err {
		case errEmptyValue:
			fallthrough
		case errNoCoveredType:
			// Do not have a coverter to this type
			return nil
		case errNoConvertible:
			err = fmt.Errorf("'%s' expected type '%s', got non-convertible type '%s'", field.Name, field.Type, value.Type())
			fallthrough
		default:
			err = fmt.Errorf("%s: %s", field.Name, err)
			return err
		}
	}

	field.Value = toValue.Interface()

	return nil
})

func toType(from reflect.Type, value reflect.Value) (to reflect.Value, err error) {
	// set received value to return
	to = value

	value = internal.Value(value, true)

	if !value.IsValid() {
		// Do not has a value
		to = reflect.Zero(from)
		err = errEmptyValue
		return
	}

	// check if type is a ptr and get real type
	fromType := internal.Type(from)

	convert, ok := convertTo[toKind(fromType)]
	if !ok {
		err = errNoCoveredType
		return
	}

	var rawValue interface{}
	if rawValue, err = convert(fromType, value); err == nil {
		to = reflect.ValueOf(rawValue)

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
