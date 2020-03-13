package env

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// Error's values.
var (
	InitializedError = fmt.Errorf("object must be initialized")
	ObjectError      = fmt.Errorf("object must be a structure")
	TypeError        = fmt.Errorf("incorrect type")
)

// Marshaller describes an interface for implementing
// a custom method for marshaling.
type Marshaller interface {
	MarshalENV() ([]string, error)
}

// Unmarshaller describes an interface for implementing
// a custom method for unmarshaling.
type Unmarshaller interface {
	UnmarshalENV() error
}

// marshalENV saves scope into environment data.
//
// Supported types: int, int8, int16, int32, int64, uint, uint8, uint16,
// uint32, uint64, bool, float32, float64, string, and slice from thous types.
func marshalENV(scope interface{}) ([]string, error) {
	var (
		rt reflect.Type  // type
		rv reflect.Value // value
		rp reflect.Value // pointer

		err    error
		result []string
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
		return result, ObjectError
	case !rv.IsValid():
		return result, InitializedError
	}

	// Implements Marshaler interface.
	if rp.Type().Implements(reflect.TypeOf((*Marshaller)(nil)).Elem()) {
		// Try to run custom MarshalENV function.
		if m := rp.MethodByName("MarshalENV"); m.IsValid() {
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
	result = make([]string, 0, rv.NumField()-1)
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
		case reflect.Array:
			value, err = getSequence(&instance, sep)
			if err != nil {
				return result, err
			}
		case reflect.Slice:
			value, err = getSequence(&instance, sep)
			if err != nil {
				return result, err
			}
		case reflect.TypeOf(&url.URL{}).Kind():
			instance = instance.Elem()
			fallthrough
		case reflect.TypeOf(url.URL{}).Kind():
			u := instance.Interface().(url.URL)
			value = u.String()
		default:
			return result, TypeError
		} // switch

		// Set into environment and add to result list.
		Set(key, value)
		result = append(result, fmt.Sprintf("%s=%s", key, value))
	} // for

	return result, nil
}

// unmarshalENV gets variables from the environment and sets them by
// pointer into scope. Returns an error if something went wrong.
//
// Supported types: int, int8, int16, int32, int64, uint, uint8, uint16,
// uint32, uint64, bool, float32, float64, string, and slice from thous types.
func unmarshalENV(scope interface{}) error {
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
		return fmt.Errorf("cannot use scope (type %s) as type *%s "+
			"in argument to decode", rt, rt)
	}

	// Scope validation.
	switch {
	case rt.Kind() != reflect.Struct:
		return ObjectError
	case !rv.IsValid():
		return InitializedError
	}

	// Implements Unmarshaler interface.
	if rp.Type().Implements(reflect.TypeOf((*Unmarshaller)(nil)).Elem()) {
		// If there is the custom method, MarshlaENV - run it.
		if m := rp.MethodByName("UnmarshalENV"); m.IsValid() {
			tmp := m.Call([]reflect.Value{})
			err := tmp[0].Interface()
			if err != nil {
				return fmt.Errorf("unmarshal: %v", err)
			}
			return nil
		}
	}

	// Walk through all the fields of the transferred object.
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Type().Field(i)
		name, sep := parseTag(field.Tag.Get("env"), field.Name, " ")

		// Change value.
		instance := rv.FieldByName(field.Name)
		kind := field.Type.Kind()
		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64:
			r, err := strToIntKind(Get(name), kind)
			if err != nil {
				return err
			}
			instance.SetInt(r)
		case reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64:
			r, err := strToUintKind(Get(name), kind)
			if err != nil {
				return err
			}
			instance.SetUint(r)
		case reflect.Float32, reflect.Float64:
			r, err := strToFloatKind(Get(name), kind)
			if err != nil {
				return err
			}
			instance.SetFloat(r)
		case reflect.Bool:
			r, err := strToBool(Get(name))
			if err != nil {
				return err
			}
			instance.SetBool(r)
		case reflect.String:
			instance.SetString(Get(name))
		case reflect.Array:
			max := instance.Type().Len()
			seq := strings.Split(Get(name), sep)
			if len(seq) > max {
				return fmt.Errorf("%d items overwhelms the [%d]array",
					len(seq), max)
			}

			err := setSequence(&instance, strings.Split(Get(name), sep))
			if err != nil {
				return err
			}
		case reflect.Slice:
			seq := strings.Split(Get(name), sep)
			tmp := reflect.MakeSlice(instance.Type(), len(seq), len(seq))
			err := setSequence(&tmp, strings.Split(Get(name), sep))
			if err != nil {
				return err
			}

			instance.Set(reflect.AppendSlice(instance, tmp))
		case reflect.TypeOf(&url.URL{}).Kind():
			u, err := url.Parse(Get(name))
			if err != nil {
				return err
			}

			instance.Set(reflect.ValueOf(u))
		case reflect.TypeOf(url.URL{}).Kind():
			u, err := url.Parse(Get(name))
			if err != nil {
				return err
			}

			instance.Set(reflect.ValueOf(*u))
		default:
			return TypeError
		} // switch
	} // for

	return nil
}

// getSequence get sequence as string.
func getSequence(instance *reflect.Value, sep string) (string, error) {
	var kind reflect.Kind
	var max int

	switch instance.Kind() {
	case reflect.Array:
		kind = instance.Index(0).Kind()
		max = instance.Type().Len()
	case reflect.Slice:
		tmp := reflect.MakeSlice(instance.Type(), 1, 1)
		kind = tmp.Index(0).Kind()
		max = instance.Len()
	default:
		return "", TypeError
	}

	switch kind {
	case reflect.TypeOf(&url.URL{}).Kind():
		var tmp = []string{}
		for i := 0; i < max; i++ {
			v := instance.Index(i).Elem().Interface().(url.URL)
			tmp = append(tmp, v.String())
		}
		str := strings.Replace(fmt.Sprint(tmp), " ", sep, -1)
		return strings.Trim(str, "[]"+sep), nil
	case reflect.TypeOf(url.URL{}).Kind():
		var tmp = []string{}
		for i := 0; i < max; i++ {
			v := instance.Index(i).Interface().(url.URL)
			tmp = append(tmp, v.String())
		}
		str := strings.Replace(fmt.Sprint(tmp), " ", sep, -1)
		return strings.Trim(str, "[]"+sep), nil
	}
	str := strings.Replace(fmt.Sprint(*instance), " ", sep, -1)
	return strings.Trim(str, "[]"+sep), nil

}

// setSlice sets slice into instance.
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
