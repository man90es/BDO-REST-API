package cache

import (
	"strings"
	"time"

	goCache "github.com/patrickmn/go-cache"
)

type cacheEntry[T any] struct {
	data   T
	date   time.Time
	status int
}

type cache[T any] struct {
	internalCache *goCache.Cache
}

func joinKeys(keys []string) string {
	return strings.Join(keys, ",")
}

func formatDate(date *time.Time) string {
	return date.Format(time.RFC1123Z)
}

func NewCache[T any]() *cache[T] {
	return &cache[T]{
		internalCache: goCache.New(5*time.Minute, 10*time.Minute),
	}
}

func (c *cache[T]) AddRecord(keys []string, data T, status int) (date string, expires string) {
	entry := cacheEntry[T]{
		data:   data,
		date:   time.Now(),
		status: status,
	}

	c.internalCache.Add(joinKeys(keys), entry, goCache.DefaultExpiration)
	expirationDate := entry.date.Add(5 * time.Minute)

	return formatDate(&entry.date), formatDate(&expirationDate)
}

func (c *cache[T]) GetRecord(keys []string) (data T, status int, date string, expires string, found bool) {
	var anyEntry interface{}
	var expirationDate time.Time

	anyEntry, expirationDate, found = c.internalCache.GetWithExpiration(joinKeys(keys))

	if !found {
		return
	}

	entry := anyEntry.(cacheEntry[T])

	return entry.data, entry.status, formatDate(&entry.date), formatDate(&expirationDate), found
}

func (c *cache[T]) GetItemCount() int {
	return c.internalCache.ItemCount()
}
