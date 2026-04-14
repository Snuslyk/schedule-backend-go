package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Snuslyk/schedule-backend-go/internal/utils"
)

type ScheduleCache struct {
	redis *RedisClient
}

func NewScheduleCache(redis *RedisClient) *ScheduleCache {
	return &ScheduleCache{redis: redis}
}

func (c *ScheduleCache) Get(ctx context.Context, key string, dest any) bool {
	if c.redis == nil || c.redis.client == nil {
		return false
	}
	val, err := c.redis.client.Get(ctx, key).Bytes()
	if err != nil {
		return false
	}
	return json.Unmarshal(val, dest) == nil
}

func (c *ScheduleCache) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	if c.redis == nil || c.redis.client == nil {
		return nil
	}
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	if ttl <= 0 {
		ttl = utils.DurationUntilEndOfDay(time.Now())
	}
	return c.redis.client.Set(ctx, key, data, ttl).Err()
}

func (c *ScheduleCache) Del(ctx context.Context, keys ...string) error {
	if c.redis == nil || c.redis.client == nil || len(keys) == 0 {
		return nil
	}
	return c.redis.client.Del(ctx, keys...).Err()
}

func (c *ScheduleCache) InvalidateGroup(ctx context.Context, groupID int, groupName string) {
	if c.redis == nil || c.redis.client == nil {
		return
	}
	keys := []string{
		fmt.Sprintf("group_week:%d", groupID),
		fmt.Sprintf("group_week:%s", groupName),
	}
	_ = c.Del(ctx, keys...)

	iter := c.redis.client.Scan(ctx, 0, fmt.Sprintf("schedule:group:*:%d:*", groupID), 0).Iterator()
	for iter.Next(ctx) {
		_ = c.redis.client.Del(ctx, iter.Val())
	}
}

func (c *ScheduleCache) InvalidateTeacher(ctx context.Context, teacherID int) {
	if c.redis == nil || c.redis.client == nil {
		return
	}
	keys := []string{
		fmt.Sprintf("teacher_week:%d", teacherID),
	}
	_ = c.Del(ctx, keys...)

	iter := c.redis.client.Scan(ctx, 0, fmt.Sprintf("schedule:teacher:*:%d:*", teacherID), 0).Iterator()
	for iter.Next(ctx) {
		_ = c.redis.client.Del(ctx, iter.Val())
	}
}
