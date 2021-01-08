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

func (s *Progress) Draw() {
	var output string
	for _, name := range s.names {
		output += fmt.Sprintf("%s %s\n", s.states[name], name)
	}
	_, _ = fmt.Fprint(s.writer, output)
	_ = s.writer.Flush()
}

func (s *Progress) add(name, state string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.states[name]; !ok {
		s.names = append(s.names, name)
	}
	s.states[name] = state
}

func (s *Progress) Init(name string) {
	s.add(name, "⌛️")
}

func (s *Progress) Error(name string) {
	s.add(name, "❌")
}

func (s *Progress) Success(name string) {
	s.add(name, "✅")
}
