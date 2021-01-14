package cmd

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/shekhirin/bionic-cli/internal/progress"
	"github.com/shekhirin/bionic-cli/providers"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

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
		defer provider.RollbackTx() //nolint:errcheck

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
