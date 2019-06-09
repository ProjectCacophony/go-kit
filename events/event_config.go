package events

import (
	"time"
)

func (e *Event) Timezone() *time.Location {
	// TODO: return timezone based on user config

	return time.FixedZone("UTC", 0)
}
