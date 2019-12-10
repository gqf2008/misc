package hfsm

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strings"
)

// StateError ....
type StateError struct {
	BadEvent     string
	CurrentState string
	Args         []interface{}
}

func (e *StateError) Error() string {
	return fmt.Sprintf("状态机发生错误: 当触发[%s]事件时当前状态[%s]没有找到转换器 Args: %+v\n", e.BadEvent, e.CurrentState, e.Args)
}

type transition struct {
	From   string
	Event  string
	Action string
	To     string
	Child  *StateMachine
}

//StateEnterFunc ....
type StateEnterFunc = func(ctx context.Context, event, to string, args ...interface{}) error

//ActionFunc ....
type ActionFunc = func(ctx context.Context, from, event string, action string, to string, args ...interface{}) error

//StateExitFunc ....
type StateExitFunc = func(ctx context.Context, from, event string, args ...interface{}) error

//StateMachine ....
type StateMachine struct {
	Name        string
	sef         StateEnterFunc
	se          StateExitFunc
	af          ActionFunc
	Transitions []transition
	actions     map[string]ActionFunc
}

//NewFromFile ....
func NewFromFile(file string) *StateMachine {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return NewFromJSON(b)
}

//NewFromJSON ....
func NewFromJSON(b []byte) *StateMachine {
	var sm StateMachine
	err := json.Unmarshal(b, &sm)
	if err != nil {
		panic(err)
	}
	return &sm
}

//WithActionFunc ....
func (m *StateMachine) WithActionFunc(action string, f ActionFunc) *StateMachine {
	m.actions[action] = f
	return m
}

//WithActionFuncs ....
func (m *StateMachine) WithActionFuncs(actions map[string]ActionFunc) *StateMachine {
	m.actions = actions
	return m
}

//WithStateExitFunc ....
func (m *StateMachine) WithStateExitFunc(f StateExitFunc) *StateMachine {
	m.se = f
	return m
}

//WithStateEnterFunc ....
func (m *StateMachine) WithStateEnterFunc(f StateEnterFunc) *StateMachine {
	m.sef = f
	return m
}

// //Event ....
// func (m *StateMachine) Event(ctx context.Context, currentState, ev string, args ...interface{}) error {
// 	return m.event(ctx, currentState, ev, args...)
// }

// Event ....
func (m *StateMachine) Event(ctx context.Context, current, event string, args ...interface{}) error {
	for _, trans := range m.Transitions {
		if trans.From == current && trans.Event == event {
			changingStates := current != trans.To
			if changingStates && m.se != nil {
				if err := m.se(ctx, current, event, args...); err != nil {
					return err
				}
			}
			if trans.Action != "" {
				err := m.af(ctx, current, event, trans.Action, trans.To, args...)
				if err != nil {
					return err
				}
			}
			if changingStates && m.sef != nil {
				if err := m.sef(ctx, event, trans.To, args...); err != nil {
					return err
				}
			}
			return nil
		}
	}
	return &StateError{event, current, args}
}

// ExportPNG 导出状态图
func (m *StateMachine) ExportPNG(outfile string) error {
	if !strings.HasSuffix(outfile, ".png") {
		outfile = outfile + ".png"
	}
	return m.ExportWithDetails(outfile, "png", "dot", "72", "-Gsize=10,5 -Gdpi=200")
}

//ExportJPG ....
func (m *StateMachine) ExportJPG(outfile string) error {
	if !strings.HasSuffix(outfile, ".jpg") {
		outfile = outfile + ".jpg"
	}
	return m.ExportWithDetails(outfile, "jpg", "dot", "72", "-Gsize=10,5 -Gdpi=200")
}

// ExportWithDetails  导出状态图
func (m *StateMachine) ExportWithDetails(outfile string, format string, layout string, scale string, more string) error {
	dot := `digraph StateMachine {

		rankdir=LR
		node[width=1 fixedsize=true shape=ellipse style=filled fillcolor="darkorchid1" ]
		
		`

	for _, t := range m.Transitions {
		var link string
		if t.Action == "" {
			link = fmt.Sprintf(`"%s" -> "%s" [label="%s"]`, t.From, t.To, t.Event)
		} else {
			link = fmt.Sprintf(`"%s" -> "%s" [label="%s | %s"]`, t.From, t.To, t.Event, t.Action)
		}

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
