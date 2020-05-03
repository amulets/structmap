package structmap

import "strings"

// FieldTags options
type FieldTags []string

// Has returns true if a option is available in FieldTags
func (tags FieldTags) Has(opt string) bool {
	for _, option := range tags {
		if option == opt {
			return true
		}
	}

	return false
}

// ParseTag return field name and options
// tag is one of followings:
// ""
// "name"
// "name,opt"
// "name,opt,opt2"
// ",opt"
func ParseTag(tag string) (string, FieldTags) {
	parts := strings.Split(tag, ",")
	return parts[0], parts[1:]
}
