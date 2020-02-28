package env

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

const FloatAccuracy = 1e-7

// strToIntKind convert string to int64 type with checking for conversion
// to intX type. Returns default value for int type if value is empty.
//
// P.s. The intX determined by reflect.Kind.
func strToIntKind(value string, kind reflect.Kind) (r int64, err error) {
	// For empty string returns zero.
	if len(value) == 0 {
		return 0, nil
	}

	// Convert string to int64.
	r, err = strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	switch kind {
	case reflect.Int:
		// For 32-bit platform it is necessary to check overflow.
		if strconv.IntSize == 32 {
			if r < math.MinInt32 || r > math.MaxInt32 {
				return 0, fmt.Errorf("%d overflows int (int32)", r)
			}
		}
	case reflect.Int8:
		if r < math.MinInt8 || r > math.MaxInt8 {
			return 0, fmt.Errorf("%d overflows int8", r)
		}
	case reflect.Int16:
		if r < math.MinInt16 || r > math.MaxInt16 {
			return 0, fmt.Errorf("%d overflows int16", r)
		}
	case reflect.Int32:
		if r < math.MinInt32 || r > math.MaxInt32 {
			return 0, fmt.Errorf("%d overflows int32", r)
		}
	case reflect.Int64:
		// pass
	default:
		r, err = 0, fmt.Errorf("incorrect kind %v", kind)
	}

	return
}

// strToUintKind convert string to uint64 type with checking for conversion
// to uintX type. Returns default value for uint type if value is empty.
//
// P.s. The uintX determined by reflect.Kind.
func strToUintKind(value string, kind reflect.Kind) (r uint64, err error) {
	// For empty string returns zero.
	if len(value) == 0 {
		return 0, nil
	}

	// Convert string to uint64.
	r, err = strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}

	switch kind {
	case reflect.Uint:
		// For 32-bit platform it is necessary to check overflow.
		if strconv.IntSize == 32 && r > math.MaxUint32 {
			return 0, fmt.Errorf("%d overflows uint (uint32)", r)
		}
	case reflect.Uint8:
		if r > math.MaxUint8 {
			return 0, fmt.Errorf("%d overflows uint8", r)
		}
	case reflect.Uint16:
		if r > math.MaxUint16 {
			return 0, fmt.Errorf("%d overflows uint16", r)
		}
	case reflect.Uint32:
		if r > math.MaxUint32 {
			return 0, fmt.Errorf("strToUint32: %d overflows uint32", r)
		}
	case reflect.Uint64:
		// pass
	default:
		r, err = 0, fmt.Errorf("incorrect kind %v", kind)
	}

	return
}

// strToFloatKind convert string to float64 type with checking for conversion
// to floatX type. Returns default value for float64 type if value is empty.
//
// P.s. The floatX determined by reflect.Kind.
func strToFloatKind(value string, kind reflect.Kind) (r float64, err error) {
	// For empty string returns zero.
	if len(value) == 0 {
		return 0.0, nil
	}

	// Convert string to Float64.
	r, err = strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0, err
	}

	switch kind {
	case reflect.Float32:
		if math.Abs(r) > math.MaxFloat32 {
			return 0.0, fmt.Errorf("%f overflows float32", r)
		}
	case reflect.Float64:
		// pass
	default:
		r, err = 0, fmt.Errorf("incorrect kind")
	}

	return
}

// strToBool convert string to bool type. Returns: result, error.
// Returns default value for bool type if value is empty.
func strToBool(value string) (bool, error) {
	// For empty string returns false.
	if len(value) == 0 {
		return false, nil
	}

	r, errB := strconv.ParseBool(value)
	if errB != nil {
		f, errF := strconv.ParseFloat(value, 64)
		if errF != nil {
			return r, errB
		}

		if math.Abs(f) > FloatAccuracy {
			r = true
		}
	}

	return bool(r), nil
}

/*
// strToInt64 convert string to int64 type with checking for conversion
// to int64 type. Returns default value for int type if value is empty.
func strToInt64(value string) (int64, error) {
	if len(value) == 0 {
		return 0, nil
	}

	return strconv.ParseInt(value, 10, 64)
}

// strToInt32 convert string to int64 type with checking for conversion
// to int32 type. Returns default value for int type if value is empty.
func strToInt32(value string) (int64, error) {
	r, err := strToInt64(value)
	if err != nil {
		return 0, err
	}

	if r >= math.MaxInt32 {
		return 0, fmt.Errorf("strToInt32: %d overflows int32", r)
	}

	return r, err
}

// strToInt16 convert string to int64 type with checking for conversion
// to int16 type. Returns default value for int type if value is empty.
func strToInt16(value string) (int64, error) {
	r, err := strToInt64(value)
	if err != nil {
		return 0, err
	}

	if r >= math.MaxInt16 {
		return 0, fmt.Errorf("strToInt16: %d overflows int16", r)
	}

	return r, err
}

// strToInt8 convert string to int64 type with checking for conversion
// to int8 type. Returns default value for int type if value is empty.
func strToInt8(value string) (int64, error) {
	r, err := strToInt64(value)
	if err != nil {
		return 0, err
	}

	if r >= math.MaxInt8 {
		return 0, fmt.Errorf("strToInt8: %d overflows int8", r)
	}

	return r, err
}

// strToInt convert string to int64 type with checking for conversion
// to int type. Returns default value for int type if value is empty.
func strToInt(value string) (int64, error) {
	r, err := strToInt64(value)
	if err != nil {
		return 0, err
	}

	// For 32-bit platform it is necessary to check overflow. Overflow for
	// 64-bit platform will be generated by the strToInt64 function.
	if strconv.IntSize == 32 && r >= math.MaxInt32 {
		return 0, fmt.Errorf("strToInt: %d overflows int (int32)", r)
	}

	return r, err
}

// strToIntKind convert string to int64 type with checking for conversion
// to intX type. Returns default value for int type if value is empty.
//
// The intX determined by reflect.Kind.
func strToIntKind(value string, kind reflect.Kind) (r int64, err error) {
	switch kind {
	case reflect.Int:
		r, err = strToInt(value)
	case reflect.Int8:
		r, err = strToInt8(value)
	case reflect.Int16:
		r, err = strToInt16(value)
	case reflect.Int32:
		r, err = strToInt32(value)
	case reflect.Int64:
		r, err = strToInt64(value)
	default:
		r, err = 0, fmt.Errorf("incorrect kind")
	}

	return
}
*/

