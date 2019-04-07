package state

import (
	"errors"
)

// defines various errors used
var (
	ErrPresenceStateNotFound    = errors.New("presence not found in state cache")
	ErrRoleStateNotFound        = errors.New("role not found in state cache")
	ErrEmojiStateNotFound       = errors.New("emoji not found in state cache")
	ErrBotForGuildStateNotFound = errors.New("bot for guild not found in state cache")

	ErrTargetWrongServer = errors.New("target is on wrong server")
	ErrTargetWrongType   = errors.New("target is of wrong type")

	ErrUserNotFound    = errors.New("user not found")
	ErrChannelNotFound = errors.New("channel not found")
	ErrRoleNotFound    = errors.New("role not found")
)
