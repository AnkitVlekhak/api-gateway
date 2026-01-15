package gateway

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisTokenBucketRateLimiter struct {
	client *redis.Client
	script *redis.Script
}

func NewRedisTokenBucketRateLimiter(client *redis.Client) *RedisTokenBucketRateLimiter {
	return &RedisTokenBucketRateLimiter{
		client: client,
		script: redis.NewScript(tokenBucketLua),
	}
}

const tokenBucketLua = `
-- KEYS[1] = token bucket key
-- ARGV[1] = capacity
-- ARGV[2] = refill_rate (tokens per second)
-- ARGV[3] = current_time (unix seconds)

local bucket = redis.call("HMGET", KEYS[1], "tokens", "last_refill")
local tokens = tonumber(bucket[1])
local last_refill = tonumber(bucket[2])

if tokens == nil then
  tokens = tonumber(ARGV[1])
  last_refill = tonumber(ARGV[3])
end

local delta = tonumber(ARGV[3]) - last_refill
local refill = delta * tonumber(ARGV[2])
tokens = math.min(tonumber(ARGV[1]), tokens + refill)

if tokens < 1 then
  redis.call("HMSET", KEYS[1], "tokens", tokens, "last_refill", ARGV[3])
  redis.call("EXPIRE", KEYS[1], math.ceil(ARGV[1] / ARGV[2]))
  return 0
end

tokens = tokens - 1
redis.call("HMSET", KEYS[1], "tokens", tokens, "last_refill", ARGV[3])
redis.call("EXPIRE", KEYS[1], math.ceil(ARGV[1] / ARGV[2]))

return 1
`

func (r *RedisTokenBucketRateLimiter) Allow(route *Route, req *http.Request) bool {
	policy := route.RateLimitPolicy
	if policy == nil {
		return true
	}

	var identity string

	switch policy.KeyBy {
	case RateLimitKeyIP:
		identity = req.RemoteAddr
	case RateLimitByAPIKey:
		identity = req.Header.Get("X-API-Key")
	}

	if identity == "" {
		return false
	}

	key := fmt.Sprintf("rate_limit:%s:%s", route.Path, identity)

	capacity := policy.Requests
	refillRate := 1.0 / policy.Window.Seconds()
	now := float64(time.Now().Unix())

	res, err := r.script.Run(
		context.Background(),
		r.client,
		[]string{key},
		capacity,
		refillRate,
		now,
	).Int()

	if err != nil {
		return false
	}

	return res == 1
}
