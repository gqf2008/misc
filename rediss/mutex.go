package rediss

import (
	"fmt"
	"log"
	"strings"
	"time"
)

//NewMutex ....
func NewMutex(prefix string) *Mutex {
	_init()
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
		if err != nil {
			log.Println("mutex", key, err)
			time.Sleep(interval)
			continue
		}
		if !ok {
			log.Println("mutex", key, "false")
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
