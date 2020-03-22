package env

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"testing"
)

// UIFDataTestType the uint, int and float test type.
type UIFDataTestType struct {
	Value   string
	Control string
	Correct bool
	Kind    reflect.Kind
}

// BoolDataTestType the boolean test type.
type BoolDataTestType struct {
	Value   string
	Control bool
	Correct bool
}

// TestParseFieldTag tests parseFieldTag function.
func TestParseFieldTag(t *testing.T) {
	var tests = [][]string{
		//       tagValue, defaultName, defaultSep
		[]string{"", "HOST", " ", "HOST", " "},
		[]string{"HOST", "host", " ", "HOST", " "},
		[]string{"PATHS,:", "paths", " ", "PATHS", ":"},
		[]string{",:", "PORT", " ", "PORT", ":"},
		[]string{",", "PORT", ":", "PORT", ":"},
	}

	for _, test := range tests {
		name, sep := parseFieldTag(test[0], test[1], test[2])
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
		if _, _, err := parseExpression(test); err != KeyError {
			t.Errorf("For `%s` value must be KeyError.", test)
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
		if _, _, err := parseExpression(test); err != ValueError {
			t.Errorf("For `%s` value must be ValueError.", test)
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
	var tests = []BoolDataTestType{
		BoolDataTestType{"", false, true},
		BoolDataTestType{"0", false, true},
		BoolDataTestType{"1", true, true},
		BoolDataTestType{"1.1", true, true},
		BoolDataTestType{"-1.1", true, true},
		BoolDataTestType{"0.0", false, true},
		BoolDataTestType{"true", true, true},
		BoolDataTestType{"True", true, true},
		BoolDataTestType{"TRUE", true, true},
		BoolDataTestType{"false", false, true},
		BoolDataTestType{"False", false, true},
		BoolDataTestType{"FALSE", false, true},
		BoolDataTestType{"string", false, false},
		BoolDataTestType{"a:b:c", false, false},
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
		tests    []UIFDataTestType
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
	tests = []UIFDataTestType{
		UIFDataTestType{"", "0", true, reflect.Int},
		UIFDataTestType{"0", "0", true, reflect.Int},
		UIFDataTestType{"-3", "-3", true, reflect.Int},
		UIFDataTestType{"3", "3", true, reflect.Int},

		UIFDataTestType{"-128", "-128", true, reflect.Int8},
		UIFDataTestType{"127", "127", true, reflect.Int8},

		UIFDataTestType{maxInt, maxInt, true, reflect.Int},
		UIFDataTestType{maxInt8, maxInt8, true, reflect.Int8},
		UIFDataTestType{maxInt16, maxInt16, true, reflect.Int16},
		UIFDataTestType{maxInt32, maxInt32, true, reflect.Int32},
		UIFDataTestType{maxInt64, maxInt64, true, reflect.Int64},

		UIFDataTestType{"string", "0", false, reflect.Int},
		UIFDataTestType{"3" + maxInt, "0", false, reflect.Int},
		UIFDataTestType{"3" + maxInt8, "0", false, reflect.Int8},
		UIFDataTestType{"-129", "0", false, reflect.Int8},
		UIFDataTestType{"128", "0", false, reflect.Int8},
		UIFDataTestType{"3" + maxInt16, "0", false, reflect.Int16},
		UIFDataTestType{"3" + maxInt32, "0", false, reflect.Int32},
		UIFDataTestType{"3" + maxInt64, "0", false, reflect.Int64},
		UIFDataTestType{"0", "0", false, reflect.Slice},
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
		tests     []UIFDataTestType
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
	tests = []UIFDataTestType{
		UIFDataTestType{"", "0", true, reflect.Uint},
		UIFDataTestType{"0", "0", true, reflect.Uint},
		UIFDataTestType{"3", "3", true, reflect.Uint},
		UIFDataTestType{maxUint, maxUint, true, reflect.Uint},
		UIFDataTestType{maxUint8, maxUint8, true, reflect.Uint8},
		UIFDataTestType{maxUint16, maxUint16, true, reflect.Uint16},
		UIFDataTestType{maxUint32, maxUint32, true, reflect.Uint32},
		UIFDataTestType{maxUint64, maxUint64, true, reflect.Uint64},

		UIFDataTestType{"string", "0", false, reflect.Uint},
		UIFDataTestType{"-3", "0", false, reflect.Uint},
		UIFDataTestType{"9" + maxUint, "0", false, reflect.Uint},
		UIFDataTestType{"9" + maxUint8, "0", false, reflect.Uint8},
		UIFDataTestType{"9" + maxUint16, "0", false, reflect.Uint16},
		UIFDataTestType{"9" + maxUint32, "0", false, reflect.Uint32},
		UIFDataTestType{"9" + maxUint64, "0", false, reflect.Uint64},
		UIFDataTestType{"0", "0", false, reflect.Slice},
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
		tests      []UIFDataTestType
		maxFloat32 string = fmt.Sprintf("%.2f", math.MaxFloat32-1)
		maxFloat64 string = fmt.Sprintf("%.2f", math.MaxFloat64-1)
	)

	// Test data.
	tests = []UIFDataTestType{
		UIFDataTestType{"", "0.00", true, reflect.Float64},
		UIFDataTestType{"0.0", "0.00", true, reflect.Float64},
		UIFDataTestType{"3.0", "3.00", true, reflect.Float64},
		UIFDataTestType{"-3.1", "-3.10", true, reflect.Float64},
		UIFDataTestType{maxFloat32, maxFloat32, true, reflect.Float32},
		UIFDataTestType{maxFloat64, maxFloat64, true, reflect.Float64},

		UIFDataTestType{"string", "0.00", false, reflect.Float64},
		UIFDataTestType{"9" + maxFloat32, "0.00", false, reflect.Float32},
		UIFDataTestType{"9" + maxFloat64, "0.00", false, reflect.Float64},
		UIFDataTestType{"0.00", "0.00", false, reflect.Slice},
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
