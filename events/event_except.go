package events

import (
	"github.com/bwmarrin/discordgo"
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

		errorMessage = errD.Message.Message
	}
	// // do not log state errors
	// if err == state.ErrStateNotFound ||
	// 	err == state.ErrTargetWrongServer ||
	// 	err == state.ErrTargetWrongType {
	// 	dontLog = true
	// }

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
		e.Logger().Error("received error while executing event", zap.Error(err))
		// TODO: send to sentry…
	}
}