/*
// strToUint64 convert string to uint64 type with checking for conversion
// to uint64 type. Returns default value for uint type if value is empty.
func strToUint64(value string) (uint64, error) {
	if len(value) == 0 {
		return 0, nil
	}

	return strconv.ParseUint(value, 10, 64)
}

// strToUint64 convert string to uint64 type with checking for conversion
// to uint32 type. Returns default value for uint type if value is empty.
func strToUint32(value string) (uint64, error) {
	r, err := strToUint64(value)
	if err != nil {
		return 0, err
	}

	if r >= math.MaxUint32 {
		return 0, fmt.Errorf("strToUint32: %d overflows uint32", r)
	}

	return r, err
}

// strToUint64 convert string to uint64 type with checking for conversion
// to uint16 type. Returns default value for uint type if value is empty.
func strToUint16(value string) (uint64, error) {
	r, err := strToUint64(value)
	if err != nil {
		return 0, err
	}

	if r >= math.MaxUint16 {
		return 0, fmt.Errorf("strToUint16: %d overflows uint16", r)
	}

	return r, err
}

// strToUint64 convert string to uint64 type with checking for conversion
// to uint8 type. Returns default value for uint type if value is empty.
func strToUint8(value string) (uint64, error) {
	r, err := strToUint64(value)
	if err != nil {
		return 0, err
	}

	if r >= math.MaxUint8 {
		return 0, fmt.Errorf("strToUint8: %d overflows uint8", r)
	}

	return r, err
}

// strToUint convert string to uint64 type with checking for conversion
// to uint type. Returns default value for uint type if value is empty.
func strToUint(value string) (uint64, error) {
	r, err := strToUint64(value)
	if err != nil {
		return 0, err
	}

	// For 32-bit platform it is necessary to check overflow. Overflow for
	// 64-bit platform will be generated by the strToUint64 function.
	if strconv.IntSize == 32 && r >= math.MaxUint32 {
		return 0, fmt.Errorf("strToUint: %d overflows uint (uint32)", r)
	}

	return r, err
}

// strToUintKind convert string to uint64 type with checking for conversion
// to uintX type. Returns default value for uint type if value is empty.
//
// The uintX determined by reflect.Kind.
func strToUintKind(value string, kind reflect.Kind) (r uint64, err error) {
	switch kind {
	case reflect.Uint:
		r, err = strToUint(value)
	case reflect.Uint8:
		r, err = strToUint8(value)
	case reflect.Uint16:
		r, err = strToUint16(value)
	case reflect.Uint32:
		r, err = strToUint32(value)
	case reflect.Uint64:
		r, err = strToUint64(value)
	default:
		r, err = 0, fmt.Errorf("incorrect kind")
	}

	return
}
*/

/*
// strToFloat64 convert string to float64 type with checking for conversion
// to float64 type. Returns default value for float64 type if value is empty.
func strToFloat64(value string) (float64, error) {
	if len(value) == 0 {
		return 0.0, nil
	}

	return strconv.ParseFloat(value, 64)
}

// strToFloat32 convert string to float64 type with checking for conversion
// to float32 type. Returns default value for float64 type if value is empty.
func strToFloat32(value string) (float64, error) {
	// Parse value from the environment.
	r, err := strToFloat64(value)
	if err != nil {
		return r, err
	}

	if r > math.MaxFloat32 {
		return 0, fmt.Errorf("strToFloat32: %f overflows float32", r)
	}

	return r, nil
}

// strToFloatKind convert string to float64 type with checking for conversion
// to floatX type. Returns default value for float64 type if value is empty.
//
// The floatX determined by reflect.Kind.
func strToFloatKind(value string, kind reflect.Kind) (r float64, err error) {
	switch kind {
	case reflect.Float64:
		r, err = strToFloat64(value)
	case reflect.Float32:
		r, err = strToFloat32(value)
	default:
		r, err = 0, fmt.Errorf("incorrect kind")
	}

	return
}
*/
