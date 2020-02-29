package env

import (
	"fmt"
	"reflect"
	"strings"
)

// marshalENV ...
func marshalENV(scope interface{}) error {
	var (
		rt = reflect.TypeOf(scope)
		rv = reflect.ValueOf(scope)
	)

	// For pointer return real object.
	if rt.Kind() == reflect.Ptr {
		rt, rv = rt.Elem(), rv.Elem()
	}

	// Value must be a structure.
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("value must be initialized struct")
	}

	// If there is the custom method, MarshlaENV - run it.
	if m := reflect.New(rt).MethodByName("MarshalENV"); m.IsValid() {
		result := m.Call([]reflect.Value{})
		if len(result) != 0 {
			err := result[0].Interface()
			if err != nil {
				return fmt.Errorf("marshal: %v", err)
			}
		}
		return nil
	}

	// Walk through the fields.
	for i := 0; i < rv.NumField(); i++ {
		var key, value, sep string
		field := rv.Type().Field(i)
		key, sep = parseTag(field.Tag.Get("env"), field.Name, " ")

		// Set value.
		instance := rv.FieldByName(field.Name)
		kind := field.Type.Kind()
		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64:
			value = fmt.Sprintf("%d", instance.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64:
			value = fmt.Sprintf("%d", instance.Uint())
		case reflect.Float32, reflect.Float64:
			value = fmt.Sprintf("%f", instance.Float())
		case reflect.Bool:
			value = fmt.Sprintf("%t", instance.Bool())
		case reflect.String:
			value = fmt.Sprintf("%s", instance.String())
		case reflect.Slice:
			str := strings.Replace(fmt.Sprint(instance), " ", sep, -1)
			value = strings.Trim(str, "[]")
		default:
			return fmt.Errorf("incorrect type")
		} // switch

		// Set into environment.
		Set(key, value)
	} // for

	return nil
}
