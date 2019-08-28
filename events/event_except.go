package events

import (
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

			message := "**Something went wrong.** :sad:" + "\n**Error:** " + e.Translate(errorMessage) + ""
			if doLog {
				message += "I sent our top people to fix the issue as soon as possible."
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

	// state errors
	if err == state.ErrPresenceStateNotFound ||
		err == state.ErrRoleStateNotFound ||
		err == state.ErrEmojiStateNotFound ||
		err == state.ErrTargetWrongServer ||
		err == state.ErrTargetWrongType ||
		err == state.ErrUserNotFound ||
		err == state.ErrChannelNotFound ||
		err == state.ErrRoleNotFound ||
		strings.Contains(err.Error(), NoStoragePermission) ||
		strings.Contains(err.Error(), NoStorageSpace) {
		return true
	}

	return false
}
