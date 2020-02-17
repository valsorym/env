package env

import (
	"testing"
)

// TestLoad tests Load function.
func TestLoad(t *testing.T) {
	Clear()
	Set("KEY_0", "default")

	// Load env-file.
	err := Load("./examples/variables.env")
	if err != nil {
		t.Error("Unable to open file.")
	}

	// Variable update protection.
	if Get("KEY_0") != "default" {
		t.Error("The existing variable has been overwritten.")
	}

	// Setting a new variable.
	if Get("KEY_1") != "value_1" {
		t.Error("Data was don't loaded.")
	}

	// Expand test.
	if v := Get("KEY_2"); v != "default01" { // KEY_0 not overwritten
		t.Errorf("Expected value `default01` != `%s`.", v)
	}
}

// TestLoadSafe tests LoadSafe function.
func TestLoadSafe(t *testing.T) {
	Clear()
	Set("KEY_0", "default")

	// Load env-file.
	err := LoadSafe("./examples/variables.env")
	if err != nil {
		t.Error("Unable to open file.")
	}

	// Expand test.
	if v := Get("KEY_2"); v != "${KEY_0}01" { // LoadSafe not to do Expand
		t.Errorf("Expected value `${KEY_0}01` != `%s`.", v)
	}
}

// TestUpdate tests Update function.
func TestUpdate(t *testing.T) {
	Clear()
	Set("KEY_0", "default")

	// Load env-file.
	err := Update("./examples/variables.env")
	if err != nil {
		t.Error("Unable to open file.")
	}

	// Variable update protection.
	if Get("KEY_0") == "default" {
		t.Error("The existing variable hasen't overwritten.")
	}

	// Setting a new variable.
	if Get("KEY_1") != "value_1" {
		t.Error("Data was don't loaded.")
	}

	// Expand test.
	if v := Get("KEY_2"); v != "value_001" { // KEY_0 not overwritten
		t.Errorf("Expected value `value_001` != `%s`.", v)
	}
}

// TestUpdateSafe tests UpdateSafe function.
func TestUpdateSafe(t *testing.T) {
	Clear()
	Set("KEY_0", "default")

	// Load env-file.
	err := UpdateSafe("./examples/variables.env")
	if err != nil {
		t.Error("Unable to open file.")
	}

	// Expand test.
	if v := Get("KEY_2"); v != "${KEY_0}01" { // UploadSafe not to do Expand
		t.Errorf("Expected value `${KEY_0}01` != `%s`.", v)
	}
}

// TestExist tests Exist function.
func TestExist(t *testing.T) {
	Clear()
	Set("KEY_0", "default")

	if !Exist("KEY_0") {
		t.Error("Expected value `ture` != `false`.")
	}

	if Exist("KEY_1") {
		t.Error("Expected value `false` != `true`.")
	}
}
