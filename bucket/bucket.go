package bucket

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func key(suffix string) string {
	return "cacophony:bucket:" + suffix
}

// TODO: add bucket limit?

// Add adds an item to a leaky bucket, and returns the amount of items in the leaky bucket
// Source: https://engineering.classdojo.com/blog/2015/02/06/rolling-rate-limiter/
func Add(
	client *redis.Client,
	keySuffix string,
	decay time.Duration,
) (int, error) {
	key := key(keySuffix)

	now := time.Now().UTC().UnixNano()

	interval := now - decay.Nanoseconds()

	pipeline := client.TxPipeline()

	// remove all items in the bucket before interval
	pipeline.ZRemRangeByScore(key, "0", strconv.FormatInt(interval, 10))
	// add item to bucket
	pipeline.ZAdd(key, redis.Z{
		Score:  float64(now),
		Member: now,
	})
	// count number of items in the bucket
	items := pipeline.ZRangeWithScores(key, 0, -1)
	// expire bucket after interval to save space
	pipeline.Expire(key, decay)

	_, err := pipeline.Exec()
	if err != nil {
		return 0, err
	}

	return len(items.Val()), nil
}
