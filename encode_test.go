package env

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

type ConfigPlain struct {
	Host         string   `env:"HOST"`
	Port         int      `env:"PORT"`
	AllowedHosts []string `env:"ALLOWED_HOSTS,:"`
}

type ConfigExtended struct {
	Host         string   `env:"HOST"`
	Port         int      `env:"PORT"`
	AllowedHosts []string `env:"ALLOWED_HOSTS,:"`
}

func (c *ConfigExtended) MarshalENV() error {
	str := strings.Replace(fmt.Sprint(c.AllowedHosts), " ", ":", -1)
	os.Setenv("ALLOWED_HOSTS", strings.Trim(str, "[]"))
	os.Setenv("PORT", fmt.Sprintf("%d", c.Port))
	os.Setenv("HOST", c.Host)
	return nil
}

// TestMarshal ...
func TestMarshalNotStruct(t *testing.T) {
	var scope string
	_, err := Marshal(scope)
	if err == nil {
		t.Error("exception expected for an object other than structure")
	}
}

func TestMarshalPointNil(t *testing.T) {
	var scope *ConfigPlain
	_, err := Marshal(scope)
	if err == nil {
		t.Error("exception expected for an uninitialized object")
	}
}

/*
// TestMarshal ...
func TestMarshalPlain(t *testing.T) {
	var config = ConfigPlain{}
	err := Update("./examples/config.env")
	if err != nil {
		t.Error(err)
	}

	Marshal(config)
	if config.Host != "0.0.0.0" || config.Port != 8080 {
		t.Errorf("incorrect data parsing: %v", config)
	}
}

func TestMarshalPoint(t *testing.T) {
	var config = &ConfigPlain{}
	err := Update("./examples/config.env")
	if err != nil {
		t.Error(err)
	}

	Marshal(config)
	if config.Host != "0.0.0.0" || config.Port != 8080 {
		t.Error("incorrect data parsing")
	}
}
*/
