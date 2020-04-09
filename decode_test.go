package env

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

// The str convert sequence in string.
func str(seq interface{}) string {
	return strings.Trim(strings.Replace(fmt.Sprint(seq), " ", ":", -1), "[]")
}

// The dataUnmarshalENV structure with custom UnmarshalENV method.
type dataUnmarshalENV struct {
	Host         string   `env:"HOST"`
	Port         int      `env:"PORT"`
	AllowedHosts []string `env:"ALLOWED_HOSTS,:"`
}

// UnmarshalENV the custom method for unmarshalling.
func (c *dataUnmarshalENV) UnmarshalENV() error {
	c.Host = "192.168.0.1"
	c.Port = 80
	c.AllowedHosts = []string{"192.168.0.1"}
	return nil
}

// TestUnmarshalENVNotPointer tests unmarshalENV for the correct handling
// of an exception for a non-pointer value.
func TestUnmarshalENVNotPointer(t *testing.T) {
	type data struct{}
	if err := unmarshalENV(data{}, ""); err == nil {
		t.Error("An exception must be thrown on the value of the non-pointer.")
	}
}

// TestUnmarshalENVNotInitialized tests unmarshalENV for the correct handling
// of an exception for a not initialized value.
func TestUnmarshalENVNotInitialized(t *testing.T) {
	type Empty struct{}
	var e *Empty
	if err := unmarshalENV(e, ""); err == nil {
		t.Error("An exception must be thrown on the not initialized value.")
	}
}

// TestUnmarshalENVNotStruct tests unmarshalENV for the correct handling
// of an exception for a value that isn't struct.
func TestUnmarshalENVNotStruct(t *testing.T) {
	if err := unmarshalENV(new(int), ""); err == nil {
		t.Error("An exception must be thrown on the value that isn't struct.")
	}
}

// TestUnmarshalENVNumber tests unmarshalENV for Int, Uint and Float types.
func TestUnmarshalENVNumber(t *testing.T) {
	type Number struct {
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

	var (
		max   = "922337203685477580777555333"
		tests = map[string][]string{
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
	)

	// Correct value.
	for i := 0; i < 3; i++ {
		var err error
		for key, data := range tests {
			var d = &Number{}

			Clear()
			err = Set(key, data[i])
			if err != nil {
				t.Error(err)
			}

			err = unmarshalENV(d, "")
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
			}
		}
	}
}

// TestUnmarshalENVBoll tests unmarshalENV function for bool type.
func TestUnmarshalENVBool(t *testing.T) {
	type Boolean struct {
		KeyBool bool `env:"KEY_BOOL"`
	}

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
		var d = &Boolean{}

		Clear()
		Set("KEY_BOOL", value)

		err := unmarshalENV(d, "")
		if err != nil {
			t.Error(err)
		}

		if d.KeyBool != test {
			t.Errorf("KeyBool == %t but need %t", d.KeyBool, test)
		}
	}

	// Incorrect value.
	for _, value := range []string{"string", "0.d", "true/false"} {
		var d = &Boolean{}

		Clear()
		Set("KEY_BOOL", value)

		err := unmarshalENV(d, "")
		if err == nil {
			t.Error("didn't handle the error")
		}
	}
}

// TestUnmarshalENVString tests unmarshalENV function for string type.
func TestUnmarshalENVString(t *testing.T) {
	type String struct {
		KeyString string `env:"KEY_STRING"`
	}
	var tests = []interface{}{
		8080,
		"Hello World",
		"true",
		true,
		3.14,
	}

	// Test correct values.
	for _, test := range tests {
		var d = &String{}
		var s = fmt.Sprintf("%v", test)

		Clear()
		Set("KEY_STRING", s)

		err := unmarshalENV(d, "")
		if err != nil {
			t.Error(err)
		}

		if d.KeyString != s {
			t.Errorf("KeyString == `%s` but need `%s`", d.KeyString, s)
		}
	}
}

