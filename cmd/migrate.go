package cmd

import (
	"log"

	"github.com/scalland/bitebuddy/pkg/db"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		if err := db.MigrateDB(_u); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migration successful!")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
