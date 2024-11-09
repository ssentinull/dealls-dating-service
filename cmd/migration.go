package cmd

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func init() {
	migrateCmd.PersistentFlags().String("direction", "up", "migration direction")
	RootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate db",
	Long:  `This subcommand start the server`,
	Run:   runMigration,
}

func runMigration(cmd *cobra.Command, args []string) {
	migrations, err := migrate.New("file://docs/sql/migration", conf.SQL.Leader.DSN)
	if err != nil {
		logrus.Fatal("Failed to create new migrate: ", err)
	}

	direction := cmd.Flag("direction").Value.String()
	switch direction {
	case "up":
		err = migrations.Up()
	case "down":
		err = migrations.Down()
	default:
		logrus.WithField("direction", direction).Error("invalid direction: ", direction)
		return
	}

	if err != nil {
		logrus.WithField("direction", direction).Fatal("Failed to migrate database: ", err)
	}

	logrus.Infof("Applied migrations!")
}
