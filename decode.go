package env

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

const TagName = "env"

// getInt64 to parse string and get int64 value or error.
// Returns `0` for empty string.
func getInt64(value string) (int64, error) {
	var (
		r   int64
		err error
	)

	if len(value) != 0 {
		r, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0, err
		}
	}

	return r, nil
}

// getUint64 to parse string and get uint64 value or error.
// Returns `0` for empty string.
func getUint64(value string) (uint64, error) {
	var (
		r   uint64
		err error
	)

	if len(value) != 0 {
		r, err = strconv.ParseUint(value, 10, 64)
		if err != nil {
			return 0, err
		}
	}

	return r, nil
}

// getFloat64 to parse string and get float64 value or error.
// Returns `0` for empty string.
func getFloat64(value string) (float64, error) {
	var (
		r   float64
		err error
	)

	if len(value) != 0 {
		r, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, err
		}
	}

	return r, nil
}

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
			// err := setInt(&instance, name)
			// if err != nil {
			// 	return err
			// }

			r, err := strToInt(Get(name))
			if err != nil {
				return err
			}

			instance.SetInt(int64(r))
		case reflect.Int8:
			// err := setInt8(&instance, name)
			// if err != nil {
			// 	return err
			// }
			r, err := strToInt8(Get(name))
			if err != nil {
				return err
			}

			instance.SetInt(int64(r))
		case reflect.Int16:
			// err := setInt16(&instance, name)
			// if err != nil {
			// 	return err
			// }
			r, err := strToInt16(Get(name))
			if err != nil {
				return err
			}

			instance.SetInt(int64(r))
		case reflect.Int32:
			// err := setInt32(&instance, name)
			// if err != nil {
			// 	return err
			// }
			r, err := strToInt32(Get(name))
			if err != nil {
				return err
			}

			instance.SetInt(int64(r))
		case reflect.Int64:
			// err := setInt64(&instance, name)
			// if err != nil {
			// 	return err
			// }
			r, err := strToInt64(Get(name))
			if err != nil {
				return err
			}

			instance.SetInt(int64(r))
		case reflect.Uint:
			err := setUint(&instance, name)
			if err != nil {
				return err
			}
		case reflect.Uint8:
			err := setUint8(&instance, name)
			if err != nil {
				return err
			}
		case reflect.Uint16:
			err := setUint16(&instance, name)
			if err != nil {
				return err
			}
		case reflect.Uint32:
			err := setUint32(&instance, name)
			if err != nil {
				return err
			}
		case reflect.Uint64:
			err := setUint64(&instance, name)
			if err != nil {
				return err
			}
		case reflect.Float32:
			err := setFloat32(&instance, name)
			if err != nil {
				return err
			}
		case reflect.Float64:
			err := setFloat64(&instance, name)
			if err != nil {
				return err
			}
		case reflect.Bool:
			err := setBool(&instance, name)
			if err != nil {
				return err
			}
		case reflect.String:
			err := setString(&instance, name)
			if err != nil {
				return err
			}
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

// setInt try to set `int` to instance.
func setInt(instance *reflect.Value, name string) error {
	// Parse value from the environment.
	r, err := getInt64(Get(name))
	if err != nil {
		return err
	}

	if strconv.IntSize == 32 && math.MaxInt32 < r {
		return fmt.Errorf("%s \"%d\": value out of range", name, r)
	}

	if strconv.IntSize == 64 && math.MaxInt64 < r {
		return fmt.Errorf("%s \"%d\": value out of range", name, r)
	}

	instance.SetInt(r)
	return nil
}

// setInt8 try to set `int8` to instance.
func setInt8(instance *reflect.Value, name string) error {
	// Parse value from the environment.
	r, err := getInt64(Get(name))
	if err != nil {
		return err
	}

	if math.MaxInt8 < r {
		return fmt.Errorf("%s \"%d\": value out of range", name, r)
	}

	instance.SetInt(r)
	return nil
}

// setInt16 try to set `int16` to instance.
func setInt16(instance *reflect.Value, name string) error {
	// Parse value from the environment.
	r, err := getInt64(Get(name))
	if err != nil {
		return err
	}

	if math.MaxInt16 < r {
		return fmt.Errorf("%s \"%d\": value out of range", name, r)
	}

	instance.SetInt(r)
	return nil
}

