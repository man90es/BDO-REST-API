package cache

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	goCache "github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"

	"bdo-rest-api/models"
	"bdo-rest-api/utils"
)

type CacheEntry[T any] struct {
	Data   T         `json:"data"`
	Date   time.Time `json:"date"`
	Status int       `json:"status"`
}

type Cache[T any] interface {
	AddRecord(keys []string, data T, status int, taskId string) (date string, expires string)
	GetRecord(keys []string) (data T, status int, date string, expires string, found bool)
	GetItemCount() int
	GetKeys() []string
	GetValues() []CacheEntry[T]
}

func joinKeys(keys []string) string {
	return strings.Join(keys, ",")
}

type cache[T any] struct {
	internalCache *goCache.Cache
}

func newMemoryCache[T any]() *cache[T] {
	cacheTTL := viper.GetDuration("cachettl")

	return &cache[T]{
		internalCache: goCache.New(cacheTTL, min(time.Hour, cacheTTL)),
	}
}

func (c *cache[T]) AddRecord(keys []string, data T, status int, taskId string) (date string, expires string) {
	cacheTTL := viper.GetDuration("cachettl")
	entry := CacheEntry[T]{
		Data:   data,
		Date:   time.Now(),
		Status: status,
	}

	c.internalCache.Add(joinKeys(keys), entry, cacheTTL)

	return utils.FormatDateForHeaders(entry.Date), utils.FormatDateForHeaders(entry.Date.Add(cacheTTL))
}

func (c *cache[T]) GetRecord(keys []string) (data T, status int, date string, expires string, found bool) {
	cacheTTL := viper.GetDuration("cachettl")
	anyEntry, found := c.internalCache.Get(joinKeys(keys))

	if !found {
		return
	}

	entry := anyEntry.(CacheEntry[T])

	return entry.Data, entry.Status, utils.FormatDateForHeaders(entry.Date), utils.FormatDateForHeaders(entry.Date.Add(cacheTTL)), true
}

func (c *cache[T]) GetItemCount() int {
	return c.internalCache.ItemCount()
}

func (c *cache[T]) GetKeys() []string {
	return maps.Keys(c.internalCache.Items())
}

func (c *cache[T]) GetValues() []CacheEntry[T] {
	items := c.internalCache.Items()
	result := make([]CacheEntry[T], 0, len(items))

	for _, item := range items {
		result = append(result, item.Object.(CacheEntry[T]))
	}

	return result
}

type redisCache[T any] struct {
	ctx       context.Context
	namespace string
}

func newRedisCache[T any](namespace string) *redisCache[T] {
	return &redisCache[T]{
		ctx:       context.Background(),
		namespace: namespace + ":",
	}
}

func (c *redisCache[T]) AddRecord(keys []string, data T, status int, taskId string) (string, string) {
	cacheTTL := viper.GetDuration("cachettl")

	entry := CacheEntry[T]{
		Data:   data,
		Date:   time.Now(),
		Status: status,
	}

	b, _ := json.Marshal(entry)
	redisClient.Set(c.ctx, c.namespace+joinKeys(keys), b, cacheTTL)

	return utils.FormatDateForHeaders(entry.Date), utils.FormatDateForHeaders(entry.Date.Add(cacheTTL))
}

func (c *redisCache[T]) GetRecord(keys []string) (data T, status int, date string, expires string, found bool) {
	cacheTTL := viper.GetDuration("cachettl")

	val, err := redisClient.Get(c.ctx, c.namespace+joinKeys(keys)).Bytes()
	if err != nil {
		return
	}

	var entry CacheEntry[T]
	if err := json.Unmarshal(val, &entry); err != nil {
		return
	}

	return entry.Data, entry.Status, utils.FormatDateForHeaders(entry.Date), utils.FormatDateForHeaders(entry.Date.Add(cacheTTL)), true
}

func (c *redisCache[T]) GetItemCount() int {
	keys, err := redisClient.Keys(c.ctx, c.namespace+"*").Result()
	if err != nil {
		return 0
	}
	return len(keys)
}

func (c *redisCache[T]) GetKeys() []string {
	keys, _ := redisClient.Keys(c.ctx, c.namespace+"*").Result()

	// Remove namespace from keys
	for i, k := range keys {
		keys[i] = strings.TrimPrefix(k, c.namespace)
	}

	return keys
}

func (c *redisCache[T]) GetValues() []CacheEntry[T] {
	keys, _ := redisClient.Keys(c.ctx, c.namespace+"*").Result()
	result := make([]CacheEntry[T], 0, len(keys))

	for _, k := range keys {
		val, err := redisClient.Get(c.ctx, k).Bytes()
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

var GuildProfiles Cache[models.GuildProfile]
var GuildSearch Cache[[]models.GuildProfile]
var Profiles Cache[models.Profile]
var ProfileSearch Cache[[]models.Profile]

func InitCache() {
	if redisURI := viper.GetString("redis"); redisURI != "" {
		initRedisClient(redisURI)

		GuildProfiles = newRedisCache[models.GuildProfile]("gpc")
		GuildSearch = newRedisCache[[]models.GuildProfile]("gsc")
		Profiles = newRedisCache[models.Profile]("pc")
		ProfileSearch = newRedisCache[[]models.Profile]("psc")
	} else {
		GuildProfiles = newMemoryCache[models.GuildProfile]()
		GuildSearch = newMemoryCache[[]models.GuildProfile]()
		Profiles = newMemoryCache[models.Profile]()
		ProfileSearch = newMemoryCache[[]models.Profile]()
	}
}
