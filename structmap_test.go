package structmap_test

import (
	"reflect"
	"testing"

	"github.com/dungeon-code/structmap"
	"github.com/dungeon-code/structmap/mutation"
	"github.com/dungeon-code/structmap/mutation/cast"
)

type SubSubStruct struct {
	Address string
	Number  *int
}

type SubStruct struct {
	*SubSubStruct
	Age *string `default:"15"`
}

type MyStruct struct {
	SubStruct
	MyAddress SubSubStruct `structmap:"myAddress"`
	Name      *string      `structmap:"name,omitempty"`
	Username  string       `structmap:"user"`
	UserNames []string
	MyBool    bool
	MyUint    uint32
	MyFloat   float32
	MyMap     map[string]interface{}
}

func TestDecode(t *testing.T) {
	s := &MyStruct{}
	m := map[string]interface{}{
		"name":      "Marisa",
		"user":      "{{name}}",
		"UserNames": []string{"A", "B", "C"},
		// "Age":       18,
		"Address": "Street A",
		"Number":  "1832",
		"myAddress": map[string]interface{}{
			"Address": "Street B",
			"Number":  1345,
		},
		"MyBool":  1,
		"MyUint":  true,
		"MyFloat": false,
		"MyMap": map[string]interface{}{
			"key": "value",
		},
	}

	d := structmap.New()
	d.AddMutation(mutation.SetFieldName("structmap"))
	d.AddMutation(func(field *structmap.FieldPart) error {
		name, _ := structmap.ParseTag(field.Tag.Get("default"))
		if field.Value == nil && name != "" {
			field.Value = name
		}

		return nil
	})
	d.AddMutation(cast.ToType)

	if err := d.Decode(m, s); err != nil {
		t.Error(err)
	}

	// &{SubStruct:{SubSubStruct:0xc000010240 Age:18} MyAddress:{Address:Street B} Name:0xc0000102a0 Username:{{name}} UserNames:[A B C]}
	name := "Marisa"
	age := "15"
	n1 := 1832
	n2 := 1345

	expected := &MyStruct{
		SubStruct: SubStruct{
			SubSubStruct: &SubSubStruct{
				Address: "Street A",
				Number:  &n1,
			},
			Age: &age,
		},
		MyAddress: SubSubStruct{
			Address: "Street B",
			Number:  &n2,
		},
		Name:      &name,
		Username:  "{{name}}",
		UserNames: []string{"A", "B", "C"},
		MyBool:    true,
		MyUint:    1,
		MyFloat:   0,
		MyMap: map[string]interface{}{
			"key": "value",
		},
	}

	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Expected = %+v; got = %+v", expected, s)
	}
}