// setInt32 try to set `int32` to instance.
func setInt32(instance *reflect.Value, name string) error {
	// Parse value from the environment.
	r, err := getInt64(Get(name))
	if err != nil {
		return err
	}

	if math.MaxInt32 < r {
		return fmt.Errorf("%s \"%d\": value out of range", name, r)
	}

	instance.SetInt(r)
	return nil
}

// setInt64 try to set `int64` to instance.
func setInt64(instance *reflect.Value, name string) error {
	// Parse value from the environment.
	r, err := getInt64(Get(name))
	if err != nil {
		return err
	}

	instance.SetInt(r)
	return nil
}

// setUint try to set `uint` to instance.
func setUint(instance *reflect.Value, name string) error {
	// Parse value from the environment.
	r, err := getUint64(Get(name))
	if err != nil {
		return err
	}

	if strconv.IntSize == 32 && math.MaxUint32 < r {
		return fmt.Errorf("%s \"%d\": value out of range", name, r)
	}

	if strconv.IntSize == 64 && math.MaxUint64 < r {
		return fmt.Errorf("%s \"%d\": value out of range", name, r)
	}

	instance.SetUint(r)
	return nil
}

// setUint8 try to set `uint8` to instance.
func setUint8(instance *reflect.Value, name string) error {
	// Parse value from the environment.
	r, err := getUint64(Get(name))
	if err != nil {
		return err
	}

	if math.MaxUint8 < r {
		return fmt.Errorf("%s \"%d\": value out of range", name, r)
	}

	instance.SetUint(r)
	return nil
}

// setUint16 try to set `uint16` to instance.
func setUint16(instance *reflect.Value, name string) error {
	// Parse value from the environment.
	r, err := getUint64(Get(name))
	if err != nil {
		return err
	}

	if math.MaxUint16 < r {
		return fmt.Errorf("%s \"%d\": value out of range", name, r)
	}

	instance.SetUint(r)
	return nil
}

// setUint32 try to set `uint32` to instance.
func setUint32(instance *reflect.Value, name string) error {
	// Parse value from the environment.
	r, err := getUint64(Get(name))
	if err != nil {
		return err
	}

	if math.MaxUint32 < r {
		return fmt.Errorf("%s \"%d\": value out of range", name, r)
	}

	instance.SetUint(r)
	return nil
}

// setUint64 try to set `uint64` to instance.
func setUint64(instance *reflect.Value, name string) error {
	// Parse value from the environment.
	r, err := getUint64(Get(name))
	if err != nil {
		return err
	}

	instance.SetUint(r)
	return nil
}

// setFloat32 try to set `float32` to instance.
func setFloat32(instance *reflect.Value, name string) error {
	// Parse value from the environment.
	r, err := getFloat64(Get(name))
	if err != nil {
		return err
	}

	if math.MaxFloat32 < r {
		return fmt.Errorf("%s \"%f\": value out of range", name, r)
	}

	instance.SetFloat(r)
	return nil
}

// setFloat64 try to set `float64` to instance.
func setFloat64(instance *reflect.Value, name string) error {
	// Parse value from the environment.
	r, err := getFloat64(Get(name))
	if err != nil {
		return err
	}

	instance.SetFloat(r)
	return nil
}

// setBool try to set `bool` to instance.
func setBool(instance *reflect.Value, name string) error {
	// Parse value from the environment.
	v := Get(name)
	if len(v) == 0 {
		instance.SetBool(false)
		return nil
	}

	// Try parse.
	b, err := strconv.ParseBool(v)
	if err != nil {
		f, errF := strconv.ParseFloat(v, 64)
		if errF != nil {
			return err
		}

		if math.Abs(f) > 1e-9 {
			b = true
		}
	}

	instance.SetBool(b)
	return nil
}

// setString try to set `string` to instance.
func setString(instance *reflect.Value, name string) error {
	v := Get(name)
	instance.SetString(v)
	return nil
}

// // setSlice try to set `slice` to instance.
// func setSlice(instance *reflect.Value, name, sep string) error {
// 	v := strings.Split(Get(name), sep)
//
// 	instance.SetSlice(v)
// 	return nil
// }
