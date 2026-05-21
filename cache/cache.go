package cache

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"

	"bdo-rest-api/models"
	"bdo-rest-api/utils"
)

type CacheEntry[T any] struct {
	ClientID    string        `json:"clientId,omitempty"`
	Data        T             `json:"data"`
	Date        time.Time     `json:"date"`
	InstanceID  string        `json:"instanceId,omitempty"`
	Status      int           `json:"status"`
	TimeElapsed time.Duration `json:"timeElapsed,omitempty"`
}

type Cache[T any] interface {
	AddRecord(keys []string, data T, status int, clientID string, startTime time.Time) (date string, expires string)
	GetRecord(keys []string) (data T, status int, date string, expires string, found bool)
	GetItemCount() int
	GetKeys() []string
	GetValues() []CacheEntry[T]
}

func joinKeys(keys []string) string {
	return strings.Join(keys, ",")
}

type redisCache[T any] struct {
	client    *redis.Client
	ctx       context.Context
	namespace string
	ttl       time.Duration
}

func newRedisCache[T any](client *redis.Client, namespace string) *redisCache[T] {
	return &redisCache[T]{
		client:    client,
		ctx:       context.Background(),
		namespace: namespace + ":",
		ttl:       viper.GetDuration("cachettl"),
	}
}

func (c *redisCache[T]) AddRecord(keys []string, data T, status int, clientID string, startTime time.Time) (date, expires string) {
	entry := CacheEntry[T]{
		ClientID:    clientID,
		Data:        data,
		Date:        time.Now(),
		InstanceID:  viper.GetString("instanceid"),
		Status:      status,
		TimeElapsed: time.Since(startTime),
	}

	b, _ := json.Marshal(entry)
	c.client.Set(c.ctx, c.namespace+joinKeys(keys), b, c.ttl)

	return utils.FormatDateForHeaders(entry.Date), utils.FormatDateForHeaders(entry.Date.Add(c.ttl))
}

func (c *redisCache[T]) GetRecord(keys []string) (data T, status int, date string, expires string, found bool) {
	val, err := c.client.Get(c.ctx, c.namespace+joinKeys(keys)).Bytes()
	if err != nil {
		return
	}

	var entry CacheEntry[T]
	if err := json.Unmarshal(val, &entry); err != nil {
		return
	}

	ttl := c.client.TTL(c.ctx, c.namespace+joinKeys(keys)).Val()

	return entry.Data, entry.Status, utils.FormatDateForHeaders(entry.Date), utils.FormatDateForHeaders(time.Now().Add(ttl)), true
}

func (c *redisCache[T]) GetItemCount() int {
	keys, err := c.client.Keys(c.ctx, c.namespace+"*").Result()
	if err != nil {
		return 0
	}
	return len(keys)
}

func (c *redisCache[T]) GetKeys() []string {
	keys, _ := c.client.Keys(c.ctx, c.namespace+"*").Result()

	// Remove namespace from keys
	for i, k := range keys {
		keys[i] = strings.TrimPrefix(k, c.namespace)
	}

	return keys
}

func (c *redisCache[T]) GetValues() []CacheEntry[T] {
	keys, _ := c.client.Keys(c.ctx, c.namespace+"*").Result()
	result := make([]CacheEntry[T], 0, len(keys))

	for _, k := range keys {
		val, err := c.client.Get(c.ctx, k).Bytes()
		if err != nil {
			continue
		}

		var entry CacheEntry[T]
		if err := json.Unmarshal(val, &entry); err != nil {
			continue
		}
		result = append(result, entry)
	}

	return result
}

var (
	GuildProfiles Cache[models.GuildProfile]
	GuildSearch   Cache[[]models.GuildProfile]
	Profiles      Cache[models.Profile]
	ProfileSearch Cache[[]models.Profile]
)

func InitCache() {
	redisClient, err := newRedisClient(viper.GetString("redis"))
	if err != nil {
		panic(err)
	}

	GuildProfiles = newRedisCache[models.GuildProfile](redisClient, "gpc")
	GuildSearch = newRedisCache[[]models.GuildProfile](redisClient, "gsc")
	Profiles = newRedisCache[models.Profile](redisClient, "pc")
	ProfileSearch = newRedisCache[[]models.Profile](redisClient, "psc")
}
