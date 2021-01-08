package progress

import (
	"fmt"
	"github.com/gosuri/uilive"
	"sync"
)

type Progress struct {
	mu     sync.Mutex
	names  []string
	states map[string]string
	writer *uilive.Writer
}

func New() Progress {
	return Progress{
		states: map[string]string{},
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

func (p *Progress) add(name, state string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.states[name]; !ok {
		p.names = append(p.names, name)
	}
	p.states[name] = state
}

func (p *Progress) Init(name string) {
	p.add(name, "⌛️")
}

func (p *Progress) Error(name string) {
	p.add(name, "❌")
}

func (p *Progress) Success(name string) {
	p.add(name, "✅")
}
