package structmap

import "reflect"

type Field struct {
	reflect.StructField
	Value reflect.Value
}

// IsEmbedded returns true if the given field is an anonymous field (embedded)
func (f Field) IsEmbedded() bool {
	return f.Anonymous
}

// IsExported returns true if the given field is exported.
func (f Field) IsExported() bool {
	return f.PkgPath == ""
}

// IsZero returns field is a zero value
func (f Field) IsZero() bool {
	zero := reflect.Zero(f.Type).Interface()
	current := f.Value.Interface()

	return reflect.DeepEqual(current, zero)
}
