package env

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

type Address struct {
	Country string `env:"COUNTRY"`
	Town    string `env:"TOWN"`
}

// User internal structure for StructTestType type.
type User struct {
	Name    string  `env:"NAME"`
	Email   string  `env:"EMAIL"`
	Address Address `env:"ADDRESS"`
}

// StructTestType structure for testing struct values.
type StructTestType struct {
	User     User    `env:"USER"`
	HomePage url.URL `env:"HOME_PAGE"`
}

// URLTestType structure for testing url.URL values.
type URLTestType struct {
	KeyURLPlain      url.URL     `env:"KEY_URL_PLAIN"`
	KeyURLPoint      *url.URL    `env:"KEY_URL_POINT"`
	KeyURLPlainSlice []url.URL   `env:"KEY_URL_PLAIN_SLICE,!"`
	KeyURLPointSlice []*url.URL  `env:"KEY_URL_POINT_SLICE,!"`
	KeyURLPlainArray [2]url.URL  `env:"KEY_URL_PLAIN_ARRAY,!"`
	KeyURLPointArray [2]*url.URL `env:"KEY_URL_POINT_ARRAY,!"`
}

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

// PlainTestType simple complex type with array (not slice).
type PlainArrayTestType struct {
	Host         string    `env:"HOST"`
	Port         int       `env:"PORT"`
	AllowedHosts [3]string `env:"ALLOWED_HOSTS,:"`
}

// Extended simple complex type that implements
// Marshaller and Unmarshaller interfaces.
type ExtendedTestType struct {
	Host         string   `env:"HOST"`
	Port         int      `env:"PORT"`
	AllowedHosts []string `env:"ALLOWED_HOSTS,:"`
}

// MarshalENV the custom method for marshalling.
func (c *ExtendedTestType) MarshalENV() ([]string, error) {
	// Test data set manually.
	Set("HOST", "192.168.0.1")
	Set("PORT", "80")
	Set("ALLOWED_HOSTS", "192.168.0.1")
	return []string{
		"HOST=192.168.0.1",
		"PORT=80",
		"ALLOWED_HOSTS=192.168.0.1",
	}, nil
}

// UnmarshalENV the custom method for unmarshalling.
func (c *ExtendedTestType) UnmarshalENV() error {
	// Test data set manually.
	c.Host = "192.168.0.1"
	c.Port = 80
	c.AllowedHosts = []string{"192.168.0.1"}
	return nil
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

			err := unmarshalENV(d, "")
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
		var d = &BoolTestType{}

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

		err := unmarshalENV(d, "")
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

		err := unmarshalENV(d, "")
		if err != nil {
			t.Error(err)
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

		err := unmarshalENV(d, "")
		if err == nil {
			t.Error("must be error")
		}
	}
}

// TestMarshalENVNotStruct tests marshalENV function for not struct values.
func TestMarshalNotStruct(t *testing.T) {
	var scope string
	_, err := marshalENV(scope, "")
	if err == nil {
		t.Error("exception expected for an object other than structure")
	}
}

// TestMarshalENVPointerNil tests marshalENV function
// for uninitialized pointer.
func TestMarshalENVPointerNil(t *testing.T) {
	var scope *PlainTestType
	_, err := marshalENV(scope, "")
	if err == nil {
		t.Error("exception expected for an uninitialized object")
	}
}

// TestMarshalENVObj tests marshalENV function with struct value.
func TestMarshalENVObj(t *testing.T) {
	var scope = PlainTestType{
		"localhost",
		8080,
		[]string{"localhost", "127.0.0.1"},
	}

	Clear()
	_, err := marshalENV(scope, "")
	if err != nil {
		t.Error(err)
	}

	// Test marshalling.
	if v := Get("HOST"); v != "localhost" {
		t.Errorf("Incorrect value set for HOST: %s", v)
	}

	if v := Get("PORT"); v != "8080" {
		t.Errorf("Incorrect value set for PORT: %s", v)
	}

	if v := Get("ALLOWED_HOSTS"); v != "localhost:127.0.0.1" {
		t.Errorf("Incorrect value set for PORT: %s", v)
	}
}

