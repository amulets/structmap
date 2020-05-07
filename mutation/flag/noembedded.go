package flag

import "github.com/dungeon-code/structmap"

// NoEmbedded treat embedded struct as non embedded struct
// and get it's value from root of map
// Only for fields marked with the flag: noembedded
// Example: `tagName:",noembedded"`
func NoEmbedded(tagName string) structmap.MutationFunc {
	return func(field *structmap.FieldPart) error {
		if !field.IsEmbedded {
			return nil
		}

		_, flags := structmap.ParseTag(field.Tag.Get(tagName))

		if flags.Has("noembedded") {
			field.IsEmbedded = false
		}

		return nil
	}
}
