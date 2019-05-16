package fsm

import (
	"fmt"
	"os/exec"
	"strings"
)

//TransitionListener ....
type TransitionListener interface {
	OnStateEnter(event, toState string, args ...interface{}) error
	OnAction(fromState, event string, action Action, toState string, args ...interface{}) error
	OnStateExit(fromState, event string, args ...interface{}) error
}

//Action ....
type Action struct {
	Action  string
	Handler func(args ...interface{}) error
}

// Transition ....
type Transition struct {
	From   string
	Event  string
	Action Action
	To     string
}

// StateMachine ....
type StateMachine struct {
	listener    TransitionListener
	transitions []Transition
}

// StatemError ....
type StatemError struct {
	BadEvent     string
	CurrentState string
}

func (e *StatemError) Error() string {
	return fmt.Sprintf("状态机发生错误: 当触发[%s]事件时当前状态[%s]没有找到转换器\n", e.BadEvent, e.CurrentState)
}

// New ....
func New() *StateMachine {
	return &StateMachine{}
}

//WithTransition ....
func (m *StateMachine) WithTransition(transition Transition) *StateMachine {
	if transition.Action.Handler != nil && transition.Action.Action == "" {
		panic("当处理函数不为空时，Action不能为空")
	}
	m.transitions = append(m.transitions, transition)
	return m
}

//WithTransitions ....
func (m *StateMachine) WithTransitions(transitions []Transition) *StateMachine {
	m.transitions = transitions
	for _, t := range transitions {
		m.WithTransition(t)
	}
	return m
}

//WithTransitionListener ....
func (m *StateMachine) WithTransitionListener(l TransitionListener) *StateMachine {
	m.listener = l
	return m
}

// Event ....
func (m *StateMachine) Event(currentState, event string, args ...interface{}) error {
	trans := m.findTransMatching(currentState, event)
	if trans == nil {
		return &StatemError{event, currentState}
	}
	changingStates := currentState != trans.To
	if changingStates && m.listener != nil {
		if err := m.listener.OnStateExit(currentState, event, args...); err != nil {
			return err
		}
	}

	if m.listener != nil && trans.Action.Action != "" {
		err := m.listener.OnAction(currentState, event, trans.Action, trans.To, args...)
		if err != nil {
			return err
		}
	}

	if changingStates && m.listener != nil {
		if err := m.listener.OnStateEnter(event, trans.To, args...); err != nil {
			return err
		}
	}
	return nil
}

func (m *StateMachine) findTransMatching(fromState string, event string) *Transition {
	for _, v := range m.transitions {
		if v.From == fromState && v.Event == event {
			return &v
		}
	}
	return nil
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
	// dot := `digraph StateMachine {
	// rankdir=LR
	// node[width=1 fixedsize=false shape=ellipse style=filled fillcolor="darkorchid1" ]
	// `

	// for _, t := range m.transitions {
	// 	link := fmt.Sprintf(`%s -> %s [label="%s"]`, t.From, t.To, t.Event)
	// 	if t.Action != "" {
	// 		link = fmt.Sprintf(`%s -> %s [label="%s | %s"]`, t.From, t.To, t.Event, t.Action)
	// 	}
	// 	dot = dot + "\r\n" + link
	// }

	// dot = dot + "\r\n}"
	// cmd := fmt.Sprintf("dot -o%s -T%s -K%s -s%s %s", outfile, format, layout, scale, more)
	dot := `digraph StateMachine {
		rankdir=LR
		node[width=1 fixedsize=true shape=circle style=filled fillcolor="darkorchid1" ]
		
		`

	for _, t := range m.transitions {
		link := fmt.Sprintf(`%s -> %s [label="%s | %s"]`, t.From, t.To, t.Event, t.Action.Action)
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
