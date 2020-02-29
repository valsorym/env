package env

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"testing"
)

// ...
type TestIntUintFloat struct {
	Value   string
	Control string
	Correct bool
	Kind    reflect.Kind
}

// ...
type TestBool struct {
	Value   string
	Control bool
	Correct bool
}

// TestParseTag tests parseTag function.
func TestParseTag(t *testing.T) {
	var tests = [][]string{
		//       tagValue, defaultName, defaultSep
		[]string{"", "HOST", " ", "HOST", " "},
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

// TestIsEmpty tests isEmpty function.
// Function returns true for empty or comment strings.
func TestIsEmpty(t *testing.T) {
	var tests = []string{
		"",             // just empty string
		"  ",           // spaces only
		"\t \n",        // separators only
		"# Comment.",   // comment line
		"\t #Comment.", // mix: separators and comment
	}

	for _, test := range tests {
		if !isEmpty(test) {
			t.Errorf("The `%s` value marked as a non-empty.", test)
		}
	}
}

// TestRemoveInlineComment tests removeCommentInline function.
func TestRemoveInlineComment(t *testing.T) {
	var items = [][]string{
		[]string{`="value" # comment`, `="value"`},
		[]string{`="value's" # comment`, `="value's"`},
		[]string{`="value # here" # comment`, `="value # here"`},
	}

	for _, item := range items {
		test, result := item[0], item[1]
		if v := removeInlineComment(test, "\""); v != result {
			t.Errorf("The `%s` doesn't match  `%s`.", v, result)
		}
	}
}

// TestParseExpressionIncorrectKey tests parseExpression function.
// The function returns an error on the wrong key.
func TestParseExpressionIncorrectKey(t *testing.T) {
	var tests []string = []string{
		`2KEY="value"`, // incorrect key
		`K EY="value"`, // broken key
		`KEY ="value"`, // space before equal sign
		`="value"`,     // without key
		`# Comment.`,   // comment
		``,             // empty string
	}

	for _, test := range tests {
		if _, _, err := parseExpression(test); err != incorrectKeyError {
			t.Errorf("For `%s` value must be incorrectKeyError.", test)
		}
	}
}

// TestParseExpressionIncorrectValue tests parseExpression function.
// The function returns an error on the wrong value.
func TestParseExpressionIncorrectValue(t *testing.T) {
	var tests []string = []string{
		`export KEY='value`, // not end-quote
		`KEY="value`,        // not end-quote
		`KEY='value"`,       // end-quote does not match
		`KEY="value\"`,      // end-quote part of the string \"
		`KEY='value\'`,      // end-quote part of the string \'
		`KEY= "value"`,      // space after equal sign
	}

	for _, test := range tests {
		if _, _, err := parseExpression(test); err != incorrectValueError {
			t.Errorf("For `%s` value must be incorrectValueError.", test)
		}
	}
}

// TestParseExpression tests parseExpression function.
func TestParseExpression(t *testing.T) {
	var tests = []string{
		`export KEY="value"`,
		`KEY="value"`,
		`KEY="value" # comment`,
	}

	for _, test := range tests {
		if k, v, _ := parseExpression(test); k != "KEY" || v != "value" {
			t.Errorf("Incorrect parsing for `%s` value, "+
				"whre KEY=`%s` and VALUE=`%s`", test, k, v)
		}
	}
}

// TestStrToBool tests strToBool function.
func TestStrToBool(t *testing.T) {
	var tests = []TestBool{
		TestBool{"", false, true},
		TestBool{"0", false, true},
		TestBool{"1", true, true},
		TestBool{"1.1", true, true},
		TestBool{"-1.1", true, true},
		TestBool{"0.0", false, true},
		TestBool{"true", true, true},
		TestBool{"True", true, true},
		TestBool{"TRUE", true, true},
		TestBool{"false", false, true},
		TestBool{"False", false, true},
		TestBool{"FALSE", false, true},
		TestBool{"string", false, false},
		TestBool{"a:b:c", false, false},
	}

	// Test correct values.
	for _, test := range tests {
		r, err := strToBool(test.Value)
		if test.Correct && err != nil {
			t.Error(err)
		} else if !test.Correct && err == nil {
			t.Errorf("value %s should throw an exception", test.Value)
		}

		if r != test.Control {
			t.Errorf("expected %s but the result %t", test.Value, r)
		}
	}
}

// TestStrToIntKind tests strToIntKind function.
func TestStrToIntKind(t *testing.T) {
	var (
		tests    []TestIntUintFloat
		maxInt   string = fmt.Sprintf("%d", math.MaxInt64-1)
		maxInt8  string = fmt.Sprintf("%d", math.MaxInt8-1)
		maxInt16 string = fmt.Sprintf("%d", math.MaxInt16-1)
		maxInt32 string = fmt.Sprintf("%d", math.MaxInt32-1)
		maxInt64 string = fmt.Sprintf("%d", math.MaxInt64-1)
	)

	// For 32-bit platform.
	if strconv.IntSize == 32 {
		maxInt = maxInt32
	}

	// Test data.
	tests = []TestIntUintFloat{
		TestIntUintFloat{"", "0", true, reflect.Int},
		TestIntUintFloat{"0", "0", true, reflect.Int},
		TestIntUintFloat{"-3", "-3", true, reflect.Int},
		TestIntUintFloat{"3", "3", true, reflect.Int},

		TestIntUintFloat{"-128", "-128", true, reflect.Int8},
		TestIntUintFloat{"127", "127", true, reflect.Int8},

		TestIntUintFloat{maxInt, maxInt, true, reflect.Int},
		TestIntUintFloat{maxInt8, maxInt8, true, reflect.Int8},
		TestIntUintFloat{maxInt16, maxInt16, true, reflect.Int16},
		TestIntUintFloat{maxInt32, maxInt32, true, reflect.Int32},
		TestIntUintFloat{maxInt64, maxInt64, true, reflect.Int64},

		TestIntUintFloat{"string", "0", false, reflect.Int},
		TestIntUintFloat{"3" + maxInt, "0", false, reflect.Int},
		TestIntUintFloat{"3" + maxInt8, "0", false, reflect.Int8},
		TestIntUintFloat{"-129", "0", false, reflect.Int8},
		TestIntUintFloat{"128", "0", false, reflect.Int8},
		TestIntUintFloat{"3" + maxInt16, "0", false, reflect.Int16},
		TestIntUintFloat{"3" + maxInt32, "0", false, reflect.Int32},
		TestIntUintFloat{"3" + maxInt64, "0", false, reflect.Int64},
		TestIntUintFloat{"0", "0", false, reflect.Slice},
	}

	// Test correct values.
	for _, data := range tests {
		r, err := strToIntKind(data.Value, data.Kind)
		if data.Correct && err != nil {
			t.Error(err)
		} else if !data.Correct && err == nil {
			t.Errorf("the value %s should throw an exception", data.Value)
		} else if err != nil && r != 0 {
			t.Errorf("any error should return zero but returns %v", r)
		}

		control := fmt.Sprintf("%d", int64(r))
		if control != data.Control {
			t.Errorf("expected %s but returns %s", data.Control, control)
		}
	}
}

// TestStrToUintKind tests strToUintKind function.
func TestStrToUintKind(t *testing.T) {
	var (
		tests     []TestIntUintFloat
		maxUint   string = "18446744073709551614"
		maxUint8  string = fmt.Sprintf("%d", math.MaxUint8-1)
		maxUint16 string = fmt.Sprintf("%d", math.MaxUint16-1)
		maxUint32 string = fmt.Sprintf("%d", math.MaxUint32-1)
		maxUint64 string = "18446744073709551614"
	)

	// For 32-bit platform.
	if strconv.IntSize == 32 {
		maxUint = maxUint32
	}

	// Test data.
	tests = []TestIntUintFloat{
		TestIntUintFloat{"", "0", true, reflect.Uint},
		TestIntUintFloat{"0", "0", true, reflect.Uint},
		TestIntUintFloat{"3", "3", true, reflect.Uint},
		TestIntUintFloat{maxUint, maxUint, true, reflect.Uint},
		TestIntUintFloat{maxUint8, maxUint8, true, reflect.Uint8},
		TestIntUintFloat{maxUint16, maxUint16, true, reflect.Uint16},
		TestIntUintFloat{maxUint32, maxUint32, true, reflect.Uint32},
		TestIntUintFloat{maxUint64, maxUint64, true, reflect.Uint64},

		TestIntUintFloat{"string", "0", false, reflect.Uint},
		TestIntUintFloat{"-3", "0", false, reflect.Uint},
		TestIntUintFloat{"9" + maxUint, "0", false, reflect.Uint},
		TestIntUintFloat{"9" + maxUint8, "0", false, reflect.Uint8},
		TestIntUintFloat{"9" + maxUint16, "0", false, reflect.Uint16},
		TestIntUintFloat{"9" + maxUint32, "0", false, reflect.Uint32},
		TestIntUintFloat{"9" + maxUint64, "0", false, reflect.Uint64},
		TestIntUintFloat{"0", "0", false, reflect.Slice},
	}

	// Test correct values.
	for _, data := range tests {
		r, err := strToUintKind(data.Value, data.Kind)
		if data.Correct && err != nil {
			t.Error(err)
		} else if !data.Correct && err == nil {
			t.Errorf("the value %s should throw an exception", data.Value)
		} else if err != nil && r != 0 {
			t.Errorf("any error should return zero but returns %v", r)
		}

		control := fmt.Sprintf("%d", uint64(r))
		if control != data.Control {
			t.Errorf("expected %s but generated %s", data.Control, control)
		}
	}
}

// TestStrToFloatKind tests strToFloatKind function.
func TestStrToFloatKind(t *testing.T) {
	var (
		tests      []TestIntUintFloat
		maxFloat32 string = fmt.Sprintf("%.2f", math.MaxFloat32-1)
		maxFloat64 string = fmt.Sprintf("%.2f", math.MaxFloat64-1)
	)

	// Test data.
	tests = []TestIntUintFloat{
		TestIntUintFloat{"", "0.00", true, reflect.Float64},
		TestIntUintFloat{"0.0", "0.00", true, reflect.Float64},
		TestIntUintFloat{"3.0", "3.00", true, reflect.Float64},
		TestIntUintFloat{"-3.1", "-3.10", true, reflect.Float64},
		TestIntUintFloat{maxFloat32, maxFloat32, true, reflect.Float32},
		TestIntUintFloat{maxFloat64, maxFloat64, true, reflect.Float64},

		TestIntUintFloat{"string", "0.00", false, reflect.Float64},
		TestIntUintFloat{"9" + maxFloat32, "0.00", false, reflect.Float32},
		TestIntUintFloat{"9" + maxFloat64, "0.00", false, reflect.Float64},
		TestIntUintFloat{"0.00", "0.00", false, reflect.Slice},
	}

	// Test correct values.
	for _, data := range tests {
		r, err := strToFloatKind(data.Value, data.Kind)
		if data.Correct && err != nil {
			t.Error(err)
		} else if !data.Correct && err == nil {
			t.Errorf("the value %s should throw an exception", data.Value)
		} else if err != nil && r != 0 {
			t.Errorf("any error should return zero but returns %v", r)
		}

		control := fmt.Sprintf("%.2f", float64(r))
		if control != data.Control {
			t.Errorf("expected %s but generated %s", data.Control, control)
		}
	}
}
