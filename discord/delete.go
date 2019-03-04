package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

func Delete(
	redis *redis.Client,
	session *discordgo.Session,
	channelID string,
	messageID string,
	dm bool,
) error {
	var err error

	if dm {
		if redis == nil {
			return errors.New("reacting in DMs requires redis")
		}

		channelID, err = DMChannel(redis, session, channelID)
		if err != nil {
			return err
		}
	}

	return session.ChannelMessageDelete(channelID, messageID)
}
