package cast

import (
	"reflect"
	"time"
)

type typeConvertStruct = map[reflect.Type]func(source reflect.Type, value reflect.Value) (result interface{}, err error)

// Covered struct types to convert
var convertToStruct typeConvertStruct

func init() {
	convertToStruct = typeConvertStruct{
		reflect.TypeOf(time.Time{}): toStructTime,
	}
}

func toStruct(source reflect.Type, value reflect.Value) (result interface{}, err error) {
	// Is the same type
	if source == value.Type() {
		return
	}

	convert, ok := convertToStruct[source]
	if !ok {
		err = errNoCoveredType
		return
	}

	return convert(source, value)
}
