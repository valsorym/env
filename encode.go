package env

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// Marshaler describes an interface for implementing
// a custom method for marshaling.
type Marshaler interface {
	MarshalENV() ([]string, error)
}

// marshalENV saves obj into environment data.
//
// Supported types: int, int8, int16, int32, int64, uint, uint8, uint16,
// uint32, uint64, bool, float32, float64, string, and slice from thous types.
func marshalENV(obj interface{}, prefix string) ([]string, error) {
	var (
		err    error
		result []string
	)

	inst := instance{}
	inst.Init(obj)

	// The object must be an initialized of the struct.
	switch {
	case !inst.IsValid:
		return []string{}, IsNotInitializedError
	case !inst.IsStruct:
		return []string{}, IsNotStructError
	}

	// Implements Marshaler interface.
	if inst.Implements((*Marshaler)(nil)) {
		// Try to run custom MarshalENV function.
		if m := inst.Ptr.MethodByName("MarshalENV"); m.IsValid() {
			tmp := m.Call([]reflect.Value{})
			value := tmp[0].Interface()
			err := tmp[1].Interface()
			if err != nil {
				return []string{}, fmt.Errorf("marshal: %v", err)
			}
			return value.([]string), nil
		}
	}

	// Walk through the fields.
	result = make([]string, 0, inst.Value.NumField()) // -1
	for i := 0; i < inst.Value.NumField(); i++ {
		var key, value, sep string
		field := inst.Value.Type().Field(i)
		item := inst.Value.FieldByName(field.Name)

		key, sep = parseTag(field.Tag.Get("env"), field.Name, " ")
		kind := field.Type.Kind()
		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64:
			value = fmt.Sprintf("%d", item.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64:
			value = fmt.Sprintf("%d", item.Uint())
		case reflect.Float32, reflect.Float64:
			value = fmt.Sprintf("%f", item.Float())
		case reflect.Bool:
			value = fmt.Sprintf("%t", item.Bool())
		case reflect.String:
			value = fmt.Sprintf("%s", item.String())
		case reflect.Array:
			value, err = getSequence(&item, sep)
			if err != nil {
				return result, err
			}
		case reflect.Slice:
			value, err = getSequence(&item, sep)
			if err != nil {
				return result, err
			}
		case reflect.Ptr:
			item = item.Elem()
			if item.Kind() != reflect.Struct {
				return result, TypeError
			}
			fallthrough
		case reflect.Struct:
			if u, ok := item.Interface().(url.URL); ok {
				// Type of url.URL.
				value = u.String()
			} else {
				// Recursive parsing
				p := fmt.Sprintf("%s%s_", prefix, key)
				value, err := marshalENV(item.Interface(), p)
				if err != nil {
					return result, err
				}

				// Expand the result.
				for _, v := range value {
					result = append(result, v)
				}
				continue // object doesn't save
			}
		default:
			return result, TypeError
		} // switch

		// Set into environment and add to result list.
		key = fmt.Sprintf("%s%s", prefix, key)
		Set(key, value)
		result = append(result, fmt.Sprintf("%s=%s", key, value))
	} // for

	return result, nil
}

// getSequence get sequence as string.
func getSequence(item *reflect.Value, sep string) (string, error) {
	var kind reflect.Kind
	var max int

	switch item.Kind() {
	case reflect.Array:
		kind = item.Index(0).Kind()
		max = item.Type().Len()
	case reflect.Slice:
		tmp := reflect.MakeSlice(item.Type(), 1, 1)
		kind = tmp.Index(0).Kind()
		max = item.Len()
	default:
		return "", TypeError
	}

	switch kind {
	case reflect.Ptr:
		var tmp = []string{}
		for i := 0; i < max; i++ {
			elem := item.Index(i).Elem()
			if v, ok := elem.Interface().(url.URL); ok {
				tmp = append(tmp, v.String())
			} else {
				return "", TypeError
			}
		}
		str := strings.Replace(fmt.Sprint(tmp), " ", sep, -1)
		return strings.Trim(str, "[]"+sep), nil
	case reflect.Struct:
		var tmp = []string{}
		for i := 0; i < max; i++ {
			elem := item.Index(i)
			if v, ok := elem.Interface().(url.URL); ok {
				tmp = append(tmp, v.String())
			} else {
				return "", TypeError
			}
		}
		str := strings.Replace(fmt.Sprint(tmp), " ", sep, -1)
		return strings.Trim(str, "[]"+sep), nil
	}
	str := strings.Replace(fmt.Sprint(*item), " ", sep, -1)
	return strings.Trim(str, "[]"+sep), nil

}
