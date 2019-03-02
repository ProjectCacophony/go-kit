package paginator

import (
	"bytes"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (p *Paginator) setPage(message *PagedEmbedMessage, page int) error {
	// normalise page
	if page > message.TotalNumOfPages {
		page = 1
	}
	if page <= 0 {
		page = message.TotalNumOfPages
	}

	message.CurrentPage = page

	session, err := p.getSession(message.GuildID)
	if err != nil {
		return err
	}

	tempEmbed := &discordgo.MessageEmbed{}
	*tempEmbed = *message.FullEmbed

	switch message.Type {

	case FieldType:
		// get start and end fields based on current page and fields per page
		startField := (message.CurrentPage - 1) * message.FieldsPerPage
		endField := startField + message.FieldsPerPage
		if endField > len(message.FullEmbed.Fields) {
			endField = len(message.FullEmbed.Fields)
		}

		tempEmbed.Fields = tempEmbed.Fields[startField:endField]
		tempEmbed.Footer = p.getEmbedFooter(message)
		_, err = p.editComplex(message.GuildID, &discordgo.MessageEdit{
			Embed:   tempEmbed,
			ID:      message.MessageID,
			Channel: message.ChannelID,
		})
		if err != nil {
			return err
		}

		err = setPagedMessage(p.redis, message.MessageID, message)
		if err != nil {
			return err
		}

	case ImageType:
		// image embeds can't be edited, need to delete and remake it
		session.ChannelMessageDelete(message.ChannelID, message.MessageID) // nolint: errcheck

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

		// change the url of the embed to point to the new image
		tempEmbed.Image.URL = fmt.Sprintf("attachment://%s", message.Files[message.CurrentPage-1].Name)

		// update footer and send message
		tempEmbed.Footer = p.getEmbedFooter(message)
		sentMessage, err := p.sendComplex(
			message.GuildID, message.ChannelID, &discordgo.MessageSend{
				Embed: tempEmbed,
				Files: []*discordgo.File{{
					Name:        message.Files[message.CurrentPage-1].Name,
					ContentType: message.Files[message.CurrentPage-1].ContentType,
					Reader:      bytes.NewReader(message.Files[message.CurrentPage-1].Data),
				}},
			})
		if err != nil {
			return err
		}

		// update map with new message id since
		originalmsgID := message.MessageID
		message.MessageID = sentMessage[0].ID
		p.addReactionsToMessage(message) // nolint: errcheck
		err = setPagedMessage(p.redis, sentMessage[0].ID, message)
		if err != nil {
			return err
		}
		err = deletePagedMessage(p.redis, originalmsgID)
		if err != nil {
			return err
		}

	case EmbedType:

		if len(message.Embeds) < message.CurrentPage {
			return nil
		}

		tempEmbed = message.Embeds[message.CurrentPage-1]
		tempEmbed.Footer = p.getEmbedFooter(message)

		_, err = p.editComplex(message.GuildID, &discordgo.MessageEdit{
			Embed:   tempEmbed,
			ID:      message.MessageID,
			Channel: message.ChannelID,
		})
		if err != nil {
			return err
		}

		err = setPagedMessage(p.redis, message.MessageID, message)
		if err != nil {
			return err
		}

	}

	return nil
}
