package delaytask

import (
	"log"
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
	Addr         string `default:"redis://127.0.0.1:6379/1?poolsize=200&retries=3&pool_timeout=30"`
	User         string
	Passwd       string
	MaxConns     int `default:"200"`
	MaxIdleConns int `default:"50"`
}

//NewDelayTask ....
func NewDelayTask(name string) *DelayTask {
	once.Do(func() {
		conf := config{}
		err := misc.Fill("DELAYTASK", &conf)
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
	return &DelayTask{
		name: name,
	}
}

//DelayTask ....
type DelayTask struct {
	name      string
	poolSize  int
	pool      *misc.WorkerPool
	afterFunc func(interface{}, time.Duration)
}

//WithAfterFunc ....
func (t *DelayTask) WithAfterFunc(f func(interface{}, time.Duration)) *DelayTask {
	t.afterFunc = f
	return t
}

//WithPoolSize ....
func (t *DelayTask) WithPoolSize(ps int) *DelayTask {
	t.poolSize = ps
	return t
}

//Start ....
func (t *DelayTask) Start() {
	t.pool = &misc.WorkerPool{
		MaxWorkersCount:       t.poolSize,
		MaxIdleWorkerDuration: time.Second * 60 * 5,
	}
	t.pool.Start()
}

//Stop ....
func (t *DelayTask) Stop() {
	if t.pool != nil {
		t.pool.Stop()
	}
}

//Add ....
func (t *DelayTask) Add(task interface{}, deadline time.Duration) error {
	return cli.ZAdd(t.name, redis.Z{Member: task, Score: float64(deadline)}).Err()
}

func (t *DelayTask) loop() {
	defer func() {
		log.Println("延迟任务", t.name, "退出了")
	}()
	if t.afterFunc == nil {
		log.Println("没有设置延迟任务函数")
		return
	}
	for {
		ret := cli.ZRangeByScoreWithScores(t.name, redis.ZRangeBy{
			Min:    "0",
			Max:    strconv.Itoa(int(time.Now().Unix())),
			Offset: 0,
			Count:  1000,
		})
		if err := ret.Err(); err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}
		for _, val := range ret.Val() {
			res := cli.ZRem(t.name, val.Member)
			if err := res.Err(); err != nil {
				log.Println(err)
				break
			}
			if res.Val() == 0 {
				continue
			}
			for {
				if t.pool.Serve(func() {
					t.afterFunc(val.Member, time.Duration(val.Score))
				}) {
					break
				}
				time.Sleep(time.Millisecond * 500)
			}
		}
		time.Sleep(time.Second)
	}
}
