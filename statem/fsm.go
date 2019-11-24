package statem

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gqf2008/misc"
)

// Transition ....
type transition struct {
	from   string
	event  string
	action string
	to     string
}

//StateEnterFunc ....
type StateEnterFunc = func(event, to string, args ...interface{}) error

//ActionFunc ....
type ActionFunc = func(from, event string, action string, to string, args ...interface{}) error

//StateExitFunc ....
type StateExitFunc = func(from, event string, args ...interface{}) error

//ErrorFunc ....
type ErrorFunc = func(err error, from, event string, args ...interface{})

type event struct {
	ev   string
	args []interface{}
}

// FSM ....
type FSM struct {
	sef         StateEnterFunc
	af          ActionFunc
	se          StateExitFunc
	ef          ErrorFunc
	transitions []transition
	pool        *misc.WorkerPool
}

// StateError ....
type StateError struct {
	BadEvent     string
	CurrentState string
}

func (e *StateError) Error() string {
	return fmt.Sprintf("状态机发生错误: 当触发[%s]事件时当前状态[%s]没有找到转换器\n", e.BadEvent, e.CurrentState)
}

// New ....
func New(poolSize int) *FSM {
	fsm := &FSM{
		pool: &misc.WorkerPool{
			MaxWorkersCount: poolSize,
		},
	}
	return fsm
}

//Start ....
func (m *FSM) Start() {
	m.pool.Start()
}

//Stop ....
func (m *FSM) Stop() {
	m.pool.Stop()
}

//WithTransition ....
func (m *FSM) WithTransition(from, event, action, to string) *FSM {
	m.transitions = append(m.transitions, transition{from, event, action, to})
	return m
}

//WithStateEnterFunc ....
func (m *FSM) WithStateEnterFunc(f StateEnterFunc) *FSM {
	m.sef = f
	return m
}

//WithErrorFunc ....
func (m *FSM) WithErrorFunc(f ErrorFunc) *FSM {
	m.ef = f
	return m
}

//Event ....
func (m *FSM) Event(currentState, ev string, args ...interface{}) {
	if !m.pool.Serve(func() {
		err := m.event(currentState, ev, args...)
		if err != nil && m.ef != nil {
			m.ef(err, currentState, ev, args...)
		}
	}) && m.ef != nil {
		m.ef(errors.New("没有足够的工作进程"), currentState, ev, args...)
	}
}

// Event ....
func (m *FSM) event(current, event string, args ...interface{}) error {
	for _, trans := range m.transitions {
		if trans.from == current && trans.event == event {
			changingStates := current != trans.to
			if changingStates && m.se != nil {
				if err := m.se(current, event, args...); err != nil {
					return err
				}
			}
			if m.af != nil && trans.action != "" {
				err := m.af(current, event, trans.action, trans.to, args...)
				if err != nil {
					return err
				}
			}
			if changingStates && m.se != nil {
				if err := m.se(event, trans.to, args...); err != nil {
					return err
				}
			}
			return nil
		}
	}
	return &StateError{event, current}
}

// ExportPNG 导出状态图
func (m *FSM) ExportPNG(outfile string) error {
	if !strings.HasSuffix(outfile, ".png") {
		outfile = outfile + ".png"
	}
	return m.ExportWithDetails(outfile, "png", "dot", "72", "-Gsize=10,5 -Gdpi=200")
}

//ExportJPG ....
func (m *FSM) ExportJPG(outfile string) error {
	if !strings.HasSuffix(outfile, ".jpg") {
		outfile = outfile + ".jpg"
	}
	return m.ExportWithDetails(outfile, "jpg", "dot", "72", "-Gsize=10,5 -Gdpi=200")
}

// ExportWithDetails  导出状态图
func (m *FSM) ExportWithDetails(outfile string, format string, layout string, scale string, more string) error {
	dot := `digraph FSM {
		rankdir=LR
		node[width=1 fixedsize=true shape=circle style=filled fillcolor="darkorchid1" ]
		
		`

	for _, t := range m.transitions {
		link := fmt.Sprintf(`%s -> %s [label="%s | %s"]`, t.from, t.to, t.event, t.action)
		dot = dot + "\r\n" + link
	}

	dot = dot + "\r\n}"
	cmd := fmt.Sprintf("dot -o%s -T%s -K%s -s%s %s", outfile, format, layout, scale, more)
	return system(cmd, dot)
}

func system(c string, dot string) error {
	cmd := exec.Command(`/bin/sh`, `-c`, c)
	cmd.Stdin = strings.NewReader(dot)
	return cmd.Run()

}
