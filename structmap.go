package structmap

import (
	"fmt"
	"reflect"

	"github.com/amulets/structmap/internal"
)

type (
	// FieldPart is a field representation
	FieldPart struct {
		Name       string
		Value      interface{}
		Type       reflect.Type
		Tag        reflect.StructTag
		Skip       bool
		IsEmbedded bool
	}

	// Behavior implementation
	Behavior interface {
		Do(*FieldPart) error
	}

	// StructMap is a structmap
	StructMap struct {
		behaviors []Behavior
	}
)

// New instance of StructMap
func New() *StructMap {
	return &StructMap{}
}

// AddBehavior a new behavior logic
func (sm *StructMap) AddBehavior(behavior Behavior) {
	sm.behaviors = append(sm.behaviors, behavior)
}

// Decode map to struct
func (sm *StructMap) Decode(from map[string]interface{}, to interface{}) (err error) {
	defer func() {
		if err == nil {
			if recovered := recover(); recovered != nil {
				err = fmt.Errorf("%v", recovered)
			}
		}
	}()

	s, err := newStruct(to)
	if err != nil {
		return err
	}

	// Struct is configurable?
	if !s.CanSet() {
		return ErrNotIsToPointer
	}

	for _, field := range s.Fields() {
		fieldType := internal.Type(field.Type)

		fp := &FieldPart{
			Name:       field.Name,
			Tag:        field.Tag,
			Type:       fieldType,
			IsEmbedded: field.IsEmbedded(),
		}

		// run behaviors
		for i, behavior := range sm.behaviors {
			if err := behavior.Do(fp); err != nil {
				return err
			}

			// expects there first behavior get field name to get field value
			if i == 0 {
				if rawValue, ok := from[fp.Name]; ok {
					value := internal.Value(reflect.ValueOf(rawValue), false)

					if value.IsValid() {
						fp.Value = value.Interface()
					}
				}
			}
		}

		if fp.Skip {
			continue
		}

		if fieldType.Kind() == reflect.Struct {
			value := field.Value

			if field.Value.Kind() == reflect.Ptr && field.IsZero() {
				pv := reflect.New(fieldType)

				internal.SetValue(field.Value, pv.Elem())
			} else {
				value = field.Value.Addr()
			}

			mapFrom := from
			mapNeedDecode := true

			if !fp.IsEmbedded {
				var ok bool
				if mapFrom, ok = fp.Value.(map[string]interface{}); !ok {
					mapNeedDecode = false
				}
			}

			if mapNeedDecode {
				if err := sm.Decode(mapFrom, value.Interface()); err != nil {
					return err
				}

				continue
			}
		}

		value := reflect.ValueOf(fp.Value)

		// Ignore if no have a value
		if !value.IsValid() {
			continue
		}

		if value.Type().ConvertibleTo(fieldType) {
			value = value.Convert(fieldType)
		}

		if value.Kind() != fieldType.Kind() {
			return fmt.Errorf("field %s value of type %s is not assignable to type %s", field.Name, value.Type(), fieldType)
		}

		internal.SetValue(field.Value, value)
	}

	return nil
}
