package cmd

import (
	"github.com/shekhirin/bionic-cli/providers"

	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset [service]",
	Short: "Reset service data stored in local db",
	RunE: func(cmd *cobra.Command, args []string) error {
		providerName := args[0]

		dbPath := rootCmd.PersistentFlags().Lookup("db").Value.String()

		manager, err := providers.NewManager(dbPath)
		if err != nil {
			return err
		}

		provider, err := manager.GetByName(providerName)
		if err != nil {
			return err
		}

		return manager.Reset(provider)
	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
