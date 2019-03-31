package discord

import (
	"encoding/json"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
)

func Invite(
	redisClient *redis.Client,
	session *Session,
	inviteCode string,
) (*discordgo.Invite, error) {
	var err error
	var cached []byte

	if redisClient != nil {
		cached, err = redisClient.Get(inviteCacheKey(inviteCode)).Bytes()
		if err != nil && err != redis.Nil {
			return nil, err
		}
	}

	var invite *discordgo.Invite
	if len(cached) > 0 {
		err = json.Unmarshal(cached, invite)
		if err != nil {
			return nil, err
		}

		if invite != nil && invite.Code == inviteCode {
			return invite, nil
		}
	}

	invite, err = session.Client.InviteWithCounts(inviteCode)
	if err != nil {
		return nil, err
	}

	if redisClient != nil && invite != nil && invite.Code == inviteCode {
		cached, err = json.Marshal(invite)
		if err == nil {
			redisClient.Set(inviteCacheKey(inviteCode), cached, 24*time.Hour)
		}
	}

	return invite, nil
}

func inviteCacheKey(inviteCode string) string {
	return "cacophony:go-kit:discord:invite-" + inviteCode
}
