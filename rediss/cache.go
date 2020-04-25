package rediss

import (
	"time"

	"encoding/json"

	"github.com/go-redis/redis"
)

//NewCache ....
func NewCache(name string) *Cache {
	_init()
	return &Cache{name}
}

//Cache .....
type Cache struct {
	name string
}

//Set ....
func (c *Cache) Set(k string, v interface{}, exp time.Duration) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return cli.Set(k, b, exp).Err()
}

//Delete ....
func (c *Cache) Delete(keys ...string) error {
	return cli.Del(keys...).Err()
}

//Get .....
func (c *Cache) Get(k string, v interface{}) (bool, error) {
	b, err := cli.Get(k).Bytes()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, json.Unmarshal(b, v)
}

//Exist ....
func (c *Cache) Exist(k string) (bool, error) {
	err := cli.Get(k).Err()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
