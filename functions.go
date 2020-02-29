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

// Unmarshal extracts the contents of the environment and populates
// the scope data structure.
func Unmarshal(scope interface{}) error {
	return unmarshalENV(scope)
}

// Marshal converts the scope in to key/value and put it into environment
// with update old data.
//
// Returns nil or error if there are problems with marshaling.
//
//
//
// Supports the following field types: int, int8, int16, int32, int64, uin,
// uint8, uin16, uint32, uin64, float32, float64, string, bool and slice
// from thous types.
//
// If object has MarshalENV and isn't a nil pointer, Marshal calls its
// MarshalENV method to scope convertation.
func Marshal(scope interface{}) error {
	return marshalENV(scope)
}
