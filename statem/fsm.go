package statem

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strings"
)

// Transition ....
type transition struct {
	from  string
	event string
	//action string
	to string
	f  ActionFunc
}

//StateEnterFunc ....
type StateEnterFunc = func(ctx context.Context, event, to string, args ...interface{}) error

//ActionFunc ....
type ActionFunc = func(ctx context.Context, from, event string, to string, args ...interface{}) error

//StateExitFunc ....
type StateExitFunc = func(ctx context.Context, from, event string, args ...interface{}) error

type event struct {
	ev   string
	args []interface{}
}

// FSM ....
type FSM struct {
	sef         StateEnterFunc
	af          ActionFunc
	se          StateExitFunc
	transitions []transition
	actions     map[string]ActionFunc
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
		actions: map[string]ActionFunc{},
	}
	fsm.af = func(ctx context.Context, from, event, to string, args ...interface{}) error {
		if f, has := fsm.actions[event]; has {
			return f(ctx, from, event, to, args...)
		}
		return nil
		// return fmt.Errorf("action %s not found,state(%s)->event(%s)->state(%s)", from, event, to)
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
func (m *FSM) WithTransition(from, event, to string, f ...ActionFunc) *FSM {
	var af ActionFunc
	if len(f) > 0 {
		af = f[0]
	}
	m.transitions = append(m.transitions, transition{from, event, to, af})
	return m
}

// //WithTransitionActionName ....
// func (m *FSM) WithTransitionActionName(from, event, to string, action ...string) *FSM {
// 	var a = ""
// 	if len(action) > 0 {
// 		a = action[0]
// 	}
// 	m.transitions = append(m.transitions, transition{from, event, a, to, nil})
// 	return m
// }

//WithTransitionFromFile ....
func (m *FSM) WithTransitionFromFile(file string) *FSM {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return m.WithTransitionFromJSON(b)
}

//WithTransitionFromJSON ....
func (m *FSM) WithTransitionFromJSON(b []byte) *FSM {
	var transitions = [][]string{}
	err := json.Unmarshal(b, &transitions)
	if err != nil {
		panic(err)
	}
	for _, a := range transitions {
		l := len(a)
		if l < 3 {
			panic(errors.New("error transition"))
		}
		m.WithTransition(a[0], a[1], a[2])
	}
	return m
}

//WithActionFunc ....
func (m *FSM) WithActionFunc(action string, f ActionFunc) *FSM {
	m.actions[action] = f
	return m
}

//WithActionFuncs ....
func (m *FSM) WithActionFuncs(actions map[string]ActionFunc) *FSM {
	m.actions = actions
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
			if trans.f != nil {
				err := trans.f(ctx, current, event, trans.to, args...)
				if err != nil {
					return err
				}
			} else {
				err := m.af(ctx, current, event, trans.to, args...)
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
	dot := `digraph StateMachine {

		rankdir=LR
		node[width=1 fixedsize=true shape=ellipse style=filled fillcolor="darkorchid1" ]
		
		`

	for _, t := range m.transitions {
		var link string
		link = fmt.Sprintf(`"%s" -> "%s" [label="%s"]`, t.from, t.to, t.event)
		dot = dot + "\r\n" + link
	}

	dot = dot + "\r\n}"
	cmd := fmt.Sprintf("dot -o%s -T%s -K%s -s%s %s", outfile, format, layout, scale, more)
	return system(cmd, dot)
}

func system(c string, dot string) error {
	fmt.Println(c)
	fmt.Println(dot)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command(`cmd`, `/C`, c)
	} else {
		cmd = exec.Command(`/bin/sh`, `-c`, c)
	}
	cmd.Stdin = strings.NewReader(dot)
	return cmd.Run()

}
