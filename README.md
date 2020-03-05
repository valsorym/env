# env

The `env` it's simple lib for manage environment's variables, parse and convert it into Go-structures.

*P.s. Environment variables are a universal mechanism for conveying configuration information to Unix programs. They are part of the environment in which a process runs.*

The env package allows to control environment's variables: 

* set new variables;
* delete variables;
* get value by key;
* load data from the env-file into environment;
* store environment values into Go- structure.

## Quick example

Let's imagine the task. There is a web-project that is develop and tests on the local computer and runs on the production server. On the production server some settings (for example host and port) enforced seated in an environment but on the local computer data must be loaded from the file.

The configuration file `config.env` with env-variables look likes:

```
HOST=0.0.0.0
PORT=8080
ALLOWED_HOSTS="localhost:127.0.0.1"

# In configurations and in the environment there may be
# a many of variables that willn't be parsed.
SECRET_KEY="N0XRABLZ5ZZY6dCNgD7pLjTIlx8v4G5d"
```

Make test project `main.go`. We need to load the missing configurations from the file into the environment. Convert data from the environment into Go-structure. Use this data to start the server.

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

// Addr returns the server's address.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Home it's example of the homepage handler.
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	// Create a pointer to the configuration object.
	var config = &Config{}

	// Load configurations in an environment, loads from the env-file.
	// P.s. We use the Load (not Update) method for as not to overwrite
	// the data in the environment because on the production server this
	// data can be set forcibly.
	env.Load("./config.env")

	// Parses the environment data and stores the result in the value
	// pointed to by config.
	env.Unmarshal(config)

	// Routing.
	http.HandleFunc("/", Home)

	// Run.
	log.Printf("Server started on %s\n", config.Addr())
	log.Printf("Allowed hosts: %v\n", config.AllowedHosts)
	http.ListenAndServe(config.Addr(), nil)
}
```

## Types

Parser supports the following field types: `int`, `int8`, `int16`, `int32`, `int64`, `uin`, `uint8`, `uin16`, `uint32`, `uin64`, `float32`, `float64`, `string`, `bool` and `array` or `slice` <ins>from thous types</ins> *(i.e. `[]int`, `[]int8`, ..., `[]bool`)*. For other filed's types will be returned an error - **`env.TypeError`**.

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

The `Load` to loads data from env-file into environment without replacing
existing values. During loading replaces `${var}` or `$var` in the string
based on the data in the environment.

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
// The method prints environment variables starting with 'KEY_'.
echo := func() {
    for _, item := range env.Environ() {
        if strings.HasPrefix(item, "KEY_") {
            fmt.Println(item)
        }
    }
}

// Printed only:
//  KEY_0=VALUE_X
echo()

// Load values without replacement.
err := env.Load(".env")
if err != nil {
    // something went wrong
}

// Printed all variables:
//  KEY_0=VALUE_X    // not replaced by VALUE_001;
//  KEY_1=VALUE_001  // add new value;
//  KEY_2=VALUE_002  // add new value and replaced ${LAST_ID}
                     // to the plain text.
echo()
```
## LoadSafe

The `LoadSafe` to loads data from env-file into environment without replacing
existing values. Ignores the replacing of a `${var}` or `$var` in a string.

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
// The method prints environment variables starting with 'KEY_'.
echo := func() {
    for _, item := range env.Environ() {
        if strings.HasPrefix(item, "KEY_") {
            fmt.Println(item)
        }
    }
}

// Printed only:
//  KEY_0=VALUE_X
echo()

// Load values without replacement.
err := env.LoadSafe(".env")
if err != nil {
    // something went wrong
}

