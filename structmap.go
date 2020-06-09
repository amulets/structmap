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

	// BehaviorFunc that deals with the behavior of the field
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
func (sm *StructMap) Decode(from interface{}, to interface{}) (err error) {
	defer func() {
		if err == nil {
			if recovered := recover(); recovered != nil {
				err = fmt.Errorf("%v", recovered)
			}
		}
	}()

	if _, ok := from.(map[string]interface{}); !ok {
		from, err = newStruct(from)
		if err != nil {
			return err
		}
	}

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
				switch fromValue := from.(type) {
				case map[string]interface{}:
					if value, ok := fromValue[fp.Name]; ok {
						fp.Value = value
					}
				case *strct:
					if _, ok := fromValue.Type().FieldByName(fp.Name); ok {
						fp.Value = fromValue.FieldByName(fp.Name).Interface()
					}
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

			fromFace := from
			needDecode := true

			if !fp.IsEmbedded {
				switch from.(type) {
				case map[string]interface{}:
					if mapValue, ok := fp.Value.(map[string]interface{}); ok {
						needDecode = len(mapValue) > 0

						if needDecode {
							fromFace = mapValue
						}
					} else {
						needDecode = false
					}
				case *strct:
					if strctValue, err := newStruct(fp.Value); err == nil {
						needDecode = strctValue.IsValid()

						if needDecode {
							fromFace = strctValue.Interface()
						}
					} else {
						needDecode = false
					}
				}
			}

			if needDecode {
				if strctValue, ok := fromFace.(*strct); ok {
					fromFace = strctValue.Interface()
				}

				if err := sm.Decode(fromFace, value.Interface()); err != nil {
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
