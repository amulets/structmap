package structmap_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/dungeon-code/structmap"
	"github.com/dungeon-code/structmap/mutations"
)

type SubSubStruct struct {
	Address string
}

type SubStruct struct {
	*SubSubStruct
	Age string
}

type MyStruct struct {
	SubStruct
	MyAddress SubSubStruct `structmap:"myAddress"`
	Name      *string      `structmap:"name,omitempty"`
	Username  string       `structmap:"user"`
	UserNames []string
}

func intToString(field *structmap.FieldPart) error {
	fieldValue := reflect.ValueOf(field.Value)

	if field.Type.Kind() == reflect.String && fieldValue.Kind() == reflect.Int {
		field.Value = strconv.Itoa(field.Value.(int))
	}

	return nil
}

func TestDecode(t *testing.T) {
	s := &MyStruct{}
	m := map[string]interface{}{
		"name":      "Marisa",
		"user":      "{{name}}",
		"UserNames": []string{"A", "B", "C"},
		"Age":       18,
		"Address":   "Street A",
		"myAddress": map[string]interface{}{
			"Address": "Street B",
		},
	}

	d := structmap.NewDecoder()
	d.AddMutation(mutations.SetFieldName("structmap"))
	d.AddMutation(intToString)

	if err := d.Decode(m, s); err != nil {
		t.Error(err)
	}

	// &{SubStruct:{SubSubStruct:0xc000010240 Age:18} MyAddress:{Address:Street B} Name:0xc0000102a0 Username:{{name}} UserNames:[A B C]}
	name := "Marisa"

	expected := &MyStruct{
		SubStruct: SubStruct{
			SubSubStruct: &SubSubStruct{
				Address: "Street A",
			},
			Age: "18",
		},
		MyAddress: SubSubStruct{
			Address: "Street B",
		},
		Name:      &name,
		Username:  "{{name}}",
		UserNames: []string{"A", "B", "C"},
	}

	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Expected = %+v; got = %+v", expected, s)
	}
}
