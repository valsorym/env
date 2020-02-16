package env

import (
	"testing"
)

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

// TestGetVariables tests getVariables function.
func TestGetVariables(t *testing.T) {
	var tests = map[string][]string{
		"The ${KEY_0}, and $KEY_1 ...":    []string{"KEY_0", "KEY_1"},
		"The ${KEY_0}01, and $KEY_10 ...": []string{"KEY_0", "KEY_10"},
	}

	for value, test := range tests {
		r := getVariables(value)
		for _, key := range test {
			if _, ok := r[key]; !ok {
				t.Errorf("The `%s` key not found.", key)
			}
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
