package env

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	emptyRegex = regexp.MustCompile(`^(\s*)$|^(\s*[#].*)$`)
	valueRegex = regexp.MustCompile(`^=[^\s].*`)
	keyRegex   = regexp.MustCompile(
		`^(?:\s*)?(?:export\s+)?(?P<key>[a-zA-Z_][a-zA-Z_0-9]*)=`,
	)
)

// isEmpty returns true if string contains separators or comment only.
func isEmpty(str string) bool {
	return emptyRegex.Match([]byte(str))
}

// removeInlineComment removes the comment in the string.
// Only in strings where the value is enclosed in quotes.
func removeInlineComment(str, quote string) string {
	// If the comment is in the string.
	if strings.Contains(str, "#") {
		chunks := strings.Split(str, "#")
		for i := range chunks {
			str := strings.Join(chunks[:i], "#")
			if len(str) > 0 && strings.Count(str, quote)%2 == 0 {
				return strings.TrimSpace(str)
			}
		}
	}
	return str
}

// parseExpression breaks expression into key and value, ignore
// comments and any spaces.
//
// Note: value must be an expression.
func parseExpression(exp string) (key, value string, err error) {
	var (
		quote  string = "\""
		marker string = fmt.Sprintf("<::%d::>", time.Now().Unix())
	)

	// Get key.
	// Remove `export` prefix, `=` suffix and trim spaces.
	tmp := keyRegex.FindStringSubmatch(exp)
	if len(tmp) < 2 {
		err = KeyError
		return
	}
	key = tmp[1]

	// Get value.
	// ... the `=` sign in the string.
	value = exp[strings.Index(exp, "="):]
	if !valueRegex.Match([]byte(value)) {
		err = ValueError
		return
	}
	value = strings.TrimSpace(value[1:])

	switch {
	case strings.HasPrefix(value, "'"):
		quote = "'"
		fallthrough
	case strings.HasPrefix(value, "\""):
		// Replace escaped quotes, remove comment in the string,
		// check begin- and end- quotes and back escaped quotes.
		value = strings.Replace(value, fmt.Sprintf("\\%s", quote), marker, -1)
		value = removeInlineComment(value, quote)
		if strings.Count(value, quote)%2 != 0 { // begin- and end- quotes
			err = ValueError
			return
		}
		value = value[1 : len(value)-1] // remove begin- and end- quotes
		// ... change `\"` and `\'` to `"` and `'`.
		value = strings.Replace(value, marker, fmt.Sprintf("%s", quote), -1)
	default:
		if strings.Contains(value, "#") {
			// Split by sharp sign and for string without quotes -
			// the first element has the meaning only.
			chunks := strings.Split(value, "#")
			chunks = strings.Split(chunks[0], " ")
			value = strings.TrimSpace(chunks[0])
		}
	}

	return
}

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
	var epsilon = math.Nextafter(1, 2) - 1

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

		if math.Abs(f) > epsilon {
			r = true
		}
	}

	return bool(r), nil
}
