package events

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// FindUser finds any kind of target user in the command
func (e *Event) FindUser() (*discordgo.User, error) {
	if e.Type != MessageCreateType {
		return nil, errors.New("event has to be MessageCreate")
	}

	// try any mentions in the command
	for _, mention := range e.MessageCreate.Mentions {
		user, err := e.State().User(mention.ID)
		if err == nil {
			return user, nil
		}
	}

	for _, field := range e.Fields() {
		user, err := e.State().UserFromMention(field)
		if err == nil {
			return user, nil
		}
	}

	return e.State().User(e.UserID)
}
