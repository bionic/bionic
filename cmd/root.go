package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var dbPath string

var rootCmd = &cobra.Command{
	Use:   "bionic",
	Short: "Load personal data exports to an SQLite database",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("verbose") {
			logrus.SetLevel(logrus.DebugLevel)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", "", "db path")
	panicOnErr(rootCmd.MarkPersistentFlagRequired("db"))

	rootCmd.PersistentFlags().Bool("verbose", false, "verbose output")
	panicOnErr(viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")))

	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	rootCmd.AddCommand(importCmd, resetCmd, generateViewsCmd, providersCmd)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
