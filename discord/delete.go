package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

func Delete(
	redis *redis.Client,
	session *Session,
	channelID string,
	messageID string,
	dm bool,
) error {
	var err error

	if dm {
		if redis == nil {
			return errors.New("acting in DMs requires redis")
		}

		channelID, err = DMChannel(redis, session, channelID)
		if err != nil {
			return err
		}
	}

	return session.Client.ChannelMessageDelete(channelID, messageID)
}

func DeleteBulk(
	redis *redis.Client,
	session *Session,
	channelID string,
	messageIDs []string,
	dm bool,
) error {
	var err error

	if len(messageIDs) > 100 {
		return errors.New("can only delete up to 100 messages at once")
	}
	if len(messageIDs) == 1 {
		return Delete(redis, session, channelID, messageIDs[0], dm)
	}
	if len(messageIDs) == 0 {
		return nil
	}

	if dm {
		if redis == nil {
			return errors.New("acting in DMs requires redis")
		}

		channelID, err = DMChannel(redis, session, channelID)
		if err != nil {
			return err
		}
	}

	return session.Client.ChannelMessagesBulkDelete(channelID, messageIDs)
}

func DeleteSmart(
	redis *redis.Client,
	session *Session,
	channelID string,
	messages []*discordgo.Message,
	dm bool,
) error {
	var err error

	if len(messages) > 100 {
		return errors.New("can only delete up to 100 messages at once")
	}
	if len(messages) == 1 {
		return Delete(redis, session, channelID, messages[0].ID, dm)
	}
	if len(messages) == 0 {
		return nil
	}

	if dm {
		if redis == nil {
			return errors.New("acting in DMs requires redis")
		}

		channelID, err = DMChannel(redis, session, channelID)
		if err != nil {
			return err
		}
	}

	var messageIDsToBulkDelete []string // nolint: prealloc
	for _, message := range messages {
		// delete one by one, if older than 14 Days
		if time.Since(message.Timestamp) > 24*time.Hour*14 {
			err = Delete(
				redis,
				session,
				channelID,
				message.ID,
				dm,
			)
			if err != nil {
				return err
			}
			continue
		}

		messageIDsToBulkDelete = append(messageIDsToBulkDelete, message.ID)
	}

	return DeleteBulk(
		redis,
		session,
		channelID,
		messageIDsToBulkDelete,
		dm,
	)
}
