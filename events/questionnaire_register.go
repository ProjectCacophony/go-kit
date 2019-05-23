package events

import (
	"encoding/json"
	"time"
)

const (
	defaultExpiration = 15 * time.Minute
)

func (q *Questionnaire) Register(
	key string,
	filter QuestionnaireFilter,
	payload map[string]interface{},
) error {
	return q.RegisterWithExpiration(key, filter, payload, defaultExpiration)
}

func (q *Questionnaire) RegisterWithExpiration(
	key string,
	filter QuestionnaireFilter,
	payload map[string]interface{},
	expiration time.Duration,
) error {
	body, err := json.Marshal(QuestionnaireMatch{
		Key:     key,
		Payload: payload,
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
