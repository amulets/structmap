package structmap

import "reflect"

type field struct {
	reflect.StructField
	Value reflect.Value
}

// IsEmbedded returns true if the given field is an anonymous field (embedded)
func (f field) IsEmbedded() bool {
	return f.Anonymous
}

// IsExported returns true if the given field is exported.
func (f field) IsExported() bool {
	return f.PkgPath == ""
}

// IsZero returns field is a zero value
func (f field) IsZero() bool {
	zero := reflect.Zero(f.Type).Interface()
	current := f.Value.Interface()

	return reflect.DeepEqual(current, zero)
}
