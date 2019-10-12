package discord

import (
	"errors"
	"fmt"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/go-redis/redis"
)

var (
	dmChannelExpirySuccessful = time.Hour * 24 * 7
	dmChannelExpiryError      = time.Hour * 1
)

func DMChannel(
	redisClient *redis.Client,
	session *Session,
	userID string,
) (string, error) {

	if redisClient == nil {
		return "", errors.New("sending DMs requires redis")
	}

	key := "cacophony:discord:dm-channel:" + session.BotID + ":" + userID

	res, err := redisClient.Get(key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	if err == nil {
		if res == "" {
			return "", errors.New("channel not set, DMs might be disabled")
		}
		return res, nil
	}

	channel, err := session.Client.UserChannelCreate(userID)
	if err != nil {
		fmt.Printf("setting %s to empty\n", key) // TODO: remove me
		if err := redisClient.Set(key, "", dmChannelExpiryError).Err(); err != nil && raven.DefaultClient != nil {
			raven.CaptureError(err, map[string]string{"key": key})
		}

		return "", err
	}

	fmt.Printf("setting %s to %s\n", key, channel.ID) // TODO: remove me
	if err := redisClient.Set(key, channel.ID, dmChannelExpirySuccessful).Err(); err != nil && raven.DefaultClient != nil {
		raven.CaptureError(err, map[string]string{"key": key})
	}

	return channel.ID, nil
}
