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

// Exist returns true if key sets in the environment.
func Exist(key string) bool {
	_, status := os.LookupEnv(key)
	return status
}
