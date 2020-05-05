package structmap

import "strings"

// FieldFlags of StructTag
type FieldFlags []string

// Has returns true if a flag is available in FieldFlags
func (flags FieldFlags) Has(flagName string) bool {
	for _, flag := range flags {
		if flag == flagName {
			return true
		}
	}

	return false
}

// ParseTag return field tag value and flags
// tag is one of followings:
// ""
// "value"
// "value,flags"
// "value,flags,flags2"
// ",flags"
func ParseTag(tag string) (string, FieldFlags) {
	parts := strings.Split(tag, ",")
	return parts[0], parts[1:]
}
