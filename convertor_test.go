package env

import (
	"fmt"
	"strings"
	"testing"
)

// NumberTestType structure for testing conversion of numeric types.
type NumberTestType struct {
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

// BoolTestType structure for testing conversion of boolean types.
type BoolTestType struct {
	KeyBool bool `env:"KEY_BOOL"`
}

// StringTestType structure for testing conversion of string types.
type StringTestType struct {
	KeyString string `env:"KEY_STRING"`
}

// SliceTestType structure for testing conversion of
// different type slices types.
type SliceTestType struct {
	KeyInt   []int   `env:"KEY_INT,:"`
	KeyInt8  []int8  `env:"KEY_INT8,:"`
	KeyInt16 []int16 `env:"KEY_INT16,:"`
	KeyInt32 []int32 `env:"KEY_INT32,:"`
	KeyInt64 []int64 `env:"KEY_INT64,:"`

	KeyUint   []uint   `env:"KEY_UINT,:"`
	KeyUint8  []uint8  `env:"KEY_UINT8,:"`
	KeyUint16 []uint16 `env:"KEY_UINT16,:"`
	KeyUint32 []uint32 `env:"KEY_UINT32,:"`
	KeyUint64 []uint64 `env:"KEY_UINT64,:"`

	KeyFloat32 []float32 `env:"KEY_FLOAT32,:"`
	KeyFloat64 []float64 `env:"KEY_FLOAT64,:"`

	KeyString []string `env:"KEY_STRING,:"`
	KeyBool   []bool   `env:"KEY_BOOL,:"`
}

// PlainTestType simple complex type.
type PlainTestType struct {
	Host         string   `env:"HOST"`
	Port         int      `env:"PORT"`
	AllowedHosts []string `env:"ALLOWED_HOSTS,:"`
}

// Extended simple complex type that implements
// Marshaller and Unmarshaller interfaces.
type ExtendedTestType struct {
	Host         string   `env:"HOST"`
	Port         int      `env:"PORT"`
	AllowedHosts []string `env:"ALLOWED_HOSTS,:"`
}

func (c *ExtendedTestType) MarshalENV() ([]string, error) {
	str := strings.Replace(fmt.Sprint(c.AllowedHosts), " ", ":", -1)
	Set("ALLOWED_HOSTS", strings.Trim(str, "[]"))
	Set("PORT", fmt.Sprintf("%d", c.Port))
	Set("HOST", c.Host)
	return []string{}, nil
}

// TestMarshalNotStruct tests unmarshalENV function for not struct values.
func TestMarshalNotStruct(t *testing.T) {
	var scope string
	_, err := Marshal(scope)
	if err == nil {
		t.Error("exception expected for an object other than structure")
	}
}

// TestMarshalPointNil tests unmarshalENV function for uninitialized pointer.
func TestMarshalPointNil(t *testing.T) {
	var scope *PlainTestType
	_, err := Marshal(scope)
	if err == nil {
		t.Error("exception expected for an uninitialized object")
	}
}

// TestUnmarshalENVNumber tests unmarshalENV function
// for Int, Uint and Float types.
func TestUnmarshalENVNumber(t *testing.T) {
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
			var d = &NumberTestType{}

			Clear()
			Set(key, data[i])

			err := unmarshalENV(d)
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

// TestUnmarshalENVBoll tests unmarshalENV function for bool type.
func TestUnmarshalENVBool(t *testing.T) {
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
		var d = &BoolTestType{}

		Clear()
		Set("KEY_BOOL", value)

		err := unmarshalENV(d)
		if err != nil {
			t.Error(err)
		}

		if d.KeyBool != test {
			t.Errorf("KeyBool == %t but need %t", d.KeyBool, test)
		}
	}

	// Incorrect value.
	for _, value := range []string{"string", "0.d", "true/false"} {
		var d = &BoolTestType{}

		Clear()
		Set("KEY_BOOL", value)

		err := unmarshalENV(d)
		if err == nil {
			t.Error("didn't handle the error")
		}
	}
}

// TestUnmarshalENVString tests unmarshalENV function for string type.
func TestUnmarshalENVString(t *testing.T) {
	var tests = []interface{}{
		8080,
		"Hello World",
		"true",
		true,
		3.14,
	}

	// Test correct values.
	for _, test := range tests {
		var d = &StringTestType{}
		var s = fmt.Sprintf("%v", test)

		Clear()
		Set("KEY_STRING", s)

		err := unmarshalENV(d)
		if err != nil {
			t.Error(err)
		}

		if d.KeyString != s {
			t.Errorf("KeyString == `%s` but need `%s`", d.KeyString, s)
		}
	}
}

// TestUnmarshalENVSliceCorrect tests unmarshalENV function
// for slice type with correct values.
func TestUnmarshalENVSliceCorrect(t *testing.T) {
	var tests = map[string]string{
		"KEY_INT":   "-30:-20:-10:0:10:20:30",
		"KEY_INT8":  "-30:-20:-10:0:10:20:30",
		"KEY_INT16": "-30:-20:-10:0:10:20:30",
		"KEY_INT32": "-30:-20:-10:0:10:20:30",
		"KEY_INT64": "-30:-20:-10:0:10:20:30",

		"KEY_UINT":   "0:10:20:30",
		"KEY_UINT8":  "0:10:20:30",
		"KEY_UINT16": "0:10:20:30",
		"KEY_UINT32": "0:10:20:30",
		"KEY_UINT64": "0:10:20:30",

		"KEY_FLOAT32": "-3.1:-1.27:0:1.27:3.3",
		"KEY_FLOAT64": "-3.1:-1.27:0:1.27:3.3",

		"KEY_STRING": "one:two:three:four:five",
		"KEY_BOOL":   "1:true:True:TRUE:0:false:False:False",
	}

	// Convert slice into string.
	toStr := func(v interface{}) string {
		return strings.Trim(strings.Replace(fmt.Sprint(v), " ", ":", -1), "[]")
	}

	// Testing.
	for key, value := range tests {
		var d = &SliceTestType{}

		Clear()
		Set(key, value)

		err := unmarshalENV(d)
		if err != nil {
			t.Error("xxxxxxxxxxxxxxxxxxxx:", err)
		}

		switch key {
		case "KEY_INT":
			if r := toStr(d.KeyInt); r != value {
				t.Errorf("KeyInt == `%s` but need `%s`", r, value)
			}
		case "KEY_INT8":
			if r := toStr(d.KeyInt8); r != value {
				t.Errorf("KeyInt8 == `%s` but need `%s`", r, value)
			}
		case "KEY_INT16":
			if r := toStr(d.KeyInt16); r != value {
				t.Errorf("KeyInt16 == `%s` but need `%s`", r, value)
			}
		case "KEY_INT32":
			if r := toStr(d.KeyInt32); r != value {
				t.Errorf("KeyInt32 == `%s` but need `%s`", r, value)
			}
		case "KEY_INT64":
			if r := toStr(d.KeyInt64); r != value {
				t.Errorf("KeyInt64 == `%s` but need `%s`", r, value)
			}
		case "KEY_UINT":
			if r := toStr(d.KeyUint); r != value {
				t.Errorf("KeyUint == `%s` but need `%s`", r, value)
			}
		case "KEY_UINT8":
			if r := toStr(d.KeyUint8); r != value {
				t.Errorf("KeyUint8 == `%s` but need `%s`", r, value)
			}
		case "KEY_UINT16":
			if r := toStr(d.KeyUint16); r != value {
				t.Errorf("KeyUint16 == `%s` but need `%s`", r, value)
			}
		case "KEY_UINT32":
			if r := toStr(d.KeyUint32); r != value {
				t.Errorf("KeyUint32 == `%s` but need `%s`", r, value)
			}
		case "KEY_UINT64":
			if r := toStr(d.KeyUint64); r != value {
				t.Errorf("KeyUint64 == `%s` but need `%s`", r, value)
			}
		case "KEY_FLOAT32":
			if r := toStr(d.KeyFloat32); r != value {
				t.Errorf("KeyFloat32 == `%s` but need `%s`", r, value)
			}
		case "KEY_FLOAT64":
			if r := toStr(d.KeyFloat64); r != value {
				t.Errorf("KeyFloat64 == `%s` but need `%s`", r, value)
			}
		case "KEY_STRING":
			if r := toStr(d.KeyString); r != value {
				t.Errorf("KeyString == `%s` but need `%s`", r, value)
			}
		case "KEY_BOOL":
			value = "true:true:true:true:false:false:false:false"
			if r := toStr(d.KeyBool); r != value {
				t.Errorf("KeyBoll == `%s` but need `%s`", r, value)
			}
		}
	}
}

// TestUnmarshalENVSliceIncorrect tests unmarshalENV function
// for slice type with correct values.
func TestUnmarshalENVSliceIncorrect(t *testing.T) {
	var tests = map[string]string{
		"KEY_INT":   "-30:-20:-10:A:10:20:30",
		"KEY_INT8":  "-30:-20:-10:A:10:20:30",
		"KEY_INT16": "-30:-20:-10:A:10:20:30",
		"KEY_INT32": "-30:-20:-10:A:10:20:30",
		"KEY_INT64": "-30:-20:-10:A:10:20:30",

		"KEY_UINT":   "0:10:-20:30",
		"KEY_UINT8":  "0:10:-20:30",
		"KEY_UINT16": "0:10:-20:30",
		"KEY_UINT32": "0:10:-20:30",
		"KEY_UINT64": "0:10:-20:30",

		"KEY_FLOAT32": "-3.1:-1.27:A:1.27:3.3",
		"KEY_FLOAT64": "-3.1:-1.27:A:1.27:3.3",
	}

	// Testing.
	for key, value := range tests {
		var d = &SliceTestType{}

		Clear()
		Set(key, value)

		err := unmarshalENV(d)
		if err == nil {
			t.Error("must be error")
		}
	}
}
