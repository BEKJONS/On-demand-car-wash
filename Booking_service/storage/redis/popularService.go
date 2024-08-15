package Redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func IncrementServiceOrderCount(redisClient *redis.Client, serviceID string) error {
	return redisClient.ZIncrBy(context.Background(), "popular_services", 1, serviceID).Err()
}

func GetPopularServices(redisClient *redis.Client, top int64) ([]redis.Z, error) {
	return redisClient.ZRevRangeWithScores(context.Background(), "popular_services", 0, top-1).Result()
}
