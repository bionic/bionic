package cmd

import (
	"fmt"
	"github.com/bionic-dev/bionic/database"
	"github.com/bionic-dev/bionic/exports"
	"github.com/bionic-dev/bionic/internal/provider/describer"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export data from local db",
}

func init() {
	for _, p := range exports.DefaultProviders(nil) {
		providerName := p.Name()

		command := &cobra.Command{
			Use: fmt.Sprintf("%s [output path]", providerName),
			RunE: func(cmd *cobra.Command, args []string) error {
				outputPath := args[0]

				db, err := database.New(dbPath)
				if err != nil {
					return err
				}

				manager, err := exports.NewManager(db, exports.DefaultProviders(db))
				if err != nil {
					return err
				}

				p, err := manager.GetByName(providerName)
				if err != nil {
					return err
				}

				return p.Export(outputPath)
			},
			Args: cobra.MinimumNArgs(1),
		}

		if d, ok := p.(describer.Export); ok {
			command.Short = d.ExportDescription()
		}

		exportCmd.AddCommand(command)
	}
}
