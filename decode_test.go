package env

/*
import (
	"net/url"
	"testing"
)

// Address it's a test structure with arbitrary fields.
type Address struct {
	Country string `env:"COUNTRY"`
	City    string `env:"CITY"`
}

// User it's a test structure with arbitrary fields.
type User struct {
	FirstName string `env:"FIRST_NAME"`
	LastName  string `env:"LAST_NAME"`
}

// UserPtr it's a test structure with structure pointer.
type UserPtr struct {
	User
	Address *Address `env:"ADDRESS"`
}

// UserPlain it's a test structure with nested structure.
type UserPlain struct {
	User
	Address Address `env:"ADDRESS"`
}

// Client it's arbitrary test structure.
type Client struct {
	Email string `env:"EMAIL"`
}

// ClientPtr it's test structure with a pointer to a structure that
// has a pointer to an another structure.
type ClientPtr struct {
	Client
	User     *UserPtr `env:"USER"`
	HomePage *url.URL `env:"HOME_PAGE"`
}

// ClientPlain it's test structure with a nested structure that
// has an another nested structure.
type ClientPlain struct {
	Client
	User     UserPlain `env:"USER"`
	HomePage url.URL   `env:"HOME_PAGE"`
}

// TestIsNotPointerError tests an exception for a value other than a pointer.
func TestIsNotPointerError(t *testing.T) {
	var client = ClientPlain{}
	if err := unmarshalENV(client, ""); err != IsNotPointerError {
		t.Error("need to raise an exception: IsNotPointerError")
	}
}

// TestIsNotPointerError tests an exception for an uninitialized value.
func TestIsNotInitializedError(t *testing.T) {
	var client *ClientPlain
	if err := unmarshalENV(client, ""); err != IsNotInitializedError {
		t.Error("need to raise an exception: IsNotInitializedError")
	}
}

// TestIsNotStructError tests an exception for a value that isn't structure.
func TestIsNotStructError(t *testing.T) {
	var client = new(int)
	if err := unmarshalENV(client, ""); err != IsNotStructError {
		t.Error("need to raise an exception: IsNotStructError")
	}
}
*/
