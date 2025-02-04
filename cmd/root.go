package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bitebuddy",
	Short: "BiteBuddy is the Food/Restaurant discovery dashboard",
	Long:  `BiteBuddy is a fully responsive, mobile‑friendly web dashboard written in Go for managing your food/restaurant data.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
