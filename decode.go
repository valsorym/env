package env

import (
	"fmt"
	"reflect"
	"strings"
)

const TagName = "env"

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

// decodeEnviron gets variables from the environment and sets them by
// pointer into scope. Returns an error if something went wrong.
//
// Supported types: Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16,
// Uint32, Uint64, Bool, Float32, Float64, String, Array, Slice.
func decodeEnviron(scope interface{}) error {
	var rv reflect.Value

	// The object must be a pointer.
	rv = reflect.ValueOf(scope)
	if rv.Type().Kind() != reflect.Ptr {
		t := rv.Type()
		return fmt.Errorf("cannot use scope (type %s) as type *%s "+
			"in argument to decode", t, t)
	}

	// Get the value of an object.
	rv = rv.Elem()

	// Walk through all the fields of the transferred object.
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Type().Field(i)
		name, _ := parseTag(field.Tag.Get(TagName), field.Name, " ")

		// Change value.
		instance := rv.FieldByName(field.Name)
		switch field.Type.Kind() {
		case reflect.Int:
			r, err := strToInt(Get(name))
			if err != nil {
				return err
			}
			instance.SetInt(int64(r))
		case reflect.Int8:
			r, err := strToInt8(Get(name))
			if err != nil {
				return err
			}
			instance.SetInt(int64(r))
		case reflect.Int16:
			r, err := strToInt16(Get(name))
			if err != nil {
				return err
			}
			instance.SetInt(int64(r))
		case reflect.Int32:
			r, err := strToInt32(Get(name))
			if err != nil {
				return err
			}
			instance.SetInt(int64(r))
		case reflect.Int64:
			r, err := strToInt64(Get(name))
			if err != nil {
				return err
			}
			instance.SetInt(r)
		case reflect.Uint:
			r, err := strToUint(Get(name))
			if err != nil {
				return err
			}
			instance.SetUint(uint64(r))
		case reflect.Uint8:
			r, err := strToUint8(Get(name))
			if err != nil {
				return err
			}
			instance.SetUint(uint64(r))
		case reflect.Uint16:
			r, err := strToUint16(Get(name))
			if err != nil {
				return err
			}
			instance.SetUint(uint64(r))
		case reflect.Uint32:
			r, err := strToUint32(Get(name))
			if err != nil {
				return err
			}
			instance.SetUint(uint64(r))
		case reflect.Uint64:
			r, err := strToUint64(Get(name))
			if err != nil {
				return err
			}
			instance.SetUint(r)
		case reflect.Float32:
			r, err := strToFloat32(Get(name))
			if err != nil {
				return err
			}
			instance.SetFloat(float64(r))
		case reflect.Float64:
			r, err := strToFloat64(Get(name))
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

			// case reflect.Array:
			// 	t := instance.Index(0).Kind() // get type of the array
			// 	err := setArray(&instance, name, sep)
			// 	if err != nil {
			// 		return err
			// 	}
			// case reflect.Slice:
			// 	tmp := reflect.MakeSlice(instance.Type(), 1, 1)
			// 	t := tmp.Index(0).Kind() // get type of the slice
			// 	err := setSlice(&instance, name, sep)
			// 	if err != nil {
			// 		return err
			// 	}
		} // switch
	} // for

	return nil
}
