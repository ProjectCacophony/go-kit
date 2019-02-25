package events

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/discord"
)

// Require calls the callback if the user has all required permissions
func (e *Event) Require(callback func(), firstPermission int, permissions ...int) {
	if !discord.UserHasPermission(
		e.state, e.UserID, e.ChannelID, firstPermission, permissions...,
	) {
		var permissionsText string
		for _, permission := range append(permissions, firstPermission) {
			permissionsText += discord.PermissionName(permission) + ", "
		}
		permissionsText = strings.TrimRight(permissionsText, ", ")

		message := "**Looks like you do not have the right permissions to do this.**" +
			"\nYou need the following permission:" +
			"\n```\n" + permissionsText + "\n```"

		e.Respond(message) // nolint: errcheck
		return
	}

	callback()
}

// Require calls the callback if the user has all required permissions
func (e *Event) RequireOr(callback func(), firstPermission int, permissions ...int) {
	if !discord.UserHasPermissionOr(
		e.state, e.UserID, e.ChannelID, firstPermission, permissions...,
	) {
		var permissionsText string
		for _, permission := range append(permissions, firstPermission) {
			permissionsText += discord.PermissionName(permission) + ", "
		}
		permissionsText = strings.TrimRight(permissionsText, ", ")

		message := "**Looks like you do not have the right permissions to do this.**" +
			"\nYou need the following permission:" +
			"\n```\n" + permissionsText + "\n```"

		e.Respond(message) // nolint: errcheck
		return
	}

	callback()
}
