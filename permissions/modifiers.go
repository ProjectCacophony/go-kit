package permissions

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/interfaces"
)

type NotModifier struct {
	permission interfaces.Permission
}

func Not(permission interfaces.Permission) *NotModifier {
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
	permissions []interfaces.Permission
}

func Or(
	firstPermission interfaces.Permission,
	secondPermission interfaces.Permission,
	additionalPermissions ...interfaces.Permission,
) *OrModifier {
	return &OrModifier{
		permissions: append([]interfaces.Permission{firstPermission, secondPermission}, additionalPermissions...),
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
	permissions []interfaces.Permission
}

func And(
	firstPermission interfaces.Permission,
	secondPermission interfaces.Permission,
	additionalPermissions ...interfaces.Permission,
) *AndModifier {
	return &AndModifier{
		permissions: append([]interfaces.Permission{firstPermission, secondPermission}, additionalPermissions...),
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
