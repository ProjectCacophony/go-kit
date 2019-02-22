package discord

import (
	"strconv"
	"time"
)

const (
	discordEpoch int64 = 1420070400000
)

// TimeFromID extracts a time from a Snowflake ID, used at the Discord API
func TimeFromID(id string) (*time.Time, error) {
	iid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}

	idTime := time.Unix(((iid>>22)+discordEpoch)/1000, 0).UTC()

	return &idTime, nil
}
