package permissions

import (
	"gitlab.com/Cacophony/go-kit/state"
)

type NotModifier struct {
	permission PermissionInterface
}

func Not(permission PermissionInterface) *NotModifier {
	return &NotModifier{
		permission: permission,
	}
}

func (p *NotModifier) Name() string {
	return "Not " + p.permission.Name()
}

func (p *NotModifier) Match(state *state.State, botOwnerIDs []string, userID, channelID string, dm bool) bool {
	return !p.permission.Match(state, botOwnerIDs, userID, channelID, dm)
}

type OrModifier struct {
	permissions []PermissionInterface
}

func Or(
	firstPermission PermissionInterface,
	secondPermission PermissionInterface,
	additionalPermissions PermissionInterface,
) *OrModifier {
	return &OrModifier{
		permissions: append([]PermissionInterface{firstPermission, secondPermission}, additionalPermissions),
	}
}

func (p *OrModifier) Match(state *state.State, botOwnerIDs []string, userID, channelID string, dm bool) bool {
	for _, permission := range p.permissions {
		if !permission.Match(state, botOwnerIDs, userID, channelID, dm) {
			continue
		}

		return true
	}

	return false
}

type AndModifier struct {
	permissions []PermissionInterface
}

func And(
	firstPermission PermissionInterface,
	secondPermission PermissionInterface,
	additionalPermissions PermissionInterface,
) *AndModifier {
	return &AndModifier{
		permissions: append([]PermissionInterface{firstPermission, secondPermission}, additionalPermissions),
	}
}

func (p *AndModifier) Match(state *state.State, botOwnerIDs []string, userID, channelID string, dm bool) bool {
	for _, permission := range p.permissions {
		if permission.Match(state, botOwnerIDs, userID, channelID, dm) {
			continue
		}

		return false
	}

	return true
}
