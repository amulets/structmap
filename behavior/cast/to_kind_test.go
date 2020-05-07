package cast

import (
	"reflect"
	"testing"
)

func TestToKind(t *testing.T) {
	values := []struct {
		raw  interface{}
		kind reflect.Kind
	}{
		{
			raw:  "Lorem",
			kind: reflect.String,
		},
		{
			raw:  false,
			kind: reflect.Bool,
		},
		{
			raw:  struct{}{},
			kind: reflect.Struct,
		},
		{
			raw:  int(1),
			kind: reflect.Int,
		},
		{
			raw:  int8(1),
			kind: reflect.Int,
		},
		{
			raw:  int16(1),
			kind: reflect.Int,
		},
		{
			raw:  int32(1),
			kind: reflect.Int,
		},
		{
			raw:  int64(1),
			kind: reflect.Int,
		},
		{
			raw:  uint(1),
			kind: reflect.Uint,
		},
		{
			raw:  uint8(1),
			kind: reflect.Uint,
		},
		{
			raw:  uint16(1),
			kind: reflect.Uint,
		},
		{
			raw:  uint32(1),
			kind: reflect.Uint,
		},
		{
			raw:  uint64(1),
			kind: reflect.Uint,
		},
		{
			raw:  float32(1),
			kind: reflect.Float32,
		},
		{
			raw:  float64(1),
			kind: reflect.Float32,
		},
	}

	for _, value := range values {
		typeOf := reflect.TypeOf(value.raw)
		typeOfPtr := reflect.PtrTo(typeOf)

		for _, typ := range []reflect.Type{typeOf, typeOfPtr} {
			if kind := toKind(typ); kind != value.kind {
				t.Errorf("expected %s; got %s", value.kind, kind)
			}
		}
	}
}
