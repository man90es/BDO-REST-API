package cache

import (
	"net/http"
	"strings"
	"sync"
	"time"

	goCache "github.com/patrickmn/go-cache"
	messagebus "github.com/vardius/message-bus"
	"golang.org/x/exp/maps"

	"bdo-rest-api/config"
	"bdo-rest-api/models"
	"bdo-rest-api/utils"
)

type CacheEntry[T any] struct {
	Data   T
	Date   time.Time
	Status int
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

func (c *cache[T]) AddRecord(keys []string, data T, status int, taskId string) (date string, expires string) {
	ttl := config.GetCacheTTL()
	entry := CacheEntry[T]{
		Data:   data,
		Date:   time.Now(),
		Status: status,
	}

	c.internalCache.Add(joinKeys(keys), entry, ttl)
	c.Bus.Publish(taskId, entry)

	return utils.FormatDateForHeaders(entry.Date), utils.FormatDateForHeaders(entry.Date.Add(ttl))
}

func (c *cache[T]) SignalBypassCache(status int, taskId string) {
	var data T
	entry := CacheEntry[T]{
		Data:   data,
		Date:   time.Now(),
		Status: status,
	}

	c.Bus.Publish(taskId, entry)
}

func (c *cache[T]) GetRecord(keys []string) (data T, status int, date string, expires string, found bool) {
	cacheTTL := config.GetCacheTTL()
	anyEntry, found := c.internalCache.Get(joinKeys(keys))

	if !found {
		return
	}

	entry := anyEntry.(CacheEntry[T])

	return entry.Data, entry.Status, utils.FormatDateForHeaders(entry.Date), utils.FormatDateForHeaders(entry.Date.Add(cacheTTL)), found
}

func (c *cache[T]) WaitForRecord(taskId string) (data T, status int, date string, expires string) {
	var wg sync.WaitGroup
	wg.Add(1)

	c.Bus.Subscribe(taskId, func(entry CacheEntry[T]) {
		data = entry.Data
		status = entry.Status
		date = utils.FormatDateForHeaders(entry.Date)

		if entry.Status == http.StatusInternalServerError {
			expires = date
		} else if entry.Status == http.StatusServiceUnavailable {
			expires = utils.FormatDateForHeaders(entry.Date.Add(config.GetMaintenanceStatusTTL()))
		} else {
			expires = utils.FormatDateForHeaders(entry.Date.Add(config.GetCacheTTL()))
		}

		wg.Done()
	})

	wg.Wait()
	return
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
