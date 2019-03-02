package paginator

import (
	"strconv"

	"github.com/go-redis/redis"

	"github.com/bwmarrin/discordgo"
)

func (p *Paginator) HandleMessageReactionAdd(messageReactionAdd *discordgo.MessageReactionAdd) error {
	if !validReactions[messageReactionAdd.Emoji.Name] {
		return nil
	}

	pagedMessage, err := getPagedMessage(
		p.redis, messageReactionAdd.MessageID,
	)
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	if messageReactionAdd.UserID != pagedMessage.UserID {
		return err
	}

	return p.handleReaction(
		pagedMessage,
		messageReactionAdd,
	)
}

func (p *Paginator) HandleMessageCreate(messageCreate *discordgo.MessageCreate) error {
	if !p.messageRegexp.MatchString(messageCreate.Content) {
		return nil
	}

	page, err := strconv.Atoi(messageCreate.Content)
	if err != nil {
		return err
	}

	if !isNumbersListening(
		p.redis, messageCreate.ChannelID, messageCreate.Author.ID,
	) {
		return nil
	}

	listener, err := getNumbersListeningMessageDelete(
		p.redis, messageCreate.ChannelID, messageCreate.Author.ID,
	)
	if err != nil {
		return err
	}

	message, err := getPagedMessage(p.redis, listener.PagedEmbedMessageID)
	if err != nil {
		return err
	}

	err = p.setPage(message, page)
	if err != nil {
		return err
	}

	session, err := p.getSession(messageCreate.GuildID)
	if err != nil {
		return err
	}

	// clean up
	err = session.ChannelMessageDelete(messageCreate.ChannelID, listener.MessageID)
	if err != nil {
		return err
	}
	return session.ChannelMessageDelete(messageCreate.ChannelID, messageCreate.ID)
}
