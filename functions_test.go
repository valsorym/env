package env

import (
	"testing"
)

// TestLoad tests Load function.
func TestLoad(t *testing.T) {
	var err error

	Clear()
	err = Set("KEY_0", "default")
	if err != nil {
		t.Error(err)
	}

	// Load env-file.
	err = Load("./fixtures/variables.env")
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
	var err error

	Clear()
	err = Set("KEY_0", "default")
	if err != nil {
		t.Error(err)
	}

	// Load env-file.
	err = LoadSafe("./fixtures/variables.env")
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
	var err error

	Clear()
	err = Set("KEY_0", "default")
	if err != nil {
		t.Error(err)
	}

	// Load env-file.
	err = Update("./fixtures/variables.env")
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
	var err error

	Clear()
	err = Set("KEY_0", "default")
	if err != nil {
		t.Error(err)
	}

	// Load env-file.
	err = UpdateSafe("./fixtures/variables.env")
	if err != nil {
		t.Error("Unable to open file.")
	}

	// Expand test.
	if v := Get("KEY_2"); v != "${KEY_0}01" { // UploadSafe not to do Expand
		t.Errorf("Expected value `${KEY_0}01` != `%s`.", v)
	}
}

// TestExists tests Exists function.
func TestExist(t *testing.T) {
	var (
		err   error
		tests = [][]string{
			{"KEY_0", "default"},
			{"KEY_1", "default"},
		}
	)

	Clear()
	for _, item := range tests {
		err = Set(item[0], item[1])
		if err != nil {
			t.Error(err)
		}
	}

	// Variables is exists.
	if !Exists("KEY_0") || !Exists("KEY_0", "KEY_1") {
		t.Error("Expected value `ture` != `false`.")
	}

	// Variables doesn't exists.
	if Exists("KEY_2") || Exists("KEY_0", "KEY_1", "KEY_2") {
		t.Error("Expected value `false` != `true`.")
	}
}
