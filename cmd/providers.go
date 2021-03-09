package cmd

import (
	"fmt"
	"github.com/bionic-dev/bionic/providers"
	"github.com/bionic-dev/bionic/providers/provider"
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

		for _, p := range providers.DefaultProviders(nil) {
			if describer, ok := p.(provider.ExportDescriber); ok {
				lines = append(lines, fmt.Sprintf("%s\t%s\n", p.Name(), describer.ExportDescription()))
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
