package cast

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/amulets/structmap"
)

var (
	errNoCoveredType = errors.New("no covered type")
	errNoConvertible = errors.New("non-convertible type")
)

type typeConvert = map[reflect.Kind]func(source reflect.Type, value reflect.Value) (result interface{}, err error)

// Covered types to convert
var convertTo typeConvert

func init() {
	convertTo = typeConvert{
		reflect.String:  toString,
		reflect.Int:     toInt,
		reflect.Uint:    toUint,
		reflect.Bool:    toBool,
		reflect.Float32: toFloat,
		reflect.Map:     toMap,
		reflect.Struct:  toStruct,
	}
}

// ToType cast value to field type value
func ToType(field *structmap.FieldPart) error {
	value := reflect.Indirect(reflect.ValueOf(field.Value))
	if !value.IsValid() {
		// Do not has a value
		return nil
	}

	toValue, err := toType(field.Type, value)
	if err != nil {
		switch err {
		case errNoCoveredType:
			// Do not have a coverter to this type
			return nil
		case errNoConvertible:
			err = fmt.Errorf("'%s' expected type '%s', got non-convertible type '%s'", field.Name, field.Type, value.Type())
			fallthrough
		default:
			return err
		}
	}

	// If null, is the same type
	if toValue != nil {
		field.Value = toValue
	}

	return nil
}

func toType(from reflect.Type, value reflect.Value) (to interface{}, err error) {
	convert, ok := convertTo[toKind(from)]
	if !ok {
		err = errNoCoveredType
		return
	}

	return convert(from, value)
}
