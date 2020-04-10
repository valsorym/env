package env

import (
	"testing"
)

// TestReadParseStoreOpen tests function when try to open a nonexistent file.
func TestLoadReadParseStoreOpen(t *testing.T) {
	err := ReadParseStore("./fixtures/nonexist.env", false, false, false)
	if err == nil {
		t.Error("Reading from a nonexistent file.")
	}
}

// TestReadParseStoreExported checks the parsing of the
// env-file with the `export` command.
func TestReadParseStoreExported(t *testing.T) {
	var tests = map[string]string{
		"KEY_0": "value 0",
		"KEY_1": "value 1",
		"KEY_2": "value_2",
		"KEY_3": "value_0:value_1:value_2:value_3",
	}

	// Load env-file.
	Clear()
	err := ReadParseStore("./fixtures/exported.env", false, false, false)
	if err != nil {
		t.Error(err.Error())
	}

	// Compare with sample.
	for key, value := range tests {
		if v := Get(key); value != v {
			t.Errorf("Incorrect value for `%s` key: `%s`!=`%s`", key, value, v)
		}
	}
}

// TestReadParseStoreComments checks the parsing of the
// env-file with the comments and empty strings.
func TestReadParseStoreComments(t *testing.T) {
	var tests = map[string]string{
		"KEY_0": "value 0",
		"KEY_1": "value 1",
		"KEY_2": "value_2",
		"KEY_3": "value_3",
		"KEY_4": "value_4:value_4:value_4",
		"KEY_5": `some text with # sharp sign and "escaped quotation" mark`,
	}

	// Load env-file.
	Clear()
	err := ReadParseStore("./fixtures/comments.env", false, false, false)
	if err != nil {
		t.Error(err.Error())
	}

	// Compare with sample.
	for key, value := range tests {
		if v := Get(key); value != v {
			t.Errorf("Incorrect value for `%s` key: `%s`!=`%s`", key, value, v)
		}
	}
}

// TestReadParseStoreWorngEqualKey tests problem with
// spaces before the equal sign.
func TestReadParseStoreWorngEqualKey(t *testing.T) {
	err := ReadParseStore("./fixtures/wrongequalkey.env", false, false, false)
	if err == nil {
		t.Error("Must be an error.")
	}

}

// TestReadParseStoreWorngEqualValue tests problem with
// space after the equal sign.
func TestReadParseStoreWorngEqualValue(t *testing.T) {
	err := ReadParseStore("./fixtures/wrongequalvalue.env", false, true, false)
	if err == nil {
		t.Error("Must be an error")
	}
}

// TestReadParseStoreIgnoreWorngEntry tests to force loading with
// the incorrect lines.
func TestReadParseStoreIgnoreWorngEntry(t *testing.T) {
	var forced = true
	var tests = map[string]string{
		"KEY_0": "value_0",
		"KEY_1": "value_1",
		"KEY_4": "value_4",
		"KEY_5": "value",
		"KEY_6": "777",
		"KEY_7": "${KEY_1}",
	}

	// Load env-file.
	Clear()
	err := ReadParseStore("./fixtures/wrongentries.env", false, false, forced)
	if err != nil {
		t.Error(err.Error())
	}

	// Compare with sample.
	for key, value := range tests {
		if v := Get(key); value != v {
			t.Errorf("Incorrect value for `%s` key: `%s`!=`%s`", key, value, v)
		}
	}
}

// TestReadParseStoreVariables tests replacing variables with real values.
func TestReadParseStoreVariables(t *testing.T) {
	var expand = true
	var tests = map[string]string{
		"KEY_0": "value_0",
		"KEY_1": "value_1",
		"KEY_2": "value_001",
		"KEY_3": "value_001->correct value",
		"KEY_4": "value_0value_001",
	}

	// Load env-file.
	Clear()
	err := ReadParseStore("./fixtures/variables.env", expand, false, false)
	if err != nil {
		t.Error(err.Error())
	}

	// Compare with sample.
	for key, value := range tests {
		if v := Get(key); value != v {
			t.Errorf("Incorrect value for `%s` key: `%s`!=`%s`", key, value, v)
		}
	}
}

// TestReadParseStoreNotUpdate tests variable update protection.
func TestReadParseStoreNotUpdate(t *testing.T) {
	var (
		update = false
		err    error
	)

	// Set test data.
	Clear()
	err = Set("KEY_0", "") // set empty string
	if err != nil {
		t.Error(err)
	}

	// Read simple env-file with KEY_0.
	err = ReadParseStore("./fixtures/simple.env", false, update, false)
	if err != nil {
		t.Error(err.Error())
	}

	// Tests.
	if v := Get("KEY_0"); v != "" {
		t.Error("The value has been updated")
	}
}

// TestReadParseStoreUpdate tests variable update.
func TestReadParseStoreUpdate(t *testing.T) {
	var (
		update = true
		err    error
	)

	// Set test data.
	Clear()
	err = Set("KEY_0", "") // set empty string
	if err != nil {
		t.Error(err)
	}

	// Read simple env-file with KEY_0.
	err = ReadParseStore("./fixtures/simple.env", false, update, false)
	if err != nil {
		t.Error(err.Error())
	}

	// Tests.
	if v := Get("KEY_0"); v != "value 0" {
		t.Error("Variable not updated.")
	}
}
