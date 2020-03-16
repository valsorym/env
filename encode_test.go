package env

import (
	"net/url"
	"testing"
)

// CustomMarshal structure with custom MarshalENV method.
type CustomMarshal struct {
	Host         string   `env:"HOST"`
	Port         int      `env:"PORT"`
	AllowedHosts []string `env:"ALLOWED_HOSTS,:"`
}

// MarshalENV the custom method for marshalling.
func (c *CustomMarshal) MarshalENV() ([]string, error) {
	// Test data set manually.
	Set("HOST", "192.168.0.1")
	Set("PORT", "80")
	Set("ALLOWED_HOSTS", "192.168.0.1")
	return []string{
		"HOST=192.168.0.1",
		"PORT=80",
		"ALLOWED_HOSTS=192.168.0.1",
	}, nil
}

// TestMarshalENVNilPointer tests marshalENV function
// for uninitialized pointer.
func TestMarshalENVNilPointer(t *testing.T) {
	type Empty struct{}
	var value *Empty
	if _, err := marshalENV(value, ""); err == nil {
		t.Error("exception expected for an uninitialized object")
	}
}

// TestMarshalENVNotStruct tests marshalENV function for not struct.
func TestMarshalNotStruct(t *testing.T) {
	var value string
	if _, err := marshalENV(value, ""); err == nil {
		t.Error("exception expected for an object other than structure")
	}
}

// TestMarshalENV tests marshalENV function with struct value.
func TestMarshalENV(t *testing.T) {
	type Struct struct {
		Host         string    `env:"HOST"`
		Port         int       `env:"PORT"`
		AllowedHosts []string  `env:"ALLOWED_HOSTS,!"`
		AllowedUsers [2]string `env:"ALLOWED_USERS,:"`
	}
	var value = Struct{
		"localhost",
		8080,
		[]string{"localhost", "127.0.0.1"},
		[2]string{"John", "Bob"},
	}

	Clear()
	_, err := marshalENV(value, "")
	if err != nil {
		t.Error(err)
	}

	// Test marshalling.
	if v := Get("HOST"); v != "localhost" {
		t.Errorf("Incorrect value set for HOST: %s", v)
	}

	if v := Get("PORT"); v != "8080" {
		t.Errorf("Incorrect value set for PORT: %s", v)
	}

	if v := Get("ALLOWED_HOSTS"); v != "localhost!127.0.0.1" {
		t.Errorf("Incorrect value set for ALLOWED_HOSTS: %s", v)
	}

	if v := Get("ALLOWED_USERS"); v != "John:Bob" {
		t.Errorf("Incorrect value set for ALLOWED_USERS: %s", v)
	}
}

// TestMarshalENVPtr tests marshalENV function for pointer of the struct value.
func TestMarshalENVPtr(t *testing.T) {
	type Struct struct {
		Host         string    `env:"HOST"`
		Port         int       `env:"PORT"`
		AllowedHosts []string  `env:"ALLOWED_HOSTS,!"`
		AllowedUsers [2]string `env:"ALLOWED_USERS,:"`
	}
	var value = &Struct{
		"localhost",
		8080,
		[]string{"localhost", "127.0.0.1"},
		[2]string{"John", "Bob"},
	}

	Clear()
	_, err := marshalENV(value, "")
	if err != nil {
		t.Error(err)
	}

	// Test marshalling.
	if v := Get("HOST"); v != "localhost" {
		t.Errorf("Incorrect value set for HOST: %s", v)
	}

	if v := Get("PORT"); v != "8080" {
		t.Errorf("Incorrect value set for PORT: %s", v)
	}

	if v := Get("ALLOWED_HOSTS"); v != "localhost!127.0.0.1" {
		t.Errorf("Incorrect value set for ALLOWED_HOSTS: %s", v)
	}

	if v := Get("ALLOWED_USERS"); v != "John:Bob" {
		t.Errorf("Incorrect value set for ALLOWED_USERS: %s", v)
	}
}

// TestMarshalENVCustom tests marshalENV function for object
// with custom MarshalENV method.
func TestMarshalENVCustom(t *testing.T) {
	var scope = CustomMarshal{
		"localhost",                        // default: 192.168.0.1
		8080,                               // default: 80
		[]string{"localhost", "127.0.0.1"}, // default: 192.168.0.1
	}

	Clear()
	_, err := marshalENV(scope, "")
	if err != nil {
		t.Error(err)
	}

	// Test marshalling.
	if v := Get("HOST"); v != "192.168.0.1" {
		t.Errorf("Incorrect value set for HOST: %s", v)
	}

	if v := Get("PORT"); v != "80" {
		t.Errorf("Incorrect value set for PORT: %s", v)
	}

	if v := Get("ALLOWED_HOSTS"); v != "192.168.0.1" {
		t.Errorf("Incorrect value set for ALLOWED_HOSTS: %s", v)
	}
}

// TestMarshalENVCustomPtr tests marshalENV function for pointer
// with custom MarshalENV method.
func TestMarshalENVCustomPtr(t *testing.T) {
	var scope = &CustomMarshal{
		"localhost",                        // default: 192.168.0.1
		8080,                               // default: 80
		[]string{"localhost", "127.0.0.1"}, // default: 192.168.0.1
	}

	Clear()
	_, err := marshalENV(scope, "")
	if err != nil {
		t.Error(err)
	}

	// Test marshalling.
	if v := Get("HOST"); v != "192.168.0.1" {
		t.Errorf("Incorrect value set for HOST: %s", v)
	}

	if v := Get("PORT"); v != "80" {
		t.Errorf("Incorrect value set for PORT: %s", v)
	}

	if v := Get("ALLOWED_HOSTS"); v != "192.168.0.1" {
		t.Errorf("Incorrect value set for ALLOWED_HOSTS: %s", v)
	}
}

