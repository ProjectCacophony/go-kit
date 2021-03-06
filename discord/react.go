package discord

import (
	"math/rand"
	"time"

	"gitlab.com/Cacophony/go-kit/discord/emoji"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

func React(
	redis *redis.Client,
	session *Session,
	channelID string,
	messageID string,
	dm bool,
	emojiID string, emojiIDs ...string,
) error {
	var err error

	if len(emojiIDs) > 0 {
		emojiID = append(emojiIDs, emojiID)[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(emojiIDs)+1)]
	}

	if dm {
		if redis == nil {
			return errors.New("reacting in DMs requires redis")
		}

		channelID, err = DMChannel(redis, session, channelID)
		if err != nil {
			return err
		}
	}

	return session.Client.MessageReactionAdd(channelID, messageID, emoji.GetWithout(emojiID))
}

func RemoveReact(
	redis *redis.Client,
	session *Session,
	channelID string,
	messageID string,
	userID string,
	dm bool,
	emojiID string,
) error {
	var err error

	if dm {
		if dm {
			return errors.New("cannot perform action in DMs")
		}

		if redis == nil {
			return errors.New("reacting in DMs requires redis")
		}

		channelID, err = DMChannel(redis, session, channelID)
		if err != nil {
			return err
		}
	}

	return session.Client.MessageReactionRemove(channelID, messageID, emoji.GetWithout(emojiID), userID)
}
