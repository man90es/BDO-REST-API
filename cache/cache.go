package cache

import (
	"strings"
	"time"

	goCache "github.com/patrickmn/go-cache"
	messagebus "github.com/vardius/message-bus"
	"golang.org/x/exp/maps"

	"bdo-rest-api/config"
	"bdo-rest-api/models"
	"bdo-rest-api/utils"
)

type CacheEntry[T any] struct {
	Data    T
	Date    time.Time
	Expires time.Time // FIXME: Expiration date is also stored in the cache library but it's harder to reach
	Status  int
}

type cache[T any] struct {
	Bus           messagebus.MessageBus
	internalCache *goCache.Cache
}

func joinKeys(keys []string) string {
	return strings.Join(keys, ",")
}

func newCache[T any]() *cache[T] {
	cacheTTL := config.GetCacheTTL()

	return &cache[T]{
		Bus:           messagebus.New(100), // Idk what buffer size is optimal
		internalCache: goCache.New(cacheTTL, min(time.Hour, cacheTTL)),
	}
}

func (c *cache[T]) AddRecord(keys []string, data T, status int) (date string, expires string) {
	cacheTTL := config.GetCacheTTL()
	entry := CacheEntry[T]{
		Data:    data,
		Date:    time.Now(),
		Expires: time.Now().Add(cacheTTL),
		Status:  status,
	}

	c.internalCache.Add(joinKeys(keys), entry, cacheTTL)
	c.Bus.Publish(joinKeys(keys), entry)

	return utils.FormatDateForHeaders(entry.Date), utils.FormatDateForHeaders(entry.Expires)
}

func (c *cache[T]) GetRecord(keys []string) (data T, status int, date string, expires string, found bool) {
	anyEntry, found := c.internalCache.Get(joinKeys(keys))

	if !found {
		return
	}

	entry := anyEntry.(CacheEntry[T])

	return entry.Data, entry.Status, utils.FormatDateForHeaders(entry.Date), utils.FormatDateForHeaders(entry.Expires), found
}

func (c *cache[T]) GetItemCount() int {
	return c.internalCache.ItemCount()
}

func (c *cache[T]) GetKeys() []string {
	return maps.Keys(c.internalCache.Items())
}

var GuildProfiles = newCache[models.GuildProfile]()
var GuildSearch = newCache[[]models.GuildProfile]()
var Profiles = newCache[models.Profile]()
var ProfileSearch = newCache[[]models.Profile]()
