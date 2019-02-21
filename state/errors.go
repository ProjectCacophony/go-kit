package state

import (
	"errors"
)

// defines various errors used
var (
	ErrStateNotFound     = errors.New("shared state cache not found")
	ErrTargetWrongServer = errors.New("target is on wrong server")
	ErrTargetWrongType   = errors.New("target is on wrong server")

	ErrUserNotFound    = errors.New("user not found")
	ErrChannelNotFound = errors.New("channel not found")
	ErrRoleNotFound    = errors.New("role not found")
)
