package events

import (
	"errors"
	"strconv"
	"strings"

	"gitlab.com/Cacophony/go-kit/discord"

	"github.com/bwmarrin/discordgo"
	raven "github.com/getsentry/raven-go"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

// TODO: ratelimit error sending

func (e *Event) Except(err error) {
	if err == nil {
		return
	}

	doLog := true

	errorMessage := err.Error()

	if ignoreError(err) {
		doLog = false
	}

	// do not log discord permission errors
	if errD, ok := err.(*discordgo.RESTError); ok && errD != nil && errD.Message != nil {
		if errD.Message.Message != "" {
			errorMessage = errD.Message.Message
		}
	}

	if e.Type == MessageCreateType {
		if e.DM() ||
			discord.UserHasPermission(e.State(), e.BotUserID, e.ChannelID, discordgo.PermissionSendMessages) {

			message := "**Something went wrong.** :sad:" + "\n**Error:** " + e.Translate(errorMessage) + "\n" +
				"I sent our top people to fix the issue as soon as possible.\n"
			if !doLog {
				message = e.Translate(errorMessage)
			}

			e.Respond(
				message,
			)

		} else if discord.UserHasPermission(e.State(), e.BotUserID, e.ChannelID, discordgo.PermissionAddReactions) {

			e.React(
				":stop:", ":shh:", ":nogood:", ":speaknoevil:",
			)
		}

	}

	if doLog {
		e.Logger().Error("error occurred while executing event", zap.Error(err))

		if raven.DefaultClient != nil {
			raven.CaptureError(
				err,
				generateRavenTags(e, false),
				&raven.User{
					ID: e.UserID,
				},
			)
		}
	}
}

func (e *Event) ExceptSilent(err error) {
	if ignoreError(err) {
		return
	}

	if e.logger != nil {
		e.Logger().Error("silent occurred error while executing event", zap.Error(err))
	}

	if raven.DefaultClient != nil {
		raven.CaptureError(
			err,
			generateRavenTags(e, true),
			&raven.User{
				ID: e.UserID,
			},
		)
	}
}

func generateRavenTags(event *Event, silent bool) map[string]string {
	tags := map[string]string{
		"event_id":    event.ID,
		"event_type:": string(event.Type),
		"bot_id":      event.BotUserID,
		"guild_id":    event.GuildID,
		"silent":      strconv.FormatBool(silent),
	}

	if event.Type == MessageCreateType {
		tags["message_content"] = event.MessageCreate.Content
	}

	return tags
}

func ignoreError(err error) bool {
	if err == nil {
		return true
	}

	// discord permission errors
	if errD, ok := err.(*discordgo.RESTError); ok && errD != nil && errD.Message != nil {
		if errD.Message.Code == discordgo.ErrCodeMissingPermissions ||
			errD.Message.Code == discordgo.ErrCodeMissingAccess ||
			errD.Message.Code == discordgo.ErrCodeCannotSendMessagesToThisUser {
			return true
		}
	}

	if errors.Is(err, state.ErrPresenceStateNotFound) ||
		errors.Is(err, state.ErrRoleStateNotFound) ||
		errors.Is(err, state.ErrEmojiStateNotFound) ||
		errors.Is(err, state.ErrTargetWrongServer) ||
		errors.Is(err, state.ErrTargetWrongType) ||
		errors.Is(err, state.ErrUserNotFound) ||
		errors.Is(err, state.ErrChannelNotFound) ||
		errors.Is(err, state.ErrRoleNotFound) ||
		errors.Is(err, discord.ErrInvalidMessageLink) ||
		errors.Is(err, discord.ErrMessageOnWrongServer) ||
		strings.Contains(err.Error(), NoStoragePermission) ||
		strings.Contains(err.Error(), NoStorageSpace) ||
		strings.Contains(err.Error(), FileTooBig) ||
		strings.Contains(err.Error(), CouldNotExtractFilename) {
		return true
	}

	var userError *UserError
	return errors.As(err, &userError)
}
