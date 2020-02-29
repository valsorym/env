package env

import "os"

// Load to loads data from env-file into environment without replacing
// existing values. During loading replaces ${var} or $var in the string
// based on the data in the environment.
//
// Returns an error in case of failure.
//
// Examples:
//
// Suppose that the some value was set into environment as:
//
//    goloop$ export KEY_0=VALUE_X
//
// And there is .env file with data:
//
//    LAST_ID=002
//    KEY_0=VALUE_000
//    KEY_1=VALUE_001
//    KEY_2=VALUE_${LAST_ID}
//
// Make code to loads custom values into environment:
//
//    // The method prints environment variables starting with 'KEY_'.
//    print := func() {
//        for _, item := range env.Environ() {
//            if strings.HasPrefix(item, "KEY_") {
//                fmt.Println(item)
//            }
//        }
//    }
//
//    // Printed only:
//    //  KEY_0=VALUE_X
//    print()
//
//    // Load values without replacement.
//    err := env.Load(".env")
//    if err != nil {
//        // something went wrong
//    }
//
//    // Printed all variables:
//    //  KEY_0=VALUE_X    // not replaced by VALUE_001;
//    //  KEY_1=VALUE_001  // add new value;
//    //  KEY_2=VALUE_002  // add new value and replaced ${LAST_ID}
//                         // to the plain text.
//    print()
func Load(filename string) error {
	var expand, update, forced = true, false, false
	return ReadParseStore(filename, expand, update, forced)
}

// LoadSafe to loads data from env-file into environment without replacing
// existing values. Ignores the replecing of a ${var} or $var in a string.
//
// Returns an error in case of failure.
//
// Examples:
//
// Suppose that the some value was set into environment as:
//
//    goloop$ export KEY_0=VALUE_X
//
// And there is .env file with data:
//
//    LAST_ID=002
//    KEY_0=VALUE_000
//    KEY_1=VALUE_001
//    KEY_2=VALUE_${LAST_ID}
//
// Make code to loads custom values into environment:
//
//    // The method prints environment variables starting with 'KEY_'.
//    print := func() {
//        for _, item := range env.Environ() {
//            if strings.HasPrefix(item, "KEY_") {
//                fmt.Println(item)
//            }
//        }
//    }
//
//    // Printed only:
//    //  KEY_0=VALUE_X
//    print()
//
//    // Load values without replacement.
//    err := env.LoadSafe(".env")
//    if err != nil {
//        // something went wrong
//    }
//
//    // Printed all variables:
//    //  KEY_0=VALUE_X            // not replaced by VALUE_001;
//    //  KEY_1=VALUE_001         // add new value;
//    //  KEY_2=VALUE_${LAST_ID}  // add new value without replecing $var.
//    print()
func LoadSafe(filename string) error {
	var expand, update, forced = false, false, false
	return ReadParseStore(filename, expand, update, forced)
}

// Update to loads data from env-file into environment with replacing
// existing values. During loading replaces ${var} or $var in the string
// based on the data in the environment.
//
// Returns an error in case of failure.
//
// Examples:
//
// Suppose that the some value was set into environment as:
//
//    goloop$ export KEY_0=VALUE_X
//
// And there is .env file with data:
//
//    LAST_ID=002
//    KEY_0=VALUE_000
//    KEY_1=VALUE_001
//    KEY_2=VALUE_${LAST_ID}
//
// Make code to loads custom values into environment:
//
//    // The method prints environment variables starting with 'KEY_'.
//    print := func() {
//        for _, item := range env.Environ() {
//            if strings.HasPrefix(item, "KEY_") {
//                fmt.Println(item)
//            }
//        }
//    }
//
//    // Printed only:
//    //  KEY_0=VALUE_X
//    print()
//
//    // Load values with replacement.
//    err := env.Update(".env")
//    if err != nil {
//        // something went wrong
//    }
//
//    // Printed all variables:
//    //  KEY_0=VALUE_000  // data has been updated;
//    //  KEY_1=VALUE_001  // add new value;
//    //  KEY_2=VALUE_002  // add new value and replaced ${LAST_ID}
//                         // to the plain text.
//    print()
func Update(filename string) error {
	var expand, update, forced = true, true, false
	return ReadParseStore(filename, expand, update, forced)
}

// UpdateSafe loads keys with replacing existing ones.
func UpdateSafe(filename string) error {
	var expand, update, forced = false, true, false
	return ReadParseStore(filename, expand, update, forced)
}

// Exists returns true if all keys sets in the environment.
func Exists(keys ...string) bool {
	for _, key := range keys {
		if _, ok := os.LookupEnv(key); !ok {
			return false
		}
	}
	return true
}

// Unmarshal parses the environment data and stores the result in the value
// pointed to by scope. If scope is nil or not a pointer, Unmarshal returns
// an error.
//
// Supports the following field types: int, int8, int16, int32, int64, uin,
// uint8, uin16, uint32, uin64, float32, float64, string, bool and slice
// from thous types.
//
// To unmarshal environment into a value implementing the
// Unmarshaler interface, Unmarshal can to calls the
// custom UnmarshalENV method.
func Unmarshal(scope interface{}) error {
	return unmarshalENV(scope)
}

// Marshal converts the scope in to key/value and put it into environment
// with update old data.
//
// Returns nil or error if there are problems with marshaling.
//
//    type Config struct {
//        Host          string   `env:"HOST"`
//        Port          int      `env:"PORT"`
//        AllowedHosts  []string `env:"ALLOWED_HOSTS,:"`
//    }
//
//    // ...
//    var c = &Config{"0.0.0.0", 8080, []string{"localhost", "127.0.0.1"}}
//    _, err := env.Marshal(c)
//    if err != nil {
//        // problem with marshaling
//    }
//
//    // ...
//    // Result:
//    // HOST="0.0.0.0"
//    // PORT=8080
//    // ALLOWED_HOSTS="localhost:127.0.0.1"
//
// Supports the following field types: int, int8, int16, int32, int64, uin,
// uint8, uin16, uint32, uin64, float32, float64, string, bool and slice
// from thous types.
//
// If object has MarshalENV and isn't a nil pointer, Marshal calls its
// MarshalENV method to scope convertation.
//
//    type Config struct {
//        Host         string   `env:"HOST"`
//        Port         int      `env:"PORT"`
//        AllowedHosts []string `env:"ALLOWED_HOSTS,:"`
//    }
//
//    func (c *Config) MarshalENV() error {
//        str := strings.Replace(fmt.Sprint(c.AllowedHosts), " ", ":", -1)
//        os.Setenv("SERVER_ALLOWED_HOSTS", strings.Trim(str, "[]"))
//        os.Setenv("SERVER_PORT", fmt.Sprintf("%d", c.Port))
//        os.Setenv("SERVER_HOST", c.Host)
//        return nil
//    }
//
//    // ...
//    var c = &Config{"0.0.0.0", 8080, []string{"localhost", "127.0.0.1"}}
//    _, err := env.Marshal(c)
//    if err != nil {
//        // problem with marshaling
//    }
//
//    // ...
//    // Result:
//    // SERVER_HOST="0.0.0.0"
//    // SERVER_PORT=8080
//    // SERVER_ALLOWED_HOSTS="localhost:127.0.0.1"
func Marshal(scope interface{}) (map[string]string, error) {
	return marshalENV(scope)
}
