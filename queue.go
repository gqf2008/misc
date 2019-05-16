package misc

import (
	"errors"
	"log"
	"time"

	"github.com/beeker1121/goque"
)

//NewPersistQueue 创建一个持久队列，path为队列目录，bufsize为缓冲区大小，ssd硬盘单队列测试下来每秒20000左右
func NewPersistQueue(path string, bufsize int) (*PersistQueue, error) {
	q, err := goque.OpenQueue(path)
	if err != nil {
		return nil, err
	}
	pq := PersistQueue{
		q:     q,
		ch:    make(chan *goque.Item, 1),
		close: make(chan struct{}),
	}
	go pq.worker()

	return &pq, nil
}

//PersistQueue ....
type PersistQueue struct {
	q     *goque.Queue
	ch    chan *goque.Item
	close chan struct{}
}

//Put ....
func (q *PersistQueue) Put(msg interface{}) error {
	_, err := q.q.EnqueueObject(msg)
	return err
}

//Poll ...
func (q *PersistQueue) Poll(msg interface{}) error {
	item := <-q.ch
	err := item.ToObject(msg)
	if err != nil {
		return err
	}
	return nil
}

//ErrorTimeout ....
var ErrorTimeout = errors.New("timeout")

//PollTimeout ...
func (q *PersistQueue) PollTimeout(msg interface{}, timeout time.Duration) error {
	select {
	case item := <-q.ch:
		err := item.ToObject(msg)
		if err != nil {
			return err
		}
	case <-time.After(timeout):
		return ErrorTimeout
	}
	return nil
}

//Size ....
func (q *PersistQueue) Size() uint64 {
	return q.q.Length()
}

//Empty ....
func (q *PersistQueue) Empty() error {
	q.q.Drop()
	return nil
}

//Delete ....
func (q *PersistQueue) Delete() error {
	q.q.Drop()
	return nil
}

//Close ....
func (q *PersistQueue) Close() error {
	q.close <- struct{}{}
	q.q.Close()
	return nil
}

func (q *PersistQueue) worker() {
	for {
		select {
		case <-q.close:
			log.Println("队列关闭")
			return
		default:
			item, err := q.q.Dequeue()
			if err == goque.ErrEmpty {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if err != nil {
				log.Println(err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			q.ch <- item
		}
	}
}
