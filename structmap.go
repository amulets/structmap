package structmap

import (
	"fmt"
	"reflect"
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

	// BehaviorFunc that's change field information
	BehaviorFunc func(*FieldPart) error

	// StructMap is a structmap
	StructMap struct {
		behaviors []BehaviorFunc
	}
)

// New instance of StructMap
func New() *StructMap {
	return &StructMap{}
}

// AddBehavior a new behavior logic
func (sm *StructMap) AddBehavior(behavior BehaviorFunc) {
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
		fp := &FieldPart{
			Name:       field.Name,
			Tag:        field.Tag,
			Type:       field.Type,
			IsEmbedded: field.IsEmbedded(),
		}

		// run behaviors
		for i, behavior := range sm.behaviors {
			if err := behavior(fp); err != nil {
				return err
			}

			// expects there first behavior get field name to get field value
			if i == 0 {
				if value, ok := from[fp.Name]; ok {
					fp.Value = value
				}
			}
		}

		if fp.Skip {
			continue
		}

		if field.Type.Kind() == reflect.Struct || (field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct) {
			value := field.Value

			if field.Value.Kind() == reflect.Ptr && field.IsZero() {
				pv := reflect.New(field.Type.Elem())
				field.Value.Set(pv)
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
		fieldValue := field.Value

		// Get value element
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}

		// Ignore if no have a value
		if !value.IsValid() {
			continue
		}

		// Get field value element
		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsZero() {
				fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
			}

			fieldValue = fieldValue.Elem()
		}

		if value.Type().ConvertibleTo(fieldValue.Type()) {
			value = value.Convert(fieldValue.Type())
		}

		if value.Kind() != fieldValue.Kind() {
			return fmt.Errorf("field %s value of type %s is not assignable to type %s", field.Name, value.Type(), fieldValue.Type())
		}

		if field.Value.Kind() == reflect.Ptr {
			field.Value.Elem().Set(value)
		} else {
			field.Value.Set(value)
		}
	}

	return nil
}
