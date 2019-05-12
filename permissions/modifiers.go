package permissions

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/interfaces"
)

type NotModifier struct {
	permission interfaces.PermissionInterface
}

func Not(permission interfaces.PermissionInterface) *NotModifier {
	return &NotModifier{
		permission: permission,
	}
}

func (p *NotModifier) Name() string {
	return "Not " + p.permission.Name()
}

func (p *NotModifier) Match(
	state interfaces.State,
	userID string,
	channelID string,
	dm bool,
) bool {
	return !p.permission.Match(state, userID, channelID, dm)
}

type OrModifier struct {
	permissions []interfaces.PermissionInterface
}

func Or(
	firstPermission interfaces.PermissionInterface,
	secondPermission interfaces.PermissionInterface,
	additionalPermissions ...interfaces.PermissionInterface,
) *OrModifier {
	return &OrModifier{
		permissions: append([]interfaces.PermissionInterface{firstPermission, secondPermission}, additionalPermissions...),
	}
}

func (p *OrModifier) Name() string {
	names := make([]string, len(p.permissions))
	for i, permission := range p.permissions {
		names[i] = permission.Name()
	}

	return "(" + strings.Join(names, " or ") + ")"
}

func (p *OrModifier) Match(
	state interfaces.State,
	userID string,
	channelID string,
	dm bool,
) bool {
	for _, permission := range p.permissions {
		if !permission.Match(state, userID, channelID, dm) {
			continue
		}

		return true
	}

	return false
}

type AndModifier struct {
	permissions []interfaces.PermissionInterface
}

func And(
	firstPermission interfaces.PermissionInterface,
	secondPermission interfaces.PermissionInterface,
	additionalPermissions ...interfaces.PermissionInterface,
) *AndModifier {
	return &AndModifier{
		permissions: append([]interfaces.PermissionInterface{firstPermission, secondPermission}, additionalPermissions...),
	}
}

func (p *AndModifier) Name() string {
	names := make([]string, len(p.permissions))
	for i, permission := range p.permissions {
		names[i] = permission.Name()
	}

	return "(" + strings.Join(names, " and ") + ")"
}

func (p *AndModifier) Match(
	state interfaces.State,
	userID string,
	channelID string,
	dm bool,
) bool {
	for _, permission := range p.permissions {
		if permission.Match(state, userID, channelID, dm) {
			continue
		}

		return false
	}

	return true
}
