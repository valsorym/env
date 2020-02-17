package env

import "os"

// Get retrieves the value of the environment variable named by the key.
// It returns the value, which will be empty if the variable is not present.
//
// P.s. It is synonym for the os.Getenv.
func Get(key string) string {
	return os.Getenv(key)
}

// Set sets the value of the environment variable named by the key.
// It returns an error, if any.
//
// P.s. It is synonym for the os.Setenv.
func Set(key, value string) error {
	return os.Setenv(key, value)
}

// Unset unsets a single environment variable.
//
// P.s. It is synonym for the os.Unsetenv.
func Unset(key string) error {
	return os.Unsetenv(key)
}

// Clear deletes all environment variables.
//
// P.s. It is synonym for the os.Clearenv.
func Clear() {
	os.Clearenv()
}

// Environ returns a copy of strings representing the environment,
// in the form "key=value".
//
// P.s. It is synonym for the os.Environ.
func Environ() []string {
	return os.Environ()
}

// Expand replaces ${var} or $var in the string according to the values
// of the current environment variables. References to undefined
// variables are replaced by the empty string.
//
// P.s. It is synonym for the os.ExpandEnv.
func Expand(value string) string {
	return os.Expand(value, os.Getenv)
}
