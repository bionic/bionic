package cmd

import (
	"bufio"
	"fmt"
	"github.com/bionic-dev/bionic/providers"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

var providersCmd = &cobra.Command{
	Use:   "providers",
	Short: fmt.Sprintf("List providers available for importing using \"bionic %s\"", importCmd.Use),
	RunE: func(cmd *cobra.Command, args []string) error {
		githubProviders := map[string]string{}

		resp, err := http.Get("https://raw.githubusercontent.com/bionic-dev/how-to-export-personal-data/main/readme.md")
		if err == nil {
			scanner := bufio.NewScanner(resp.Body)
			for scanner.Scan() {
				line := scanner.Text()
				if name := strings.TrimPrefix(line, "### "); name != line {
					line = ""
					for line == "" && scanner.Scan() {
						line = scanner.Text()
					}
					if line != "" {
						githubProviders[strings.ToLower(name)] = line
					}
				}
			}

			if err := resp.Body.Close(); err != nil {
				return err
			}
		}

		var lines []string

		for _, provider := range providers.DefaultProviders(nil) {
			if desc, ok := githubProviders[provider.Name()]; ok {
				lines = append(lines, fmt.Sprintf("%s\t%s\n", provider.Name(), desc))
			} else {
				lines = append(lines, fmt.Sprintf("%s\t\n", provider.Name()))
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

func init() {
	rootCmd.AddCommand(providersCmd)
}
