package events

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/interfaces"
	"go.opentelemetry.io/otel/api/global"
)

// Has returns true if the user has all the required permissions
func (e *Event) Has(
	firstPermission interfaces.Permission,
	permissions ...interfaces.Permission,
) bool {
	_, span := global.Tracer("cacophony.dev/kit").Start(e.Context(), "event.Has")
	defer span.End()

	if e.UserID == "" {
		return false
	}

	for _, permission := range append(permissions, firstPermission) {
		if permission.Match(e.state, e.db, e.UserID, e.ChannelID, e.DM(), e.SuperUser) {
			continue
		}

		return false
	}

	return true
}

// HasOr returns true if the user has one of the required permissions
func (e *Event) HasOr(
	firstPermission interfaces.Permission,
	permissions ...interfaces.Permission,
) bool {
	_, span := global.Tracer("cacophony.dev/kit").Start(e.Context(), "event.HasOr")
	defer span.End()

	if e.UserID == "" {
		return false
	}

	for _, permission := range append(permissions, firstPermission) {
		if !permission.Match(e.state, e.db, e.UserID, e.ChannelID, e.dm, e.SuperUser) {
			continue
		}

		return true
	}

	return false
}

// Require calls the callback if the user has all of the required permissions
func (e *Event) Require(
	callback func(),
	firstPermission interfaces.Permission,
	permissions ...interfaces.Permission,
) {
	_, span := global.Tracer("cacophony.dev/kit").Start(e.Context(), "event.Require")
	defer span.End()

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

	e.Respond(message)
}

// RequireOr calls the callback if the user has one of the required permissions
func (e *Event) RequireOr(
	callback func(),
	firstPermission interfaces.Permission,
	permissions ...interfaces.Permission,
) {
	_, span := global.Tracer("cacophony.dev/kit").Start(e.Context(), "event.RequireOr")
	defer span.End()

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

	e.Respond(message)
}