// TestMarshalENVPointer tests marshalENV function
// with pointer of the struct.
func TestMarshalENVPointer(t *testing.T) {
	var scope = &PlainTestType{
		"localhost",
		8080,
		[]string{"localhost", "127.0.0.1"},
	}

	Clear()
	_, err := marshalENV(scope, "")
	if err != nil {
		t.Error(err)
	}

	// Test marshalling.
	if v := Get("HOST"); v != "localhost" {
		t.Errorf("Incorrect value set for HOST: %s", v)
	}

	if v := Get("PORT"); v != "8080" {
		t.Errorf("Incorrect value set for PORT: %s", v)
	}

	if v := Get("ALLOWED_HOSTS"); v != "localhost:127.0.0.1" {
		t.Errorf("Incorrect value set for PORT: %s", v)
	}
}

// TestMarshalENVObjCustom tests marshalENV function for object
// with custom MarshalENV method.
func TestMarshalENVObjCustom(t *testing.T) {
	var scope = ExtendedTestType{
		"localhost",                        // default: 192.168.0.1
		8080,                               // default: 80
		[]string{"localhost", "127.0.0.1"}, // default: 192.168.0.1
	}

	Clear()
	_, err := marshalENV(scope, "")
	if err != nil {
		t.Error(err)
	}

	// Test marshalling.
	if v := Get("HOST"); v != "192.168.0.1" {
		t.Errorf("Incorrect value set for HOST: %s", v)
	}

	if v := Get("PORT"); v != "80" {
		t.Errorf("Incorrect value set for PORT: %s", v)
	}

	if v := Get("ALLOWED_HOSTS"); v != "192.168.0.1" {
		t.Errorf("Incorrect value set for PORT: %s", v)
	}
}

// TestMarshalENVPointerCustom tests marshalENV function for pointer
// with custom MarshalENV method.
func TestMarshalENVPointerCustom(t *testing.T) {
	var scope = &ExtendedTestType{
		"localhost",                        // default: 192.168.0.1
		8080,                               // default: 80
		[]string{"localhost", "127.0.0.1"}, // default: 192.168.0.1
	}

	Clear()
	_, err := marshalENV(scope, "")
	if err != nil {
		t.Error(err)
	}

	// Test marshalling.
	if v := Get("HOST"); v != "192.168.0.1" {
		t.Errorf("Incorrect value set for HOST: %s", v)
	}

	if v := Get("PORT"); v != "80" {
		t.Errorf("Incorrect value set for PORT: %s", v)
	}

	if v := Get("ALLOWED_HOSTS"); v != "192.168.0.1" {
		t.Errorf("Incorrect value set for PORT: %s", v)
	}
}

// TestUnmarshalENVNotStruct tests unmarshalENV function for not struct values.
func TestUnmarshalNotStruct(t *testing.T) {
	var scope string
	err := unmarshalENV(scope, "")
	if err == nil {
		t.Error("exception expected for an object other than structure")
	}
}

// TestUnmarshalENVNotPointer tests unmarshalENV function
// for not pointer value.
func TestUnmarshalNotPointer(t *testing.T) {
	var scope PlainTestType
	err := unmarshalENV(scope, "")
	if err == nil {
		t.Error("exception expected for not pointer")
	}
}

// TestUnmarshalENVPointerNil tests unmarshalENV function
// for uninitialized pointer.
func TestUnmarshalPointerNil(t *testing.T) {
	var scope *PlainTestType
	err := unmarshalENV(scope, "")
	if err == nil {
		t.Error("exception expected for an uninitialized object")
	}
}

// TestUnmarshalENV tests unmarshalENV function.
func TestUnmarshalENV(t *testing.T) {
	var scope = &PlainTestType{}

	// Set test data.
	Clear()
	Set("HOST", "localhost")
	Set("PORT", "8080")
	Set("ALLOWED_HOSTS", "localhost:127.0.0.1")

	err := unmarshalENV(scope, "")
	if err != nil {
		t.Error(err)
	}

	// Test marshalling.
	if scope.Host != "localhost" {
		t.Errorf("Incorrect value set for HOST: %s", scope.Host)
	}

	if scope.Port != 8080 {
		t.Errorf("Incorrect value set for PORT: %d", scope.Port)
	}

	str := strings.Replace(fmt.Sprint(scope.AllowedHosts), " ", ":", -1)
	if value := strings.Trim(str, "[]:"); value != "localhost:127.0.0.1" {
		t.Errorf("Incorrect value set for ALLOWED_HOSTS: %s", value)
	}
}

