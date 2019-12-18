package misc

import (
	"fmt"
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
	//defer p.lock.Unlock()
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

//ToString 把ID格式化成字符串
func (p *IDGenerator) ToString(id uint64) string {
	seq := id & maxSeq
	node := (id >> 20) & maxNode
	timestamp := (id >> 28) & maxTimestamp
	prefix := (id >> 60) & maxPrefix
	return fmt.Sprintf("%d%s%03d%07d", prefix, time.Unix(0, (int64(timestamp)+epoch)*int64(time.Second)).Format("20060102150405"), node, seq)
}

//FormatID ....
func FormatID(id uint64) string {
	seq := id & maxSeq
	node := (id >> 20) & maxNode
	timestamp := (id >> 28) & maxTimestamp
	prefix := (id >> 60) & maxPrefix
	return fmt.Sprintf("%d%s%03d%07d", prefix, time.Unix(0, (int64(timestamp)+epoch)*int64(time.Second)).Format("20060102150405"), node, seq)
}

//State 获取当前内部状态
func (p *IDGenerator) State() string {
	return fmt.Sprintf("prefix=%d,node=%d,last=%d,seq=%d,step=%d", p.prefix, p.nodeID, p.last, p.seq, p.step)
}
