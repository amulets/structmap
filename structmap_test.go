package structmap_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/amulets/structmap"
	"github.com/amulets/structmap/behavior/cast"
	"github.com/amulets/structmap/behavior/flag"
	"github.com/amulets/structmap/behavior/name"
)

type SubSubStruct struct {
	Address string
	Number  *int
}

type SubStruct struct {
	*SubSubStruct
	Age *string `structmap:",required" default:"15"`
}

type MyStruct struct {
	SubStruct
	MyAddress SubSubStruct `structmap:"myAddress"`
	Name      *string      `structmap:"name"`
	Username  string       `structmap:"user,required"`
	UserNames []string
	MyBool    bool
	MyUint    uint32
	MyFloat   float32
	MyMap     map[string]interface{}
	Headers   map[string]string `structmap:"headers"`
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
		"headers": map[string]interface{}{
			"a": "b",
			"b": "c",
		},
	}

	defaultTag := "structmap"

	d := structmap.New()
	d.AddBehavior(name.FromTag(defaultTag))
	d.AddBehavior(func(field *structmap.FieldPart) error {
		if field.Value != nil {
			return nil
		}

		value, _ := structmap.ParseTag(field.Tag.Get("default"))
		if value != "" {
			field.Value = value
		}

		return nil
	})
	d.AddBehavior(flag.Required(defaultTag))
	d.AddBehavior(cast.ToType)

	if err := d.Decode(m, s); err != nil {
		t.Error(err)
		t.FailNow()
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
		Headers: map[string]string{
			"a": "b",
			"b": "c",
		},
	}

	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Expected = %+v; got = %+v", expected, s)
	}
}

type DefaultTypes struct {
	Tstring    string `structmap:"tstring"`
	Tint       int    `structmap:"tint"`
	Tint8      int8   `structmap:"tint8"`
	Tint16     int16  `structmap:"tint16"`
	Tint32     int32  `structmap:"tint32"`
	Tint64     int64  `structmap:"tint64"`
	Tuint      uint
	Tbool      bool `structmap:"tbool"`
	Tfloat     float64
	unexported bool
	Tdata      interface{} `structmap:"tdata"`
}

type DefaultTypesPointer struct {
	Tstring    *string `structmap:"tstring"`
	Tint       *int    `structmap:"tint"`
	Tuint      *uint
	Tbool      *bool `structmap:"tbool"`
	Tfloat     *float64
	unexported *bool
	Tdata      *interface{} `structmap:"tdata"`
}

func TestDefaultTypes(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"tstring":    "foo",
		"tint":       20,
		"tint8":      20,
		"tint16":     20,
		"tint32":     20,
		"tint64":     20,
		"Tuint":      20,
		"tbool":      true,
		"Tfloat":     20.20,
		"unexported": true,
		"tdata":      20,
	}

	var result DefaultTypes

	sm := structmap.New()
	sm.AddBehavior(name.FromTag("structmap"))

	err := sm.Decode(input, &result)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if result.Tstring != "foo" {
		t.Errorf("tstring value should be 'foo': %#v", result.Tstring)
	}

	if result.Tint != 20 {
		t.Errorf("tint value should be 20: %#v", result.Tint)
	}
	if result.Tint8 != 20 {
		t.Errorf("tint8 value should be 20: %#v", result.Tint)
	}
	if result.Tint16 != 20 {
		t.Errorf("tint16 value should be 20: %#v", result.Tint)
	}
	if result.Tint32 != 20 {
		t.Errorf("tint32 value should be 20: %#v", result.Tint)
	}
	if result.Tint64 != 20 {
		t.Errorf("tint64 value should be 20: %#v", result.Tint)
	}

	if result.Tuint != 20 {
		t.Errorf("tuint value should be 20: %#v", result.Tuint)
	}

	if result.Tbool != true {
		t.Errorf("tbool value should be true: %#v", result.Tbool)
	}

	if result.Tfloat != 20.20 {
		t.Errorf("tfloat value should be 20.20: %#v", result.Tfloat)
	}

	if result.unexported != false {
		t.Error("unexported should not be set, it is unexported")
	}

	if result.Tdata != 20 {
		t.Error("tdata should be valid")
	}
}