// Printed all variables:
//  KEY_0=VALUE_X            // not replaced by VALUE_001;
//  KEY_1=VALUE_001         // add new value;
//  KEY_2=VALUE_${LAST_ID}  // add new value without replecing $var.
echo()
```
## Update

The `Update` to loads data from env-file into environment with replacing
existing values. During loading replaces `${var}` or `$var` in the string
based on the data in the environment.

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
// The method prints environment variables starting with 'KEY_'.
echo := func() {
    for _, item := range env.Environ() {
        if strings.HasPrefix(item, "KEY_") {
            fmt.Println(item)
        }
    }
}

// Printed only:
//  KEY_0=VALUE_X
echo()

// Load values with replacement.
err := env.Update(".env")
if err != nil {
    // something went wrong
}

// Printed all variables:
//  KEY_0=VALUE_000  // data has been updated;
//  KEY_1=VALUE_001  // add new value;
//  KEY_2=VALUE_002  // add new value and replaced ${LAST_ID}
                     // to the plain text.
echo()
```

## UpdateSafe

The `UpdateSafe` to loads data from env-file into environment with replacing
existing values. Ignores the replecing of a `${var}` or `$var` in a string.

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
// The method prints environment variables starting with 'KEY_'.
echo := func() {
    for _, item := range env.Environ() {
        if strings.HasPrefix(item, "KEY_") {
            fmt.Println(item)
        }
    }
}

// Printed only:
//  KEY_0=VALUE_X
echo()

// Load values with replacement.
err := env.Update(".env")
if err != nil {
    // something went wrong
}

// Printed all variables:
//  KEY_0=VALUE_000         // data has been updated;
//  KEY_1=VALUE_001         // add new value;
//  KEY_2=VALUE_${LAST_ID}  // add new value without replecing $var.
echo()
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

The `Unmarshal` to parses the environment data and stores the result in the
value pointed to by scope. If scope isn't struct, not a pointer or is nil -
returns an error.

Supports the following field types: `int`, `int8`, `int16`, `int32`, `int64`, `uin`,
`uint8`, `uin16`, `uint32`, `uin64`, `float32`, `float64`, `string`, `bool` and `array` or `slice`
from thous types. For other filed's types will be returned an error.

If the structure implements Unmarshaller interface - the custom UnmarshalENV
method will be called.

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
// Important: pointer to initialized structure!
var config = &Config{}

err := env.Unmarshal(config)
if err != nil {
    // something went wrong
}

config.Host         // "0.0.0.0"
config.Port         // 8080
config.AllowedHosts // []string{"localhost", "127.0.0.1"}
```

If the structure will havs custom UnmarshalENV - it will be called:

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

The `Marshal` converts the structure in to key/value and put it into
environment with update old values. The first return value returns a map
of the data that was correct set into environment. The seconden -
error or nil.

Supports the following field types: `int`, `int8`, `int16`, `int32`, `int64`, `uin`,
`uint8`, `uin16`, `uint32`, `uin64`, `float32`, `float64`, `string`, `bool` and `array` or `slice`
from thous types. For other filed's types will be returned an error.

If the structure implements Marshaller interface - the custom MarshalENV
method - will be called.

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
// It can be a structure or a pointer to a structure.
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

If object has MarshalENV and isn't a nil pointer - will be calls it
method to scope convertation.

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

There are synonyms for the  'os. * env' functions.

## Get

The `Get` retrieves the value of the environment variable named by the key.
It returns the value, which will be empty if the variable is not present.

*P.s. It is synonym for the `os.Getenv`*.

## Set

The `Set` sets the value of the environment variable named by the key.
It returns an error, if any.

*P.s. It is synonym for the `os.Setenv`*.

## Unset

The `Unset` unsets a single environment variable.

*P.s. It is synonym for the `os.Unsetenv`*.

## Clear

The `Clear` deletes all environment variables.

*P.s. It is synonym for the `os.Clearenv`*.

## Environ

The `Environ` returns a copy of strings representing the environment,
in the form "key=value".

*P.s. It is synonym for the `os.Environ`*.

## Expand

The `Expand` replaces `${var}` or `$var` in the string according to
the values of the current environment variables. References to undefined
variables are replaced by the empty string.

*P.s. It is synonym for the `os.ExpandEnv`*.
