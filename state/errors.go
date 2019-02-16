package state

import (
	"errors"
)

// ErrStateNotFound will be returned if the item was not found in the shared state
var ErrStateNotFound = errors.New("shared state cache not found")

// ErrTargetWrongServer will be returned if the target is on the wrong server
var ErrTargetWrongServer = errors.New("target is on wrong server")

// ErrTargetWrongType will be returned if the target is on the wrong server
var ErrTargetWrongType = errors.New("target is on wrong server")