// TestUnmarshalENVSlice tests unmarshalENV function for slice.
func TestUnmarshalENVSlice(t *testing.T) {
	type Slice struct {
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

	var (
		corretc = map[string]string{
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
		incorrect = map[string]string{
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
	)

	// Testing correct values.
	for key, value := range corretc {
		var d = &Slice{}

		Clear()
		Set(key, value)

		err := unmarshalENV(d, "")
		if err != nil {
			t.Error(err)
		}

		switch key {
		case "KEY_INT":
			if r := str(d.KeyInt); r != value {
				t.Errorf("KeyInt == `%s` but need `%s`", r, value)
			}
		case "KEY_INT8":
			if r := str(d.KeyInt8); r != value {
				t.Errorf("KeyInt8 == `%s` but need `%s`", r, value)
			}
		case "KEY_INT16":
			if r := str(d.KeyInt16); r != value {
				t.Errorf("KeyInt16 == `%s` but need `%s`", r, value)
			}
		case "KEY_INT32":
			if r := str(d.KeyInt32); r != value {
				t.Errorf("KeyInt32 == `%s` but need `%s`", r, value)
			}
		case "KEY_INT64":
			if r := str(d.KeyInt64); r != value {
				t.Errorf("KeyInt64 == `%s` but need `%s`", r, value)
			}
		case "KEY_UINT":
			if r := str(d.KeyUint); r != value {
				t.Errorf("KeyUint == `%s` but need `%s`", r, value)
			}
		case "KEY_UINT8":
			if r := str(d.KeyUint8); r != value {
				t.Errorf("KeyUint8 == `%s` but need `%s`", r, value)
			}
		case "KEY_UINT16":
			if r := str(d.KeyUint16); r != value {
				t.Errorf("KeyUint16 == `%s` but need `%s`", r, value)
			}
		case "KEY_UINT32":
			if r := str(d.KeyUint32); r != value {
				t.Errorf("KeyUint32 == `%s` but need `%s`", r, value)
			}
		case "KEY_UINT64":
			if r := str(d.KeyUint64); r != value {
				t.Errorf("KeyUint64 == `%s` but need `%s`", r, value)
			}
		case "KEY_FLOAT32":
			if r := str(d.KeyFloat32); r != value {
				t.Errorf("KeyFloat32 == `%s` but need `%s`", r, value)
			}
		case "KEY_FLOAT64":
			if r := str(d.KeyFloat64); r != value {
				t.Errorf("KeyFloat64 == `%s` but need `%s`", r, value)
			}
		case "KEY_STRING":
			if r := str(d.KeyString); r != value {
				t.Errorf("KeyString == `%s` but need `%s`", r, value)
			}
		case "KEY_BOOL":
			value = "true:true:true:true:false:false:false:false"
			if r := str(d.KeyBool); r != value {
				t.Errorf("KeyBoll == `%s` but need `%s`", r, value)
			}
		}
	}

	// Testing incorrect values.
	for key, value := range incorrect {
		var d = &Slice{}

		Clear()
		Set(key, value)

		err := unmarshalENV(d, "")
		if err == nil {
			t.Error("must be error")
		}
	}
}

// TestUnmarshalENVArray tests unmarshalENV function with array.
func TestUnmarshalENVArray(t *testing.T) {
	type Array struct {
		KeyInt   [7]int   `env:"KEY_INT,:"`
		KeyInt8  [7]int8  `env:"KEY_INT8,:"`
		KeyInt16 [7]int16 `env:"KEY_INT16,:"`
		KeyInt32 [7]int32 `env:"KEY_INT32,:"`
		KeyInt64 [7]int64 `env:"KEY_INT64,:"`

		KeyUint   [4]uint   `env:"KEY_UINT,:"`
		KeyUint8  [4]uint8  `env:"KEY_UINT8,:"`
		KeyUint16 [4]uint16 `env:"KEY_UINT16,:"`
		KeyUint32 [4]uint32 `env:"KEY_UINT32,:"`
		KeyUint64 [4]uint64 `env:"KEY_UINT64,:"`

		KeyFloat32 [5]float32 `env:"KEY_FLOAT32,:"`
		KeyFloat64 [5]float64 `env:"KEY_FLOAT64,:"`

		KeyString [5]string `env:"KEY_STRING,:"`
		KeyBool   [8]bool   `env:"KEY_BOOL,:"`
	}

	var (
		corretc = map[string]string{
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
		incorrect = map[string]string{
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
		overflow = map[string]string{
			"KEY_INT":   "-30:-20:-10:0:10:20:30:100",
			"KEY_INT8":  "-30:-20:-10:0:10:20:30:100",
			"KEY_INT16": "-30:-20:-10:0:10:20:30:100",
			"KEY_INT32": "-30:-20:-10:0:10:20:30:100",
			"KEY_INT64": "-30:-20:-10:0:10:20:30:100",

			"KEY_UINT":   "0:10:20:30:100",
			"KEY_UINT8":  "0:10:20:30:100",
			"KEY_UINT16": "0:10:20:30:100",
			"KEY_UINT32": "0:10:20:30:100",
			"KEY_UINT64": "0:10:20:30:100",

			"KEY_FLOAT32": "-3.1:-1.27:0:1.27:3.3:100.0",
			"KEY_FLOAT64": "-3.1:-1.27:0:1.27:3.3:100.0",

			"KEY_STRING": "one:two:three:four:five:one hundred",
			"KEY_BOOL":   "1:true:True:TRUE:0:false:False:False:true",
		}
	)

	// Convert slice into string.

	// Test correct values.
	for key, value := range corretc {
		var d = &Array{}

		Clear()
		Set(key, value)

		err := unmarshalENV(d, "")
		if err != nil {
			t.Error(err)
		}

		switch key {
		case "KEY_INT":
			if r := str(d.KeyInt); r != value {
				t.Errorf("KeyInt == `%s` but need `%s`", r, value)
			}
		case "KEY_INT8":
			if r := str(d.KeyInt8); r != value {
				t.Errorf("KeyInt8 == `%s` but need `%s`", r, value)
			}
		case "KEY_INT16":
			if r := str(d.KeyInt16); r != value {
				t.Errorf("KeyInt16 == `%s` but need `%s`", r, value)
			}
		case "KEY_INT32":
			if r := str(d.KeyInt32); r != value {
				t.Errorf("KeyInt32 == `%s` but need `%s`", r, value)
			}
		case "KEY_INT64":
			if r := str(d.KeyInt64); r != value {
				t.Errorf("KeyInt64 == `%s` but need `%s`", r, value)
			}
		case "KEY_UINT":
			if r := str(d.KeyUint); r != value {
				t.Errorf("KeyUint == `%s` but need `%s`", r, value)
			}
		case "KEY_UINT8":
			if r := str(d.KeyUint8); r != value {
				t.Errorf("KeyUint8 == `%s` but need `%s`", r, value)
			}
		case "KEY_UINT16":
			if r := str(d.KeyUint16); r != value {
				t.Errorf("KeyUint16 == `%s` but need `%s`", r, value)
			}
		case "KEY_UINT32":
			if r := str(d.KeyUint32); r != value {
				t.Errorf("KeyUint32 == `%s` but need `%s`", r, value)
			}
		case "KEY_UINT64":
			if r := str(d.KeyUint64); r != value {
				t.Errorf("KeyUint64 == `%s` but need `%s`", r, value)
			}
		case "KEY_FLOAT32":
			if r := str(d.KeyFloat32); r != value {
				t.Errorf("KeyFloat32 == `%s` but need `%s`", r, value)
			}
		case "KEY_FLOAT64":
			if r := str(d.KeyFloat64); r != value {
				t.Errorf("KeyFloat64 == `%s` but need `%s`", r, value)
			}
		case "KEY_STRING":
			if r := str(d.KeyString); r != value {
				t.Errorf("KeyString == `%s` but need `%s`", r, value)
			}
		case "KEY_BOOL":
			value = "true:true:true:true:false:false:false:false"
			if r := str(d.KeyBool); r != value {
				t.Errorf("KeyBoll == `%s` but need `%s`", r, value)
			}
		}
	}

	// Test incorrect values.
	for key, value := range incorrect {
		var d = &Array{}

		Clear()
		Set(key, value)

		err := unmarshalENV(d, "")
		if err == nil {
			t.Error("There should be an exception due to an invalid value.")
		}
	}

	// Test array overflow.
	for key, value := range overflow {
		var d = &Array{}

		Clear()
		Set(key, value)

		err := unmarshalENV(d, "")
		if err == nil {
			t.Error("There should be an exception due to array overflow.")
		}
	}
}

// TestUnmarshalURL tests unmarshalENV for url.URL type.
func TestUnmarshalURL(t *testing.T) {
	type URL struct {
		KeyURLPlain      url.URL     `env:"KEY_URL_PLAIN"`
		KeyURLPoint      *url.URL    `env:"KEY_URL_POINT"`
		KeyURLPlainSlice []url.URL   `env:"KEY_URL_PLAIN_SLICE,,!"`
		KeyURLPointSlice []*url.URL  `env:"KEY_URL_POINT_SLICE,,!"`
		KeyURLPlainArray [2]url.URL  `env:"KEY_URL_PLAIN_ARRAY,,!"`
		KeyURLPointArray [2]*url.URL `env:"KEY_URL_POINT_ARRAY,,!"`
	}
	var (
		slice []string
		str   string

		data = URL{}
	)

	// Set tests data.
	Set("KEY_URL_PLAIN", "http://plain.example.com")
	Set("KEY_URL_POINT", "http://point.example.com")
	Set("KEY_URL_PLAIN_SLICE",
		"http://a.plain.example.com!http://b.plain.example.com")
	Set("KEY_URL_POINT_SLICE",
		"http://a.point.example.com!http://b.point.example.com")
	Set("KEY_URL_PLAIN_ARRAY",
		"http://c.plain.example.com!http://d.plain.example.com")
	Set("KEY_URL_POINT_ARRAY",
		"http://c.point.example.com!http://d.point.example.com")

	// Unmarshaling.
	unmarshalENV(&data, "")

	// Tests results.
	if v := data.KeyURLPlain.String(); v != "http://plain.example.com" {
		t.Errorf("Incorrect unmarshaling plain url.URL: %s", v)
	}

	if v := data.KeyURLPoint.String(); v != "http://point.example.com" {
		t.Errorf("Incorrect unmarshaling point url.URL: %s", v)
	}

	// Plain slice.
	slice = []string{}
	for _, v := range data.KeyURLPlainSlice {
		slice = append(slice, v.String())
	}
	str = strings.Trim(strings.Replace(fmt.Sprint(slice), " ", "!", -1), "[]")
	if str != "http://a.plain.example.com!http://b.plain.example.com" {
		t.Errorf("Incorrect unmarshaling plain slice []url.URL: %s", str)
	}

	// Point slice.
	slice = []string{}
	for _, v := range data.KeyURLPointSlice {
		slice = append(slice, v.String())
	}
	str = strings.Trim(strings.Replace(fmt.Sprint(slice), " ", "!", -1), "[]")
	if str != "http://a.point.example.com!http://b.point.example.com" {
		t.Errorf("Incorrect unmarshaling point alice []*url.URL: %s", str)
	}

	// Plain array.
	slice = []string{}
	for _, v := range data.KeyURLPlainArray {
		slice = append(slice, v.String())
	}
	str = strings.Trim(strings.Replace(fmt.Sprint(slice), " ", "!", -1), "[]")
	if str != "http://c.plain.example.com!http://d.plain.example.com" {
		t.Errorf("Incorrect unmarshaling plain array [2]url.URL: %s", str)
	}

	// Point array.
	slice = []string{}
	for _, v := range data.KeyURLPointArray {
		slice = append(slice, v.String())
	}
	str = strings.Trim(strings.Replace(fmt.Sprint(slice), " ", "!", -1), "[]")
	if str != "http://c.point.example.com!http://d.point.example.com" {
		t.Errorf("Incorrect unmarshaling point array [2]*url.URL: %s", str)
	}
}

// TestUnmarshalStruct tests unmarshalENV for the struct.
func TestUnmarshalStruct(t *testing.T) {
	type Address struct {
		Country string `env:"COUNTRY"`
	}

	type User struct {
		Name    string  `env:"NAME"`
		Address Address `env:"ADDRESS"`
	}

	type Client struct {
		User     User    `env:"USER"`
		HomePage url.URL `env:"HOME_PAGE"`
	}

	var c = Client{}

	Set("USER_NAME", "John")
	Set("USER_ADDRESS_COUNTRY", "USA")
	Set("HOME_PAGE", "http://example.com")

	// Unmarshaling.
	err := unmarshalENV(&c, "")
	if err != nil {
		t.Error("Incorrect ummarshaling.")
	}

	// Tests.
	if c.User.Address.Country != "USA" {
		t.Errorf("Incorrect ummarshaling User.Address: %v", c.User.Address)
	}

	if c.User.Name != "John" {
		t.Errorf("Incorrect ummarshaling User: %v", c.User)
	}

	if c.HomePage.String() != "http://example.com" {
		t.Errorf("Incorrect ummarshaling url.URL: %v", c.HomePage)
	}
}

// TestUnmarshalStructPtr tests unmarshalENV for the pointerf of the struct.
func TestUnmarshalStructPtr(t *testing.T) {
	type Address struct {
		Country string `env:"COUNTRY"`
	}

	type User struct {
		Name    string   `env:"NAME"`
		Address *Address `env:"ADDRESS"`
	}

	type Client struct {
		User     *User    `env:"USER"`
		HomePage *url.URL `env:"HOME_PAGE"`
	}

	var c = Client{}

	Set("USER_NAME", "John")
	Set("USER_ADDRESS_COUNTRY", "USA")
	Set("HOME_PAGE", "http://example.com")

	// Unmarshaling.
	err := unmarshalENV(&c, "")
	if err != nil {
		t.Error("Incorrect ummarshaling.")
	}

	// Tests.
	if c.User.Address.Country != "USA" {
		t.Errorf("Incorrect ummarshaling User.Address: %v", c.User.Address)
	}

	if c.User.Name != "John" {
		t.Errorf("Incorrect ummarshaling User: %v", c.User)
	}

	if c.HomePage.String() != "http://example.com" {
		t.Errorf("Incorrect ummarshaling url.URL: %v", c.HomePage)
	}
}

// TestUnmarshalENVCustom tests unmarshalENV function
// with custom UnmarshalENV method.
func TestUnmarshalENVCustom(t *testing.T) {
	var c = &dataUnmarshalENV{}

	// Set test data.
	Clear()
	Set("HOST", "localhost")                    // default: 192.168.0.1
	Set("PORT", "8080")                         // default: 80
	Set("ALLOWED_HOSTS", "localhost:127.0.0.1") // default: 192.168.0.1

	err := unmarshalENV(c, "")
	if err != nil {
		t.Error(err)
	}

	// Test marshalling.
	if c.Host != "192.168.0.1" {
		t.Errorf("Incorrect value set for HOST: %s", c.Host)
	}

	if c.Port != 80 {
		t.Errorf("Incorrect value set for PORT: %d", c.Port)
	}

	str := strings.Replace(fmt.Sprint(c.AllowedHosts), " ", ":", -1)
	if value := strings.Trim(str, "[]"); value != "192.168.0.1" {
		t.Errorf("Incorrect value set for ALLOWED_HOSTS: %v", value)
	}
}

// TestUnmarshalENVStringPtr tests unmarshalENV function
// for pointer on the string type.
func TestUnmarshalENVStringPtr(t *testing.T) {
	type String struct {
		KeyString *string `env:"KEY_STRING"`
	}
	var (
		keyString string

		d = String{KeyString: &keyString}
	)

	Set("KEY_STRING", "Hello World")
	err := unmarshalENV(&d, "")
	if err != nil {
		t.Error(err)
	}

	if *d.KeyString != "Hello World" {
		t.Errorf("Incorrect value set for KEY_STRING: %v", *d.KeyString)
	}

}

// TestUnmarshalDefaultValue tests unmarshalENV for default value.
func TestUnmarshalDefaultValue(t *testing.T) {
	type data struct {
		Host         string    `env:"HOST,0.0.0.0"`
		AllowedHosts []string  `env:"ALLOWED_HOSTS,{localhost:0.0.0.0},:"`
		Names        [3]string `env:"NAME_LIST,'John,Bob,Smit',,"` // sep `,`
	}

	var (
		d   data
		err error
	)

	Clear() // make empty environment

	// Unmarshaling wit default values.
	d = data{}
	err = unmarshalENV(&d, "")
	if err != nil {
		t.Error("Incorrect ummarshaling.")
	}

	if d.Host != "0.0.0.0" {
		t.Errorf("incorrect Host %s", d.Host)
	}

	if str(d.AllowedHosts) != "localhost:0.0.0.0" {
		t.Errorf("incorrect AllowedHosts %s", d.AllowedHosts)
	}

	if str(d.Names) != "John:Bob:Smit" {
		t.Errorf("incorrect AllowedHosts %s", d.AllowedHosts)
	}

	// Set any values.
	Set("HOST", "localhost")
	Set("ALLOWED_HOSTS", "127.0.0.1:localhost")
	Set("NAME_LIST", "John")

	// Unmarshaling wit environment values.
	d = data{}
	err = unmarshalENV(&d, "")
	if err != nil {
		t.Error("Incorrect ummarshaling.")
	}

	if d.Host == "0.0.0.0" {
		t.Errorf("Host sets as default %s", d.Host)
	}

	if str(d.AllowedHosts) == "localhost:0.0.0.0" {
		t.Errorf("AllowedHosts sets as default %s", d.AllowedHosts)
	}

	if str(d.Names) == "John:Bob:Smit" {
		t.Errorf("Names setas as default %s", d.Names)
	}
}
