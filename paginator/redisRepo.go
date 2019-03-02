package paginator

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

// nolint: gochecknoglobals
var expiration = time.Hour

func pagedMessageKey(messageID string) string {
	return "cacophony:paginator:paged-message:" + messageID
}

func numbersListenerKey(channelID, userID string) string {
	return "cacophony:paginator:numbers-listener:" + channelID + ":" + userID
}

func getPagedMessage(redis *redis.Client, messageID string) (*PagedEmbedMessage, error) {
	data, err := redis.Get(pagedMessageKey(messageID)).Bytes()
	if err != nil {
		return nil, err
	}

	var message PagedEmbedMessage
	err = json.Unmarshal(data, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func setPagedMessage(redis *redis.Client, messageID string, message *PagedEmbedMessage) error {
	if message == nil {
		return nil
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return redis.Set(pagedMessageKey(messageID), data, expiration).Err()
}

func deletePagedMessage(redis *redis.Client, messageID string) error {
	return redis.Del(pagedMessageKey(messageID)).Err()
}

func addNumbersListener(redis *redis.Client, channelID, userID string, listener *numbersListener) error {
	if listener == nil {
		return nil
	}

	data, err := json.Marshal(listener)
	if err != nil {
		return err
	}

	return redis.Set(numbersListenerKey(channelID, userID), data, expiration).Err()
}

func isNumbersListening(redis *redis.Client, channelID, userID string) bool {
	res, _ := redis.Exists(numbersListenerKey(channelID, userID)).Result() // nolint: errcheck
	return res > 0
}

func getNumbersListeningMessageDelete(redis *redis.Client, channelID, userID string) (*numbersListener, error) {
	pipeline := redis.TxPipeline()

	get := pipeline.Get(numbersListenerKey(channelID, userID))
	pipeline.Del(numbersListenerKey(channelID, userID))

	_, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}

	data, err := get.Bytes()
	if err != nil {
		return nil, err
	}

	var listener numbersListener
	err = json.Unmarshal(data, &listener)
	if err != nil {
		return nil, err
	}

	return &listener, err
}