// TestMarshalURL tests marshaling of the URL.
func TestMarshalURL(t *testing.T) {
	type URLTestType struct {
		KeyURLPlain      url.URL     `env:"KEY_URL_PLAIN"`
		KeyURLPoint      *url.URL    `env:"KEY_URL_POINT"`
		KeyURLPlainSlice []url.URL   `env:"KEY_URL_PLAIN_SLICE,!"`
		KeyURLPointSlice []*url.URL  `env:"KEY_URL_POINT_SLICE,!"`
		KeyURLPlainArray [2]url.URL  `env:"KEY_URL_PLAIN_ARRAY,!"`
		KeyURLPointArray [2]*url.URL `env:"KEY_URL_POINT_ARRAY,!"`
	}

	var test string
	var data = URLTestType{
		KeyURLPlain: url.URL{Scheme: "http", Host: "plain.example.com"},
		KeyURLPoint: &url.URL{Scheme: "http", Host: "point.example.com"},
		KeyURLPlainSlice: []url.URL{
			url.URL{Scheme: "http", Host: "a.plain.example.com"},
			url.URL{Scheme: "http", Host: "b.plain.example.com"},
		},
		KeyURLPointSlice: []*url.URL{
			&url.URL{Scheme: "http", Host: "a.point.example.com"},
			&url.URL{Scheme: "http", Host: "b.point.example.com"},
		},
		KeyURLPlainArray: [2]url.URL{
			url.URL{Scheme: "http", Host: "c.plain.example.com"},
			url.URL{Scheme: "http", Host: "d.plain.example.com"},
		},
		KeyURLPointArray: [2]*url.URL{
			&url.URL{Scheme: "http", Host: "c.point.example.com"},
			&url.URL{Scheme: "http", Host: "d.point.example.com"},
		},
	}

	Marshal(data)

	// Tests results.
	if v := Get("KEY_URL_PLAIN"); v != "http://plain.example.com" {
		t.Errorf("Incorrect marshaling plain url.URL: %s", v)
	}

	if v := Get("KEY_URL_POINT"); v != "http://point.example.com" {
		t.Errorf("Incorrect marshaling poin url.URL: %s", v)
	}

	// Plain slice.
	test = "http://a.plain.example.com!http://b.plain.example.com"
	if v := Get("KEY_URL_PLAIN_SLICE"); v != test {
		t.Errorf("Incorrect marshaling poin slice []url.URL: %s", v)
	}

	// Point slice.
	test = "http://a.point.example.com!http://b.point.example.com"
	if v := Get("KEY_URL_POINT_SLICE"); v != test {
		t.Errorf("Incorrect marshaling point slice []*url.URL: %s", v)
	}

	// Plain array.
	test = "http://c.plain.example.com!http://d.plain.example.com"
	if v := Get("KEY_URL_PLAIN_ARRAY"); v != test {
		t.Errorf("Incorrect marshaling plain array []url.URL: %s", v)
	}

	// Point array.
	test = "http://c.point.example.com!http://d.point.example.com"
	if v := Get("KEY_URL_POINT_ARRAY"); v != test {
		t.Errorf("Incorrect marshaling point array []*url.URL: %s", v)
	}
}

// TestMarshalStruct tests marshaling of the struct.
func TestMarshalStruct(t *testing.T) {
	type Address struct {
		Country string `env:"COUNTRY"`
	}

	type User struct {
		Name    string  `env:"NAME"`
		Address Address `env:"ADDRESS"`
	}

	type Client struct {
		User     User    `env:"USER"`
		HomePage url.URL `env:"HOME_PAGE"`
	}

	var data = Client{
		User: User{
			Name: "John",
			Address: Address{
				Country: "USA",
			},
		},
		HomePage: url.URL{Scheme: "http", Host: "example.com"},
	}

	// Marshaling.
	result, _ := Marshal(data)

	// Tests.
	if v := Get("USER_NAME"); v != "John" {
		t.Errorf("Incorrect marshaling (Name): %s\n%v", v, result)
	}

	if v := Get("USER_ADDRESS_COUNTRY"); v != "USA" {
		t.Errorf("Incorrect marshaling (Cuontry): %s\n%v", v, result)
	}

	if v := Get("HOME_PAGE"); v != "http://example.com" {
		t.Errorf("Incorrect marshaling url.URL (HomePage):%s", v)
	}
}

// TestMarshalStructPtr tests marshaling of the pointer on the struct.
func TestMarshalStructPtr(t *testing.T) {
	type Address struct {
		Country string `env:"COUNTRY"`
	}

	type User struct {
		Name    string   `env:"NAME"`
		Address *Address `env:"ADDRESS"`
	}

	type Client struct {
		User     *User    `env:"USER"`
		HomePage *url.URL `env:"HOME_PAGE"`
	}

	var data = Client{
		User: &User{
			Name: "John",
			Address: &Address{
				Country: "USA",
			},
		},
		HomePage: &url.URL{Scheme: "http", Host: "example.com"},
	}

	// Marshaling.
	result, _ := Marshal(data)

	// Tests.
	if v := Get("USER_NAME"); v != "John" {
		t.Errorf("Incorrect marshaling (Name): %s\n%v", v, result)
	}

	if v := Get("USER_ADDRESS_COUNTRY"); v != "USA" {
		t.Errorf("Incorrect marshaling (Cuontry): %s\n%v", v, result)
	}

	if v := Get("HOME_PAGE"); v != "http://example.com" {
		t.Errorf("Incorrect marshaling url.URL (HomePage):%s", v)
	}
}
