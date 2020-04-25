package rediss

import (
	"sync"

	"github.com/go-redis/redis"
)

var once = sync.Once{}
var cli *redis.Client

type config struct {
	Addr string `default:"redis://127.0.0.1:6379/1?poolsize=200&retries=3&pool_timeout=30"`
}
