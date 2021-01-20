package cmd

import (
	"github.com/BionicTeam/bionic/database"
	"github.com/BionicTeam/bionic/internal/progress"
	"github.com/BionicTeam/bionic/views"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var viewsCmd = &cobra.Command{
	Use:   "views",
	Short: "Create derivative SQL tables and materialized views based on raw imported data",
	RunE: func(cmd *cobra.Command, args []string) error {
		dbPath := rootCmd.PersistentFlags().Lookup("db").Value.String()

		db, err := database.New(dbPath)
		if err != nil {
			return err
		}

		manager, err := views.NewManager(db, views.DefaultViews())
		if err != nil {
			return err
		}

		errs, _ := errgroup.WithContext(cmd.Context())

		viewProgress := progress.New()

		for _, view := range manager.Views {
			name := view.TableName()
			viewProgress.Init(name)
		}

		viewProgress.Draw()

		for _, view := range manager.Views {
			name := view.TableName()
			fn := view.Update

			errs.Go(func() error {
				defer viewProgress.Draw()

				err := fn(db)

				if err != nil {
					viewProgress.Error(name)
					return err
				}

				viewProgress.Success(name)

				return nil
			})
		}

		if err := errs.Wait(); err != nil {
			return err
		}

		return nil
	},
	Args: cobra.MinimumNArgs(0),
}

func init() {
	rootCmd.AddCommand(viewsCmd)
}
