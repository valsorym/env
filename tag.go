package env

import (
	"fmt"
	"regexp"
	"strings"
)

//"regexp"
//"strings"

// The correctKeyRgx is regular expression
// to verify the correctness of the key name.
var correctKeyRgx = regexp.MustCompile(`^[A-Za-z_]{1}\w*$`)

// The splitFieldTag parse the field's tag and returns separate
// elements like: key, value, separator and error ID.
//
// Input:
//   - ft is field's tag as string;
//
// Output:
//   - key is environment variable name;
//   - value is default value (if key not exists in environment);
//   - sep is item separator (for lists only);
//   - err is error id.
func splitFieldTag(ft string) (key, value, sep string, err error) {
	var (
		scope  = []*string{&key, &value, &sep}
		covers = map[string]string{"'": "'", "\"": "\"", "{": "}"}
	)

	// Set default values.
	sep = ":" // classic element separator for environment variables

	// Get key and put right part (default value and separator) into value.
	for i, item := range strings.SplitN(ft, ",", 2) {
		*scope[i] = strings.Trim(item, " _")
	}

	// Checking key for correctness.
	if len(key) != 0 && !correctKeyRgx.Match([]byte(key)) {
		err = fmt.Errorf("incorrect key %s", key)
	}

	// The right part is missing or the key is incorrect.
	if len(value) == 0 || err != nil {
		return
	}

	// Now the value contains value data and help text.
	// The value can contains `,` symbol if is string or slice/array.
	// Note: If the value is a string or slice/array with a comma symbol,
	//       the value should be written as:
	//       - 'value here' or \"value here\" for string;
	//       - {s,l,i,c,e} for slice or array.

	// The value doesn't start from a special character: `'`, `"` or `{`.
	begin := string(value[0])
	end, ok := covers[begin]
	if !ok {
		// Value doesn't contains `,` symbol.
		for i, item := range strings.SplitN(value, ",", 2) {
			*scope[i+1] = strings.Trim(item, " _")
		}
		return
	}

	// Separate value and separator symbol.
	// Check for closing a special character.
	index := strings.Index(value[1:], end) + 1
	if index < 1 {
		err = fmt.Errorf("missing `%s` - closing character to value", end)
		return
	}

	// Define a separate string.
	tmp := strings.SplitN(value[index+1:], ",", 2)
	if p := index + 2 + len(tmp[0]); p <= len(value) {
		sep = strings.Trim(value[p:], " _")
	}

	// Define value string.
	value = strings.Trim(value[1:index]+tmp[0], " _")

	// Undefined remainder in the tail of the value data.
	if len(strings.Trim(tmp[0], " _")) != 0 {
		err = fmt.Errorf("undefined remainder in the tail "+
			"of the value: %s", tmp[0])
	}

	return
}

/*
// Rules for handling tags.
var (
	nameRegex = regexp.MustCompile(`^[A-Za-z_]{1}\w*$`)
)

// The splitFieldTag parse the field's tag and returns separate
// elements like: name, value, separator and error ID.
func splitFieldTag(ft, n, v, s string) (name, value, sep string, err error) {
	var (
		r = []*string{&name, &value, &sep}
		d = map[string]string{"'": "'", "\"": "\"", "{": "}"}
	)

	// Set default separator.
	name, value, sep = n, v, s

	// Split into two part only because value can have a ',' symbol.
	// Note: value contains right parts of the ft.
	for i, item := range strings.SplitN(ft, ",", 2) {
		*r[i] = strings.Trim(item, " _")
	}

	// Checking the name for correctness.
	if !nameRegex.Match([]byte(name)) {
		err = fmt.Errorf("incorrect variable name %s", name)
	}

	// The right part is missing or equal to value.
	if value == v || err != nil {
		return
	}

	// Now the value contains value data and help text. Value data can
	// contains `,` symbol if is string or slice/array.
	// Note: If the value is a string or slice/array with a comma symbol,
	//       the value should be written as:
	//       - 'value here' or \"value here\" for string;
	//       - {s,l,i,c,e} for slice or array.

	// The value doesn't start from a special character: ', " or {.
	begin := string(value[0])
	end, ok := d[begin]
	if !ok {
		for i, item := range strings.SplitN(value, ",", 2) {
			*r[i+1] = strings.Trim(item, " _")
		}
		return
	}

	// Separate value and separator.
	// Check for closing a special character.
	index := strings.Index(value[1:], end) + 1
	if index < 1 {
		err = fmt.Errorf("missing `%s` - closing character to value", end)
		return
	}

	// Define a separate string.
	tmp := strings.SplitN(value[index+1:], ",", 2)
	if p := index + 2 + len(tmp[0]); p <= len(value) {
		sep = strings.Trim(value[p:], " _")
	}

	// Define value string.
	value = strings.Trim(value[1:index]+tmp[0], " _")

	// Undefined remainder in the tail of the value data.
	if len(strings.Trim(tmp[0], " _")) != 0 {
		err = fmt.Errorf("undefined remainder in the tail "+
			"of the value: %s", tmp[0])
	}

	return
}
*/
