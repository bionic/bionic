package cmd

import (
	"fmt"
	"github.com/bionic-dev/bionic/providers"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"text/tabwriter"
)

var providerDescriptions = map[string]string{
	"airbnb":    "https://www.airbnb.com/privacy/manage-your-data",
	"amazon":    "https://www.amazon.com/gp/privacycentral/dsar/preview.html",
	"apple":     "https://privacy.apple.com/ => \"Get a copy of your data\"",
	"duolingo":  "https://drive-thru.duolingo.com/",
	"ebay":      "https://www.sarweb.ebay.com/sar",
	"facebook":  "https://www.facebook.com/dyi",
	"google":    "https://takeout.google.com/",
	"instagram": "https://www.instagram.com/download/request/",
	"reddit":    "https://www.reddit.com/settings/data-request",
	"snapchat":  "https://accounts.snapchat.com/accounts/downloadmydata",
	"telegram":  "Desktop App (https://desktop.telegram.org/ only): Settings => Advanced => Export Telegram data",
	"tiktok":    "Mobile App: Settings => Privacy => Personalization & Data",
	"tinder":    "https://account.gotinder.com/data",
}

var providersCmd = &cobra.Command{
	Use:   "providers",
	Short: fmt.Sprintf("List providers available for importing using \"bionic %s\"", importCmd.Use),
	RunE: func(cmd *cobra.Command, args []string) error {
		var lines []string

		for _, provider := range providers.DefaultProviders(nil) {
			if desc, ok := providerDescriptions[provider.Name()]; ok {
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
