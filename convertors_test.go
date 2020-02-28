package env

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"testing"
)

type IntUintFloat struct {
	Value   string
	Control string
	Correct bool
	Kind    reflect.Kind
}

// TestStrToIntKind ...
func TestStrToIntKind(t *testing.T) {
	var (
		tests    []IntUintFloat
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
	tests = []IntUintFloat{
		IntUintFloat{"", "0", true, reflect.Int},
		IntUintFloat{"0", "0", true, reflect.Int},
		IntUintFloat{"-3", "-3", true, reflect.Int},
		IntUintFloat{"3", "3", true, reflect.Int},

		IntUintFloat{"-128", "-128", true, reflect.Int8},
		IntUintFloat{"127", "127", true, reflect.Int8},

		IntUintFloat{maxInt, maxInt, true, reflect.Int},
		IntUintFloat{maxInt8, maxInt8, true, reflect.Int8},
		IntUintFloat{maxInt16, maxInt16, true, reflect.Int16},
		IntUintFloat{maxInt32, maxInt32, true, reflect.Int32},
		IntUintFloat{maxInt64, maxInt64, true, reflect.Int64},

		IntUintFloat{"string", "0", false, reflect.Int},
		IntUintFloat{maxInt + "1", "0", false, reflect.Int},
		IntUintFloat{maxInt8 + "1", "0", false, reflect.Int8},
		IntUintFloat{"-129", "0", false, reflect.Int8},
		IntUintFloat{"128", "0", false, reflect.Int8},
		IntUintFloat{maxInt16 + "1", "0", false, reflect.Int16},
		IntUintFloat{maxInt32 + "1", "0", false, reflect.Int32},
		IntUintFloat{maxInt64 + "1", "0", false, reflect.Int64},
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

// TestStrToUintKind ...
func TestStrToUintKind(t *testing.T) {
	var (
		tests     []IntUintFloat
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
	tests = []IntUintFloat{
		IntUintFloat{"", "0", true, reflect.Uint},
		IntUintFloat{"0", "0", true, reflect.Uint},
		IntUintFloat{"3", "3", true, reflect.Uint},
		IntUintFloat{maxUint, maxUint, true, reflect.Uint},
		IntUintFloat{maxUint8, maxUint8, true, reflect.Uint8},
		IntUintFloat{maxUint16, maxUint16, true, reflect.Uint16},
		IntUintFloat{maxUint32, maxUint32, true, reflect.Uint32},
		IntUintFloat{maxUint64, maxUint64, true, reflect.Uint64},

		IntUintFloat{"string", "0", false, reflect.Uint},
		IntUintFloat{"-3", "0", false, reflect.Uint},
		IntUintFloat{maxUint + "1", "0", false, reflect.Uint},
		IntUintFloat{maxUint8 + "1", "0", false, reflect.Uint8},
		IntUintFloat{maxUint16 + "1", "0", false, reflect.Uint16},
		IntUintFloat{maxUint32 + "1", "0", false, reflect.Uint32},
		IntUintFloat{maxUint64 + "1", "0", false, reflect.Uint64},
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

// TestStrToFloatKind ...
func TestStrToFloatKind(t *testing.T) {
	var (
		tests      []IntUintFloat
		maxFloat32 string = fmt.Sprintf("%.2f", math.MaxFloat32-1)
		maxFloat64 string = fmt.Sprintf("%.2f", math.MaxFloat64-1)
	)

	// Test data.
	tests = []IntUintFloat{
		IntUintFloat{"", "0.00", true, reflect.Float64},
		IntUintFloat{"0.0", "0.00", true, reflect.Float64},
		IntUintFloat{"3.0", "3.00", true, reflect.Float64},
		IntUintFloat{"-3.1", "-3.10", true, reflect.Float64},
		IntUintFloat{maxFloat32, maxFloat32, true, reflect.Float32},
		IntUintFloat{maxFloat64, maxFloat64, true, reflect.Float64},

		IntUintFloat{"string", "0.00", false, reflect.Float64},
		//IntUintFloat{maxFloat32 + "1", "0.00", false, reflect.Float32},
		//IntUintFloat{maxFloat64 + "1", "0.00", false, reflect.Float64},
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

/*
// TestStrToInt64Correct tests strToInt64 function for correct values.
func TestStrToInt64Coreect(t *testing.T) {
	var tests = map[string]int64{
		"":   0,
		"0":  0,
		"-3": -3,
		"3":  3,
	}

	tests[fmt.Sprintf("%d", math.MaxInt64-1)] = math.MaxInt64 - 1

	// Test correct values.
	for value, test := range tests {
		r, err := strToInt64(value)
		if err != nil {
			t.Error(err)
		}

		if r != test {
			t.Errorf("Expected value %d but the result %d", test, r)
		}
	}
}

// TestStrToInt64Incorrect tests strToInt64 function for incorrect values.
func TestStrToInt64Incorrect(t *testing.T) {
	var tests = []string{
		"922337203685477580777",
		"3.14",
		"3,14",
		".14",
		"string",
	}

	// Test incorrect values.
	for _, value := range tests {
		_, err := strToInt64(value)
		if err == nil {
			t.Errorf("The value %s should throw an exception.", value)
		}
	}
}

// TestStrToInt32Correct tests strToInt32 function for correct values.
func TestStrToInt32Coreect(t *testing.T) {
	var tests = map[string]int32{
		"":   0,
		"0":  0,
		"-3": -3,
		"3":  3,
	}

	tests[fmt.Sprintf("%d", math.MaxInt32-1)] = math.MaxInt32 - 1

	// Test correct values.
	for value, test := range tests {
		r, err := strToInt32(value)
		if err != nil {
			t.Error(err)
		}

		if r != int64(test) {
			t.Errorf("Expected value %d but the result %d", test, r)
		}
	}
}

// TestStrToInt32Incorrect tests strToInt32 function for incorrect values.
func TestStrToInt32Incorrect(t *testing.T) {
	var tests = []string{
		"3.14",
		"3,14",
		".14",
		"string",
	}

	tests = append(tests, fmt.Sprintf("%d", math.MaxInt32+1))

	// Test incorrect values.
	for _, value := range tests {
		_, err := strToInt32(value)
		if err == nil {
			t.Errorf("The value %s should throw an exception.", value)
		}
	}
}

// TestStrToInt16Correct tests strToInt16 function for correct values.
func TestStrToInt16Coreect(t *testing.T) {
	var tests = map[string]int16{
		"":   0,
		"0":  0,
		"-3": -3,
		"3":  3,
	}

	tests[fmt.Sprintf("%d", math.MaxInt16-1)] = math.MaxInt16 - 1

	// Test correct values.
	for value, test := range tests {
		r, err := strToInt16(value)
		if err != nil {
			t.Error(err)
		}

		if r != int64(test) {
			t.Errorf("Expected value %d but the result %d", test, r)
		}
	}
}

// TestStrToInt16Incorrect tests strToInt16 function for incorrect values.
func TestStrToInt16Incorrect(t *testing.T) {
	var tests = []string{
		"3.14",
		"3,14",
		".14",
		"string",
	}

	tests = append(tests, fmt.Sprintf("%d", math.MaxInt16+1))

	// Test incorrect values.
	for _, value := range tests {
		_, err := strToInt16(value)
		if err == nil {
			t.Errorf("The value %s should throw an exception.", value)
		}
	}
}

// TestStrToInt8Correct tests strToInt8 function for correct values.
func TestStrToInt8Coreect(t *testing.T) {
	var tests = map[string]int8{
		"":   0,
		"0":  0,
		"-3": -3,
		"3":  3,
	}

	tests[fmt.Sprintf("%d", math.MaxInt8-1)] = math.MaxInt8 - 1

	// Test correct values.
	for value, test := range tests {
		r, err := strToInt8(value)
		if err != nil {
			t.Error(err)
		}

		if r != int64(test) {
			t.Errorf("Expected value %d but the result %d", test, r)
		}
	}
}

// TestStrToInt8Incorrect tests strToInt8 function for incorrect values.
func TestStrToInt8Incorrect(t *testing.T) {
	var tests = []string{
		"3.14",
		"3,14",
		".14",
		"string",
	}

	tests = append(tests, fmt.Sprintf("%d", math.MaxInt8+1))

	// Test incorrect values.
	for _, value := range tests {
		_, err := strToInt8(value)
		if err == nil {
			t.Errorf("The value %s should throw an exception.", value)
		}
	}
}

// TestStrToIntCorrect tests strToInt function for correct values.
func TestStrToIntCoreect(t *testing.T) {
	var tests = map[string]int{
		"":   0,
		"0":  0,
		"-3": -3,
		"3":  3,
	}

	// Set the maximum value for a specific platform.
	if strconv.IntSize == 32 {
		tests[fmt.Sprintf("%d", math.MaxInt32-1)] = math.MaxInt32 - 1
	} else if strconv.IntSize == 64 {
		tests[fmt.Sprintf("%d", math.MaxInt64-1)] = math.MaxInt64 - 1
	}

	// Test correct values.
	for value, test := range tests {
		r, err := strToInt(value)
		if err != nil {
			t.Error(err)
		}

		if r != int64(test) {
			t.Errorf("Expected value %d but the result %d", test, r)
		}
	}
}

// TestStrToIntIncorrect tests strToInt function for incorrect values.
func TestStrToIntIncorrect(t *testing.T) {
	var tests = []string{
		"922337203685477580777",
		"3.14",
		"3,14",
		".14",
		"string",
	}

	// Test incorrect values.
	for _, value := range tests {
		_, err := strToInt(value)
		if err == nil {
			t.Errorf("The value %s should throw an exception.", value)
		}
	}
}

// TestStrToUint64Correct tests strToUint64 function for correct values.
func TestStrToUint64Coreect(t *testing.T) {
	var tests = map[string]uint64{
		"":  0,
		"0": 0,
		"3": 3,

		"18446744073709551614": 18446744073709551614,
	}

	// Test correct values.
	for value, test := range tests {
		r, err := strToUint64(value)
		if err != nil {
			t.Error(err)
		}

		if r != test {
			t.Errorf("Expected value %d but the result %d", test, r)
		}
	}
}

// TestStrToUint64Incorrect tests strToUint64 function for incorrect values.
func TestStrToUint64Incorrect(t *testing.T) {
	var tests = []string{
		"922337203685477580777",
		"3.14",
		"3,14",
		".14",
		"string",
		"-1",
	}

	// Test incorrect values.
	for _, value := range tests {
		_, err := strToUint64(value)
		if err == nil {
			t.Errorf("The value %s should throw an exception.", value)
		}
	}
}

// TestStrToUint32Correct tests strToUint32 function for correct values.
func TestStrToUint32Coreect(t *testing.T) {
	var tests = map[string]uint32{
		"":  0,
		"0": 0,
		"3": 3,
	}

	tests[fmt.Sprintf("%d", math.MaxUint32-1)] = math.MaxUint32 - 1

	// Test correct values.
	for value, test := range tests {
		r, err := strToUint32(value)
		if err != nil {
			t.Error(err)
		}

		if r != uint64(test) {
			t.Errorf("Expected value %d but the result %d", test, r)
		}
	}
}

// TestStrToUint32Incorrect tests strToUint32 function for incorrect values.
func TestStrToUint32Incorrect(t *testing.T) {
	var tests = []string{
		"3.14",
		"3,14",
		".14",
		"string",
		"-1",
	}

	tests = append(tests, fmt.Sprintf("%d", math.MaxUint32+1))

	// Test incorrect values.
	for _, value := range tests {
		_, err := strToUint32(value)
		if err == nil {
			t.Errorf("The value %s should throw an exception.", value)
		}
	}
}

// TestStrToUint16Correct tests strToUint16 function for correct values.
func TestStrToUint16Coreect(t *testing.T) {
	var tests = map[string]uint16{
		"":  0,
		"0": 0,
		"3": 3,
	}

	tests[fmt.Sprintf("%d", math.MaxUint16-1)] = math.MaxUint16 - 1

	// Test correct values.
	for value, test := range tests {
		r, err := strToUint16(value)
		if err != nil {
			t.Error(err)
		}

		if r != uint64(test) {
			t.Errorf("Expected value %d but the result %d", test, r)
		}
	}
}

// TestStrToUint16Incorrect tests strToUint16 function for incorrect values.
func TestStrToUint16Incorrect(t *testing.T) {
	var tests = []string{
		"3.14",
		"3,14",
		".14",
		"string",
		"-1",
	}

	tests = append(tests, fmt.Sprintf("%d", math.MaxUint16+1))

	// Test incorrect values.
	for _, value := range tests {
		_, err := strToUint16(value)
		if err == nil {
			t.Errorf("The value %s should throw an exception.", value)
		}
	}
}

// TestStrToUint8Correct tests strToUint8 function for correct values.
func TestStrToUint8Coreect(t *testing.T) {
	var tests = map[string]uint8{
		"":  0,
		"0": 0,
		"3": 3,
	}

	tests[fmt.Sprintf("%d", math.MaxUint8-1)] = math.MaxUint8 - 1

	// Test correct values.
	for value, test := range tests {
		r, err := strToUint8(value)
		if err != nil {
			t.Error(err)
		}

		if r != uint64(test) {
			t.Errorf("Expected value %d but the result %d", test, r)
		}
	}
}

// TestStrToUint8Incorrect tests strToUint8 function for incorrect values.
func TestStrToUint8Incorrect(t *testing.T) {
	var tests = []string{
		"3.14",
		"3,14",
		".14",
		"string",
		"-1",
	}

	tests = append(tests, fmt.Sprintf("%d", math.MaxUint8+1))

	// Test incorrect values.
	for _, value := range tests {
		_, err := strToUint8(value)
		if err == nil {
			t.Errorf("The value %s should throw an exception.", value)
		}
	}
}

// TestStrToUintCorrect tests strToUint function for correct values.
func TestStrToUintCoreect(t *testing.T) {
	var tests = map[string]uint{
		"":  0,
		"0": 0,
		"3": 3,

		"18446744073709551614": 18446744073709551614,
	}

	// Test correct values.
	for value, test := range tests {
		r, err := strToUint(value)
		if err != nil {
			t.Error(err)
		}

		if r != uint64(test) {
			t.Errorf("Expected value %d but the result %d", test, r)
		}
	}
}

// TestStrToUintIncorrect tests strToUint function for incorrect values.
func TestStrToUintIncorrect(t *testing.T) {
	var tests = []string{
		"922337203685477580777",
		"3.14",
		"3,14",
		".14",
		"string",
		"-1",
	}

	// Test incorrect values.
	for _, value := range tests {
		_, err := strToUint(value)
		if err == nil {
			t.Errorf("The value %s should throw an exception.", value)
		}
	}
}

// TestStrToFloat64Correct tests strToFloat64 function for correct values.
func TestStrToFloat64Coreect(t *testing.T) {
	var tests = map[string]float64{
		"":     0.0,
		"0.0":  0.0,
		"3.3":  3.3,
		"-3.3": -3.3,
	}

	tests[fmt.Sprintf("%f", math.MaxFloat64-1)] = math.MaxFloat64 - 1

	// Test correct values.
	for value, test := range tests {
		r, err := strToFloat64(value)
		if err != nil {
			t.Error(err)
		}

		if math.Abs(r-test) > FloatAccuracy {
			t.Errorf("Expected value %f but the result %f", test, r)
		}
	}
}

// TestStrToFloat64Incorrect tests strToFloat64 function for incorrect values.
func TestStrToFloat64Incorrect(t *testing.T) {
	var tests = []string{
		"3,14",
		"string",
		"1-1",
	}

	// Test incorrect values.
	for _, value := range tests {
		_, err := strToFloat64(value)
		if err == nil {
			t.Errorf("The value %s should throw an exception.", value)
		}
	}
}

// TestStrToFloat32Correct tests strToFloat32 function for correct values.
func TestStrToFloat32Coreect(t *testing.T) {
	var tests = map[string]float32{
		"":     0.0,
		"0.0":  0.0,
		"3.3":  3.3,
		"-3.3": -3.3,
	}

	tests[fmt.Sprintf("%f", math.MaxFloat32-1)] = math.MaxFloat32 - 1

	// Test correct values.
	for value, test := range tests {
		r, err := strToFloat32(value)
		if err != nil {
			t.Error(err)
		}

		if math.Abs(r-float64(test)) > FloatAccuracy {
			t.Errorf("Expected value %f but the result %f", test, r)
		}
	}
}

// TestStrToFloat32Incorrect tests strToFloat32 function for incorrect values.
func TestStrToFloat32Incorrect(t *testing.T) {
	var tests = []string{
		"3,14",
		"string",
		"1-1",
	}

	// Test incorrect values.
	for _, value := range tests {
		_, err := strToFloat32(value)
		if err == nil {
			t.Errorf("The value %s should throw an exception.", value)
		}
	}
}

// TestStrToBoolCorrect tests strToBool function for correct values.
func TestStrToBoolCoreect(t *testing.T) {
	var tests = map[string]bool{
		"":      false,
		"0":     false,
		"1":     true,
		"1.1":   true,
		"True":  true,
		"TRUE":  true,
		"False": false,
		"FALSE": false,
	}

	// Test correct values.
	for value, test := range tests {
		r, err := strToBool(value)
		if err != nil {
			t.Error(err)
		}

		if r != test {
			t.Errorf("Expected value %t but the result %t", test, r)
		}
	}
}

// TestStrToBoolIncorrect tests strToBool function for incorrect values.
func TestStrToBoolIncorrect(t *testing.T) {
	var tests = []string{
		"a:b:c",
		"string",
	}

	// Test incorrect values.
	for _, value := range tests {
		_, err := strToBool(value)
		if err == nil {
			t.Errorf("The value %s should throw an exception.", value)
		}
	}
}
*/
