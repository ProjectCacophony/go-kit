package discord

import (
	"errors"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Session struct {
	Client *discordgo.Session
	BotID  string
}

// NewSession creates a new DiscordGo Client for the given BotID
// from the tokens map (key: BotID, value: token)
func NewSession(tokens map[string]string, botID string) (*Session, error) {
	token := tokens[botID]
	if token == "" {
		return nil, errors.New("token for bot ID is not configured")
	}

	newSession, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, errors.New("error creating bot session")
	}
	newSession.MaxRestRetries = 0

	newSession.Client = &http.Client{
		Timeout: 1 * time.Hour,
	}

	return &Session{
		Client: newSession,
		BotID:  botID,
	}, nil
}
