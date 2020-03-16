package env

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

var (
	IsNotPointerError     = errors.New("object isn't a pointer")
	IsNotInitializedError = errors.New("object must be initialized")
	IsNotStructError      = errors.New("object isn't a struct")
	TypeError             = errors.New("incorrect type")
)

// instance is an auxiliary structure for performing reflection.
type instance struct {
	Ptr   reflect.Value
	Type  reflect.Type
	Kind  reflect.Kind
	Value reflect.Value

	IsPtr    bool
	IsStruct bool
	IsValid  bool
}

// Init defines the main reflect's parameters.
func (inst *instance) Init(obj interface{}) {
	inst.Type = reflect.TypeOf(obj)
	inst.Value = reflect.ValueOf(obj)
	inst.Kind = inst.Type.Kind()

	if inst.Kind == reflect.Ptr {
		inst.IsPtr = true
		inst.Ptr = inst.Value
		inst.Type = inst.Type.Elem()
		inst.Kind = inst.Type.Kind()
		inst.Value = inst.Value.Elem()
	} else {
		inst.Ptr = reflect.New(inst.Type)
		tmp := inst.Ptr.Elem()
		tmp.Set(inst.Value)
	}

	inst.IsStruct = inst.Kind == reflect.Struct
	inst.IsValid = inst.Value.IsValid()
}

// Implements returns true if instance implements interface.
//
// Usage:
//	if inst.Implements((*CustomInterface)(nil)) { ... }
func (inst *instance) Implements(ifc interface{}) bool {
	return inst.Ptr.Type().Implements(reflect.TypeOf(ifc).Elem())
}

// Marshaler describes an interface for implementing
// a custom method for marshaling.
type Marshaler interface {
	MarshalENV() ([]string, error)
}

// Unmarshaler describes an interface for implementing
// a custom method for unmarshaling.
type Unmarshaler interface {
	UnmarshalENV() error
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

// unmarshalENV gets variables from the environment and sets them by
// pointer into obj. Returns an error if something went wrong.
//
// Supported types: int, int8, int16, int32, int64, uint, uint8, uint16,
// uint32, uint64, bool, float32, float64, string, and slice from thous types.
func unmarshalENV(obj interface{}, pfx string) error {
	var inst instance = instance{}
	inst.Init(obj)

	// The object must be an initialized pointer of the struct.
	switch {
	case !inst.IsPtr:
		return IsNotPointerError
	case !inst.IsValid:
		return IsNotInitializedError
	case !inst.IsStruct:
		return IsNotStructError
	}

	// If objects implements Unmarshaler interface try to calling
	// a custom Unmarshal method.
	if inst.Implements((*Unmarshaler)(nil)) {
		if m := inst.Ptr.MethodByName("UnmarshalENV"); m.IsValid() {
			tmp := m.Call([]reflect.Value{})
			err := tmp[0].Interface()
			if err != nil {
				return fmt.Errorf("env: unmarshal: %v", err)
			}
			return nil
		}
	}

	// Walk through all the fields of the struct.
	for i := 0; i < inst.Value.NumField(); i++ {
		// Get item.
		field := inst.Type.Field(i)
		item := inst.Value.FieldByName(field.Name)

		// Get key and sep for sequences.
		key, sep := parseTag(field.Tag.Get("env"), field.Name, " ")
		key = fmt.Sprintf("%s%s", pfx, key)

		// Set values of the desired type.
		kind := field.Type.Kind()
		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64:
			r, err := strToIntKind(Get(key), kind)
			if err != nil {
				return err
			}
			item.SetInt(r)
		case reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64:
			r, err := strToUintKind(Get(key), kind)
			if err != nil {
				return err
			}
			item.SetUint(r)
		case reflect.Float32, reflect.Float64:
			r, err := strToFloatKind(Get(key), kind)
			if err != nil {
				return err
			}
			item.SetFloat(r)
		case reflect.Bool:
			r, err := strToBool(Get(key))
			if err != nil {
				return err
			}
			item.SetBool(r)
		case reflect.String:
			item.SetString(Get(key))
		case reflect.Array:
			max := item.Type().Len()
			seq := strings.Split(Get(key), sep)
			if len(seq) > max {
				return errors.New(fmt.Sprintf(
					"%d items overwhelms the [%d]array",
					len(seq), max,
				))
			}
			err := setSequence(&item, strings.Split(Get(key), sep))
			if err != nil {
				return err
			}
		case reflect.Slice:
			seq := strings.Split(Get(key), sep)
			tmp := reflect.MakeSlice(item.Type(), len(seq), len(seq))
			err := setSequence(&tmp, strings.Split(Get(key), sep))
			if err != nil {
				return err
			}
			item.Set(reflect.AppendSlice(item, tmp))
		case reflect.Ptr:
			if item.Type() == reflect.TypeOf((*url.URL)(nil)) {
				// The *url.URL pointer.
				u, err := url.Parse(Get(key))
				if err != nil {
					return err
				}
				item.Set(reflect.ValueOf(u))
			} else {
				// Another type of pointer's struct.
				tmp := reflect.New(item.Type().Elem()).Interface()
				err := unmarshalENV(tmp, fmt.Sprintf("%s_", key))
				if err != nil {
					return err
				}
				item.Set(reflect.ValueOf(tmp))
			}
		case reflect.Struct:
			if _, ok := item.Interface().(url.URL); ok {
				// Parse url.URL.
				u, err := url.Parse(Get(key))
				if err != nil {
					return err
				}
				item.Set(reflect.ValueOf(*u))
			} else {
				// Another type of struct.
				tmp := reflect.New(item.Type()).Interface()
				err := unmarshalENV(tmp, fmt.Sprintf("%s_", key))
				if err != nil {
					return err
				}
				item.Set(reflect.ValueOf(tmp).Elem())
			}
		default:
			return TypeError
		}
	}

	return nil
}

