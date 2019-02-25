package discord

import (
	"gitlab.com/Cacophony/go-kit/state"
)

// UserHasPermission returns true if the User has all of th egiven permissions in the given channel
func UserHasPermission(
	state *state.State, userID, channelID string, firstPermission int, permissions ...int,
) bool {
	if userID == "" || channelID == "" {
		return false
	}

	if state == nil {
		return false
	}

	userChannelPermissions, err := state.UserChannelPermissions(userID, channelID)
	if err != nil {
		return false
	}

	for _, permission := range append(permissions, firstPermission) {

		if userChannelPermissions&permission != permission {
			return false
		}
	}

	return true
}

// UserHasPermissionOr returns true if the User has any of the given permissions in the given channel
func UserHasPermissionOr(
	state *state.State, userID, channelID string, firstPermission int, permissions ...int,
) bool {
	if userID == "" || channelID == "" {
		return false
	}

	if state == nil {
		return false
	}

	userChannelPermissions, err := state.UserChannelPermissions(userID, channelID)
	if err != nil {
		return false
	}

	for _, permission := range append(permissions, firstPermission) {

		if userChannelPermissions&permission == permission {
			return true
		}
	}

	return false
}
