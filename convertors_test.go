package env

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

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
