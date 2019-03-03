package permissions

import (
	"gitlab.com/Cacophony/go-kit/state"
)

type PermissionInterface interface {
	Name() string
	Match(
		state *state.State, botOwnerIDs []string,
		userID, channelID string, dm bool,
	) bool
}
