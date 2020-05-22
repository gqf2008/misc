package event

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/streadway/amqp"
)

//NewMediatorForAliyun ....
func NewMediatorForAliyun(uid, accessKey, secretKey, addr, vhost string) (*Mediator, error) {
	return NewMediator(buildURLForAliyun(uid, accessKey, secretKey, addr, vhost))
}

//NewMediator ....
func NewMediator(url string) (*Mediator, error) {
	driver := Mediator{
		url:           url,
		prefetchCount: 1,
		funcs:         map[string]HandleFunc{},
	}
	err := driver.connect()
	if err != nil {
		return nil, err
	}
	return &driver, nil
}

//WithQos ....
func (c *Mediator) WithQos(prefetchCount, prefetchSize int, global bool) {
	c.prefetchCount = prefetchCount
	c.prefetchSize = prefetchSize
	c.global = global
}

//Close ....
func (c *Mediator) Close() {
	if c.ch != nil {
		_ = c.ch.Close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
}

func (c *Mediator) connect() error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return err
	}
	c.conn = conn
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	ch.Qos(1, 0, true)
	c.ch = ch
	c.err = make(chan *amqp.Error, 10)
	c.block = make(chan amqp.Blocking, 10)
	conn.NotifyClose(c.err)
	conn.NotifyBlocked(c.block)
	return nil
}

//BindQueue 绑定队列
func (c *Mediator) BindQueue(exchange, queue, key string, nowait bool, header Header) error {
	return c.ch.QueueBind(queue, key, exchange, nowait, header)
}

//HandleFunc ....
func (c *Mediator) HandleFunc(pattern string, f HandleFunc) {
	c.funcs[pattern] = f
}

//Consumer ....
func (c *Mediator) Consumer(queue string, nowait, multiple bool) error {
	q, err := c.ch.QueueInspect(queue)
	if err != nil {
		return err
	}
	log.Println(queue, q.Consumers, q.Messages)
	for {
		deliver, err := c.ch.Consume(queue, "EventDriver", false, false, false, nowait, nil)
		if err != nil {
			log.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}
		for msg := range deliver {
			if h, has := c.funcs[msg.RoutingKey]; has {
				if err := h(msg.RoutingKey, msg.Body); err != nil {
					log.Println(err)
				} else {
					err = msg.Ack(multiple)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
		for {
			log.Println("EventDriver连接异常，5秒后重新连接")
			c.Close()
			time.Sleep(5 * time.Second)
			err = c.connect()
			if err == nil {
				//c.ch.QueueBind(queue, key, exchange, nowait, nil)
				break
			}
			log.Println(err)
		}
	}
}

//Tx ....
func (c *Mediator) Tx() error {
	return c.ch.Tx()
}

//TxCommit ....
func (c *Mediator) TxCommit() error {
	return c.ch.TxCommit()
}

//TxRollback ....
func (c *Mediator) TxRollback() error {
	return c.ch.TxRollback()
}

//Forward ....
func (c *Mediator) Forward(exchange, ev string, body []byte, header Header) error {
	select {
	case b := <-c.block:
		return errors.New(b.Reason)
	case err := <-c.err:
		return err
	default:
	}
	return c.ch.Publish(
		exchange, // exchange
		ev,
		false, // mandatory
		true,  // immediate
		amqp.Publishing{
			Headers:      header,
			DeliveryMode: 2,
			Timestamp:    time.Now(),
			Body:         body,
		})
}

//Mediator ....
type Mediator struct {
	url           string
	conn          *amqp.Connection
	ch            *amqp.Channel
	err           chan *amqp.Error
	block         chan amqp.Blocking
	prefetchCount int
	prefetchSize  int
	global        bool
	funcs         map[string]HandleFunc
}

func buildURLForAliyun(uid, accessKey, secretKey, addr, vhost string) string {
	var uname = func() string {
		return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("0:%s:%s", uid, accessKey)))
	}()
	var passwd = func() string {
		key := fmt.Sprintf("%d", time.Now().Unix()*1000)
		mac := hmac.New(sha1.New, []byte(key))
		mac.Write([]byte(secretKey))
		sign := fmt.Sprintf("%02X", mac.Sum(nil))
		return base64.StdEncoding.EncodeToString([]byte(sign + ":" + key))
	}()
	return fmt.Sprintf("amqp://%s:%s@%s.%s/%s", uname, passwd, uid, addr, vhost)
}

//Header ....
type Header = amqp.Table

var json = jsoniter.ConfigCompatibleWithStandardLibrary

//HandleFunc ....
type HandleFunc func(ev string, body []byte) error
