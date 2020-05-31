package permissions

import (
	"strings"

	"github.com/jinzhu/gorm"
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
	db *gorm.DB,
	userID string,
	channelID string,
	dm bool,
	superUser bool,
) bool {
	return !p.permission.Match(state, db, userID, channelID, dm, superUser)
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
	db *gorm.DB,
	userID string,
	channelID string,
	dm bool,
	superUser bool,
) bool {
	for _, permission := range p.permissions {
		if !permission.Match(state, db, userID, channelID, dm, superUser) {
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

func (p *AndModifier) Match(state interfaces.State, db *gorm.DB, userID string, channelID string, dm bool, superUser bool) bool {
	for _, permission := range p.permissions {
		if permission.Match(state, db, userID, channelID, dm, false) {
			continue
		}

		return false
	}

	return true
}
