package state

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var guildBansLock sync.Mutex

func (s *State) initGuildBans(session *discordgo.Session, guildID string) (err error) {
	guildBansLock.Lock()
	defer guildBansLock.Unlock()

	l := zap.L().With(zap.String("guild_id", guildID))

	// check if bot is allowed to see bans
	apermissions, err := s.UserPermissions(session.State.User.ID, guildID)
	if err != nil {
		return err
	}
	if apermissions&discordgo.PermissionBanMembers != discordgo.PermissionBanMembers {
		// reset ban list if not allowed
		l.Debug("resetting guild bans, because missing permissions")
		err = deleteStateObject(s.client, guildBanIDsSetKey(guildID))
		if err != nil {
			return err
		}
		return deleteStateObject(s.client, guildBanInitializedKey(guildID))
	}

	exists, err := s.client.Exists(guildBanInitializedKey(guildID)).Result()
	if exists > 0 {
		l.Debug("skipping initializing bans, already initialized")
		return
	}

	// reset guild bans
	err = deleteStateObject(s.client, guildBanIDsSetKey(guildID))
	if err != nil {
		return err
	}

	// cache new guild bans
	bans, err := session.GuildBans(guildID)
	if err != nil {
		return err
	}

	newSet := make([]string, 0, len(bans))
	for _, ban := range bans {
		newSet = append(newSet, ban.User.ID)
	}

	l.Debug("initializing bans", zap.Int("amount", len(newSet)))
	err = addToStateSet(s.client, guildBanIDsSetKey(guildID), newSet...)
	if err != nil {
		return err
	}

	return updateStateObject(s.client, guildBanInitializedKey(guildID), true)
}
