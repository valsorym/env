package env

import (
	"fmt"
	"reflect"
	"strings"
)

// parseTag returns tag parameters as [NAME[, SEP]] where
//     NAME variable name in the environment;
//     SEP  separator for the list (only for arrays and slices).
func parseTag(tagValue, defaultName, defaultSep string) (name, sep string) {
	var data = strings.Split(tagValue, ",")

	switch len(data) {
	case 0:
		name, sep = defaultName, defaultSep
	case 1:
		name, sep = strings.TrimSpace(data[0]), defaultSep
	default: // more then 1
		name, sep = strings.TrimSpace(data[0]), strings.TrimSpace(data[1])
	}

	if len(name) == 0 { // the name must be at least one character
		name = defaultName
	}

	if len(sep) == 0 { // the sep must be at least one character
		sep = defaultSep
	}

	return
}

// unmarshalENV gets variables from the environment and sets them by
// pointer into scope. Returns an error if something went wrong.
//
// Supported types: Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16,
// Uint32, Uint64, Bool, Float32, Float64, String, Array, Slice.
func unmarshalENV(scope interface{}) error {
	var rv reflect.Value

	// The object must be a pointer.
	rv = reflect.ValueOf(scope)
	if rv.Type().Kind() != reflect.Ptr {
		t := rv.Type()
		return fmt.Errorf("cannot use scope (type %s) as type *%s "+
			"in argument to decode", t, t)
	}

	// Call custom UnmarshalENV method.
	if rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("recipient must be initialized struct")
	} else if cue := rv.MethodByName("UnmarshalENV"); cue.IsValid() {
		// If the structure has a custom MethodByName method.
		cue.Call([]reflect.Value{})
	}

	// Walk through all the fields of the transferred object.
	rv = rv.Elem()
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
		case reflect.Slice:
			tmp := reflect.MakeSlice(instance.Type(), 1, 1)
			err := setSlice(
				&instance,
				strings.Split(Get(name), sep),
				tmp.Index(0).Kind(),
			)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("incorrect type")
		} // switch
	} // for

	return nil
}

// setSlice sets slice into instance.
func setSlice(instance *reflect.Value,
	seq []string, kind reflect.Kind) (err error) {
	var (
		intSeq    []int64
		uintSeq   []uint64
		floatSeq  []float64
		stringSeq []string
		boolSeq   []bool
	)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	// Convert to correct type slice.
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		intSeq = make([]int64, 0, len(seq))
		for _, value := range seq {
			r, err := strToIntKind(value, kind)
			if err != nil {
				return err
			}
			intSeq = append(intSeq, r)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		uintSeq = make([]uint64, 0, len(seq))
		for _, value := range seq {
			r, err := strToUintKind(value, kind)
			if err != nil {
				return err
			}
			uintSeq = append(uintSeq, r)
		}
	case reflect.Float32, reflect.Float64:
		floatSeq = make([]float64, 0, len(seq))
		for _, value := range seq {
			r, err := strToFloatKind(value, kind)
			if err != nil {
				return err
			}
			floatSeq = append(floatSeq, r)
		}
	case reflect.Bool:
		boolSeq = make([]bool, 0, len(seq))
		for _, value := range seq {
			r, err := strToBool(value)
			if err != nil {
				return err
			}
			boolSeq = append(boolSeq, r)
		}
	case reflect.String:
		stringSeq = seq
	default:
		return fmt.Errorf("incorrect type %v\n", kind)
	}

	// Set correct value.
	switch kind {
	case reflect.Int:
		for _, v := range intSeq {
			value := reflect.ValueOf(int(v))
			instance.Set(reflect.Append(*instance, value))
		}
	case reflect.Int8:
		for _, v := range intSeq {
			value := reflect.ValueOf(int8(v))
			instance.Set(reflect.Append(*instance, value))
		}
	case reflect.Int16:
		for _, v := range intSeq {
			value := reflect.ValueOf(int16(v))
			instance.Set(reflect.Append(*instance, value))
		}
	case reflect.Int32:
		for _, v := range intSeq {
			value := reflect.ValueOf(int32(v))
			instance.Set(reflect.Append(*instance, value))
		}
	case reflect.Int64:
		for _, v := range intSeq {
			value := reflect.ValueOf(int64(v))
			instance.Set(reflect.Append(*instance, value))
		}
	case reflect.Uint:
		for _, v := range uintSeq {
			value := reflect.ValueOf(uint(v))
			instance.Set(reflect.Append(*instance, value))
		}
	case reflect.Uint8:
		for _, v := range uintSeq {
			value := reflect.ValueOf(uint8(v))
			instance.Set(reflect.Append(*instance, value))
		}
	case reflect.Uint16:
		for _, v := range uintSeq {
			value := reflect.ValueOf(uint16(v))
			instance.Set(reflect.Append(*instance, value))
		}
	case reflect.Uint32:
		for _, v := range uintSeq {
			value := reflect.ValueOf(uint32(v))
			instance.Set(reflect.Append(*instance, value))
		}
	case reflect.Uint64:
		for _, v := range uintSeq {
			value := reflect.ValueOf(uint64(v))
			instance.Set(reflect.Append(*instance, value))
		}
	case reflect.Float32:
		for _, v := range floatSeq {
			value := reflect.ValueOf(float32(v))
			instance.Set(reflect.Append(*instance, value))
		}
	case reflect.Float64:
		for _, v := range floatSeq {
			value := reflect.ValueOf(float64(v))
			instance.Set(reflect.Append(*instance, value))
		}
	case reflect.Bool:
		for _, v := range boolSeq {
			value := reflect.ValueOf(bool(v))
			instance.Set(reflect.Append(*instance, value))
		}
	case reflect.String:
		for _, v := range stringSeq {
			value := reflect.ValueOf(string(v))
			instance.Set(reflect.Append(*instance, value))
		}
	}

	return nil
}
