package middleware

import (
	"time"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"

	"github.com/spf13/viper"
)

func GetRateLimitMiddleware() Middleware {
	var rate = limiter.Rate{
		Limit:  viper.GetInt64("ratelimit"),
		Period: time.Minute,
	}
	var store = memory.NewStore()
	var instance = limiter.New(store, rate, limiter.WithClientIPHeader("CF-Connecting-IP"))

	var middleware = stdlib.NewMiddleware(instance)
	return middleware.Handler
}
