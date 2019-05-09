package events

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/permissions"
)

// Has returns true if the user has all the required permissions
func (e *Event) Has(
	firstPermission permissions.PermissionInterface,
	permissions ...permissions.PermissionInterface,
) bool {
	if e.UserID == "" {
		return false
	}

	for _, permission := range append(permissions, firstPermission) {
		if permission.Match(e.state, e.UserID, e.ChannelID, e.DM()) {
			continue
		}

		return false
	}

	return true
}

// HasOr returns true if the user has one of the required permissions
func (e *Event) HasOr(
	firstPermission permissions.PermissionInterface,
	permissions ...permissions.PermissionInterface,
) bool {
	if e.UserID == "" {
		return false
	}

	for _, permission := range append(permissions, firstPermission) {
		if !permission.Match(e.state, e.UserID, e.ChannelID, e.dm) {
			continue
		}

		return true
	}

	return false
}

// Require calls the callback if the user has all of the required permissions
func (e *Event) Require(
	callback func(),
	firstPermission permissions.PermissionInterface,
	permissions ...permissions.PermissionInterface,
) {
	if e.Has(firstPermission, permissions...) {
		callback()
		return
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
}

// RequireOr calls the callback if the user has one of the required permissions
func (e *Event) RequireOr(
	callback func(),
	firstPermission permissions.PermissionInterface,
	permissions ...permissions.PermissionInterface,
) {
	if e.HasOr(firstPermission, permissions...) {
		callback()
		return
	}

	var permissionsText string
	for _, permission := range append(permissions, firstPermission) {
		permissionsText += permission.Name() + ", "
	}
	permissionsText = strings.TrimRight(permissionsText, ", ")

	message := "**Looks like you do not have the right permissions to do this.**" +
		"\nYou need one of the following permissions:" +
		"\n```\n" + permissionsText + "\n```"

	e.Respond(message) // nolint: errcheck
}
