package env

import (
	"fmt"
	"strings"
	"testing"
)

type Data struct {
	Host  string `env:"HOST"`
	Port  uint   `env:"PORT,"`
	Paths string `env:"PATH,:"`
}

type DataNumber struct {
	KeyInt     int     `env:"KEY_INT"`
	KeyInt8    int8    `env:"KEY_INT8"`
	KeyInt16   int16   `env:"KEY_INT16"`
	KeyInt32   int32   `env:"KEY_INT32"`
	KeyInt64   int64   `env:"KEY_INT64"`
	KeyUint    uint    `env:"KEY_UINT"`
	KeyUint8   uint8   `env:"KEY_UINT8"`
	KeyUint16  uint16  `env:"KEY_UINT16"`
	KeyUint32  uint32  `env:"KEY_UINT32"`
	KeyUint64  uint64  `env:"KEY_UINT64"`
	KeyFloat32 float32 `env:"KEY_FLOAT32"`
	KeyFloat64 float64 `env:"KEY_FLOAT64"`
}

type DataBool struct {
	KeyBool bool `env:"KEY_BOOL"`
}

type DataString struct {
	KeyString string `env:"KEY_STRING"`
}

// TestParseTag tests parseTag function.
func TestParseTag(t *testing.T) {
	var tests = [][]string{
		[]string{"HOST", "host", " ", "HOST", " "},
		[]string{"PATHS,:", "paths", " ", "PATHS", ":"},
		[]string{",:", "PORT", " ", "PORT", ":"},
		[]string{",", "PORT", ":", "PORT", ":"},
	}

	for _, test := range tests {
		name, sep := parseTag(test[0], test[1], test[2])
		if test[3] != name {
			t.Errorf("incorrect value for name `%s`!=`%s`", test[3], name)
		}

		if test[4] != sep {
			t.Errorf("incorrect value for sep `%s`!=`%s`", test[4], sep)
		}
	}
}

// TestDecodeEnvironNumber tests decodeEnviron function
// for Int, Uint and Float types.
func TestDeocdeEnvironNumber(t *testing.T) {
	var max = "922337203685477580777"
	var tests = map[string][]string{
		"KEY_INT":     []string{"2", "-2", max},
		"KEY_INT8":    []string{"8", "-8", max},
		"KEY_INT16":   []string{"16", "-16", max},
		"KEY_INT32":   []string{"32", "-32", max},
		"KEY_INT64":   []string{"64", "-64", max},
		"KEY_UINT":    []string{"2", "-2", max},
		"KEY_UINT8":   []string{"8", "-8", max},
		"KEY_UINT16":  []string{"16", "-16", max},
		"KEY_UINT32":  []string{"32", "-32", max},
		"KEY_UINT64":  []string{"64", "-64", max},
		"KEY_FLOAT32": []string{"32.0", "-32.0", max},
		"KEY_FLOAT64": []string{"64.0", "-64.0", max},
	}

	// Correct value.
	for i := 0; i < 3; i++ {
		for key, data := range tests {
			var d = &DataNumber{}

			Clear()
			Set(key, data[i])

			err := decodeEnviron(d)
			switch i {
			case 0:
				if err != nil {
					t.Error(err)
				}
			case 1:
				if !strings.Contains(key, "UINT") {
					// Int and float types.
					if err != nil {
						t.Error(err)
					}
				} else {
					// Uint cannot be negative.
					if err == nil {
						t.Errorf("uint cannot be negative: %s", data[i])
					}
					continue
				}
			case 2:
				// Ignore FloatX to check for `value out of range`.
				if !strings.Contains(key, "FLOAT") {
					if err == nil {
						t.Errorf("for %s must be `value out of "+
							"range` exception", key)
					}
				}
				continue
			}

			switch key {
			case "KEY_INT":
				if v := fmt.Sprintf("%d", d.KeyInt); v != data[i] {
					t.Errorf("value isn't correct `%s`!=`%s`", v, data[i])
				}
			case "KEY_INT8":
				if v := fmt.Sprintf("%d", d.KeyInt8); v != data[i] {
					t.Errorf("value isn't correct `%s`!=`%s`", v, data[i])
				}
			case "KEY_INT16":
				if v := fmt.Sprintf("%d", d.KeyInt16); v != data[i] {
					t.Errorf("value isn't correct `%s`!=`%s`", v, data[i])
				}
			case "KEY_INT32":
				if v := fmt.Sprintf("%d", d.KeyInt32); v != data[i] {
					t.Errorf("value isn't correct `%s`!=`%s`", v, data[i])
				}
			case "KEY_INT64":
				if v := fmt.Sprintf("%d", d.KeyInt64); v != data[i] {
					t.Errorf("value isn't correct `%s`!=`%s`", v, data[i])
				}
			case "KEY_UINT":
				if v := fmt.Sprintf("%d", d.KeyUint); v != data[i] {
					t.Errorf("value isn't correct `%s`!=`%s`", v, data[i])
				}
			case "KEY_UINT8":
				if v := fmt.Sprintf("%d", d.KeyUint8); v != data[i] {
					t.Errorf("value isn't correct `%s`!=`%s`", v, data[i])
				}
			case "KEY_UINT16":
				if v := fmt.Sprintf("%d", d.KeyUint16); v != data[i] {
					t.Errorf("value isn't correct `%s`!=`%s`", v, data[i])
				}
			case "KEY_UINT32":
				if v := fmt.Sprintf("%d", d.KeyUint32); v != data[i] {
					t.Errorf("value isn't correct `%s`!=`%s`", v, data[i])
				}
			case "KEY_UINT64":
				if v := fmt.Sprintf("%d", d.KeyUint64); v != data[i] {
					t.Errorf("value isn't correct `%s`!=`%s`", v, data[i])
				}
			} // switch
		}
	}
}

// TestDecodeEnvironBoll tests decodeEnviron function for bool type.
func TestDeocdeEnvironBool(t *testing.T) {
	var tests = map[string]bool{
		"true":  true,
		"false": false,
		"0":     false,
		"1":     true,
		"":      false,
		"True":  true,
		"TRUE":  true,
		"False": false,
		"FALSE": false,
	}

	// Test correct values.
	for value, test := range tests {
		var d = &DataBool{}

		Clear()
		Set("KEY_BOOL", value)

		err := decodeEnviron(d)
		if err != nil {
			t.Error(err)
		}

		if d.KeyBool != test {
			t.Errorf("KeyBool == %t but need %t", d.KeyBool, test)
		}
	}

	// Incorrect value.
	for _, value := range []string{"string", "0.d", "true/false"} {
		var d = &DataBool{}

		Clear()
		Set("KEY_BOOL", value)

		err := decodeEnviron(d)
		if err == nil {
			t.Error("didn't handle the error")
		}
	}
}

// TestDecodeEnvironString tests decodeEnviron function for string type.
func TestDeocdeEnvironString(t *testing.T) {
	var tests = []interface{}{
		8080,
		"Hello World",
		"true",
		true,
		3.14,
	}

	// Test correct values.
	for _, test := range tests {
		var d = &DataString{}
		var s = fmt.Sprintf("%v", test)

		Clear()
		Set("KEY_STRING", s)

		err := decodeEnviron(d)
		if err != nil {
			t.Error(err)
		}

		if d.KeyString != s {
			t.Errorf("KeyString == `%s` but need `%s`", d.KeyString, s)
		}
	}
}