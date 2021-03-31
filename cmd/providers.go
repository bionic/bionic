package cmd

import (
	"fmt"
	"github.com/bionic-dev/bionic/imports"
	"github.com/bionic-dev/bionic/internal/provider/describer"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"text/tabwriter"
)

var providersCmd = &cobra.Command{
	Use:   "providers",
	Short: fmt.Sprintf("List providers available for importing using \"bionic %s\"", importCmd.Use),
	RunE: func(cmd *cobra.Command, args []string) error {
		var lines []string

		for _, p := range imports.DefaultProviders(nil) {
			if d, ok := p.(describer.Export); ok {
				lines = append(lines, fmt.Sprintf("%s\t%s\n", p.Name(), d.ExportDescription()))
			} else {
				lines = append(lines, fmt.Sprintf("%s\t\n", p.Name()))
			}
		}

		sort.Strings(lines)

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 3, ' ', 0)
		for _, line := range lines {
			if _, err := fmt.Fprint(w, line); err != nil {
				return err
			}
		}

		return w.Flush()
	},
	DisableFlagParsing: true,
}