// TestUnmarshalENVArray tests unmarshalENV function with array.
func TestUnmarshalENVArray(t *testing.T) {
	var scope = &PlainArrayTestType{}

	// Set test data.
	Clear()
	Set("HOST", "localhost")
	Set("PORT", "8080")
	Set("ALLOWED_HOSTS", "localhost:127.0.0.1")

	err := unmarshalENV(scope, "")
	if err != nil {
		t.Error(err)
	}

	// Test marshalling.
	if scope.Host != "localhost" {
		t.Errorf("Incorrect value set for HOST: %s", scope.Host)
	}

	if scope.Port != 8080 {
		t.Errorf("Incorrect value set for PORT: %d", scope.Port)
	}

	str := strings.Replace(fmt.Sprint(scope.AllowedHosts), " ", ":", -1)
	if value := strings.Trim(str, "[]:"); value != "localhost:127.0.0.1" {
		t.Errorf("Incorrect value set for ALLOWED_HOSTS: %s", value)
	}
}

// TestUnmarshalENVArrayOverflow tests unmarshalENV function
// with array overflow.
func TestUnmarshalENVArrayOverflow(t *testing.T) {
	var scope = &PlainArrayTestType{}

	// Set test data.
	Clear()
	Set("HOST", "localhost")
	Set("PORT", "8080")
	Set("ALLOWED_HOSTS", "localhost:127.0.0.1:0.0.0.0:192.168.0.1") // 4 items

	err := unmarshalENV(scope, "")
	if err == nil {
		t.Error("there must be an array overflow error")
	}
}

// TestUnmarshalENVCustom tests unmarshalENV function
// with custom UnmarshalENV method.
func TestUnmarshalENVCustom(t *testing.T) {
	var scope = &ExtendedTestType{}

	// Set test data.
	Clear()
	Set("HOST", "localhost")                    // default: 192.168.0.1
	Set("PORT", "8080")                         // default: 80
	Set("ALLOWED_HOSTS", "localhost:127.0.0.1") // default: 192.168.0.1

	err := unmarshalENV(scope, "")
	if err != nil {
		t.Error(err)
	}

	// Test marshalling.
	if scope.Host != "192.168.0.1" {
		t.Errorf("Incorrect value set for HOST: %s", scope.Host)
	}

	if scope.Port != 80 {
		t.Errorf("Incorrect value set for PORT: %d", scope.Port)
	}

	str := strings.Replace(fmt.Sprint(scope.AllowedHosts), " ", ":", -1)
	if value := strings.Trim(str, "[]"); value != "192.168.0.1" {
		t.Errorf("Incorrect value set for ALLOWED_HOSTS: %v", value)
	}
}

