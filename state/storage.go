package state

import (
	"sync"
	"time"

	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
)

const (
	redisKeyPrefix = "cacophony:state:"
)

var stateLock sync.Mutex

var stateExpire = time.Duration(0)

func allUserIDsSetKey() string {
	return redisKeyPrefix + "user-ids"
}

func allGuildIDsSetKey() string {
	return redisKeyPrefix + "guild-ids"
}

func allChannelIDsSetKey() string {
	return redisKeyPrefix + "channel-ids"
}

func guildBotIDsSetKey(guildID string) string {
	return redisKeyPrefix + "guild-" + guildID + ":bot-ids"
}

func guildBannedUserIDsSetKey(guildID string) string {
	return redisKeyPrefix + "guild-" + guildID + ":banned-userids"
}

func guildBannedUserIDInitializedGuildIDsSetKey() string {
	return redisKeyPrefix + "banned-userids-initialized-guild-ids"
}

func userKey(userID string) string {
	return redisKeyPrefix + "user-" + userID
}

func memberKey(guildID, userID string) string {
	return redisKeyPrefix + "guild-" + guildID + ":member-" + userID
}

func guildKey(guildID string) string {
	return redisKeyPrefix + "guild-" + guildID
}

func guildMembersSetKey(guildID string) string {
	return redisKeyPrefix + "guild-" + guildID + ":members"
}

func guildChannelsSetKey(guildID string) string {
	return redisKeyPrefix + "guild-" + guildID + ":channels"
}

func guildRolesSetKey(guildID string) string {
	return redisKeyPrefix + "guild-" + guildID + ":roles"
}

func guildEmojiSetKey(guildID string) string {
	return redisKeyPrefix + "guild-" + guildID + ":emoji"
}

func channelKey(channelID string) string {
	return redisKeyPrefix + "channel-" + channelID
}

func roleKey(guildID, roleID string) string {
	return redisKeyPrefix + "guild-" + guildID + ":role-" + roleID
}

func emojiKey(guildID, emojiID string) string {
	return redisKeyPrefix + "guild-" + guildID + ":emoji-" + emojiID
}

func messagesListKey(channelID string) string {
	return redisKeyPrefix + "channel-" + channelID + ":messages"
}

func updateStateObject(client *redis.Client, key string, object interface{}) error {
	marshalled, err := jsoniter.Marshal(object)
	if err != nil {
		return err
	}

	err = client.Set(key, marshalled, stateExpire).Err()
	return err
}

func deleteStateObject(client *redis.Client, key string) error {
	err := client.Del(key).Err()
	return err
}

func readStateObject(client *redis.Client, key string) (data []byte, err error) {
	data, err = client.Get(key).Bytes()

	return data, err
}

func addToStateSet(client *redis.Client, key string, items ...string) (err error) {
	interfaceItems := make([]interface{}, 0)
	for _, item := range items {
		interfaceItems = append(interfaceItems, item)
	}
	if len(interfaceItems) == 0 {
		return
	}

	err = client.SAdd(key, interfaceItems...).Err()
	return err
}

func removeFromStateSet(client *redis.Client, key, item string) (err error) {
	err = client.SRem(key, item).Err()
	return err
}

func readStateSet(client *redis.Client, key string) (items []string, err error) {
	items, err = client.SMembers(key).Result()
	return items, err
}

func addToStateList(client *redis.Client, key string, items ...string) (err error) {
	interfaceItems := make([]interface{}, 0)
	for _, item := range items {
		interfaceItems = append(interfaceItems, item)
	}
	if len(interfaceItems) == 0 {
		return
	}

	err = client.LPush(key, interfaceItems...).Err()
	return err
}

func trimStateList(client *redis.Client, key string, limit int64) (err error) {
	err = client.LTrim(key, 0, limit).Err()
	return err
}

func readStateList(client *redis.Client, key string) (items []string, err error) {
	items, err = client.LRange(key, 0, -1).Result()
	return items, err
}
