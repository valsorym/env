package env

import (
	"os"
	"strings"
	"testing"
)

// TestGet tests Get function.
func TestGet(t *testing.T) {
	os.Clearenv()
	os.Setenv("KEY_0", "A")
	os.Setenv("KEY_1", "B")

	if a, b := Get("KEY_0"), Get("KEY_1"); a != "A" || b != "B" {
		t.Errorf("Get function doesn't work.")
	}
}

// TestSet tests Set function.
func TestSet(t *testing.T) {
	os.Clearenv()
	Set("KEY_0", "A")
	Set("KEY_1", "B")

	if a, b := os.Getenv("KEY_0"), os.Getenv("KEY_1"); a != "A" || b != "B" {
		t.Errorf("Set function doesn't work.")
	}
}

// TestUnset tests Unset function.
func TestUnset(t *testing.T) {
	os.Clearenv()
	Set("KEY_0", "A")
	Set("KEY_1", "B")

	Unset("KEY_0")

	if a, b := Get("KEY_0"), Get("KEY_1"); !(a != "A" || b == "B") {
		t.Errorf("Unset function doesn't work.")
	}
}

// TestClear tests Clear function.
func TestClear(t *testing.T) {
	os.Clearenv()
	Set("KEY_0", "A")
	Set("KEY_1", "B")

	Clear()

	if a, b := Get("KEY_0"), Get("KEY_1"); a == "A" || b == "B" {
		t.Errorf("Clear function doesn't work.")
	}
}

// TestEnviron tests Environ function.
func TestEnviron(t *testing.T) {
	var tests = map[string]string{
		"KEY_0": "A",
		"KEY_1": "B",
		"KEY_2": "C",
	}

	Clear()
	for key, value := range tests {
		Set(key, value)
	}

	for _, value := range Environ() {
		p := strings.Split(value, "=")
		if tests[p[0]] != p[1] {
			t.Errorf("Environ function doesn't work for `%s` key: `%s`!=`%s`",
				p[0],
				tests[p[0]],
				p[1])
		}
	}
}

// TestExpand tests Expand function.
func TestExpand(t *testing.T) {
	var tests = [][]string{
		[]string{"${KEY_0}$KEY_0$KEY_0", "777"},
		[]string{"${KEY_2}$KEY_1$KEY_0", "357"},
	}

	Clear()
	Set("KEY_0", "7")
	Set("KEY_1", "5")
	Set("KEY_2", "3")

	// Tests.
	for _, item := range tests {
		test, expect := item[0], item[1]
		if v := Expand(test); v != expect {
			t.Errorf("Expand error: `%s`!=`%s`", v, expect)
		}
	}
}
