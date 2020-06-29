package name

import (
	"github.com/amulets/structmap"
	"github.com/amulets/structmap/behavior"
	"github.com/huandu/xstrings"
)

var (
	// FromSnake transforms field name in snake_case format
	FromSnake = behavior.New(func(field *structmap.FieldPart) error {
		field.Name = xstrings.ToSnakeCase(field.Name)

		return nil
	})

	// FromCamel transforms field name in camelCase format
	FromCamel = behavior.New(func(field *structmap.FieldPart) error {
		field.Name = xstrings.FirstRuneToLower(xstrings.ToCamelCase(xstrings.ToSnakeCase(field.Name)))

		return nil
	})

	// FromPascal transforms field name in PascalCase format
	FromPascal = behavior.New(func(field *structmap.FieldPart) error {
		field.Name = xstrings.ToCamelCase(xstrings.ToSnakeCase(field.Name))

		return nil
	})

	// FromKebab transforms field name in kebab-case format
	FromKebab = behavior.New(func(field *structmap.FieldPart) error {
		field.Name = xstrings.ToKebabCase(field.Name)

		return nil
	})
)
