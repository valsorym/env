package env

import (
	"testing"
)

// TestReadParseStoreOpen tests exception handling when
// opening a nonexistent file.
func TestLoadReadParseStoreOpen(t *testing.T) {
	err := ReadParseStore("./examples/nonexist.env", true, false)
	if err == nil {
		t.Error("File descriptor leak.")
	}
}

// TestReadParseStoreExported checks the parsing of the
// file with the `export` command.
func TestReadParseStoreExported(t *testing.T) {
	var tests = map[string]string{
		`KEY_0`: `value 0`,
		`KEY_1`: `value 1`,
		`KEY_2`: `value_2`,
	}

	Clear()
	err := ReadParseStore("./examples/exported.env", true, false)
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

// TestReadParseStoreComments checks the parsing of the file with the
// comments and empty strings.
func TestReadParseStoreComments(t *testing.T) {
	var tests = map[string]string{
		`KEY_0`: `value 0`,
		`KEY_1`: `value 1`,
		`KEY_2`: `value_2`,
		`KEY_3`: `value_3`,
		`KEY_4`: `value_4:value_4:value_4`,
		`KEY_5`: `some text with # sharp sign and "escaped quotation" mark`,
	}

	Clear()
	err := ReadParseStore("./examples/comments.env", true, false)
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

// TestReadParseStoreWorngEqualKey tests problem with space
// before the equal sign.
func TestReadParseStoreWorngEqualKey(t *testing.T) {
	err := ReadParseStore("./examples/wrongequalkey.env", true, false)
	if err != incorrectKeyError {
		t.Error("Key error ignored")
	}

}

// TestReadParseStoreWorngEqualValue tests problem with space
// after the equal sign.
func TestReadParseStoreWorngEqualValue(t *testing.T) {
	err := ReadParseStore("./examples/wrongequalvalue.env", true, false)
	if err != incorrectValueError {
		t.Error("Value error ignored")
	}
}

// TestReadParseStoreIgnoreWorngEntry tests problem with space
// before and after the equal sign, and not correct lines.
func TestReadParseStoreIgnoreWorngEntry(t *testing.T) {
	var wrongentry = true
	var tests = map[string]string{
		`KEY_0`: `value_0`,
		`KEY_1`: `value_1`,
		`KEY_4`: `value_4`,
		`KEY_5`: `value`,
		`KEY_6`: `777`,
		`KEY_7`: `value_1`,
	}

	err := ReadParseStore("./examples/wrongentries.env", true, wrongentry)
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
	var tests = map[string]string{
		`KEY_0`: `value_0`,
		`KEY_1`: `value_001`,
		`KEY_2`: `value_001->correct value`,
	}
	err := ReadParseStore("./examples/variables.env", true, false)
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
