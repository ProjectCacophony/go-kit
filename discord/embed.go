package discord

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/regexp"
)

func EmbedCodeFromMessage(message *discordgo.Message) string {
	if message == nil {
		return ""
	}

	var embedCode string

	if len(message.Content) > 0 {
		embedCode += "ptext=" + cleanEmbedValue(message.Content) + " | "
	}

	if message.Embeds == nil || len(message.Embeds) <= 0 {
		return embedCode
	}

	targetEmbed := message.Embeds[0]

	if targetEmbed.Author != nil && targetEmbed.Author.Name != "" {
		if targetEmbed.Author.URL == "" && targetEmbed.Author.IconURL == "" {
			embedCode += "author=" + cleanEmbedValue(targetEmbed.Author.Name) + " | "
		} else {
			embedCode += "author=name=" + cleanEmbedValue(targetEmbed.Author.Name)
			if targetEmbed.Author.URL != "" {
				embedCode += " url=" + cleanEmbedValue(targetEmbed.Author.URL)
			}
			if targetEmbed.Author.IconURL != "" {
				embedCode += " icon=" + cleanEmbedValue(targetEmbed.Author.IconURL)
			}
			embedCode += " | "
		}
	}
	if targetEmbed.Title != "" {
		embedCode += "title=" + cleanEmbedValue(targetEmbed.Title) + " | "
	}
	if targetEmbed.Description != "" {
		embedCode += "description=" + cleanEmbedValue(targetEmbed.Description) + " | "
	}
	if targetEmbed.Thumbnail != nil && targetEmbed.Thumbnail.URL != "" {
		embedCode += "thumbnail=" + cleanEmbedValue(targetEmbed.Thumbnail.URL) + " | "
	}
	if targetEmbed.Image != nil && targetEmbed.Image.URL != "" {
		embedCode += "image=" + cleanEmbedValue(targetEmbed.Image.URL) + " | "
	}
	if targetEmbed.Fields != nil && len(targetEmbed.Fields) > 0 {
		for _, targetField := range targetEmbed.Fields {
			if targetField.Inline {
				embedCode += "field=name=" + cleanEmbedValue(targetField.Name) +
					" value=" + cleanEmbedValue(targetField.Value) + " | "
			} else {
				embedCode += "field=name=" + cleanEmbedValue(targetField.Name) +
					" value=" + cleanEmbedValue(targetField.Value) + " inline=no | "
			}
		}
	}
	if targetEmbed.Footer != nil && targetEmbed.Footer.Text != "" {
		if targetEmbed.Footer.IconURL == "" {
			embedCode += "footer=" + cleanEmbedValue(targetEmbed.Footer.Text)
		} else {
			embedCode += "footer=name=" + cleanEmbedValue(targetEmbed.Footer.Text) +
				" icon=" + cleanEmbedValue(targetEmbed.Footer.IconURL)
		}
		embedCode += " | "
	}
	if targetEmbed.Color > 0 {
		embedCode += "color=#" + ColorCodeToHex(targetEmbed.Color) + " | "
	}

	embedCode = strings.TrimSuffix(embedCode, " | ")

	return replaceEmojiCodes(embedCode)
}

