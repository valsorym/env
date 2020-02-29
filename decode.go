package env

import (
	"fmt"
	"reflect"
	"strings"
)

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
		return fmt.Errorf("object must be a structure")
	case !rv.IsValid():
		return fmt.Errorf("object must be initialized")
	}

	// If there is the custom method, MarshlaENV - run it.
	if m := rp.MethodByName("UnmarshalENV"); m.IsValid() {
		tmp := m.Call([]reflect.Value{})
		if len(tmp) != 0 {
			err := tmp[0].Interface()
			if err != nil {
				return fmt.Errorf("marshal: %v", err)
			}
		}
		return nil
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
