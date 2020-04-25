package rediss

import (
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/gqf2008/misc"
)

var once = sync.Once{}
var cli *redis.Client

type config struct {
	Redis Redis
}

//Redis ....
type Redis struct {
	Addr string `default:"redis://127.0.0.1:6379/1?poolsize=200&retries=3&pool_timeout=30"`
}

func _init() {
	once.Do(func() {
		conf := config{}
		err := misc.Fill("MISC", &conf)
		if err != nil {
			panic(err)
		}
		URL, err := url.Parse(conf.Redis.Addr)
		if err != nil {
			panic(err)
		}
		var passwd string
		if URL.User != nil {
			passwd, _ = URL.User.Password()
		}
		val := URL.Query()
		poolSize, _ := strconv.ParseUint(val.Get("poolsize"), 10, 64)
		if poolSize == 0 {
			poolSize = 10
		}
		var db = 0
		if len(URL.Path) > 1 {
			db, err = strconv.Atoi(URL.Path[1:])
			if err != nil {
				panic(err)
			}
		}
		client := redis.NewClient(&redis.Options{
			Network:      "tcp",
			Addr:         URL.Host,
			Password:     passwd,
			MaxRetries:   3,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			PoolSize:     int(poolSize),
			DB:           db,
		})
		_, err = client.Ping().Result()
		if err != nil {
			panic(err)
		}
		cli = client
	})
}
