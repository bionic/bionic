package progress

import (
	"fmt"
	"github.com/gosuri/uilive"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
)

type State int

const (
	InitState State = iota
	ErrorState
	SuccessState
)

func (s State) String() string {
	switch s {
	case InitState:
		return "⌛️"
	case ErrorState:
		return "❌"
	case SuccessState:
		return "✅"
	default:
		return ""
	}
}

type Progress struct {
	mu     sync.Mutex
	names  []string
	states map[string]State
	writer *uilive.Writer
}

func New() Progress {
	return Progress{
		states: map[string]State{},
		writer: uilive.New(),
	}
}

func (p *Progress) Draw() {
	var output string
	for _, name := range p.names {
		output += fmt.Sprintf("%s %s\n", p.states[name], name)
	}
	_, _ = fmt.Fprint(p.writer, output)
	_ = p.writer.Flush()
}

func (p *Progress) add(name string, state State) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.states[name]; !ok {
		p.names = append(p.names, name)
	}
	p.states[name] = state
}

func (p *Progress) Init(name string) {
	if !viper.GetBool("verbose") {
		p.add(name, InitState)
	}
	logrus.Debugf("%s transitioned to Init state", name)
}

func (p *Progress) Error(name string) {
	if !viper.GetBool("verbose") {
		p.add(name, ErrorState)
	}
	logrus.Debugf("%s transitioned to Error state", name)
}

func (p *Progress) Success(name string) {
	if !viper.GetBool("verbose") {
		p.add(name, SuccessState)
	}
	logrus.Debugf("%s transitioned to Success state", name)
}
