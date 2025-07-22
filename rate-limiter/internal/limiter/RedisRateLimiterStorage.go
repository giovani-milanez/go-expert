package limiter

import (
	"context"
	"encoding/json"
	// "fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRateLimiterStorage struct {
	ctx          context.Context
	redisClient *redis.Client
}

func NewRedisRateLimiterStorage(ctx context.Context, addr string) *RedisRateLimiterStorage {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // No password by default for local development
		DB:       0,  // Default DB
	})
	return &RedisRateLimiterStorage{
		redisClient: redisClient,
		ctx: ctx,
	}
}

func (s *RedisRateLimiterStorage) Flush() error {
	_, err := s.redisClient.FlushAll(s.ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

func (s *RedisRateLimiterStorage) Connect() error {
	_, err := s.redisClient.Ping(s.ctx).Result()
	if err != nil {	
		return err
	}
	return nil
}

func (s *RedisRateLimiterStorage) Get(key string) (*Rate, error) {
	val, err := s.redisClient.Get(s.ctx, key).Result()
	if err == redis.Nil {
		rate := &Rate{Count: 1, FirstSeen: time.Now(), BlockedUntil: time.Now().Add(-1 * time.Hour)}
		err = s.Update(key, rate)
		if err != nil {
			return nil, err
		}
		return rate, nil
	} else if err != nil {
		return nil, err
	}
	var rate Rate
	err = json.Unmarshal([]byte(val), &rate)
	if err != nil {
		return nil, err
	}
	// fmt.Println("Got first seen:", rate.FirstSeen, "for key:", key)
	// fmt.Println("Got blocked until:", rate.BlockedUntil, "for key:", key)
	return &rate, nil
}

func (s *RedisRateLimiterStorage) Update(key string, rate *Rate) error {
	// fmt.Println("Saving first seen:", rate.FirstSeen, "for key:", key)
	// fmt.Println("Saving count:", rate.Count, "for key:", key)
	// fmt.Println("Saving blocked until:", rate.BlockedUntil, "for key:", key)
	data, err := json.Marshal(rate)
	if err != nil {
		return err
	}	
	err = s.redisClient.Set(s.ctx, key, data, 0).Err()
	if err != nil {
		return err
	}
	return nil
}