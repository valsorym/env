# ENV

The env package implements a set of functions of the environment management.

The env package allows to control environment's variables: 
* set new variables;
* delete variables;
* get value by key;
* load data from the env-file into environment;
* store environment values into go- structure.

## Example

Micro web-server loading configurations from the environment.

Environment configuration `config.env`:

```
# Web-server's configuration.
HOST="0.0.0.0"
PORT=8080

# P.s. Additional settings can be specified in the config file
#      that are not exists in the server's configuration object.
SECRET_KEY="HxLRgxqzqf"
```

Web-server `main.go`:
```
package main

import (
	"fmt"
	"net/http"

	"github.com/goloop/env"
)

// Config it's struct of the server configuration.
type Config struct {
	Host string `env:"HOST"`
	Port int    `env:"PORT"`

	// P.s. It isn't necessary to specify all the keys that
	//      are available in the environment.
}

// Home it's handler of the homepage.
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	// Create a pointer to the configuration object.
	var config = &Config{}

	// Update configurations in an environment, loads from the fiel.
	// P.s. You can skip this step If the necessary data is manually
	//      setuped in the environment.
	env.Update("./config.env")

	// Parses the environment data and stores the result in the value
	// pointed to by config.
	env.Unmarshal(config)

	// Run test web-server.
	http.HandleFunc("/", Home)
	http.ListenAndServe(fmt.Sprintf("%s:%d", config.Host, config.Port), nil)
}
```

## Get

You can use `go get`:

```
$ go get -u github.com/goloop/env
```
or use `git clone`  and add link to sorece code into $GOPATH:

```
$ mkdir -p ~/workspace && cd ~/workspace
$ git clone https://github.com/goloop/env.git
$ cd env && make link
```

## Import

To use the env package, you must import it as:

```
import "github.com/goloop/env"

```

## Quick start

For example. We have a configuration env-file: `~/workspace/config.env` like:
```
HOST="localhost"
PORT=8080
```
To load these parameters into the environment we need to do:

```
package main
import (
    "fmt"
    "log"
    "os"

    "github.com/goloop/env"
)

func main() {
    // Load env-file.
    err := Update("~/workspace/config.env")
    if err != nil {
        log.Fatal(err.Error())
    }

    // Get value.
    url = fmt.Sprintf("https://%s:%s",
        env.Get("HOST"),
        env.Get("PORT"))

    // Do something ... 
}
```

## Functions

### Load

Loads keys without replacing existing ones and make expand.

`func Load(filename string) error`

### LoadSafe

Loads keys without replacing existing ones.

`func LoadSafe(filename string) error`

### Update

Loads keys with replacing existing ones and make expand.

`func Update(filename string) error`

### UpdateSafe

UpdateSafe loads keys with replacing existing ones.

`func UpdateSafe(filename string) error`
