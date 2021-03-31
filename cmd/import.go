package cmd

import (
	"fmt"
	"github.com/bionic-dev/bionic/database"
	"github.com/bionic-dev/bionic/imports"
	"github.com/bionic-dev/bionic/internal/progress"
	"github.com/bionic-dev/bionic/internal/provider/describer"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import data to local db",
}

func init() {
	for _, p := range imports.DefaultProviders(nil) {
		providerName := p.Name()

		command := &cobra.Command{
			Use: fmt.Sprintf("%s [path]", providerName),
			RunE: func(cmd *cobra.Command, args []string) error {
				inputPath := args[0]

				db, err := database.New(dbPath)
				if err != nil {
					return err
				}

				manager, err := imports.NewManager(db, imports.DefaultProviders(db))
				if err != nil {
					return err
				}

				if err := manager.Migrate(); err != nil {
					return err
				}

				p, err := manager.GetByName(providerName)
				if err != nil {
					return err
				}

				importFns, err := p.ImportFns(inputPath)
				if err != nil {
					return err
				}

				if err := p.BeginTx(); err != nil {
					return err
				}
				defer p.RollbackTx() //nolint:errcheck

				errs, _ := errgroup.WithContext(cmd.Context())

				importProgress := progress.New()

				for _, importFn := range importFns {
					name := importFn.Name()
					importProgress.Init(name)
				}

				importProgress.Draw()

				for _, importFn := range importFns {
					name := importFn.Name()
					fn := importFn.Call

					errs.Go(func() error {
						defer importProgress.Draw()

						err := fn()

						if err != nil {
							importProgress.Error(name)
							return err
						}

						importProgress.Success(name)

						return nil
					})
				}

				if err := errs.Wait(); err != nil {
					return err
				}

				err = p.DB().
					Create(&imports.Import{
						Provider: p.Name(),
					}).
					Error
				if err != nil {
					return err
				}

				if err := p.CommitTx(); err != nil {
					return err
				}

				return nil
			},
			Args: cobra.MinimumNArgs(1),
		}

		if d, ok := p.(describer.Import); ok {
			command.Short = d.ImportDescription()
		}

		importCmd.AddCommand(command)
	}
}
