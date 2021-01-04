package cmd

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/shekhirin/bionic-cli/providers"
	"github.com/spf13/cobra"
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

		if err := manager.Migrate(provider); err != nil {
			return err
		}

		return provider.Process(inputPath)
	},
	Args: cobra.MinimumNArgs(2),
}

func init() {
	rootCmd.AddCommand(importCmd)
}

