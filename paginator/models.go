package paginator

import (
	"github.com/bwmarrin/discordgo"
)

type PagedEmbedMessage struct {
	Files           []*File
	FullEmbed       *discordgo.MessageEmbed
	TotalNumOfPages int
	CurrentPage     int
	FieldsPerPage   int
	Color           int
	MessageID       string
	ChannelID       string
	UserID          string // user who triggered the message
	Type            Type
	Embeds          []*discordgo.MessageEmbed
	DM              bool
	BotID           string
}

type File struct {
	Name        string
	ContentType string
	Data        []byte
}

type Type int

const (
	FieldType Type = iota
	ImageType
	EmbedType
)

type numbersListener struct {
	MessageID           string // message ID of the message asking the user which page to choose
	PagedEmbedMessageID string // message ID of the paged embed message
}
