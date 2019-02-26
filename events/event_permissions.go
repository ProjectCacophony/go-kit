package events

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/permissions"
)

// Require calls the callback if the user has all of the required permissions
func (e *Event) Require(
	callback func(),
	firstPermission permissions.PermissionInterface,
	permissions ...permissions.PermissionInterface,
) {
	for _, permission := range append(permissions, firstPermission) {
		if permission.Match(e.state, e.BotOwnerIDs(), e.UserID, e.ChannelID) {
			continue
		}

		var permissionsText string
		for _, permission := range append(permissions, firstPermission) {
			permissionsText += permission.Name() + ", "
		}
		permissionsText = strings.TrimRight(permissionsText, ", ")

		message := "**Looks like you do not have the right permissions to do this.**" +
			"\nYou need all of the following permissions:" +
			"\n```\n" + permissionsText + "\n```"

		e.Respond(message) // nolint: errcheck
		return
	}

	callback()
}

// RequireOr calls the callback if the user has one of the required permissions
func (e *Event) RequireOr(
	callback func(),
	firstPermission permissions.PermissionInterface,
	permissions ...permissions.PermissionInterface,
) {
	var matched bool

	for _, permission := range append(permissions, firstPermission) {
		if !permission.Match(e.state, e.BotOwnerIDs(), e.UserID, e.ChannelID) {
			continue
		}

		matched = true
		break
	}

	if !matched {

		var permissionsText string
		for _, permission := range append(permissions, firstPermission) {
			permissionsText += permission.Name() + ", "
		}
		permissionsText = strings.TrimRight(permissionsText, ", ")

		message := "**Looks like you do not have the right permissions to do this.**" +
			"\nYou need one of the following permissions:" +
			"\n```\n" + permissionsText + "\n```"

		e.Respond(message) // nolint: errcheck
		return
	}

	callback()
}
