package paginator

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"
)

func (p *Paginator) handleReaction(message *PagedEmbedMessage, reaction *discordgo.MessageReactionAdd) error {
	session, err := p.getSession(message.BotID)
	if err != nil {
		return err
	}

	switch reaction.Emoji.Name {

	case CloseEmoji:
		err = deletePagedMessage(p.redis, reaction.MessageID)
		if err != nil {
			return err
		}

		err = discord.Delete(
			p.redis, session, message.ChannelID, message.MessageID, message.DM)
		if err != nil {
			return err
		}

	case NumbersEmoji:
		if isNumbersListening(p.redis, message.ChannelID, message.UserID) {
			return nil
		}

		resp, err := p.sendComplex(
			message.BotID, message.ChannelID, &discordgo.MessageSend{
				Content: "<@" + message.UserID + "> Which page would you like to open? <:blobidea:317047867036663809>",
			},
			message.DM,
		)
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

		discord.RemoveReact(
			p.redis, session, message.ChannelID, message.MessageID, reaction.UserID, message.DM, NumbersEmoji,
		)

	case LeftArrowEmoji:
		err = p.setPage(message, message.CurrentPage-1)
		if err != nil {
			return err
		}

		if message.Type != ImageType {
			discord.RemoveReact(
				p.redis, session, message.ChannelID, reaction.MessageID, reaction.UserID, message.DM, reaction.Emoji.Name,
			)
		}

	case RightArrowEmoji:
		err = p.setPage(message, message.CurrentPage+1)
		if err != nil {
			return err
		}

		if message.Type != ImageType {
			discord.RemoveReact(
				p.redis, session, message.ChannelID, reaction.MessageID, reaction.UserID, message.DM, reaction.Emoji.Name,
			)
		}

	}

	return nil
}
