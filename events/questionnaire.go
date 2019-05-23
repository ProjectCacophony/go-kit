package events

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type Questionnaire struct {
	logger    *zap.Logger
	redis     *redis.Client
	publisher *Publisher
}

func NewQuestionnaire(
	logger *zap.Logger,
	redis *redis.Client,
	publisher *Publisher,
) *Questionnaire {
	return &Questionnaire{
		logger:    logger,
		redis:     redis,
		publisher: publisher,
	}
}

func (q *Questionnaire) Do(ctx context.Context, event *Event) (bool, error) {
	exists, err := q.redis.WithContext(ctx).SCard(questionnaireEventRedisKey(event)).Result()
	if err != nil {
		return false, nil
	}
	if exists <= 0 {
		return false, nil
	}

	key := questionnaireEventRedisKey(event)

	pipe := q.redis.WithContext(ctx).TxPipeline()

	items := pipe.SMembers(key)
	pipe.Del(key)

	_, err = pipe.Exec()
	if err != nil {
		return false, err
	}

	values, err := items.Result()
	if err != nil {
		return false, err
	}

	if len(values) <= 0 {
		return false, nil
	}

	var recoverable bool

	for _, value := range values {
		var match QuestionnaireMatch
		err = json.Unmarshal([]byte(value), &match)
		if err != nil {
			return false, err
		}

		newEvent, err := New(CacophonyQuestionnaireMatch)
		if err != nil {
			return false, err
		}

		newEvent.UserID = event.UserID
		newEvent.GuildID = event.GuildID
		newEvent.ChannelID = event.ChannelID
		newEvent.BotUserID = event.BotUserID
		newEvent.QuestionnaireMatch = &match

		// TODO(snake): temp fix for message content being carried to new event
		newEvent.MessageCreate = event.MessageCreate

		err, recoverable = q.publisher.Publish(ctx, newEvent)
		if err != nil {
			if !recoverable {
				q.logger.Fatal(
					"received unrecoverable error while publishing questionnaire",
					zap.Error(err),
				)
			}
			return false, err
		}
	}

	return true, nil
}
