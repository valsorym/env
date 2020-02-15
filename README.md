# ENV

The env package implements a set of functions of the environment management.

The env package allows to control environment's variables: set new
variables, delete variables, get value by variable name, load env-file into
environment, store environment values in a data structure.

## Get

```
    go get -u github.com/goloop/env
```

## Import

```
    import "github.com/goloop/env"

```

## Quick start

```
    package main

    import "github.com/goloop/env"

    func main() {
        err := ReadParseStore("/path/to/file.env", true, false)
        if err != nil {
            panic(err.Error())
        }
        // ...
    }
```