// setSequence sets slice into instance.
func setSequence(item *reflect.Value, seq []string) (err error) {
	var kind = item.Index(0).Kind()

	defer func() {
		// Catch the panic and return an exception as a value.
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()

	// Ignore empty containers.
	switch {
	case kind == reflect.Array && item.Type().Len() == 0:
		fallthrough
	case kind == reflect.Slice && item.Len() == 0:
		return nil
	}

	// Convert to correct type slice.
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		for i, value := range seq {
			r, err := strToIntKind(value, kind)
			if err != nil {
				return err
			}
			item.Index(i).SetInt(r)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		for i, value := range seq {
			r, err := strToUintKind(value, kind)
			if err != nil {
				return err
			}
			item.Index(i).SetUint(r)
		}
	case reflect.Float32, reflect.Float64:
		for i, value := range seq {
			r, err := strToFloatKind(value, kind)
			if err != nil {
				return err
			}
			item.Index(i).SetFloat(r)
		}
	case reflect.Bool:
		for i, value := range seq {
			r, err := strToBool(value)
			if err != nil {
				return err
			}
			item.Index(i).SetBool(r)
		}
	case reflect.String:
		for i, value := range seq {
			item.Index(i).SetString(value)
		}
	case reflect.Ptr:
		// The *url.URL pointer only.
		if len(seq) == 0 {
			break
		}

		if item.Index(0).Type() != reflect.TypeOf((*url.URL)(nil)) {
			return TypeError
		}

		for i, value := range seq {
			u, err := url.Parse(value)
			if err != nil {
				return err
			}

			item.Index(i).Set(reflect.ValueOf(u))
		}
	case reflect.Struct:
		// The url.URL struct only.
		if len(seq) == 0 {
			break
		}

		if _, ok := item.Index(0).Interface().(url.URL); !ok {
			return TypeError
		}

		for i, value := range seq {
			u, err := url.Parse(value)
			if err != nil {
				return err
			}

			item.Index(i).Set(reflect.ValueOf(*u))
		}
	default:
		return TypeError
	}

	return nil
}
