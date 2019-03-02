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
	guildID, channelID, userID string, embed *discordgo.MessageEmbed, fieldsPerPage int,
) error {

	// if there aren't multiple fields to be paged through,
	// or if the amount of fields is less than the requested fields per page
	// just send a normal embed
	if len(embed.Fields) < 2 || len(embed.Fields) <= fieldsPerPage {
		_, err := p.sendComplex(guildID, channelID, &discordgo.MessageSend{
			Embed: embed,
		})
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
		GuildID:         guildID,
		CurrentPage:     1,
		FieldsPerPage:   fieldsPerPage,
		TotalNumOfPages: int(math.Ceil(float64(len(embed.Fields)) / float64(fieldsPerPage))),
		UserID:          userID,
		Type:            FieldType,
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
	guildID, channelID, userID string, embed *discordgo.MessageEmbed, files []*File,
) error {
	if embed == nil || len(files) == 0 {
		return nil
	}

	if embed.Image == nil {
		embed.Image = &discordgo.MessageEmbedImage{}
	}
	embed.Image.URL = fmt.Sprintf("attachment://%s", files[0].Name)

	if len(files) < 2 {
		var _, err = p.sendComplex(guildID, channelID, &discordgo.MessageSend{
			Content: "",
			Embed:   embed,
			Tts:     false,
			Files: []*discordgo.File{
				{
					Name:        files[0].Name,
					ContentType: files[0].ContentType,
					Reader:      bytes.NewReader(files[0].Data),
				},
			},
			File: nil,
		})
		return err
	}

	// create paged message
	pagedMessage := &PagedEmbedMessage{
		FullEmbed:       embed,
		ChannelID:       channelID,
		GuildID:         guildID,
		CurrentPage:     1,
		FieldsPerPage:   20,
		TotalNumOfPages: len(files),
		Files:           files,
		UserID:          userID,
		Type:            ImageType,
	}

	err := p.setupAndSendFirstMessage(pagedMessage)
	if err != nil {
		return err
	}

	err = setPagedMessage(p.redis, pagedMessage.MessageID, pagedMessage)
	return err
}

func (p *Paginator) EmbedPaginator(
	guildID, channelID, userID string, embeds ...*discordgo.MessageEmbed,
) error {
	if len(embeds) == 0 {
		return nil
	}

	if len(embeds) < 2 {
		_, err := p.sendComplex(guildID, channelID, &discordgo.MessageSend{
			Embed: embeds[0],
		})
		return err
	}

	pagedMessage := &PagedEmbedMessage{
		FullEmbed:       embeds[0],
		ChannelID:       channelID,
		GuildID:         guildID,
		CurrentPage:     1,
		TotalNumOfPages: len(embeds),
		UserID:          userID,
		Embeds:          embeds,
		Type:            EmbedType,
	}

	err := p.setupAndSendFirstMessage(pagedMessage)
	if err != nil {
		return err
	}

	err = setPagedMessage(p.redis, pagedMessage.MessageID, pagedMessage)
	return err
}