func TestFromDefaultTypesToPointer(t *testing.T) {
	t.Parallel()

	input := map[string]interface{}{
		"tstring":    "foo",
		"tint":       20,
		"Tuint":      20,
		"tbool":      true,
		"Tfloat":     20.20,
		"unexported": true,
		"tdata":      20,
	}

	var result DefaultTypesPointer

	sm := structmap.New()
	sm.AddBehavior(name.FromTag("structmap"))

	err := sm.Decode(input, &result)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if *result.Tstring != "foo" {
		t.Errorf("tstring value should be 'foo': %#v", result.Tstring)
	}

	if *result.Tint != 20 {
		t.Errorf("tint value should be 20: %#v", result.Tint)
	}

	if *result.Tuint != 20 {
		t.Errorf("tuint value should be 20: %#v", result.Tuint)
	}

	if *result.Tbool != true {
		t.Errorf("tbool value should be true: %#v", result.Tbool)
	}

	if *result.Tfloat != 20.20 {
		t.Errorf("tfloat value should be 20.20: %#v", result.Tfloat)
	}

	if result.unexported != nil {
		t.Error("unexported should not be set, it is unexported")
	}

	if *result.Tdata != 20 {
		t.Error("tdata should be valid")
	}
}

func TestFromPointerToDefaultTypes(t *testing.T) {
	t.Parallel()

	tstring := "foo"
	tint := 20
	tint8 := int8(20)
	tint16 := int16(20)
	tint32 := int32(20)
	tint64 := int64(20)
	tuint := uint(20)
	tbool := true
	tfloat := 20.20
	unexported := true
	tdata := 20

	input := map[string]interface{}{
		"tstring":    &tstring,
		"tint":       &tint,
		"tint8":      &tint8,
		"tint16":     &tint16,
		"tint32":     &tint32,
		"tint64":     &tint64,
		"Tuint":      &tuint,
		"tbool":      &tbool,
		"Tfloat":     &tfloat,
		"unexported": &unexported,
		"tdata":      &tdata,
	}

	var result DefaultTypes

	sm := structmap.New()
	sm.AddBehavior(name.FromTag("structmap"))

	err := sm.Decode(input, &result)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if result.Tstring != "foo" {
		t.Errorf("tstring value should be 'foo': %#v", result.Tstring)
	}

	if result.Tint != 20 {
		t.Errorf("tint value should be 20: %#v", result.Tint)
	}
	if result.Tint8 != 20 {
		t.Errorf("tint8 value should be 20: %#v", result.Tint)
	}
	if result.Tint16 != 20 {
		t.Errorf("tint16 value should be 20: %#v", result.Tint)
	}
	if result.Tint32 != 20 {
		t.Errorf("tint32 value should be 20: %#v", result.Tint)
	}
	if result.Tint64 != 20 {
		t.Errorf("tint64 value should be 20: %#v", result.Tint)
	}

	if result.Tuint != 20 {
		t.Errorf("tuint value should be 20: %#v", result.Tuint)
	}

	if result.Tbool != true {
		t.Errorf("tbool value should be true: %#v", result.Tbool)
	}

	if result.Tfloat != 20.20 {
		t.Errorf("tfloat value should be 20.20: %#v", result.Tfloat)
	}

	if result.unexported != false {
		t.Error("unexported should not be set, it is unexported")
	}

	if result.Tdata != 20 {
		t.Error("tdata should be valid")
	}
}

func TestFromPointerToPointer(t *testing.T) {
	t.Parallel()

	tstring := "foo"
	tint := 20
	tuint := uint(20)
	tbool := true
	tfloat := 20.20
	unexported := true
	tdata := 20

	input := map[string]interface{}{
		"tstring":    &tstring,
		"tint":       &tint,
		"Tuint":      &tuint,
		"tbool":      &tbool,
		"Tfloat":     &tfloat,
		"unexported": &unexported,
		"tdata":      &tdata,
	}

	var result DefaultTypesPointer

	sm := structmap.New()
	sm.AddBehavior(name.FromTag("structmap"))

	err := sm.Decode(input, &result)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if *result.Tstring != "foo" {
		t.Errorf("tstring value should be 'foo': %#v", result.Tstring)
	}

	if *result.Tint != 20 {
		t.Errorf("tint value should be 20: %#v", result.Tint)
	}

	if *result.Tuint != 20 {
		t.Errorf("tuint value should be 20: %#v", result.Tuint)
	}

	if *result.Tbool != true {
		t.Errorf("tbool value should be true: %#v", result.Tbool)
	}

	if *result.Tfloat != 20.20 {
		t.Errorf("tfloat value should be 20.20: %#v", result.Tfloat)
	}

	if result.unexported != nil {
		t.Error("unexported should not be set, it is unexported")
	}

	if *result.Tdata != 20 {
		t.Error("tdata should be valid")
	}
}

// TODO: Remove this test (tmp)
type testStr struct {
	Headers map[string]map[string]string
}

func TestMapCast(t *testing.T) {
	s := new(testStr)
	m := map[string]interface{}{
		"Headers": map[interface{}]interface{}{
			// "a": "1",
			// "b": "2",
			// "c": "3",
			"d": map[string]int{
				"a": 1,
			},
		},
	}

	sm := structmap.New()
	sm.AddBehavior(name.FromTag("structmap"))
	sm.AddBehavior(cast.ToType)

	err := sm.Decode(m, s)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%+v", s)
}

