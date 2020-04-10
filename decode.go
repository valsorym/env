package env

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// Unmarshaler is the interface implemented by types that can unmarshal
// an environment variables of themselves.
type Unmarshaler interface {
	UnmarshalENV() error
}

// unmarshalENV gets variables from the environment and sets them
// into object by pointer. Returns an error if something went wrong.
//
// unmarshalENV method supports the following field's types: int, int8, int16,
// int32, int64, uin, uint8, uin16, uint32, in64, float32, float64, string,
// bool, url.URL and pointers, array or slice from thous types (i.e. *int, ...,
// []int, ..., []bool, ..., [2]*url.URL, etc.). The nested structures will be
// processed recursively.
//
// For other filed's types (like chan, map ...) will be returned an error.
//
// Among the supported types are: struct and pointer to struct but
// slice/array of these types is not supported (except url.URL and
// *url.URL from the net package).
func unmarshalENV(obj interface{}, pfx string) error {
	var inst instance = instance{}
	inst.Init(obj)

	// The object must be an initialized pointer of the struct.
	switch {
	case !inst.IsPtr:
		return errors.New("object isn't a pointer")
	case !inst.IsValid:
		return errors.New("object must be initialized")
	case !inst.IsStruct:
		return errors.New("object isn't a struct")
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
		key, value, sep, err := splitFieldTag(field.Tag.Get("env"))
		if err != nil {
			return err
		}

		// Create full key name.
		if len(key) == 0 {
			key = field.Name
		}

		key = fmt.Sprintf("%s%s", pfx, key)

		// If the value is defined in environment set it into value.
		if Exists(key) {
			value = Get(key)
		}

		// Set values of the desired type.
		switch item.Kind() {
		case reflect.Array:
			max := item.Type().Len()
			seq := strings.Split(value, sep)
			if len(seq) > max {
				return fmt.Errorf("%d overflows the [%d]array", len(seq), max)
			}

			err := setSequence(&item, strings.Split(value, sep))
			if err != nil {
				return err
			}
		case reflect.Slice:
			seq := strings.Split(value, sep)
			tmp := reflect.MakeSlice(item.Type(), len(seq), len(seq))
			err := setSequence(&tmp, strings.Split(value, sep))
			if err != nil {
				return err
			}
			item.Set(reflect.AppendSlice(item, tmp))
		case reflect.Ptr:
			switch {
			case item.Type().Elem().Kind() != reflect.Struct:
				// If the pointer is not to a structure.
				tmp := reflect.Indirect(item)
				err := setValue(tmp, value)
				if err != nil {
					return err
				}
			case item.Type() == reflect.TypeOf((*url.URL)(nil)):
				// If a pointer to a structure of the url.URL.
				err := setValue(item, value)
				if err != nil {
					return err
				}
			default:
				// If a pointer to a structure of the another's types.
				// P.s. Not a *url.URL.
				tmp := reflect.New(item.Type().Elem()).Interface()
				err := unmarshalENV(tmp, fmt.Sprintf("%s_", key))
				if err != nil {
					return err
				}
				item.Set(reflect.ValueOf(tmp))
			}
		case reflect.Struct:
			switch {
			case item.Type() == reflect.TypeOf(url.URL{}):
				// If a url.URL structure.
				err := setValue(item, value)
				if err != nil {
					return err
				}
			default:
				// If a structure of the another's types.
				// P.s. Not a url.URL.
				tmp := reflect.New(item.Type()).Interface()
				err := unmarshalENV(tmp, fmt.Sprintf("%s_", key))
				if err != nil {
					return err
				}
				item.Set(reflect.ValueOf(tmp).Elem())
			}
		default:
			// Try to set correct value.
			err := setValue(item, value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// setSequence sets slice into item.
func setSequence(item *reflect.Value, seq []string) (err error) {
	var kind = item.Index(0).Kind()

	defer func() {
		// Catch the panic and return an exception as a value.
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	// Ignore empty containers.
	switch {
	case kind == reflect.Array && item.Type().Len() == 0:
		fallthrough
	case kind == reflect.Slice && item.Len() == 0:
		fallthrough
	case len(seq) == 0:
		return nil
	}

	// Set values from sequence.
	for i, value := range seq {
		elem := item.Index(i)
		err := setValue(elem, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// setValue sets value.
func setValue(item reflect.Value, value string) error {
	kind := item.Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		r, err := strToIntKind(value, kind)
		if err != nil {
			return err
		}
		item.SetInt(r)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		r, err := strToUintKind(value, kind)
		if err != nil {
			return err
		}
		item.SetUint(r)
	case reflect.Float32, reflect.Float64:
		r, err := strToFloatKind(value, kind)
		if err != nil {
			return err
		}
		item.SetFloat(r)
	case reflect.Bool:
		r, err := strToBool(value)
		if err != nil {
			return err
		}
		item.SetBool(r)
	case reflect.String:
		item.SetString(value)
	case reflect.Ptr:
		// The *url.URL pointer only.
		switch {
		case item.Type() == reflect.TypeOf((*url.URL)(nil)):
			u, err := url.Parse(value)
			if err != nil {
				return err
			}
			item.Set(reflect.ValueOf(u))
		default:
			return fmt.Errorf("incorrect type: %s", item.Type())
		}
	case reflect.Struct:
		// The url.URL struct only.
		switch {
		case item.Type() == reflect.TypeOf(url.URL{}):
			u, err := url.Parse(value)
			if err != nil {
				return err
			}
			item.Set(reflect.ValueOf(*u))
		default:
			return fmt.Errorf("incorrect type: %s", item.Type())
		}
	default:
		return fmt.Errorf("incorrect type: %s", item.Type())
	}

	return nil
}
