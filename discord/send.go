package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"gitlab.com/Cacophony/go-kit/discord/emoji"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/localisation"
)

// SendComplexWithVars sends a message to a given channelID using a given discord session
// it performs various on the actions:
// - pagifies the message content
// - trims the embed
// - escapes @ everyone, and @ here
func SendComplexWithVars(
	redis *redis.Client,
	session *Session,
	localisations []interfaces.Localisation,
	channelID string,
	send *discordgo.MessageSend,
	dm bool,
	values ...interface{},
) ([]*discordgo.Message, error) {
	send = TranslateMessageSend(
		localisations,
		send,
		values...,
	)

	var err error
	var message *discordgo.Message
	var messages []*discordgo.Message

	if send.Embed != nil {
		send.Embed = TrimEmbed(send.Embed)
	}

	if dm {
		if redis == nil {
			return nil, errors.New("sending DMs requires redis")
		}

		channelID, err = DMChannel(redis, session, channelID)
		if err != nil {
			return nil, err
		}
	}

	pages := Pagify(send.Content)
	if len(pages) > 0 {
		for i, page := range pages {
			if i+1 < len(pages) {
				message, err = session.Client.ChannelMessageSend(channelID, page)
			} else {
				send.Content = page
				message, err = session.Client.ChannelMessageSendComplex(channelID, send)
			}
			if err != nil {
				return messages, err
			}
			messages = append(messages, message)
		}

		return messages, nil
	}

	message, err = session.Client.ChannelMessageSendComplex(channelID, send)
	if err != nil {
		return messages, err
	}
	messages = append(messages, message)

	return messages, nil

}

func TranslateMessageSend(
	localisations []interfaces.Localisation,
	send *discordgo.MessageSend,
	values ...interface{}) *discordgo.MessageSend {
	send.Content = localisation.Translate(
		localisations,
		send.Content,
		values...,
	)

	send.Content = EscapeDiscordContent(send.Content)
	send.Content = emoji.Replace(send.Content) // TODO: replace emoji in other fields
	if send.Embed != nil {
		send.Embed.URL = localisation.Translate(
			localisations,
			send.Embed.URL,
			values...,
		)
		send.Embed.Title = localisation.Translate(
			localisations,
			send.Embed.Title,
			values...,
		)
		send.Embed.Title = emoji.Replace(send.Embed.Title)
		send.Embed.Description = localisation.Translate(
			localisations,
			send.Embed.Description,
			values...,
		)
		send.Embed.Description = emoji.Replace(send.Embed.Description)
		if send.Embed.Footer != nil {
			send.Embed.Footer.IconURL = localisation.Translate(
				localisations,
				send.Embed.Footer.IconURL,
				values...,
			)
			send.Embed.Footer.Text = localisation.Translate(
				localisations,
				send.Embed.Footer.Text,
				values...,
			)
			send.Embed.Footer.Text = emoji.Replace(send.Embed.Footer.Text)
		}
		if send.Embed.Image != nil {
			send.Embed.Image.URL = localisation.Translate(
				localisations,
				send.Embed.Image.URL,
				values...,
			)
		}
		if send.Embed.Thumbnail != nil {
			send.Embed.Thumbnail.URL = localisation.Translate(
				localisations,
				send.Embed.Thumbnail.URL,
				values...,
			)
		}
		if send.Embed.Video != nil {
			send.Embed.Video.URL = localisation.Translate(
				localisations,
				send.Embed.Video.URL,
				values...,
			)
		}
		if send.Embed.Provider != nil {
			send.Embed.Provider.URL = localisation.Translate(
				localisations,
				send.Embed.Provider.URL,
				values...,
			)
			send.Embed.Provider.Name = localisation.Translate(
				localisations,
				send.Embed.Provider.Name,
				values...,
			)
			send.Embed.Provider.Name = emoji.Replace(send.Embed.Provider.Name)
		}
		if send.Embed.Author != nil {
			send.Embed.Author.URL = localisation.Translate(
				localisations,
				send.Embed.Author.URL,
				values...,
			)
			send.Embed.Author.IconURL = localisation.Translate(
				localisations,
				send.Embed.Author.IconURL,
				values...,
			)
			send.Embed.Author.Name = localisation.Translate(
				localisations,
				send.Embed.Author.Name,
				values...,
			)
			send.Embed.Author.Name = emoji.Replace(send.Embed.Author.Name)
		}
		for i := range send.Embed.Fields {
			send.Embed.Fields[i].Name = localisation.Translate(
				localisations,
				send.Embed.Fields[i].Name,
				values...,
			)
			send.Embed.Fields[i].Name = emoji.Replace(send.Embed.Fields[i].Name)
			send.Embed.Fields[i].Value = localisation.Translate(
				localisations,
				send.Embed.Fields[i].Value,
				values...,
			)
			send.Embed.Fields[i].Value = emoji.Replace(send.Embed.Fields[i].Value)
		}
	}

	return send
}