// TODO: Remove this test (tmp)
func TestName(t *testing.T) {
	s := &struct {
		FirstName string `bson:"first_name"`
		LastName  string `json:"last_name"`
		SnakeCase string
	}{}

	m := map[string]interface{}{
		"first_name": "MyFirstName",
		"last_name":  "MyLastName",
		"snake_case": "MySnakeCase",
	}

	sm := structmap.New()
	sm.AddBehavior(name.Discovery(name.FromTag("json"), name.FromTag("bson"), name.FromSnake))

	err := sm.Decode(m, s)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%+v", s)
}

func TestNameNoop(t *testing.T) {
	s := &struct {
		ValueA string
		ValueB string
	}{}

	m := map[string]interface{}{
		"ValueA": "valA",
		"ValueB": "valB",
	}

	sm := structmap.New()
	sm.AddBehavior(name.Noop)

	err := sm.Decode(m, s)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%+v", s)
}

func TestStructConvert(t *testing.T) {
	s := &struct {
		Date time.Time
	}{}

	m := map[string]interface{}{
		"Date": 1588791963946,
	}

	sm := structmap.New()
	sm.AddBehavior(name.Noop)
	sm.AddBehavior(cast.ToType)

	err := sm.Decode(m, s)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%+v", s)
}

type Person struct {
	Name string
}

type Ladies struct {
	Person `structmap:",noembedded"`
}

func TestLadies(t *testing.T) {
	s := new(Ladies)
	m := map[string]interface{}{
		"Name": "Luana",
		"Person": map[string]interface{}{
			"Name": "Jessica",
		},
	}

	sm := structmap.New()
	sm.AddBehavior(name.Noop)
	sm.AddBehavior(flag.NoEmbedded("structmap"))

	if err := sm.Decode(m, s); err != nil {
		t.Error(err)
	}

	t.Logf("%+v", s)
}

func TestSliceToSlice(t *testing.T) {
	s := new(struct {
		Numbers []int
	})
	m := map[string]interface{}{
		"Numbers": []int{1, 2, 3},
	}

	sm := structmap.New()
	sm.AddBehavior(name.Noop)

	if err := sm.Decode(m, s); err != nil {
		t.Error(err)
	}

	t.Logf("%+v", s)
}

func TestArrayToArray(t *testing.T) {
	s := new(struct {
		Numbers [3]int
	})
	m := map[string]interface{}{
		"Numbers": [3]int{1, 2, 3},
	}

	sm := structmap.New()
	sm.AddBehavior(name.Noop)

	if err := sm.Decode(m, s); err != nil {
		t.Error(err)
	}

	t.Logf("%+v", s)
}

func TestSliceToArrayConverter(t *testing.T) {
	s := new(struct {
		Times [3]int
	})
	m := map[string]interface{}{
		"Times": []int{1588791963946, 1588791963946, 1588791963946},
	}

	sm := structmap.New()
	sm.AddBehavior(name.Noop)
	sm.AddBehavior(cast.ToType)

	if err := sm.Decode(m, s); err != nil {
		t.Error(err)
	}

	t.Logf("%+v", s)
}

func TestSliceToArrayConverterType(t *testing.T) {
	s := new(struct {
		Times [3]time.Time
	})
	m := map[string]interface{}{
		"Times": []int{1588791963946, 1588791963946, 1588791963946},
	}

	sm := structmap.New()
	sm.AddBehavior(name.Noop)
	sm.AddBehavior(cast.ToType)

	if err := sm.Decode(m, s); err != nil {
		t.Error(err)
	}

	t.Logf("%+v", s)
}

func TestToNoStructOrMap(t *testing.T) {
	to := 10
	from := map[string]interface{}{}

	sm := structmap.New()
	sm.AddBehavior(name.Noop)

	if err := sm.Decode(from, &to); err == nil {
		t.Error("expected error: to value: cannot is a struct; got nil")
	}
}

func TestFromNoStructOrMap(t *testing.T) {
	to := struct{}{}
	from := 10

	sm := structmap.New()
	sm.AddBehavior(name.Noop)

	if err := sm.Decode(from, &to); err == nil {
		t.Error("expected error: from value: cannot is map or struct; got nil")
	}
}

func TestStructToStruct(t *testing.T) {
	to := struct {
		Age int
	}{}

	from := struct {
		Age string
	}{
		Age: "20",
	}

	sm := structmap.New()
	sm.AddBehavior(name.Noop)
	sm.AddBehavior(cast.ToType)

	if err := sm.Decode(from, &to); err != nil {
		t.Error(err)
	}

	t.Logf("%+v", to)
}
