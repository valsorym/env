package env

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"testing"
)

// ...
type toIntUintFloat struct {
	Value   string
	Control string
	Correct bool
	Kind    reflect.Kind
}

// ...
type toBool struct {
	Value   string
	Control bool
	Correct bool
}

// TestStrToBool tests strToBool function.
func TestStrToBool(t *testing.T) {
	var tests = []toBool{
		toBool{"", false, true},
		toBool{"0", false, true},
		toBool{"1", true, true},
		toBool{"1.1", true, true},
		toBool{"-1.1", true, true},
		toBool{"0.0", false, true},
		toBool{"true", true, true},
		toBool{"True", true, true},
		toBool{"TRUE", true, true},
		toBool{"false", false, true},
		toBool{"False", false, true},
		toBool{"FALSE", false, true},
		toBool{"string", false, false},
		toBool{"a:b:c", false, false},
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
		tests    []toIntUintFloat
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
	tests = []toIntUintFloat{
		toIntUintFloat{"", "0", true, reflect.Int},
		toIntUintFloat{"0", "0", true, reflect.Int},
		toIntUintFloat{"-3", "-3", true, reflect.Int},
		toIntUintFloat{"3", "3", true, reflect.Int},

		toIntUintFloat{"-128", "-128", true, reflect.Int8},
		toIntUintFloat{"127", "127", true, reflect.Int8},

		toIntUintFloat{maxInt, maxInt, true, reflect.Int},
		toIntUintFloat{maxInt8, maxInt8, true, reflect.Int8},
		toIntUintFloat{maxInt16, maxInt16, true, reflect.Int16},
		toIntUintFloat{maxInt32, maxInt32, true, reflect.Int32},
		toIntUintFloat{maxInt64, maxInt64, true, reflect.Int64},

		toIntUintFloat{"string", "0", false, reflect.Int},
		toIntUintFloat{"3" + maxInt, "0", false, reflect.Int},
		toIntUintFloat{"3" + maxInt8, "0", false, reflect.Int8},
		toIntUintFloat{"-129", "0", false, reflect.Int8},
		toIntUintFloat{"128", "0", false, reflect.Int8},
		toIntUintFloat{"3" + maxInt16, "0", false, reflect.Int16},
		toIntUintFloat{"3" + maxInt32, "0", false, reflect.Int32},
		toIntUintFloat{"3" + maxInt64, "0", false, reflect.Int64},
		toIntUintFloat{"0", "0", false, reflect.Slice},
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
		tests     []toIntUintFloat
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
	tests = []toIntUintFloat{
		toIntUintFloat{"", "0", true, reflect.Uint},
		toIntUintFloat{"0", "0", true, reflect.Uint},
		toIntUintFloat{"3", "3", true, reflect.Uint},
		toIntUintFloat{maxUint, maxUint, true, reflect.Uint},
		toIntUintFloat{maxUint8, maxUint8, true, reflect.Uint8},
		toIntUintFloat{maxUint16, maxUint16, true, reflect.Uint16},
		toIntUintFloat{maxUint32, maxUint32, true, reflect.Uint32},
		toIntUintFloat{maxUint64, maxUint64, true, reflect.Uint64},

		toIntUintFloat{"string", "0", false, reflect.Uint},
		toIntUintFloat{"-3", "0", false, reflect.Uint},
		toIntUintFloat{"9" + maxUint, "0", false, reflect.Uint},
		toIntUintFloat{"9" + maxUint8, "0", false, reflect.Uint8},
		toIntUintFloat{"9" + maxUint16, "0", false, reflect.Uint16},
		toIntUintFloat{"9" + maxUint32, "0", false, reflect.Uint32},
		toIntUintFloat{"9" + maxUint64, "0", false, reflect.Uint64},
		toIntUintFloat{"0", "0", false, reflect.Slice},
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
		tests      []toIntUintFloat
		maxFloat32 string = fmt.Sprintf("%.2f", math.MaxFloat32-1)
		maxFloat64 string = fmt.Sprintf("%.2f", math.MaxFloat64-1)
	)

	// Test data.
	tests = []toIntUintFloat{
		toIntUintFloat{"", "0.00", true, reflect.Float64},
		toIntUintFloat{"0.0", "0.00", true, reflect.Float64},
		toIntUintFloat{"3.0", "3.00", true, reflect.Float64},
		toIntUintFloat{"-3.1", "-3.10", true, reflect.Float64},
		toIntUintFloat{maxFloat32, maxFloat32, true, reflect.Float32},
		toIntUintFloat{maxFloat64, maxFloat64, true, reflect.Float64},

		toIntUintFloat{"string", "0.00", false, reflect.Float64},
		toIntUintFloat{"9" + maxFloat32, "0.00", false, reflect.Float32},
		toIntUintFloat{"9" + maxFloat64, "0.00", false, reflect.Float64},
		toIntUintFloat{"0.00", "0.00", false, reflect.Slice},
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
