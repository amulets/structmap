package cast

import "reflect"

// TypeFunc returns reflect.Type and ConverterFunc to this type
type TypeFunc func() (reflect.Type, ConverterFunc)

// Type creates a TypeFunc for this type
func Type(value interface{}, convert ConverterFunc) TypeFunc {
	return func() (reflect.Type, ConverterFunc) {
		return reflect.TypeOf(value), convert
	}
}

// WithTypes packages everything in the cast
func WithTypes(types ...TypeFunc) OptionFunc {
	return func(c *Cast) {
		for _, pack := range types {
			typeOf, convert := pack()
			c.convertToType[typeOf] = convert
		}
	}
}
