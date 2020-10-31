package events

import (
	"errors"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	raven "github.com/getsentry/raven-go"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/state"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	"go.uber.org/zap"
)

// TODO: ratelimit error sending

func (e *Event) Except(err error, fields ...string) {
	if err == nil {
		return
	}

	fieldsMap := fieldsListToMap(fields)

	_, span := global.Tracer("cacophony.dev/kit").Start(e.Context(), "event.Except",
		trace.WithAttributes(label.Any("fields", fieldsMap), label.String("error", err.Error())),
	)
	defer span.End()
	span.RecordError(e.Context(), err, trace.WithErrorStatus(codes.Unknown))

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
		e.Logger().Error("error occurred while executing event", zap.Error(err), zap.Any("fields", fields))

		if raven.DefaultClient != nil {
			raven.CaptureError(
				err,
				generateRavenTags(e, false, fieldsMap),
				&raven.User{
					ID: e.UserID,
				},
			)
		}
	}
}

func (e *Event) ExceptSilent(err error, fields ...string) {
	if ignoreError(err) {
		return
	}

	fieldsMap := fieldsListToMap(fields)

	_, span := global.Tracer("cacophony.dev/kit").Start(e.Context(), "event.ExceptSilent",
		trace.WithAttributes(label.Any("fields", fieldsMap), label.String("error", err.Error())),
	)
	defer span.End()
	span.RecordError(e.Context(), err, trace.WithErrorStatus(codes.Unknown))

	if e.logger != nil {
		e.Logger().Error("silent occurred error while executing event", zap.Error(err), zap.Any("fields", fields))
	}

	if raven.DefaultClient != nil {
		raven.CaptureError(
			err,
			generateRavenTags(e, true, fieldsMap),
			&raven.User{
				ID: e.UserID,
			},
		)
	}
}

func generateRavenTags(event *Event, silent bool, fields map[string]string) map[string]string {
	if fields == nil {
		fields = make(map[string]string)
	}

	fields["event_id"] = event.ID
	fields["event_type"] = string(event.Type)
	fields["bot_id"] = event.BotUserID
	if event.GuildID != "" {
		fields["guild_id"] = event.GuildID
	}
	if event.ChannelID != "" {
		fields["channel_ud"] = event.ChannelID
	}
	if event.UserID != "" {
		fields["user_id"] = event.UserID
	}
	if event.MessageID != "" {
		fields["message_id"] = event.MessageID
	}
	fields["silent"] = strconv.FormatBool(silent)

	if event.Type == MessageCreateType {
		fields["message_content"] = event.MessageCreate.Content
	}

	return fields
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
		errors.Is(err, discord.ErrNoDMChannel) ||
		strings.Contains(err.Error(), NoStoragePermission) ||
		strings.Contains(err.Error(), NoStorageSpace) ||
		strings.Contains(err.Error(), FileTooBig) ||
		strings.Contains(err.Error(), CouldNotExtractFilename) {
		return true
	}

	var userError *UserError
	return errors.As(err, &userError)
}

func fieldsListToMap(fields []string) map[string]string {
	fieldsData := map[string]string{}

	for i := range fields {
		if i%2 == 0 && len(fields) > i+1 {
			fieldsData[fields[i]] = fields[i+1]
		}
	}

	return fieldsData
}
