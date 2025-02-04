package cmd

import (
	"github.com/scalland/bitebuddy/pkg/utils"
	"log"
	"net/http"

	"github.com/scalland/bitebuddy/internal/routes"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the BiteBuddy web server",
	Run: func(cmd *cobra.Command, args []string) {
		env := "local"
		u := utils.NewUtils()
		appName := utils.APP_NAME
		u.ViperReadConfig(env, appName, "app.yml")
		
		router := routes.SetupRoutes()
		log.Println("Starting server on :8080")
		log.Fatal(http.ListenAndServe(":8080", router))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
