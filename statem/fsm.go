package statem

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Transition ....
type transition struct {
	from   string
	event  string
	action string
	to     string
	f      ActionFunc
}

//StateEnterFunc ....
type StateEnterFunc = func(ctx context.Context, event, to string, args ...interface{}) error

//ActionFunc ....
type ActionFunc = func(ctx context.Context, from, event string, action string, to string, args ...interface{}) error

//StateExitFunc ....
type StateExitFunc = func(ctx context.Context, from, event string, args ...interface{}) error

type event struct {
	ev   string
	args []interface{}
}

// FSM ....
type FSM struct {
	sef         StateEnterFunc
	se          StateExitFunc
	transitions []transition
	//pool        *misc.WorkerPool
}

// StateError ....
type StateError struct {
	BadEvent     string
	CurrentState string
	Args         []interface{}
}

func (e *StateError) Error() string {
	return fmt.Sprintf("状态机发生错误: 当触发[%s]事件时当前状态[%s]没有找到转换器 Args: %+v\n", e.BadEvent, e.CurrentState, e.Args)
}

// New ....
func New(poolSize int) *FSM {
	fsm := &FSM{
		// pool: &misc.WorkerPool{
		// 	MaxWorkersCount: poolSize,
		// },
	}
	return fsm
}

//Start ....
// func (m *FSM) Start() {
// 	m.pool.Start()
// }

// //Stop ....
// func (m *FSM) Stop() {
// 	m.pool.Stop()
// }

//WithTransition ....
func (m *FSM) WithTransition(from, event, action, to string, f ...ActionFunc) *FSM {
	var af ActionFunc
	if len(f) > 0 {
		af = f[0]
	}
	m.transitions = append(m.transitions, transition{from, event, action, to, af})
	return m
}

//WithStateExitFunc ....
func (m *FSM) WithStateExitFunc(f StateExitFunc) *FSM {
	m.se = f
	return m
}

//WithStateEnterFunc ....
func (m *FSM) WithStateEnterFunc(f StateEnterFunc) *FSM {
	m.sef = f
	return m
}

//Event ....
func (m *FSM) Event(ctx context.Context, currentState, ev string, args ...interface{}) error {
	return m.event(ctx, currentState, ev, args...)
}

// Event ....
func (m *FSM) event(ctx context.Context, current, event string, args ...interface{}) error {
	for _, trans := range m.transitions {
		if trans.from == current && trans.event == event {
			changingStates := current != trans.to
			if changingStates && m.se != nil {
				if err := m.se(ctx, current, event, args...); err != nil {
					return err
				}
			}
			if trans.action != "" && trans.f != nil {
				err := trans.f(ctx, current, event, trans.action, trans.to, args...)
				if err != nil {
					return err
				}
			}
			if changingStates && m.sef != nil {
				if err := m.sef(ctx, event, trans.to, args...); err != nil {
					return err
				}
			}
			return nil
		}
	}
	return &StateError{event, current, args}
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
