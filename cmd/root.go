package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
)

var dbPath string

var rootCmd = &cobra.Command{
	Use:   "bionic",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	homeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	defaultDbPath := path.Join(homeDir, ".bionic", "db.sqlite")

	rootCmd.PersistentFlags().StringVar(&dbPath, "db", defaultDbPath, "db path (default is $HOME/.bionic/db.sqlite)")

	rootCmd.PersistentFlags().Bool("verbose", false, "verbose output")
	panicOnErr(viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")))
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
