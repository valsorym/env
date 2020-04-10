// Copyright (c) 2020, GoLoop. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package env implements a set of functions of the environment management.

The env package allows to control environment's variables: set new variables,
delete variables, get value by variable name, load env-file into environment,
store environment values in a data structure.

*/
package env // import "github.com/goloop/env"

import (
	"bufio"
	"os"
)

// ReadParseStore reads env-file, parse it to `key` and `value` and
// to store it into environment.
//
// Arguments:
//    filename path to the env-file;
//    expand   if true replaces ${var} or $var in the string according
//             to the values of the current environment variables;
//    update   if true to overwrites the set value in the environment
//             to the new one from the env-file;
//    forced   if true ignores wrong entries and loads all possible options,
//             without causing an exception.
//
// P.s. The function can be used to build more flexible tools.
func ReadParseStore(filename string, expand, update, forced bool) (err error) {
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
			if forced {
				continue // ignore wrong entry
			}
			return // incorrect expression
		}

		// Overwrite or add new value.
		if _, ok := os.LookupEnv(key); update || !ok {
			if expand {
				value = Expand(value)
			}
			err = Set(key, value)
			if err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}