// TestMarshalURL tests marshaling of the URL.
func TestMarshalURL(t *testing.T) {
	var test string
	var data = URLTestType{
		KeyURLPlain: url.URL{Scheme: "http", Host: "plain.example.com"},
		KeyURLPoint: &url.URL{Scheme: "http", Host: "point.example.com"},
		KeyURLPlainSlice: []url.URL{
			url.URL{Scheme: "http", Host: "a.plain.example.com"},
			url.URL{Scheme: "http", Host: "b.plain.example.com"},
		},
		KeyURLPointSlice: []*url.URL{
			&url.URL{Scheme: "http", Host: "a.point.example.com"},
			&url.URL{Scheme: "http", Host: "b.point.example.com"},
		},
		KeyURLPlainArray: [2]url.URL{
			url.URL{Scheme: "http", Host: "c.plain.example.com"},
			url.URL{Scheme: "http", Host: "d.plain.example.com"},
		},
		KeyURLPointArray: [2]*url.URL{
			&url.URL{Scheme: "http", Host: "c.point.example.com"},
			&url.URL{Scheme: "http", Host: "d.point.example.com"},
		},
	}

	Marshal(data)

	// Tests results.
	if v := Get("KEY_URL_PLAIN"); v != "http://plain.example.com" {
		t.Errorf("Incorrect marshaling plain url.URL: %s", v)
	}

	if v := Get("KEY_URL_POINT"); v != "http://point.example.com" {
		t.Errorf("Incorrect marshaling poin url.URL: %s", v)
	}

	// Plain slice.
	test = "http://a.plain.example.com!http://b.plain.example.com"
	if v := Get("KEY_URL_PLAIN_SLICE"); v != test {
		t.Errorf("Incorrect marshaling poin slice []url.URL: %s", v)
	}

	// Point slice.
	test = "http://a.point.example.com!http://b.point.example.com"
	if v := Get("KEY_URL_POINT_SLICE"); v != test {
		t.Errorf("Incorrect marshaling point slice []*url.URL: %s", v)
	}

	// Plain array.
	test = "http://c.plain.example.com!http://d.plain.example.com"
	if v := Get("KEY_URL_PLAIN_ARRAY"); v != test {
		t.Errorf("Incorrect marshaling plain array []url.URL: %s", v)
	}

	// Point array.
	test = "http://c.point.example.com!http://d.point.example.com"
	if v := Get("KEY_URL_POINT_ARRAY"); v != test {
		t.Errorf("Incorrect marshaling point array []*url.URL: %s", v)
	}
}

// TestUnmarshalURL tests unmarshaling of the URL.
func TestUnmarshalURL(t *testing.T) {
	var (
		slice []string
		str   string

		data = URLTestType{}
	)

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

	Unmarshal(&data)

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

// TestMarshalStruct tests marshaling of the Struct.
func TestMarshalStruct(t *testing.T) {
	var data = StructTestType{
		User: User{
			Name:  "John",
			Email: "john@example.com",
			Address: Address{
				Country: "Ukraine",
				Town:    "Chernihiv",
			},
		},
		HomePage: url.URL{Scheme: "http", Host: "example.com"},
	}

	// Marshaling.
	result, _ := Marshal(data)

	// Tests.
	if v := Get("USER_NAME"); v != "John" {
		t.Errorf("Incorrect marshaling (Name): %s\n%v", v, result)
	}

	if v := Get("USER_EMAIL"); v != "john@example.com" {
		t.Errorf("Incorrect marshaling (Email): %s\n%v", v, result)
	}

	if v := Get("USER_ADDRESS_COUNTRY"); v != "Ukraine" {
		t.Errorf("Incorrect marshaling (Cuontry): %s\n%v", v, result)
	}

	if v := Get("USER_ADDRESS_TOWN"); v != "Chernihiv" {
		t.Errorf("Incorrect marshaling (Town): %s\n%v", v, result)
	}

	if v := Get("HOME_PAGE"); v != "http://example.com" {
		t.Errorf("Incorrect marshaling url.URL (HomePage):%s", v)
	}
}

// TestUnmarshalStruct tests unmarshaling of the Struct.
func TestUnmarshalStruct(t *testing.T) {
	var data = StructTestType{User: User{}}

	Set("USER_NAME", "John")
	Set("USER_EMAIL", "john@example.com")
	Set("USER_ADDRESS_COUNTRY", "Ukraine")
	Set("USER_ADDRESS_TOWN", "Chernihiv")
	Set("HOME_PAGE", "http://example.com")

	// Unmarshaling.
	err := unmarshalENV(&data, "")
	if err != nil {
		t.Error("Incorrect ummarshaling")
	}

	// Tests.
	if data.User.Address.Country != "Ukraine" ||
		data.User.Address.Town != "Chernihiv" {
		t.Errorf("Incorrect ummarshaling User.Address: %v", data.User.Address)
	}

	if data.User.Name != "John" || data.User.Email != "john@example.com" {
		t.Errorf("Incorrect ummarshaling User: %v", data.User)
	}

	if data.HomePage.String() != "http://example.com" {
		t.Errorf("Incorrect ummarshaling url.URL: %v", data.HomePage)
	}
}
