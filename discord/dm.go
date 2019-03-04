package discord

import (
	"time"

	"github.com/pkg/errors"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
)

// nolint: gochecknoglobals
var (
	dmChannelExpirySuccessful = time.Hour * 24 * 7
	dmChannelExpiryError      = time.Hour * 1
)

func DMChannel(
	redisClient *redis.Client,
	session *discordgo.Session,
	userID string,
) (string, error) {
	key := "cacophony:discord:dm-channel:" + userID

	res, err := redisClient.Get(key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	if err == nil {
		if res == "" {
			return "", errors.New("channel not set, DMs disabled?")
		}
		return res, nil
	}

	channel, err := session.UserChannelCreate(userID)
	if err != nil {
		redisClient.Set(key, "", dmChannelExpiryError)

		return "", err
	}

	redisClient.Set(key, channel.ID, dmChannelExpirySuccessful)

	return channel.ID, nil
}
