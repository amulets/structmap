package structmap

import "errors"

// All erros of StructMap
var (
	ErrNotIsToPointer = errors.New("to value cannot is a pointer")
	ErrNotIsStruct    = errors.New("value cannot is a struct")
)
