package env

import (
	"fmt"
	"reflect"
	"strings"
)

/*
 var ptr reflect.Value
    var value reflect.Value
    var finalMethod reflect.Value

    value = reflect.ValueOf(i)

    // if we start with a pointer, we need to get value pointed to
    // if we start with a value, we need to get a pointer to that value
    if value.Type().Kind() == reflect.Ptr {
        ptr = value
        value = ptr.Elem()
    } else {
        ptr = reflect.New(reflect.TypeOf(i))
        temp := ptr.Elem()
        temp.Set(value)
    }
*/

// marshalENV ...
func marshalENV(scope interface{}) error {
	var (
		rt reflect.Type  // type
		rv reflect.Value // value
		rp reflect.Value // pointer
	)

	// Define: type, value and pointer.
	rt = reflect.TypeOf(scope)
	rv = reflect.ValueOf(scope)
	if rt.Kind() == reflect.Ptr {
		rp, rt, rv = rv, rt.Elem(), rv.Elem()
	} else {
		rp = reflect.New(reflect.TypeOf(scope))
		temp := rp.Elem()
		temp.Set(rv)
	}

	// Scope validation.
	switch {
	case rt.Kind() != reflect.Struct:
		return fmt.Errorf("object must be a structure")
	case !rv.IsValid():
		return fmt.Errorf("object must be initialized")
	}

	// If there is the custom method, MarshlaENV - run it.
	if m := rp.MethodByName("MarshalENV"); m.IsValid() {
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
		//fmt.Println(key, value)
	} // for

	return nil
}
