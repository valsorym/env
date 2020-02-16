package env

import (
	"bufio"
	"os"
	"strings"
)

// ReadParseStore reads env-file, parse it to `key` and `value` and
// to store it into environment.
//
// Options:
//    filename   path to the env-file;
//    overwrite  if true to overwrites the set value in the environment
//               to the new one from the env-file;
//    wrongentry if true ignores wrong entries and loads all possible options,
//               without causing an exception.
//
// P.s. The function can be used to build more flexible tools.
func ReadParseStore(filename string, overwrite, wrongentry bool) (err error) {
	var (
		file       *os.File
		key, value string
	)

	// Open env-file.
	file, err = os.Open(filename)
	if err != nil {
		return // unable to open file
	}
	defer file.Close()

	// Parse file.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Get current line and ignore empty string or comments.
		str := scanner.Text()
		if isEmpty(str) {
			continue
		}

		// Parse expression.
		// The string containing the expression must be of the
		// format like: [export] KEY=VALUE [# Comment]
		key, value, err = parseExpression(str)
		if err != nil {
			if wrongentry {
				continue // ignore wrong entry
			}
			return // incorrect expression
		}

		// Overwrite or add new value.
		if overwrite || len(os.Getenv(key)) == 0 {
			if variables := getVariables(value); len(variables) != 0 {
				for key, item := range variables {
					value = strings.ReplaceAll(value, item, os.Getenv(key))
				}
			}
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}
