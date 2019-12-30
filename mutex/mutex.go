package mutex

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gqf2008/misc"

	"github.com/go-redis/redis"
)

var once = sync.Once{}
var cli *redis.Client

type config struct {
	Redis Redis
}

//Redis ....
type Redis struct {
	Addr string `default:"redis://127.0.0.1:6379/0?poolsize=200&retries=3&pool_timeout=30"`
}

//NewMutex ....
func NewMutex(prefix string) *Mutex {
	once.Do(func() {
		conf := config{}
		err := misc.Fill("MUTEX", &conf)
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
	return &Mutex{prefix}
}

//Mutex ....
type Mutex struct {
	// key string
	prefix string
	//ttl time.Duration
}

var intervalseq = []time.Duration{500 * time.Millisecond, 500 * time.Millisecond, time.Second, time.Second, time.Second, time.Second, time.Second, time.Second, time.Second, time.Second, time.Second, time.Second, time.Second, time.Second, time.Second, time.Second}

//Lock ....
func (m *Mutex) Lock(key string, ttl time.Duration) error {
	var sb strings.Builder
	sb.WriteString(m.prefix)
	sb.WriteString(key)
	key = sb.String()
	for _, interval := range intervalseq {
		ok, err := cli.SetNX(key, "lock", ttl).Result()
		if !ok || err != nil {
			log.Println(err)
			time.Sleep(interval)
			continue
		}
		return nil
	}
	return fmt.Errorf("lock %s %d timeout", key, ttl)
}

//Unlock ....
func (m *Mutex) Unlock(key string) error {
	var sb strings.Builder
	sb.WriteString(m.prefix)
	sb.WriteString(key)
	key = sb.String()
	return cli.Del(key).Err()
}
