package cmd

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var markdownCmd = &cobra.Command{
	Use:   "markdown",
	Short: "Export data from local db to markdown format",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(dbPath)

		return nil
	},
}

