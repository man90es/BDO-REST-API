package cache

import (
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

var redisClient *redis.Client

func initRedisClient(url string) {
	opts, err := redis.ParseURL(url)

	if err != nil {
		panic(err)
	}

	// Explicitly disable maintenance notifications
	// This prevents the client from sending CLIENT MAINT_NOTIFICATIONS ON
	opts.MaintNotificationsConfig = &maintnotifications.Config{
		Mode: maintnotifications.ModeDisabled,
	}

	redisClient = redis.NewClient(opts)
}
