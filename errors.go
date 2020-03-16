package env

import "errors"

var (
	KeyError   = errors.New("incorrect key")
	ValueError = errors.New("incorrect value")

	IsNotPointerError     = errors.New("object isn't a pointer")
	IsNotInitializedError = errors.New("object must be initialized")
	IsNotStructError      = errors.New("object isn't a struct")
	TypeError             = errors.New("incorrect type")
)
