package events

import (
	"errors"
	"time"

	// embeding timezones
	_ "time/tzdata"

	"gitlab.com/Cacophony/go-kit/config"
	"go.opentelemetry.io/otel/api/global"
)

const userTimezoneKey = "cacophony:kit:timezone"

var defaultTimezone = time.FixedZone("UTC", 0)

func (e *Event) Timezone() *time.Location {
	_, span := global.Tracer("cacophony.dev/kit").Start(e.Context(), "events.Event.Timezone")
	defer span.End()

	if e.timezone != nil {
		return e.timezone
	}

	if e.DB() == nil || e.UserID == "" {
		return defaultTimezone
	}

	timezoneText, err := config.UserGetString(e.DB(), e.UserID, userTimezoneKey)
	if err != nil {
		e.ExceptSilent(err)
		return defaultTimezone
	}

	timezone, err := time.LoadLocation(timezoneText)
	if err != nil {
		e.ExceptSilent(err)
		return defaultTimezone
	}

	e.timezone = timezone
	return timezone
}

func (e *Event) SetTimezone(timezone *time.Location) error {
	if timezone == nil {
		return errors.New("timezone cannot be empty")
	}

	if e.DB() == nil || e.UserID == "" {
		return errors.New("event is missing fields to set timezone")
	}

	return config.UserSetString(e.DB(), e.UserID, userTimezoneKey, timezone.String())
}
