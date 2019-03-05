package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"gitlab.com/Cacophony/go-kit/interfaces"
)

func EditComplexWithVars(
	redis *redis.Client,
	session *Session,
	localisations []interfaces.Localisation,
	edit *discordgo.MessageEdit,
	dm bool,
	values ...interface{},
) (*discordgo.Message, error) {
	if edit == nil {
		return nil, nil
	}

	// TODO: support translating

	var err error

	if edit.Embed != nil {
		edit.Embed = TrimEmbed(edit.Embed)
	}

	if dm {
		if redis == nil {
			return nil, errors.New("sending DMs requires redis")
		}

		edit.Channel, err = DMChannel(redis, session, edit.Channel)
		if err != nil {
			return nil, err
		}
	}

	return session.Client.ChannelMessageEditComplex(edit)

}
