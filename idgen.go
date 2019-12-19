package misc

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

const (
	maxPrefix    = 2<<3 - 1
	maxNode      = 2<<7 - 1
	maxSeq       = 2<<19 - 1
	maxTimestamp = 2<<31 - 1
	epoch        = int64(1575129600) //second
)

//IDGenerator ID发号器
type IDGenerator struct {
	prefix uint
	nodeID uint
	last   int64
	seq    uint64
	step   uint64
	lock   sync.Mutex
}

//NewIDGenerator ....
func NewIDGenerator(prefix, node uint) *IDGenerator {
	if prefix > maxPrefix {
		panic(fmt.Errorf("prefix不合法，prefix不能大于%d", maxPrefix))
	}
	if node > maxNode {
		panic(fmt.Errorf("node不合法，node不能大于%d", maxNode))
	}
	time.Sleep(time.Second)
	return &IDGenerator{
		prefix: prefix,
		nodeID: node,
		last:   time.Now().Unix(),
		step:   1,
	}
}

//WithStep 设置步长
func (p *IDGenerator) WithStep(step uint64) *IDGenerator {
	if step == 0 || step > maxSeq {
		panic(fmt.Errorf("步长必须在[1-%d]之间", maxSeq))
	}
	p.step = step
	return p
}

//NextID 获取下一个ID
func (p *IDGenerator) NextID() uint64 {
	p.lock.Lock()
	current := time.Now().Unix()
	if current < p.last {
		refused := p.last - current
		p.lock.Unlock()
		panic(fmt.Errorf("发生时钟回拨，拒绝执行%d秒", refused))
	}
	if current > p.last {
		p.seq = 0
		p.last = current
	}
	p.seq += p.step
	if p.seq > maxSeq {
		p.lock.Unlock()
		return p.NextID()
	}
	timestamp := current - epoch
	var v = uint64(p.prefix)<<60 | uint64(timestamp)<<28 | uint64(p.nodeID)<<20 | uint64(p.seq)
	p.lock.Unlock()
	return v
}

//FormatID ....
func FormatID(id uint64) string {
	seq := id & maxSeq
	node := (id >> 20) & maxNode
	timestamp := (id >> 28) & maxTimestamp
	prefix := (id >> 60) & maxPrefix
	return fmt.Sprintf("%d%s%03d%07d", prefix, time.Unix(0, (int64(timestamp)+epoch)*int64(time.Second)).Format("060102150405"), node, seq)
}

//ParserID ....
func ParserID(id string) (uint64, error) {
	l := len(id)
	if l < 24 {
		return 0, errors.New("ID不合法")
	}
	seq, err := strconv.ParseInt(id[l-7:], 10, 64)
	if err != nil || seq > maxSeq {
		return 0, errors.New("ID不合法")
	}
	node, err := strconv.ParseInt(id[l-10:l-7], 10, 64)
	if err != nil || node > maxNode {
		return 0, errors.New("ID不合法")
	}
	t, err := time.ParseInLocation("060102150405", id[l-24:l-10], time.Local)
	if err != nil {
		return 0, errors.New("ID不合法")
	}
	time := t.Unix()
	biz, err := strconv.ParseInt(id[:l-24], 10, 64)
	if err != nil || biz > maxPrefix {
		return 0, errors.New("ID不合法")
	}
	timestamp := time - epoch
	return uint64(biz)<<60 | uint64(timestamp)<<28 | uint64(node)<<20 | uint64(seq), nil
}

//State 获取当前内部状态
func (p *IDGenerator) State() string {
	return fmt.Sprintf("prefix=%d,node=%d,last=%d,seq=%d,step=%d", p.prefix, p.nodeID, p.last, p.seq, p.step)
}
