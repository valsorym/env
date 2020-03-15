package env

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// Unmarshaler describes an interface for implementing
// a custom method for unmarshaling.
type Unmarshaler interface {
	UnmarshalENV() error
}

// Unmarshal
func unmarshalENV(obj interface{}, pfx string) error {
	var inst instance = instance{}
	inst.Init(obj)

	// The object must be an initialized pointer to the structure.
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
func setSequence(instance *reflect.Value, seq []string) (err error) {
	var kind = instance.Index(0).Kind()

	defer func() {
		// Catch the panic and return an exception as a value.
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	// Ignore empty containers.
	switch {
	case kind == reflect.Array && instance.Type().Len() == 0:
		fallthrough
	case kind == reflect.Slice && instance.Len() == 0:
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
			instance.Index(i).SetInt(r)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		for i, value := range seq {
			r, err := strToUintKind(value, kind)
			if err != nil {
				return err
			}
			instance.Index(i).SetUint(r)
		}
	case reflect.Float32, reflect.Float64:
		for i, value := range seq {
			r, err := strToFloatKind(value, kind)
			if err != nil {
				return err
			}
			instance.Index(i).SetFloat(r)
		}
	case reflect.Bool:
		for i, value := range seq {
			r, err := strToBool(value)
			if err != nil {
				return err
			}
			instance.Index(i).SetBool(r)
		}
	case reflect.String:
		for i, value := range seq {
			instance.Index(i).SetString(value)
		}
	case reflect.TypeOf(&url.URL{}).Kind():
		for i, value := range seq {
			u, err := url.Parse(value)
			if err != nil {
				return err
			}

			instance.Index(i).Set(reflect.ValueOf(u))
		}
	case reflect.TypeOf(url.URL{}).Kind():
		for i, value := range seq {
			u, err := url.Parse(value)
			if err != nil {
				return err
			}

			instance.Index(i).Set(reflect.ValueOf(*u))
		}
	default:
		return TypeError
	}

	return nil
}
