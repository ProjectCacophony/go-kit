package paginator

import (
	"github.com/bwmarrin/discordgo"
)

func (p *Paginator) handleReaction(message *PagedEmbedMessage, reaction *discordgo.MessageReactionAdd) error {
	session, err := p.getSession(message.GuildID)
	if err != nil {
		return err
	}

	switch reaction.Emoji.Name {

	case CloseEmoji:
		err = deletePagedMessage(p.redis, reaction.MessageID)
		if err != nil {
			return err
		}

		err = session.ChannelMessageDelete(message.ChannelID, message.MessageID)
		if err != nil {
			return err
		}

	case NumbersEmoji:
		if isNumbersListening(p.redis, message.ChannelID, message.UserID) {
			return nil
		}

		resp, err := p.sendComplex(
			message.GuildID, message.ChannelID, &discordgo.MessageSend{
				Content: "<@" + message.UserID + "> Which page would you like to open? <:blobidea:317047867036663809>",
			})
		if err != nil {
			return err
		}

		err = addNumbersListener(p.redis, message.ChannelID, message.UserID, &numbersListener{
			MessageID:           resp[0].ID,
			PagedEmbedMessageID: message.MessageID,
		})
		if err != nil {
			return err
		}

		session.MessageReactionRemove( // nolint: errcheck
			message.ChannelID, message.MessageID, NumbersEmoji, message.UserID,
		)

	case LeftArrowEmoji:
		err = p.setPage(message, message.CurrentPage-1)
		if err != nil {
			return err
		}

		if message.Type != ImageType {
			session.MessageReactionRemove( // nolint: errcheck
				reaction.ChannelID, reaction.MessageID, reaction.Emoji.Name, reaction.UserID,
			)
		}

	case RightArrowEmoji:
		err = p.setPage(message, message.CurrentPage+1)
		if err != nil {
			return err
		}

		if message.Type != ImageType {
			session.MessageReactionRemove( // nolint: errcheck
				reaction.ChannelID, reaction.MessageID, reaction.Emoji.Name, reaction.UserID,
			)
		}

	}

	return nil
}
