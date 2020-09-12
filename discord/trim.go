package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// TrimEmbed enforces Embed Limits to the given Embed
// Source: https://discordapp.com/developers/docs/resources/channel#embed-limits
func TrimEmbed(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	if embed == nil || (&discordgo.MessageEmbed{}) == embed {
		return nil
	}

	if embed.Title != "" && len(embed.Title) > 256 {
		embed.Title = embed.Title[0:255] + "…"
	}
	if len(embed.Description) > 2048 {
		embed.Description = embed.Description[0:2047] + "…"
	}
	if embed.Footer != nil && len(embed.Footer.Text) > 2048 {
		embed.Footer.Text = embed.Footer.Text[0:2047] + "…"
	}
	if embed.Author != nil && len(embed.Author.Name) > 256 {
		embed.Author.Name = embed.Author.Name[0:255] + "…"
	}
	newFields := make([]*discordgo.MessageEmbedField, 0)
	for _, field := range embed.Fields {
		if strings.TrimSpace(field.Name) == "" ||
			strings.TrimSpace(field.Value) == "" {
			continue
		}
		if len(field.Name) > 256 {
			field.Name = field.Name[0:255] + "…"
		}
		// TODO: better cutoff (at commas and stuff)
		if len(field.Value) > 1024 {
			field.Value = field.Value[0:1023] + "…"
		}
		newFields = append(newFields, field)
		if len(newFields) >= 25 {
			break
		}
	}
	embed.Fields = newFields

	if CalculateEmbedLength(embed) > 6000 {
		if embed.Footer != nil {
			embed.Footer.Text = ""
		}
		if CalculateEmbedLength(embed) > 6000 {
			if embed.Author != nil {
				embed.Author.Name = ""
			}
			if CalculateEmbedLength(embed) > 6000 {
				embed.Fields = []*discordgo.MessageEmbedField{{}}
			}
		}
	}

	return embed
}

// CalculateEmbedLength calculates the total length of embed
func CalculateEmbedLength(embed *discordgo.MessageEmbed) int {
	var count int
	count += len(embed.Title)
	count += len(embed.Description)
	if embed.Footer != nil {
		count += len(embed.Footer.Text)
	}
	if embed.Author != nil {
		count += len(embed.Author.Name)
	}
	for _, field := range embed.Fields {
		count += len(field.Name)
		count += len(field.Value)
	}
	return count
}
