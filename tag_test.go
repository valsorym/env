package env

import "testing"

// TestSplitFieldTag tests splitFieldTag function.
func TestSplitFieldTag(t *testing.T) {
	type sample struct {
		tag, key, value, sep string
	}

	var (
		msg     = "for `%s` incorrect %s: `%s` != `%s`"
		correct = []sample{
			{
				tag:   "HOST",
				key:   "HOST",
				value: "",
				sep:   ":",
			},
			{
				tag:   ",{localhost,0.0.0.0}",
				key:   "",
				value: "localhost,0.0.0.0",
				sep:   ":",
			},
		}
		incorrect = []string{
			"KEY,'0.0.0.0",
			"KEY,\"0.0.0.0",
			"KEY,{0.0.0.0,:",
			"KEY,{8080,8081} 80,!",
		}
	)

	// Tests.
	for _, check := range correct {
		key, value, sep, err := splitFieldTag(check.tag)
		if err != nil {
			t.Error(err)
		}

		if key != check.key {
			t.Errorf(msg, check.tag, "key", key, check.key)
		}

		if value != check.value {
			t.Errorf(msg, check.tag, "value", value, check.value)
		}

		if sep != check.sep {
			t.Errorf(msg, check.tag, "sep", sep, check.sep)
		}
	}

	// Tests.
	for _, sample := range incorrect {
		_, _, _, err := splitFieldTag(sample)
		if err == nil {
			t.Error("there must be a error for expression:", sample)
		}
	}
}
