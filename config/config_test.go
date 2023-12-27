package config

import (
	"reflect"
	"testing"
	"time"
)

func TestConfigSingleton(t *testing.T) {
	// Test that the instance is the same when retrieved multiple times
	instance1 := getInstance()
	instance2 := getInstance()

	if instance1 != instance2 {
		t.Error("Multiple instances of Config were created")
	}
}

func TestConfigOperations(t *testing.T) {
	// Set values
	ttl := 5 * time.Minute
	SetCacheTTL(ttl)
	SetPort(8081)
	proxies := []string{"proxy3", "proxy4"}
	SetProxyList(proxies)
	SetVerbosity(true)

	// Test updated values
	if GetCacheTTL() != ttl {
		t.Error("Failed to set Cache TTL")
	}
	if GetPort() != 8081 {
		t.Error("Failed to set Port")
	}
	if !reflect.DeepEqual(GetProxyList(), proxies) {
		t.Error("Failed to set Proxy List")
	}
	if !GetVerbosity() {
		t.Error("Failed to set Verbosity")
	}
}

func TestConcurrency(t *testing.T) {
	// Set values concurrently
	go func() {
		ttl := 5 * time.Minute
		SetCacheTTL(ttl)
	}()
	go func() {
		proxies := []string{"proxy5", "proxy6"}
		SetProxyList(proxies)
	}()
	go func() {
		SetVerbosity(true)
	}()

	// Allow goroutines to finish
	time.Sleep(100 * time.Millisecond)

	// Test updated values
	if GetCacheTTL() == 0 {
		t.Error("Failed to set Cache TTL concurrently")
	}
	if !reflect.DeepEqual(GetProxyList(), []string{"proxy5", "proxy6"}) {
		t.Error("Failed to set Proxy List concurrently")
	}
	if !GetVerbosity() {
		t.Error("Failed to set Verbosity concurrently")
	}
}
