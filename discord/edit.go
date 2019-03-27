package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"gitlab.com/Cacophony/go-kit/discord/emoji"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/localisation"
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

	var err error

	edit = TranslateMessageEdit(
		localisations,
		edit,
		values...,
	)

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

func TranslateMessageEdit(
	localisations []interfaces.Localisation,
	edit *discordgo.MessageEdit,
	values ...interface{},
) *discordgo.MessageEdit {
	if edit.Content != nil {
		newContent := localisation.Translate(
			localisations,
			*edit.Content,
			values...,
		)
		newContent = EscapeDiscordContent(newContent)
		newContent = emoji.Replace(newContent)
		edit.Content = &newContent
	}

	if edit.Embed != nil {
		edit.Embed = TranslateEmbed(
			localisations,
			edit.Embed,
			values...,
		)
	}

	return edit
}
