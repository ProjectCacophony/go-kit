package events

import (
	"encoding/json"
	"time"
)

const (
	expiration = 15 * time.Minute
)

func (q *Questionnaire) Register(
	key string,
	filter QuestionnaireFilter,
	Payload map[string]interface{},
) error {
	body, err := json.Marshal(QuestionnaireMatch{
		Key:     key,
		Payload: Payload,
	})
	if err != nil {
		return err
	}

	redisKey := questionnaireRedisKey(
		filter.Type,
		filter.GuildID,
		filter.ChannelID,
		filter.UserID,
	)

	err = q.redis.SAdd(
		redisKey,
		body,
	).Err()
	if err != nil {
		return err
	}

	return q.redis.Expire(redisKey, expiration).Err()
}
