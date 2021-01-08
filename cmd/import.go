package cmd

import (
	"fmt"
	"github.com/gosuri/uilive"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shekhirin/bionic-cli/providers"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"sync"
)

type State struct {
	mu     sync.Mutex
	names  []string
	states map[string]string
	writer *uilive.Writer
}

func NewState(writer *uilive.Writer) State {
	return State{
		states: map[string]string{},
		writer: writer,
	}
}

func (s *State) draw() {
	var output string
	for _, name := range s.names {
		output += fmt.Sprintf("%s %s\n", s.states[name], name)
	}
	_, _ = fmt.Fprint(s.writer, output)
	_ = s.writer.Flush()
}

func (s *State) add(name, state string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.states[name]; !ok {
		s.names = append(s.names, name)
	}
	s.states[name] = state
}

func (s *State) Init(name string) {
	s.add(name, "⌛️")
}

func (s *State) Error(name string) {
	s.add(name, "❌")
	s.draw()
}

func (s *State) Success(name string) {
	s.add(name, "✅")
	s.draw()
}

var importCmd = &cobra.Command{
	Use:   "import [service] [path]",
	Short: "Import GDPR export to local db",
	RunE: func(cmd *cobra.Command, args []string) error {
		providerName, inputPath := args[0], args[1]

		dbPath := rootCmd.PersistentFlags().Lookup("db").Value.String()

		manager, err := providers.NewManager(dbPath)
		if err != nil {
			return err
		}

		provider, err := manager.GetByName(providerName)
		if err != nil {
			return err
		}

		importFns, err := provider.ImportFns(inputPath)
		if err != nil {
			return err
		}

		if err := provider.BeginTx(); err != nil {
			return err
		}
		defer provider.RollbackTx()

		errs, _ := errgroup.WithContext(cmd.Context())

		writer := uilive.New()
		state := NewState(writer)

		for _, importFn := range importFns {
			name := importFn.Name()
			state.Init(name)
		}

		state.draw()

		for _, importFn := range importFns {
			name := importFn.Name()
			fn := importFn.Call

			errs.Go(func() error {
				err := fn()

				if err != nil {
					state.Error(name)
					return err
				}

				state.Success(name)

				return nil
			})
		}

		err = errs.Wait()
		if err != nil {
			return err
		}

		if err := provider.CommitTx(); err != nil {
			return err
		}

		return nil
	},
	Args: cobra.MinimumNArgs(2),
}

func init() {
	rootCmd.AddCommand(importCmd)
}
