# env

Package `env` it's simple lib for manage [environment's variables](https://en.wikipedia.org/wiki/Environment_variable), parse and convert it into Go-structures.

The `env` package allows to control environment's variables: 

* set, read and delete variables;
* load or update environment's variables from env-file;
* to unmarshal environment variables into Go-structure;
* to marshal Go-structure's fields into environment.

## Quick example

Let’s imagine the task. There is a web-project that is develop and tests on the local computer and deploys on the production server. On the production server some settings (for example `host` and `port`) enforced seated in an environment but on the local computer data must be loaded from the file (because different team members have different launch options).

For example, the local-configuration file `.env` with variables look likes:

```
HOST=0.0.0.0
PORT=8080
ALLOWED_HOSTS="localhost:127.0.0.1"

# In configurations file and/or in the environment there can be
# a many of variables that willn't be parsed, like this:
SECRET_KEY="N0XRABLZ5ZZY6dCNgD7pLjTIlx8v4G5d"
```

So, make test project like `main.go`. We need to load the missing configurations from the file into the environment. Convert data from the environment into Go-structure. And to use this data to start the server.

**Note:** We need to load the missing variables but not update the existing ones.

```
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/goloop/env"
)

// Config it's struct of the server configuration.
type Config struct {
	Host         string   `env:"HOST"`
	Port         int      `env:"PORT"`
	AllowedHosts []string `env:"ALLOWED_HOSTS,:"` // parse by `:`.

	// P.s. It isn't necessary to specify all the keys
	// that are available in the environment.
}

// Addr returns the server's address - concatenate host
// and port into one string.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Home it is handler of the home page.
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	var config = Config{}

	// Load configurations from the env-file into an environment.
	// P.s. We use the Load (but not Update) method for as not to overwrite
	// the data in the environment because on the production server this
	// data can be set forcibly.
	env.Load(".env") // set correct the path to the file with variables

	// Parsing of the environment data and storing the result
	// into object by pointer.
	env.Unmarshal(&config)

	// Make routing.
	http.HandleFunc("/", Home)

	// Run.
	log.Printf("Server started on %s\n", config.Addr())
	log.Printf("Allowed hosts: %v\n", config.AllowedHosts)
	http.ListenAndServe(config.Addr(), nil)
}
```

## Types

Marshal/Unmarshal methods  supports the following field's types: `int`, `int8`, `int16`, `int32`, `int64`, `uin`, `uint8`, `uin16`, `uint32`, `uin64`, `float32`, `float64`, `string`, `bool`, `url.URL` and `pointers`, `array` or `slice` from thous types *(i.e. `*int`, ..., `[]int`, ..., `[]bool`, ..., `[2]*url.URL`, etc.)*. The nested structures will be processed recursively.

For other filed's types (like `chan` or `map` ...) will be returned an error.

Example:

```
package main

import (
	"log"
	"net/url"

	"github.com/goloop/env"
)

type Address struct {
	City string `env:"CITY"`
}

type User struct {
	Name        string   `env:"NAME"`
	Address     *Address `env:"ADDRESS"`       // can be as pointer
	Permissions []bool   `env:"PERMISSIONS,;"` // separator like `;`
}

type Client struct {
	ID       int      // default env-variable name
	Email    string   `env:"EMAIL"`
	HomePage *url.URL `env:"HOME_PAGE"` // can be as pointer
	User     User     `env:"USER"`      // ... or like value
}

func main() {
	var clientA, clientB Client
	clientA = Client{
		ID:       3,
		Email:    "mail@example.com",
		HomePage: &url.URL{Scheme: "http", Host: "example.com"},
		User: User{
			Name:        "valsorym",
			Address:     &Address{City: "Chernihiv"},
			Permissions: []bool{true, true, true, false, false},
		},
	}

	// Save data of the clientA.
	if _, err := env.Marshal(clientA); err != nil {
		log.Fatal(err)
	}

	// Load data into clientB.
	if err := env.Unmarshal(&clientB); err != nil { // need to use pointer
		log.Fatal(err)
	}

	// Obect clientB:
	// clientB.ID                // 3
	// clientB.Email             // mail@example.com
	// clientB.HomePage.String() // http://example.com
	// clientB.User.Name         // valsorym
	// clientB.User.Address.City // Chernihiv
	// clientB.User.Permissions  // [true true true false false]

	// Envirenment:
	// ID=3
	// EMAIL=mail@example.com
	// HOME_PAGE=http://example.com
	// USER_NAME=valsorym
	// USER_ADDRESS_CITY=Chernihiv
	// USER_PERMISSIONS=true;true;true;false;false
}
```

# Get lib

Use  `go get`:

```
$ go get -u github.com/goloop/env
```

or as  `git clone`  and make link in `$GOPATH`:

```
$ mkdir -p ~/workspace && \
  cd ~/workspace && \
  git clone https://github.com/goloop/env.git && \
  cd env && \
  make link
```
To use the `env` package import it as:

```
import (
    // ...
    "github.com/goloop/env"
)
```

# Functions

## Load

The `Load` to loads data from env-file into environment without replacing existing values. During loading replaces `${var}` or `$var` in the string based on the data in the environment.

Returns an error in case of failure.

### Examples:

Suppose that the some value was set into environment as:

```
$ export KEY_0=VALUE_X
```

And there is `.env` file with data:

```
LAST_ID=002
KEY_0=VALUE_000
KEY_1=VALUE_001
KEY_2=VALUE_${LAST_ID}
```

Make code to loads custom values into environment:

```
// Environment:
// KEY_0=VALUE_X

// Load values without replacement.
err := env.Load(".env")
if err != nil {
    // something went wrong
}

// Environment:
// KEY_0=VALUE_X    // not replaced by VALUE_001;
// KEY_1=VALUE_001  // add new value;
// KEY_2=VALUE_002  // add new value and replaced ${LAST_ID}
                    // to the plain text.
```
## LoadSafe

The `LoadSafe` to loads data from env-file into environment without replacing existing values. Ignores the replacing of a `${var}` or `$var` in a string.

Returns an error in case of failure.

### Examples:

Suppose that the some value was set into environment as:

```
$ export KEY_0=VALUE_X
```

And there is `.env` file with data:

```
LAST_ID=002
KEY_0=VALUE_000
KEY_1=VALUE_001
KEY_2=VALUE_${LAST_ID}
```

Make code to loads custom values into environment:

```
// Environment:
// KEY_0=VALUE_X

// Load values without replacement.
err := env.LoadSafe(".env")
if err != nil {
    // something went wrong
}

// Environment:
// KEY_0=VALUE_X           // not replaced by VALUE_001;
// KEY_1=VALUE_001         // add new value;
// KEY_2=VALUE_${LAST_ID}  // add new value without replecing $var.
```
## Update

The `Update` to loads data from env-file into environment with replacing existing values. During loading replaces `${var}` or `$var` in the string based on the data in the environment.

Returns an error in case of failure.

### Examples:

Suppose that the some value was set into environment as:

```
$ export KEY_0=VALUE_X
```

And there is `.env` file with data:

```
LAST_ID=002
KEY_0=VALUE_000
KEY_1=VALUE_001
KEY_2=VALUE_${LAST_ID}
```

Make code to loads custom values into environment:

```
// Environment:
// KEY_0=VALUE_X

// Load values with replacement.
err := env.Update(".env")
if err != nil {
    // something went wrong
}

// Environment:
// KEY_0=VALUE_000  // data has been updated;
// KEY_1=VALUE_001  // add new value;
// KEY_2=VALUE_002  // add new value and replaced ${LAST_ID}
                    // to the plain text.
```

## UpdateSafe

The `UpdateSafe` to loads data from env-file into environment with replacing existing values. Ignores the replecing of a `${var}` or `$var` in a string.

Returns an error in case of failure.

### Examples:

Suppose that the some value was set into environment as:

```
$ export KEY_0=VALUE_X
```

And there is `.env` file with data:

```
LAST_ID=002
KEY_0=VALUE_000
KEY_1=VALUE_001
KEY_2=VALUE_${LAST_ID}
```

Make code to loads custom values into environment:

```
// Environment:
// KEY_0=VALUE_X

// Load values with replacement.
err := env.Update(".env")
if err != nil {
    // something went wrong
}

// Environment:
// KEY_0=VALUE_000         // data has been updated;
// KEY_1=VALUE_001         // add new value;
// KEY_2=VALUE_${LAST_ID}  // add new value without replecing $var.
```
## Exists

The `Exists` returns true if all keys sets in the environment.

### Examples:

Suppose that the some value was set into environment as:

```
$ export KEY_0=VALUE_X
```

And there is `.env` file with data:

```
LAST_ID=002
KEY_0=VALUE_000
KEY_1=VALUE_001
KEY_2=VALUE_${LAST_ID}
```

Make code to check existing for a variable in the environment:

```
env.Exists("KEY_0")          // true
env.Exists("KEY_1")          // false
env.Exists("KEY_0", "KEY_1") // false

// Load values with replacement.
err := env.Update(".env")
if err != nil {
    // something went wrong
}

env.Exists("KEY_1")          // true
env.Exists("KEY_0", "KEY_1") // true
```

## Unmarshal

The `Unmarshal` to parses the environment data and stores the result in the value pointed to by scope. If scope isn't struct, not a pointer or is nil - returns an error.

Unmarshal method  supports the following field's types: `int`, `int8`, `int16`, `int32`, `int64`, `uin`, `uint8`, `uin16`, `uint32`, `uin64`, `float32`, `float64`, `string`, `bool`, `url.URL` and `pointers`, `array` or `slice` from thous types *(i.e. `*int`, ..., `[]int`, ..., `[]bool`, ..., `[2]*url.URL`, etc.)*. The nested structures will be processed recursively.

For other filed's types (like `chan` or `map` ...) will be returned an error.

If the structure implements Unmarshaler interface - the custom UnmarshalENV method will be called.

Structure fields may have a `env` tag as `env:"KEY[,SEP]"` where:

   - KEY - matches the name of the key in the environment;
   - SEP - optional argument, sets the separator for lists (default: space).

Suppose that the some values was set into environment as:

```
$ export HOST="0.0.0.0"
$ export PORT=8080
$ export ALLOWED_HOSTS=localhost:127.0.0.1
$ export SECRET_KEY=AgBsdjONL53IKa33LM9SNROvD3hZXfoz
```

Structure example:

```
// Config structure for containing values from the environment.
// P.s. No need to describe all the keys that are in the environment,
// for example, we ignore the SECRET_KEY key.
type Config struct {
    Host         string   `env:"HOST"`
    Port         int      `env:"PORT"`
    AllowedHosts []string `env:"ALLOWED_HOSTS,:"`
}
```

Unmarshal data from the environment into Config struct.

```
var config = Config{}

err := env.Unmarshal(&config)
if err != nil {
    // something went wrong
}

// Object config:
// config.Host         // "0.0.0.0"
// config.Port         // 8080
// config.AllowedHosts // []string{"localhost", "127.0.0.1"}
```

If the structure will has custom UnmarshalENV - it will be called:

```
// UnmarshalENV it's custom method for unmarshalling.
func (c *Config) UnmarshalENV() error {
    c.Host = "192.168.0.1"
    c.Port = 80
    c.AllowedHosts = []string{"192.168.0.1"}

    return nil
}
...
// The result will be the data that sets in the custom
// unmarshalling method.
config.Host         // "192.168.0.1"
config.Port         // 80
config.AllowedHosts // []string{"192.168.0.1"}
```

## Marshal

The `Marshal` converts the structure in to key/value and put it into environment with update old values. The first return value returns a map of the data that was correct set into environment. The second - error or nil.

Marshal methods  supports the following field's types: `int`, `int8`, `int16`, `int32`, `int64`, `uin`, `uint8`, `uin16`, `uint32`, `uin64`, `float32`, `float64`, `string`, `bool`, `url.URL` and `pointers`, `array` or `slice` from thous types *(i.e. `*int`, ..., `[]int`, ..., `[]bool`, ..., `[2]*url.URL`, etc.)*. The nested structures will be processed recursively.

For other filed's types (like `chan` or `map` ...) will be returned an error.

If the structure implements Marshaler interface - the custom MarshalENV method - will be called.

Structure fields may have a `env` tag as `env:"KEY[,SEP]"` where:

   - KEY - matches the name of the key in the environment;
   - SEP - optional argument, sets the separator for lists (default: space).

Structure example:

```
// Config structure.
type Config struct {
    Host         string   `env:"HOST"`
    Port         int      `env:"PORT"`
    AllowedHosts []string `env:"ALLOWED_HOSTS,:"`
}
```

Marshal data into environment from the Config.

```
// It can be a structure or a pointer to the struct.
var config = &Config{
    "localhost",
    8080,
    []string{"localhost", "127.0.0.1"},
}

// Returns:
// map[ALLOWED_HOSTS:localhost:127.0.0.1 HOST:localhost PORT:8080], nil
_, err := env.Marshal(config)
if err != nil {
    // something went wrong
}

env.Get("HOST")          // "localhost"
env.Get("PORT")          // "8080"
env.Get("ALLOWED_HOSTS") // "localhost:127.0.0.1"
```

If object has MarshalENV and isn't a nil pointer - will be calls it method to scope convertation.

```
// MarshalENV it's custom method for marshalling.
func (c *Config) MarshalENV() ([]string, error ){
    os.Setenv("HOST", "192.168.0.1")
    os.Setenv("PORT", "80")
    os.Setenv("ALLOWED_HOSTS", "192.168.0.1")

    return []string{}, nil
}
...
// The result will be the data that sets in the custom
// unmarshalling method.
env.Get("HOST")          // "192.168.0.1"
env.Get("PORT")          // "80"
env.Get("ALLOWED_HOSTS") // "192.168.0.1"
```

# Synonyms

There are synonyms for the  `os.*env` functions.

## Get

The `Get` retrieves the value of the environment variable named by the key. It returns the value, which will be empty if the variable is not present.

*P.s. It is synonym for the `os.Getenv`*.

## Set

The `Set` sets the value of the environment variable named by the key. It returns an error, if any.

*P.s. It is synonym for the `os.Setenv`*.

## Unset

The `Unset` unsets a single environment variable.

*P.s. It is synonym for the `os.Unsetenv`*.

## Clear

The `Clear` deletes all environment variables.

*P.s. It is synonym for the `os.Clearenv`*.

## Environ

The `Environ` returns a copy of strings representing the environment, in the form "key=value".

*P.s. It is synonym for the `os.Environ`*.

## Expand

The `Expand` replaces `${var}` or `$var` in the string according to the values of the current environment variables. References to undefined variables are replaced by the empty string.

*P.s. It is synonym for the `os.ExpandEnv`*.
