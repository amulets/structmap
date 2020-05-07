package structmap

import "testing"

type testStruct struct {
	A string
	B int
	c bool // unexported
}

func TestNewStruct(t *testing.T) {
	s := new(testStruct)

	strct, err := newStruct(s)
	if err != nil {
		t.Error(err)
	}

	if !strct.IsValid() {
		t.Error("struct do not is valid value")
	}
}

func TestStructIsConfigurable(t *testing.T) {
	s := new(testStruct)

	strct, err := newStruct(s)
	if err != nil {
		t.Error(err)
	}

	if !strct.CanSet() {
		t.Error("struct do not is configurable")
	}
}

func TestNotIsStruct(t *testing.T) {
	s := "new(testStruct)"

	_, err := newStruct(s)
	if err == nil {
		t.Error("cannot catch that's value cannot is a struct")
	}
}

func TestStructFieldsCountUnexported(t *testing.T) {
	s := new(struct{ a string })

	strct, err := newStruct(s)
	if err != nil {
		t.Error(err)
	}

	length := len(strct.Fields())

	if length > 0 {
		t.Errorf("unexported field expect 0; got %d", length)
	}
}

func TestStructFieldsCount(t *testing.T) {
	s := new(testStruct)

	strct, err := newStruct(s)
	if err != nil {
		t.Error(err)
	}

	length := len(strct.Fields())

	if length != 2 {
		t.Errorf("fields expect 2; got %d", length)
	}
}
