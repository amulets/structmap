package name

import (
	"github.com/dungeon-code/structmap"
	"github.com/huandu/xstrings"
)

// FromSnake transforms field name in snake_case format
func FromSnake(field *structmap.FieldPart) error {
	field.Name = xstrings.ToSnakeCase(field.Name)

	return nil
}

// FromCamel transforms field name in camelCase format
func FromCamel(field *structmap.FieldPart) error {
	field.Name = xstrings.FirstRuneToLower(xstrings.ToCamelCase(xstrings.ToSnakeCase(field.Name)))

	return nil
}

// FromPascal transforms field name in PascalCase format
func FromPascal(field *structmap.FieldPart) error {
	field.Name = xstrings.ToCamelCase(xstrings.ToSnakeCase(field.Name))

	return nil
}

// FromKebab transforms field name in kebab-case format
func FromKebab(field *structmap.FieldPart) error {
	field.Name = xstrings.ToKebabCase(field.Name)

	return nil
}
