package cmd

import (
	"log"

	"github.com/scalland/bitebuddy/pkg/db"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with dummy values",
	Run: func(cmd *cobra.Command, args []string) {
		if err := db.SeedDB(_u); err != nil {
			log.Fatalf("DB Seeding failed: %s", err.Error())
		}
		log.Println("DB Seeding successful!")
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}
