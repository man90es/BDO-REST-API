package cache

import (
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

func newRedisClient(url string) (*redis.Client, error) {
	opts, err := redis.ParseURL(url)

	if err != nil {
		return nil, err
	}

	// Explicitly disable maintenance notifications
	// This prevents the client from sending CLIENT MAINT_NOTIFICATIONS ON
	opts.MaintNotificationsConfig = &maintnotifications.Config{
		Mode: maintnotifications.ModeDisabled,
	}

	return redis.NewClient(opts), nil
}
