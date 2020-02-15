package env

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	incorrectKeyError   = errors.New("bad key")
	incorrectValueError = errors.New("bad value")

	emptyRegex = regexp.MustCompile(
		`^(\s*)$|^(\s*[#].*)$`,
	)
	keyRegex = regexp.MustCompile(
		`^(?:\s*)?(?:export\s+)?(?P<key>[a-zA-Z_][a-zA-Z_0-9]*)=`,
	)
	valueRegex = regexp.MustCompile(
		`^=[^\s].*`,
	)
)

// isEmpty returns true if string contains separators or comment only.
func isEmpty(str string) bool {
	return emptyRegex.Match([]byte(str))
}

// removeInlineComment removes the comment in the string.
// Only in strings where the value is enclosed in quotes.
func removeInlineComment(str, quote string) string {
	// If the comment is in the string.
	if strings.Contains(str, "#") {
		chunks := strings.Split(str, "#")
		for i := range chunks {
			str := strings.Join(chunks[:i], "#")
			if len(str) > 0 && strings.Count(str, quote)%2 == 0 {
				return strings.TrimSpace(str)
			}
		}
	}
	return str
}

// parseExpression breaks expression into key and value, ignore
// comments and any spaces.
// Note: value must be an expression.
func parseExpression(exp string) (key, value string, err error) {
	var (
		quote  string = "\""
		marker string = fmt.Sprintf("<::%d::>", time.Now().Unix())
	)

	// Get key.
	// Remove `export` prefix, `=` suffix and trim spaces.
	tmp := keyRegex.FindStringSubmatch(exp)
	if len(tmp) < 2 {
		err = incorrectKeyError
		return
	}
	key = tmp[1]

	// Get value.
	// ... the `=` sign in the string.
	value = exp[strings.Index(exp, "="):]
	if !valueRegex.Match([]byte(value)) {
		err = incorrectValueError
		return
	}
	value = strings.TrimSpace(value[1:])

	switch {
	case strings.HasPrefix(value, "'"):
		quote = "'"
		fallthrough
	case strings.HasPrefix(value, "\""):
		// Replace escaped quotes, remove comment in the string,
		// check begin- and end- quotes and back escaped quotes.
		value = strings.Replace(value, fmt.Sprintf("\\%s", quote), marker, -1)
		value = removeInlineComment(value, quote)
		if strings.Count(value, quote)%2 != 0 { // begin- and end- quotes
			err = incorrectValueError
			return
		}
		value = value[1 : len(value)-1] // remove begin- and end- quotes
		// ... change `\"` and `\'` to `"` and `'`.
		value = strings.Replace(value, marker, fmt.Sprintf("%s", quote), -1)
	default:
		if strings.Contains(value, "#") {
			// Split by sharp sign and for string without quotes -
			// the first element has the meaning only.
			chunks := strings.Split(value, "#")
			chunks = strings.Split(chunks[0], " ")
			value = strings.TrimSpace(chunks[0])
		}
	}

	return
}
