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
		// l.Debug("resetting guild bans, because missing permissions")
		err = deleteStateObject(s.client, guildBanIDsSetKey(guildID))
		if err != nil {
			return err
		}
		return deleteStateObject(s.client, guildBanInitializedKey(guildID))
	}

	exists, err := s.client.Exists(guildBanInitializedKey(guildID)).Result()
	if exists > 0 {
		// l.Debug("skipping initializing bans, already initialized")
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

var webhooksLock sync.Mutex

func (s *State) initGuildWebhooks(session *discordgo.Session, guildID string) (err error) {
	webhooksLock.Lock()
	defer webhooksLock.Unlock()

	l := zap.L().With(zap.String("guild_id", guildID))

	// check if bot is allowed to see bans
	apermissions, err := s.UserPermissions(session.State.User.ID, guildID)
	if err != nil {
		return err
	}
	if apermissions&discordgo.PermissionManageWebhooks != discordgo.PermissionManageWebhooks {
		// reset ban list if not allowed
		// l.Debug("resetting guild webhooks, because missing permissions")
		err = deleteStateObject(s.client, guildWebhookIDsSetKey(guildID))
		if err != nil {
			return err
		}
		return deleteStateObject(s.client, guildWebhooksInitializedKey(guildID))
	}

	exists, err := s.client.Exists(guildWebhooksInitializedKey(guildID)).Result()
	if exists > 0 {
		// l.Debug("skipping initializing webhooks, already initialized")
		return
	}

	// reset guild webhooks
	err = deleteStateObject(s.client, guildWebhookIDsSetKey(guildID))
	if err != nil {
		return err
	}

	// cache new guild bans
	webhooks, err := session.GuildWebhooks(guildID)
	if err != nil {
		return err
	}

	l.Debug("saving webhooks", zap.Int("amount", len(webhooks)))
	for _, webhook := range webhooks {
		err = s.webhookAdd(webhook)
		if err != nil {
			return err
		}
	}

	return updateStateObject(s.client, guildWebhooksInitializedKey(guildID), true)
}

var guildInvitesLock sync.Mutex

func (s *State) initGuildInvites(session *discordgo.Session, guildID string) (err error) {
	guildInvitesLock.Lock()
	defer guildInvitesLock.Unlock()

	l := zap.L().With(zap.String("guild_id", guildID))

	// check if bot is allowed to see invites
	apermissions, err := s.UserPermissions(session.State.User.ID, guildID)
	if err != nil {
		return err
	}
	if apermissions&discordgo.PermissionManageServer != discordgo.PermissionManageServer {
		// reset invites list if not allowed
		// l.Debug("resetting guild invites, because missing permissions")
		err = deleteStateObject(s.client, guildInvitesKey(guildID))
		if err != nil {
			return err
		}
		return deleteStateObject(s.client, guildInvitesInitializedKey(guildID))
	}

	exists, err := s.client.Exists(guildInvitesInitializedKey(guildID)).Result()
	if exists > 0 {
		// l.Debug("skipping initializing bans, already initialized")
		return
	}

	// cache new guild invites
	invites, err := session.GuildInvites(guildID)
	if err != nil {
		return err
	}

	l.Debug("initializing invites", zap.Int("amount", len(invites)))
	err = updateStateObject(s.client, guildInvitesKey(guildID), invites)
	if err != nil {
		return err
	}

	return updateStateObject(s.client, guildInvitesInitializedKey(guildID), true)
}

func (s *State) updateGuildInvites(session *discordgo.Session, guildID string) error {
	guildInvitesLock.Lock()
	defer guildInvitesLock.Unlock()

	l := zap.L().With(zap.String("guild_id", guildID))

	// check if bot is allowed to see invites
	apermissions, err := s.UserPermissions(session.State.User.ID, guildID)
	if err != nil {
		return err
	}
	if apermissions&discordgo.PermissionManageServer != discordgo.PermissionManageServer {
		return nil
	}

	// cache new guild invites
	invites, err := session.GuildInvites(guildID)
	if err != nil {
		return err
	}

	l.Debug("updating invites", zap.Int("amount", len(invites)))
	err = updateStateObject(s.client, guildInvitesKey(guildID), invites)
	if err != nil {
		return err
	}

	return updateStateObject(s.client, guildInvitesInitializedKey(guildID), true)
}
