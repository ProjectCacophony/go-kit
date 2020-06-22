package paginator

import (
	"bytes"
	"errors"
	"fmt"
	"math"

	"github.com/bwmarrin/discordgo"
)

// CreatePagedMessage creates the paged messages
func (p *Paginator) FieldsPaginator(
	botID string,
	channelID string,
	userID string,
	embed *discordgo.MessageEmbed,
	fieldsPerPage int,
	dm bool,
) error {
	if dm {
		channelID = userID
	}

	// if there aren't multiple fields to be paged through,
	// or if the amount of fields is less than the requested fields per page
	// just send a normal embed
	if len(embed.Fields) < 2 || len(embed.Fields) <= fieldsPerPage {
		_, err := p.sendComplex(
			botID, channelID, &discordgo.MessageSend{
				Embed: embed,
			},
			dm,
		)
		return err
	}

	// fields per page can not be less than 1
	if fieldsPerPage < 1 {
		return errors.New("FieldsPerPage can not be less than 1")
	}

	// create paged message
	pagedMessage := &PagedEmbedMessage{
		FullEmbed:       embed,
		ChannelID:       channelID,
		CurrentPage:     1,
		FieldsPerPage:   fieldsPerPage,
		TotalNumOfPages: int(math.Ceil(float64(len(embed.Fields)) / float64(fieldsPerPage))),
		UserID:          userID,
		Type:            FieldType,
		DM:              dm,
		BotID:           botID,
	}

	err := p.setupAndSendFirstMessage(pagedMessage)
	if err != nil {
		return err
	}

	err = setPagedMessage(p.redis, pagedMessage.MessageID, pagedMessage)
	return err
}

// ImagePaginator creates the paged image messages
func (p *Paginator) ImagePaginator(
	botID string,
	channelID string,
	userID string,
	embed *discordgo.MessageEmbed,
	files []*File,
	dm bool,
) error {
	if embed == nil || len(files) == 0 {
		return nil
	}

	if dm {
		channelID = userID
	}

	if embed.Image == nil {
		embed.Image = &discordgo.MessageEmbedImage{}
	}
	embed.Image.URL = fmt.Sprintf("attachment://%s", files[0].Name)

	if len(files) < 2 {
		var _, err = p.sendComplex(botID, channelID, &discordgo.MessageSend{
			Embed: embed,
			Files: []*discordgo.File{
				{
					Name:        files[0].Name,
					ContentType: files[0].ContentType,
					Reader:      bytes.NewReader(files[0].Data),
				},
			},
		},
			dm,
		)
		return err
	}

	// create paged message
	pagedMessage := &PagedEmbedMessage{
		FullEmbed:       embed,
		ChannelID:       channelID,
		CurrentPage:     1,
		FieldsPerPage:   20,
		TotalNumOfPages: len(files),
		Files:           files,
		UserID:          userID,
		Type:            ImageType,
		DM:              dm,
		BotID:           botID,
	}

	err := p.setupAndSendFirstMessage(pagedMessage)
	if err != nil {
		return err
	}

	err = setPagedMessage(p.redis, pagedMessage.MessageID, pagedMessage)
	return err
}

func (p *Paginator) EmbedPaginator(
	botID string,
	channelID string,
	userID string,
	dm bool,
	embeds ...*discordgo.MessageEmbed,
) error {
	if len(embeds) == 0 {
		return nil
	}

	if dm {
		channelID = userID
	}

	if len(embeds) < 2 {
		_, err := p.sendComplex(botID, channelID, &discordgo.MessageSend{
			Embed: embeds[0],
		},
			dm,
		)
		return err
	}

	pagedMessage := &PagedEmbedMessage{
		FullEmbed:       embeds[0],
		ChannelID:       channelID,
		CurrentPage:     1,
		TotalNumOfPages: len(embeds),
		UserID:          userID,
		Embeds:          embeds,
		Type:            EmbedType,
		DM:              dm,
		BotID:           botID,
	}

	err := p.setupAndSendFirstMessage(pagedMessage)
	if err != nil {
		return err
	}

	err = setPagedMessage(p.redis, pagedMessage.MessageID, pagedMessage)
	return err
}
