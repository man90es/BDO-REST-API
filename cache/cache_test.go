package cache

import (
	"testing"
	"time"

	"bdo-rest-api/config"
)

func init() {
	config.SetCacheTTL(time.Second)
}

func TestCache(t *testing.T) {
	// Create a cache instance for testing
	testCache := newCache[string]()

	// Test AddRecord and GetRecord
	keys := []string{"key1", "key2"}
	data := "test data"
	status := 200
	taskId := "task-id"

	date, expires := testCache.AddRecord(keys, data, status, taskId)

	// Validate AddRecord results
	if date == "" || expires == "" {
		t.Error("AddRecord should return non-empty date and expires values")
	}

	// Test GetRecord for an existing record
	returnedData, returnedStatus, returnedDate, returnedExpires, found := testCache.GetRecord(keys)

	if !found {
		t.Error("GetRecord should find the record")
	}

	// Validate GetRecord results
	if returnedData != data || returnedStatus != status || returnedDate == "" || returnedExpires == "" {
		t.Error("GetRecord returned unexpected values")
	}

	// Test GetItemCount
	itemCount := testCache.GetItemCount()
	if itemCount != 1 {
		t.Errorf("GetItemCount should return 1, but got %d", itemCount)
	}

	// Sleep for a while to allow the cache entry to expire
	time.Sleep(2 * time.Second)

	// Test GetRecord for an expired record
	_, _, _, _, found = testCache.GetRecord(keys)

	if found {
		t.Error("GetRecord should not find an expired record")
	}

	// Test GetItemCount after expiration
	itemCount = testCache.GetItemCount()
	if itemCount != 0 {
		t.Errorf("GetItemCount should return 0 after expiration, but got %d", itemCount)
	}
}