func EmbedCodeToMessage(embedText string) *discordgo.MessageSend {
	// Code ported from https://github.com/appu1232/Discord-Selfbot/blob/master/cogs/misc.py#L146
	// Reference https://github.com/Seklfreak/Robyul-Web/blob/master/src/RobyulWebBundle/Resources/public/js/main.js#L724
	var ptext, title, description, image, thumbnail, color, footer, author string
	var timestamp time.Time

	embedValues := strings.Split(embedText, "|")
	for _, embedValue := range embedValues {
		embedValue = strings.TrimSpace(embedValue)
		if strings.HasPrefix(embedValue, "ptext=") {
			ptext = strings.TrimSpace(embedValue[6:])
		} else if strings.HasPrefix(embedValue, "title=") {
			title = strings.TrimSpace(embedValue[6:])
		} else if strings.HasPrefix(embedValue, "description=") {
			description = strings.TrimSpace(embedValue[12:])
		} else if strings.HasPrefix(embedValue, "desc=") {
			description = strings.TrimSpace(embedValue[5:])
		} else if strings.HasPrefix(embedValue, "image=") {
			image = strings.TrimSpace(embedValue[6:])
		} else if strings.HasPrefix(embedValue, "thumbnail=") {
			thumbnail = strings.TrimSpace(embedValue[10:])
		} else if strings.HasPrefix(embedValue, "colour=") {
			color = strings.TrimSpace(embedValue[7:])
		} else if strings.HasPrefix(embedValue, "color=") {
			color = strings.TrimSpace(embedValue[6:])
		} else if strings.HasPrefix(embedValue, "footer=") {
			footer = strings.TrimSpace(embedValue[7:])
		} else if strings.HasPrefix(embedValue, "author=") {
			author = strings.TrimSpace(embedValue[7:])
		} else if strings.HasPrefix(embedValue, "timestamp") {
			timestamp = time.Now()
		} else if description == "" && !strings.HasPrefix(embedValue, "field=") {
			description = embedValue
		}
	}

	if ptext == "" && title == "" && description == "" && image == "" && thumbnail == "" && footer == "" &&
		author == "" && !strings.Contains("field=", embedText) {
		return &discordgo.MessageSend{Content: embedText}
	}

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
	}
	if !timestamp.IsZero() {
		embed.Timestamp = timestamp.Format(time.RFC3339)
	}
	if color != "" {
		embed.Color = HexToColorCode(color)
	}

	var fieldValues []string
	var field, fieldName, fieldValue string
	var fieldInline bool
	for _, embedValue := range embedValues {
		embedValue = strings.TrimSpace(embedValue)
		if strings.HasPrefix(embedValue, "field=") {
			fieldInline = true
			field = strings.TrimSpace(strings.TrimPrefix(embedValue, "field="))
			fieldValues = strings.SplitN(field, "value=", 2)
			if len(fieldValues) >= 2 {
				fieldName = fieldValues[0]
				fieldValue = fieldValues[1]
			} else if len(fieldValues) >= 1 {
				fieldName = fieldValues[0]
				fieldValue = ""
			}
			if strings.Contains(fieldValue, "inline=") {
				fieldValues = strings.SplitN(fieldValue, "inline=", 2)
				if len(fieldValues) >= 2 {
					fieldValue = fieldValues[0]
					if strings.Contains(strings.ToLower(fieldValues[1]), "false") ||
						strings.Contains(strings.ToLower(fieldValues[1]), "no") {
						fieldInline = false
					}
				} else if len(fieldValues) >= 1 {
					fieldValue = fieldValues[0]
				}
			}
			fieldName = strings.TrimSpace(strings.TrimPrefix(fieldName, "name="))
			fieldValue = strings.TrimSpace(fieldValue)
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   fieldName,
				Value:  fieldValue,
				Inline: fieldInline,
			})
		}
	}
	var authorValues, iconValues []string
	if author != "" {
		if strings.Contains(author, "icon=") {
			authorValues = strings.SplitN(author, "icon=", 2)
			if len(authorValues) >= 2 {
				if strings.Contains(authorValues[1], "url=") {
					iconValues = strings.SplitN(authorValues[1], "url=", 2)
					if len(iconValues) >= 2 {
						embed.Author = &discordgo.MessageEmbedAuthor{
							Name:    strings.TrimSpace(authorValues[0][5:]),
							IconURL: strings.TrimSpace(iconValues[0]),
							URL:     strings.TrimSpace(iconValues[1]),
						}
					}
				} else {
					embed.Author = &discordgo.MessageEmbedAuthor{
						Name:    strings.TrimSpace(authorValues[0][5:]),
						IconURL: strings.TrimSpace(authorValues[1]),
					}
				}
			}
		} else {
			if strings.Contains(author, "url=") {
				authorValues = strings.SplitN(author, "url=", 2)
				if len(iconValues) >= 2 {
					embed.Author = &discordgo.MessageEmbedAuthor{
						Name: strings.TrimSpace(authorValues[0][5:]),
						URL:  strings.TrimSpace(authorValues[1]),
					}
				}
			} else {
				embed.Author = &discordgo.MessageEmbedAuthor{
					Name: strings.TrimSpace(author),
				}
			}
		}
	}
	if image != "" {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: image,
		}
	}
	if thumbnail != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: thumbnail,
		}
	}
	var footerValues []string
	if footer != "" {
		if strings.Contains(footer, "icon=") {
			footerValues = strings.SplitN(footer, "icon=", 2)
			if len(footerValues) >= 2 {
				embed.Footer = &discordgo.MessageEmbedFooter{
					Text:    strings.TrimSpace(footerValues[0])[5:],
					IconURL: strings.TrimSpace(footerValues[1]),
				}
			}
		} else {
			embed.Footer = &discordgo.MessageEmbedFooter{
				Text: strings.TrimSpace(footer),
			}
		}
	}

	return &discordgo.MessageSend{
		Content: ptext,
		Embed:   embed,
	}
}

func cleanEmbedValue(input string) (output string) {
	return strings.Replace(input, "|", "-", -1)
}

func replaceEmojiCodes(content string) (result string) {
	var replaceWith string

	emojiPartsList := regexp.DiscordEmojiRegexp.FindAllStringSubmatch(content, -1)
	if len(emojiPartsList) > 0 {
		for _, emojiParts := range emojiPartsList {
			replaceWith = ":" + emojiParts[1] + ":"

			if replaceWith != "" {
				content = strings.Replace(content, emojiParts[0], replaceWith, -1)
			}
		}
	}

	return content
}
