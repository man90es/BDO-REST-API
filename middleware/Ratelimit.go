package middleware

import (
	"time"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

var rate = limiter.Rate{
	Limit:  512,
	Period: time.Minute,
}
var store = memory.NewStore()
var instance = limiter.New(store, rate, limiter.WithClientIPHeader("CF-Connecting-IP"))

var middleware = stdlib.NewMiddleware(instance)
var Ratelimit = middleware.Handler
