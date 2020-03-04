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
//    $ export KEY_0=VALUE_X
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
//    $ export KEY_0=VALUE_X
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
//    $ export KEY_0=VALUE_X
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

// UpdateSafe to loads data from env-file into environment with replacing
// existing values. Ignores the replecing of a ${var} or $var in a string.
//
// Returns an error in case of failure.
//
// Examples:
//
// Suppose that the some value was set into environment as:
//
//    $ export KEY_0=VALUE_X
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
//    //  KEY_0=VALUE_000         // data has been updated;
//    //  KEY_1=VALUE_001         // add new value;
//    //  KEY_2=VALUE_${LAST_ID}  // add new value without replecing $var.
//    print()
func UpdateSafe(filename string) error {
	var expand, update, forced = false, true, false
	return ReadParseStore(filename, expand, update, forced)
}

// Exists returns true if all keys sets in the environment.
// Examples:
//
// Suppose that the some value was set into environment as:
//
//    $ export KEY_0=VALUE_X
//
// And there is .env file with data:
//
//    LAST_ID=002
//    KEY_0=VALUE_000
//    KEY_1=VALUE_001
//    KEY_2=VALUE_${LAST_ID}
//
// Make code to check existing for a variable in the environment:
//
//    env.Exists("KEY_0")          // true
//    env.Exists("KEY_1")          // false
//    env.Exists("KEY_0", "KEY_1") // false
//
//    // Load values with replacement.
//    err := env.Update(".env")
//    if err != nil {
//        // something went wrong
//    }
//
//    env.Exists("KEY_1")          // true
//    env.Exists("KEY_0", "KEY_1") // true
func Exists(keys ...string) bool {
	for _, key := range keys {
		if _, ok := os.LookupEnv(key); !ok {
			return false
		}
	}
	return true
}

// Unmarshal to parses the environment data and stores the result in the value
// pointed to by scope. If scope isn't struct, not a pointer or is nil -
// returns an error.
//
// Supports the following field types: int, int8, int16, int32, int64, uin,
// uint8, uin16, uint32, uin64, float32, float64, string, bool and slice
// from thous types. For other filed's types will be returned an error.
//
// If the structure implements Unmarshaller interface - the custom UnmarshalENV
// method will be called.
//
// Structure fields may have a `env` tag as `env:"KEY[,SEP]"` where:
//
//    KEY - matches the name of the key in the environment;
//    SEP - optional argument, sets the separator for lists (default: space).
//
// Suppose that the some values was set into environment as:
//
//    $ export HOST="0.0.0.0"
//    $ export PORT=8080
//    $ export ALLOWED_HOSTS=localhost:127.0.0.1
//    $ export SECRET_KEY=AgBsdjONL53IKa33LM9SNROvD3hZXfoz
//
// Structure example:
//
//    // Config structure for containing values from the environment.
//    // P.s. No need to describe all the keys that are in the environment,
//    // for example, we ignore the SECRET_KEY key.
//    type Config struct {
//        Host         string   `env:"HOST"`
//        Port         int      `env:"PORT"`
//        AllowedHosts []string `env:"ALLOWED_HOSTS,:"`
//    }
//
// Unmarshal data from the environment into Config struct.
//
//    // Important: pointer to initialized structure!
//    var config = &Config{}
//
//    err := env.Unmarshal(config)
//    if err != nil {
//        // something went wrong
//    }
//
//    config.Host         // "0.0.0.0"
//    config.Port         // 8080
//    config.AllowedHosts // []string{"localhost", "127.0.0.1"}
//
// If the structure will havs custom UnmarshalENV - it will be called:
//
//    // UnmarshalENV it's custom method for unmarshalling.
//    func (c *Config) UnmarshalENV() error {
//        c.Host = "192.168.0.1"
//        c.Port = 80
//        c.AllowedHosts = []string{"192.168.0.1"}
//
//        return nil
//    }
//    ...
//    // The result will be the data that sets in the custom
//    // unmarshalling method.
//    config.Host         // "192.168.0.1"
//    config.Port         // 80
//    config.AllowedHosts // []string{"192.168.0.1"}
func Unmarshal(scope interface{}) error {
	return unmarshalENV(scope)
}

// Marshal converts the structure in to key/value and put it into environment
// with update old values. The first return value returns a map of the data
// that was correct set into environment. The seconden - error or nil.
//
// Supports the following field types: int, int8, int16, int32, int64, uin,
// uint8, uin16, uint32, uin64, float32, float64, string, bool and slice
// from thous types. For other filed's types will be returned an error.
//
// If the structure implements Marshaller interface - the custom MarshalENV
// method - will be called.
//
// Structure fields may have a `env` tag as `env:"KEY[,SEP]"` where:
//
//    KEY - matches the name of the key in the environment;
//    SEP - optional argument, sets the separator for lists (default: space).
//
// Structure example:
//
//    // Config structure.
//    type Config struct {
//        Host         string   `env:"HOST"`
//        Port         int      `env:"PORT"`
//        AllowedHosts []string `env:"ALLOWED_HOSTS,:"`
//    }
//
// Marshal data into environment from the Config.
//
//    // It can be a structure or a pointer to a structure.
//    var config = &Config{
//        "localhost",
//        8080,
//        []string{"localhost", "127.0.0.1"},
//    }
//
//    // Returns:
//    // map[ALLOWED_HOSTS:localhost:127.0.0.1 HOST:localhost PORT:8080], nil
//    _, err := env.Marshal(config)
//    if err != nil {
//        // something went wrong
//    }
//
//    env.Get("HOST")          // "localhost"
//    env.Get("PORT")          // "8080"
//    env.Get("ALLOWED_HOSTS") // "localhost:127.0.0.1"
//
// If object has MarshalENV and isn't a nil pointer - will be calls it
// method to scope convertation.
//
//    // MarshalENV it's custom method for marshalling.
//    func (c *Config) MarshalENV() ([]string, error ){
//        os.Setenv("HOST", "192.168.0.1")
//        os.Setenv("PORT", "80")
//        os.Setenv("ALLOWED_HOSTS", "192.168.0.1")
//
//        return []string{}, nil
//    }
//    ...
//    // The result will be the data that sets in the custom
//    // unmarshalling method.
//    env.Get("HOST")          // "192.168.0.1"
//    env.Get("PORT")          // "80"
//    env.Get("ALLOWED_HOSTS") // "192.168.0.1"
func Marshal(scope interface{}) ([]string, error) {
	return marshalENV(scope)
}
