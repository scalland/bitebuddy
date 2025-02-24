package cmd

import (
	"github.com/scalland/bitebuddy/pkg/utils"
	"log"

	"github.com/scalland/bitebuddy/pkg/db"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		createDB, createDBErr := cmd.Flags().GetBool("create-database")
		if createDBErr != nil {
			log.Printf("%s.cmd.migrateCmd: error reading value for --create-database: %s", utils.APP_NAME, createDBErr.Error())
			log.Printf("%s.cmd.migrateCmd: using default value of false", utils.APP_NAME)
			createDB = false
		}

		if err := db.MigrateDB(_u, createDB); err != nil {
			log.Fatalf("Migration failed: %s", err.Error())
		}
		log.Println("Migration successful!")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().BoolP("create-database", "d", false, "Set this to true if you want the database to be created before running migrations")
}
