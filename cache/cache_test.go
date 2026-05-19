package cache

import (
	"context"
	"testing"
	"time"

	"github.com/spf13/viper"
)

// Simple helper type for testing
type testStruct struct {
	Value string
}

func init() {
	// Ensure TTL is something predictable
	viper.Set("cachettl", time.Second*10)
	viper.Set("redis", "redis://localhost:6379/0")
}

func getTestRedisCache[T any](t *testing.T, namespace string) *redisCache[T] {
	client, err := newRedisClient(viper.GetString("redis"))
	if err != nil {
		t.Skip("Skipping test: Redis client could not be created")
	}
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Skip("Skipping test: Redis server not reachable")
	}
	client.FlushDB(context.Background())
	return newRedisCache[T](client, namespace)
}

func TestJoinKeys(t *testing.T) {
	keys := []string{"a", "b", "c"}
	expected := "a,b,c"

	if got := joinKeys(keys); got != expected {
		t.Fatalf("joinKeys() = %s; want %s", got, expected)
	}
}

func TestRedisCacheAddAndGetRecord(t *testing.T) {
	c := getTestRedisCache[testStruct](t, "test_add_get")

	data := testStruct{Value: "hello"}
	keys := []string{"key1", "key2"}
	status := 200

	_, _ = c.AddRecord(keys, data, status, "task123")

	gotData, gotStatus, _, _, found := c.GetRecord(keys)
	if !found {
		t.Fatal("Expected record to be found, got not found")
	}

	if gotData.Value != "hello" {
		t.Fatalf("Expected data 'hello', got '%s'", gotData.Value)
	}

	if gotStatus != status {
		t.Fatalf("Expected status %d, got %d", status, gotStatus)
	}
}

func TestRedisCacheMissingRecord(t *testing.T) {
	c := getTestRedisCache[testStruct](t, "test_missing")

	_, _, _, _, found := c.GetRecord([]string{"does", "not", "exist"})
	if found {
		t.Fatal("Expected record NOT to be found")
	}
}

func TestRedisCacheItemCount(t *testing.T) {
	c := getTestRedisCache[testStruct](t, "test_count")

	if c.GetItemCount() != 0 {
		t.Fatal("Expected empty cache")
	}

	c.AddRecord([]string{"a"}, testStruct{"x"}, 200, "task1")
	c.AddRecord([]string{"b"}, testStruct{"y"}, 200, "task2")

	if c.GetItemCount() != 2 {
		t.Fatalf("Expected 2 items, got %d", c.GetItemCount())
	}
}

func TestRedisCacheGetKeys(t *testing.T) {
	c := getTestRedisCache[testStruct](t, "test_keys")

	c.AddRecord([]string{"k1"}, testStruct{"v1"}, 200, "task1")
	c.AddRecord([]string{"k2"}, testStruct{"v2"}, 200, "task2")

	keys := c.GetKeys()

	if len(keys) != 2 {
		t.Fatalf("Expected 2 keys, got %d", len(keys))
	}

	// NOTE: key order in map is not stable
	found1, found2 := false, false
	for _, k := range keys {
		if k == "k1" {
			found1 = true
		}
		if k == "k2" {
			found2 = true
		}
	}

	if !found1 || !found2 {
		t.Fatalf("Expected keys k1 and k2, got %v", keys)
	}
}

func TestRedisCacheGetValues(t *testing.T) {
	c := getTestRedisCache[testStruct](t, "test_values")

	c.AddRecord([]string{"a"}, testStruct{"aaa"}, 200, "task1")
	c.AddRecord([]string{"b"}, testStruct{"bbb"}, 200, "task2")

	values := c.GetValues()

	if len(values) != 2 {
		t.Fatalf("Expected 2 values, got %d", len(values))
	}

	// Verify content
	foundA, foundB := false, false
	for _, v := range values {
		if v.Data.Value == "aaa" {
			foundA = true
		}
		if v.Data.Value == "bbb" {
			foundB = true
		}
	}

	if !foundA || !foundB {
		t.Fatalf("Missing expected values, got %v", values)
	}
}
