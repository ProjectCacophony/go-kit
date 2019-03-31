package events

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/text"
)

const (
	defaultPrefix = "."
)

// Parse parses the content of a message into fields
func (e *Event) Parse() {
	// only Message Create can be command
	if e.Type != MessageCreateType {
		return
	}

	// set DM field
	if e.MessageCreate.GuildID == "" {
		e.dm = true
	}

	// ignore empty messages (eg embeds)
	if e.MessageCreate.Content == "" {
		return
	}

	// ignore messages by bots
	if e.MessageCreate.Author == nil || e.MessageCreate.Author.Bot {
		return
	}

	// trim message
	content := strings.TrimSpace(e.MessageCreate.Content)

	// ignore empty messages
	if content == "" {
		return
	}

	// ignore messages without prefix
	if !strings.HasPrefix(content, e.Prefix()) {

		// if message doesn't have a prefix, check to see if it starts with bot mention
		if !strings.HasPrefix(content, "<@") || len(e.MessageCreate.Mentions) == 0 || e.MessageCreate.Mentions[0].ID != e.BotUserID {
			return
		}

		// Replace starting mentions with prefix
		content = strings.Replace(content, "<@"+e.BotUserID+">", e.Prefix(), -1)
		content = strings.Replace(content, "<@!"+e.BotUserID+">", e.Prefix(), -1)
	}

	args, err := text.ToArgv(content[len(e.Prefix()):])
	if err != nil {
		return
	}

	if len(args) >= 1 {
		args[0] = strings.ToLower(args[0])
	}

	// extract fields of command without prefix
	e.command = true
	e.fields = args
}

// Fields returns all fields of a command, includes the command
func (e *Event) Fields() []string {
	return e.fields
}

// Command returns true if the event is a command message
func (e *Event) Command() bool {
	return e.command
}

// Prefix returns the prefix of a command, if event is a command
func (e *Event) Prefix() string {
	if e.prefix == "" {
		// TODO: check for guild prefix

		e.prefix = defaultPrefix
	}

	return e.prefix
}

// DM returns true if the event is a DM message
func (e *Event) DM() bool {
	return e.dm
}
