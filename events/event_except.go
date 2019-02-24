package events

import (
	"strconv"

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

	// do not log discord permission errors
	if errD, ok := err.(*discordgo.RESTError); ok && errD != nil && errD.Message != nil {
		if errD.Message.Code == discordgo.ErrCodeMissingPermissions ||
			errD.Message.Code == discordgo.ErrCodeMissingAccess ||
			errD.Message.Code == discordgo.ErrCodeCannotSendMessagesToThisUser {
			doLog = false
		}

		if errD.Message.Message != "" {
			errorMessage = errD.Message.Message
		}
	}
	// do not log state errors
	if err == state.ErrStateNotFound ||
		err == state.ErrTargetWrongServer ||
		err == state.ErrTargetWrongType ||
		err == state.ErrUserNotFound ||
		err == state.ErrChannelNotFound ||
		err == state.ErrRoleNotFound {
		doLog = false
	}

	if e.Type == MessageCreateType {
		// TODO: send reaction instead if we are not allowed to send messages (check permissions from state)
		// TODO: better message, emoji, translation, …

		message := "**Something went wrong.** " + "\n```\nError: " + errorMessage + "\n```"
		if doLog {
			message += "I sent our top people to fix the issue as soon as possible."
		}

		e.Respond( // nolint: errcheck
			message,
		)
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
