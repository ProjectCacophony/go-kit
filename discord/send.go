package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
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
	values ...interface{},
) *discordgo.MessageSend {
	send.Content = localisation.Translate(
		localisations,
		send.Content,
		values...,
	)

	send.Content = EscapeDiscordContent(send.Content)
	send.Content = emoji.Replace(send.Content)

	if send.Embed != nil {
		send.Embed = TranslateEmbed(
			localisations,
			send.Embed,
			values...,
		)
	}

	return send
}

// TODO: replace emoji
func TranslateEmbed(
	localisations []interfaces.Localisation,
	embed *discordgo.MessageEmbed,
	values ...interface{},
) *discordgo.MessageEmbed {
	embed.URL = localisation.Translate(
		localisations,
		embed.URL,
		values...,
	)
	embed.Title = localisation.Translate(
		localisations,
		embed.Title,
		values...,
	)
	embed.Title = emoji.Replace(embed.Title)
	embed.Description = localisation.Translate(
		localisations,
		embed.Description,
		values...,
	)
	embed.Description = emoji.Replace(embed.Description)
	if embed.Footer != nil {
		embed.Footer.IconURL = localisation.Translate(
			localisations,
			embed.Footer.IconURL,
			values...,
		)
		embed.Footer.Text = localisation.Translate(
			localisations,
			embed.Footer.Text,
			values...,
		)
		embed.Footer.Text = emoji.Replace(embed.Footer.Text)
	}
	if embed.Image != nil {
		embed.Image.URL = localisation.Translate(
			localisations,
			embed.Image.URL,
			values...,
		)
	}
	if embed.Thumbnail != nil {
		embed.Thumbnail.URL = localisation.Translate(
			localisations,
			embed.Thumbnail.URL,
			values...,
		)
	}
	if embed.Video != nil {
		embed.Video.URL = localisation.Translate(
			localisations,
			embed.Video.URL,
			values...,
		)
	}
	if embed.Provider != nil {
		embed.Provider.URL = localisation.Translate(
			localisations,
			embed.Provider.URL,
			values...,
		)
		embed.Provider.Name = localisation.Translate(
			localisations,
			embed.Provider.Name,
			values...,
		)
		embed.Provider.Name = emoji.Replace(embed.Provider.Name)
	}
	if embed.Author != nil {
		embed.Author.URL = localisation.Translate(
			localisations,
			embed.Author.URL,
			values...,
		)
		embed.Author.IconURL = localisation.Translate(
			localisations,
			embed.Author.IconURL,
			values...,
		)
		embed.Author.Name = localisation.Translate(
			localisations,
			embed.Author.Name,
			values...,
		)
		embed.Author.Name = emoji.Replace(embed.Author.Name)
	}
	for i := range embed.Fields {
		embed.Fields[i].Name = localisation.Translate(
			localisations,
			embed.Fields[i].Name,
			values...,
		)
		embed.Fields[i].Name = emoji.Replace(embed.Fields[i].Name)
		embed.Fields[i].Value = localisation.Translate(
			localisations,
			embed.Fields[i].Value,
			values...,
		)
		embed.Fields[i].Value = emoji.Replace(embed.Fields[i].Value)
	}

	return embed
}
