package paginator

import (
	"bytes"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// setupAndSendFirstMessage
func (p *Paginator) setupAndSendFirstMessage(message *PagedEmbedMessage) error {
	var sentMessage []*discordgo.Message
	var err error

	// copy the embedded message so changes can be made to it
	tempEmbed := &discordgo.MessageEmbed{}
	*tempEmbed = *message.FullEmbed

	// set footer which will hold information about the page it is on
	tempEmbed.Footer = p.getEmbedFooter(message)

	switch message.Type {

	case FieldType:
		// reduce fields to the fields per page
		tempEmbed.Fields = tempEmbed.Fields[:message.FieldsPerPage]

		sentMessage, err = p.sendComplex(
			message.BotID, message.ChannelID, &discordgo.MessageSend{
				Embed: tempEmbed,
			},
			message.DM,
		)
		if err != nil {
			return err
		}

	case ImageType:

		// if fields were sent with image embed, handle those
		if len(message.FullEmbed.Fields) > 0 {

			// get start and end fields based on current page and fields per page
			startField := (message.CurrentPage - 1) * message.FieldsPerPage
			endField := startField + message.FieldsPerPage
			if endField > len(message.FullEmbed.Fields) {
				endField = len(message.FullEmbed.Fields)
			}

			tempEmbed.Fields = tempEmbed.Fields[startField:endField]
		}

		tempEmbed.Image.URL = fmt.Sprintf("attachment://%s", message.Files[message.CurrentPage-1].Name)
		sentMessage, err = p.sendComplex(
			message.BotID, message.ChannelID, &discordgo.MessageSend{
				Embed: tempEmbed,
				Files: []*discordgo.File{{
					Name:        message.Files[message.CurrentPage-1].Name,
					ContentType: message.Files[message.CurrentPage-1].ContentType,
					Reader:      bytes.NewReader(message.Files[message.CurrentPage-1].Data),
				}},
			},
			message.DM,
		)
		if err != nil {
			return err
		}

	case EmbedType:
		tempEmbed.Footer = p.getEmbedFooter(message)

		sentMessage, err = p.sendComplex(
			message.BotID, message.ChannelID, &discordgo.MessageSend{
				Embed: tempEmbed,
			},
			message.DM,
		)
		if err != nil {
			return err
		}

	}

	message.MessageID = sentMessage[0].ID
	return p.addReactionsToMessage(message)
}
