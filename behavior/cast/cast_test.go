package cast

import (
	"reflect"
	"testing"
)

func TestToType(t *testing.T) {
	stringType := reflect.TypeOf("")

	v, err := toType(reflect.PtrTo(reflect.PtrTo(stringType)), reflect.ValueOf(100))
	if err != nil {
		t.Error(err)
	}

	t.Log(v.Elem().Elem().String())
}
