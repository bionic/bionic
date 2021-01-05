package cmd

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/shekhirin/bionic-cli/providers"
	"github.com/shekhirin/bionic-cli/providers/provider"
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

		p, err := manager.GetByName(providerName)
		if err != nil {
			return err
		}

		importFns, err := p.ImportFns(inputPath)
		if err != nil {
			return err
		}

		dbProvider, isDbProvider := p.(provider.Database)

		if isDbProvider {
			if err := dbProvider.BeginTx(); err != nil {
				return err
			}
			defer dbProvider.RollbackTx()
		}

		errs, _ := errgroup.WithContext(cmd.Context())

		for _, importFn := range importFns {
			fn := importFn.Call
			errs.Go(fn)
		}

		err = errs.Wait()

		if isDbProvider {
			if err := dbProvider.CommitTx(); err != nil {
				return err
			}
		}

		return err
	},
	Args: cobra.MinimumNArgs(2),
}

func init() {
	rootCmd.AddCommand(importCmd)
}
