package permissions

import (
	"gitlab.com/Cacophony/go-kit/state"
)

type PermissionInterface interface {
	Name() string
	Match(
		state *state.State,
		userID string,
		channelID string,
		dm bool,
	) bool
}
