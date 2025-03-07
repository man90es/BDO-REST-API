package config

import (
	"fmt"
	"sync"
	"time"
)

type config struct {
	cacheTTL       time.Duration
	maintenanceTTL time.Duration
	mongoDB        string
	mu             sync.RWMutex
	port           int
	proxyList      []string
	rateLimit      int64
	verbosity      bool
}

var instance *config
var once sync.Once

func getInstance() *config {
	once.Do(func() {
		instance = &config{
			cacheTTL:       3 * time.Hour,
			maintenanceTTL: 5 * time.Minute,
			port:           8001,
			proxyList:      nil,
			rateLimit:      512,
			verbosity:      false,
		}
	})
	return instance
}

func SetCacheTTL(ttl time.Duration) {
	getInstance().mu.Lock()
	defer getInstance().mu.Unlock()
	getInstance().cacheTTL = ttl
}

func GetCacheTTL() time.Duration {
	getInstance().mu.RLock()
	defer getInstance().mu.RUnlock()
	return getInstance().cacheTTL
}

func SetMaintenanceStatusTTL(ttl time.Duration) {
	getInstance().mu.Lock()
	defer getInstance().mu.Unlock()
	getInstance().maintenanceTTL = ttl
}

func GetMaintenanceStatusTTL() time.Duration {
	getInstance().mu.RLock()
	defer getInstance().mu.RUnlock()
	return getInstance().maintenanceTTL
}

func SetPort(port int) {
	getInstance().mu.Lock()
	defer getInstance().mu.Unlock()
	getInstance().port = port
}

func GetPort() int {
	getInstance().mu.RLock()
	defer getInstance().mu.RUnlock()
	return getInstance().port
}

func SetProxyList(proxies []string) {
	getInstance().mu.Lock()
	defer getInstance().mu.Unlock()
	getInstance().proxyList = proxies
}

func GetProxyList() []string {
	getInstance().mu.RLock()
	defer getInstance().mu.RUnlock()
	return getInstance().proxyList
}

func SetVerbosity(verbosity bool) {
	getInstance().mu.Lock()
	defer getInstance().mu.Unlock()
	getInstance().verbosity = verbosity
}

func GetVerbosity() bool {
	getInstance().mu.RLock()
	defer getInstance().mu.RUnlock()
	return getInstance().verbosity
}

func SetRateLimit(rateLimit int64) {
	getInstance().mu.Lock()
	defer getInstance().mu.Unlock()
	getInstance().rateLimit = rateLimit
}

func GetRateLimit() int64 {
	getInstance().mu.RLock()
	defer getInstance().mu.RUnlock()
	return getInstance().rateLimit
}

func SetMongoDB(mongoDB string) {
	getInstance().mu.Lock()
	defer getInstance().mu.Unlock()
	getInstance().mongoDB = mongoDB
}

func GetMongoDB() string {
	getInstance().mu.RLock()
	defer getInstance().mu.RUnlock()
	return getInstance().mongoDB
}

func SprintfConfig() string {
	return fmt.Sprintf("\tPort:\t\t%v\n", GetPort()) +
		fmt.Sprintf("\tProxies:\t%v\n", GetProxyList()) +
		fmt.Sprintf("\tVerbosity:\t%v\n", GetVerbosity()) +
		fmt.Sprintf("\tCache TTL:\t%v\n", GetCacheTTL()) +
		fmt.Sprintf("\tRate limit:\t%v/min\n", GetRateLimit()) +
		fmt.Sprintf("\tMongoDB:\t%v", GetMongoDB())
}
