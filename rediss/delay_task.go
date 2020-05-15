package rediss

import (
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/gqf2008/misc"
)

//NewDelayTask ....
func NewDelayTask(name string) *DelayTask {
	_init()
	return &DelayTask{
		name:     name,
		poolSize: 200,
		afterFunc: func(task interface{}, at time.Duration) {
			log.Printf("At: %d Task: %+v\n", at, task)
		},
		stop: make(chan struct{}, 10),
	}
}

//DelayTask ....
type DelayTask struct {
	name      string
	poolSize  int
	pool      *misc.WorkerPool
	afterFunc func(interface{}, time.Duration)
	stop      chan struct{}
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
	go t.loop()
}

//Stop ....
func (t *DelayTask) Stop() {
	t.stop <- struct{}{}
	if t.pool != nil {
		t.pool.Stop()
	}
}

//Add ....
func (t *DelayTask) Add(task interface{}, deadline time.Duration) error {
	return cli.ZAdd(t.name, redis.Z{Member: task, Score: float64(deadline)}).Err()
}

//Remove ....
func (t *DelayTask) Remove(task interface{}) error {
	return cli.ZRem(t.name, task).Err()
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
		select {
		case <-t.stop:
			return
		default:
		}
		ret := cli.ZRangeByScoreWithScores(t.name, redis.ZRangeBy{
			Min:    "0",
			Max:    strconv.Itoa(int(time.Now().Unix())),
			Offset: 0,
			Count:  1000,
		})
		if err := ret.Err(); err != nil {
			log.Println("delay_task", err)
			time.Sleep(time.Second)
			continue
		}
		for _, val := range ret.Val() {
			for {
				if t.pool.Serve(func() {
					t.afterFunc(val.Member, time.Duration(val.Score))
				}) {
					if err := cli.ZRem(t.name, val.Member).Err(); err != nil {
						log.Println("zrem delay_task", t.name, val.Member, err)
					}
					break
				}
				log.Println("delay_task not enough worker")
				time.Sleep(time.Second)
			}
		}
		time.Sleep(time.Second)
	}
}
