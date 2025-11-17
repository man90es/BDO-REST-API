package cache

import (
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

var redisClient *redis.Client

func initRedisClient(addr string) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       0,
		Password: "",

		// Explicitly disable maintenance notifications
		// This prevents the client from sending CLIENT MAINT_NOTIFICATIONS ON
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})
}
