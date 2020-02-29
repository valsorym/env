package env

import "os"

// Load loads keys without replacing existing ones and make expand.
func Load(filename string) error {
	var expand, update, forced = true, false, false
	return ReadParseStore(filename, expand, update, forced)
}

// LoadSafe loads keys without replacing existing ones.
func LoadSafe(filename string) error {
	var expand, update, forced = false, false, false
	return ReadParseStore(filename, expand, update, forced)
}

// Update loads keys with replacing existing ones and make expand.
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
//    err := env.Marshal(c)
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
//    err := env.Marshal(c)
//    if err != nil {
//        // problem with marshaling
//    }
//
//    // ...
//    // Result:
//    // SERVER_HOST="0.0.0.0"
//    // SERVER_PORT=8080
//    // SERVER_ALLOWED_HOSTS="localhost:127.0.0.1"
func Marshal(scope interface{}) error {
	return marshalENV(scope)
}
