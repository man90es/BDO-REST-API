package cache

import (
	"strings"
	"time"

	goCache "github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	messagebus "github.com/vardius/message-bus"
	"golang.org/x/exp/maps"

	"bdo-rest-api/models"
	"bdo-rest-api/utils"
)

type CacheEntry[T any] struct {
	Data   T         `json:"data"`
	Date   time.Time `json:"date"`
	Status int       `json:"status"`
}

type cache[T any] struct {
	Bus           messagebus.MessageBus
	internalCache *goCache.Cache
}

func joinKeys(keys []string) string {
	return strings.Join(keys, ",")
}

func newCache[T any]() *cache[T] {
	cacheTTL := viper.GetDuration("cachettl")

	return &cache[T]{
		Bus:           messagebus.New(100), // Idk what buffer size is optimal
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
	c.Bus.Publish(taskId, entry)

	return utils.FormatDateForHeaders(entry.Date), utils.FormatDateForHeaders(entry.Date.Add(cacheTTL))
}

func (c *cache[T]) GetRecord(keys []string) (data T, status int, date string, expires string, found bool) {
	cacheTTL := viper.GetDuration("cachettl")
	anyEntry, found := c.internalCache.Get(joinKeys(keys))

	if !found {
		return
	}

	entry := anyEntry.(CacheEntry[T])

	return entry.Data, entry.Status, utils.FormatDateForHeaders(entry.Date), utils.FormatDateForHeaders(entry.Date.Add(cacheTTL)), found
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

var GuildProfiles = newCache[models.GuildProfile]()
var GuildSearch = newCache[[]models.GuildProfile]()
var Profiles = newCache[models.Profile]()
var ProfileSearch = newCache[[]models.Profile]()
