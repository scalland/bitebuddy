package cmd

import (
	"fmt"
	"github.com/scalland/bitebuddy/internal/handlers"
	"github.com/scalland/bitebuddy/pkg/log"
	"github.com/scalland/bitebuddy/pkg/utils"
	"github.com/spf13/viper"
	"net/http"

	"github.com/scalland/bitebuddy/internal/routes"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the BiteBuddy web server",
	Run: func(cmd *cobra.Command, args []string) {
		env := "local"
		appName := utils.APP_NAME
		_u.ViperReadConfig(env, appName, "app.yml")
		l := log.New()
		var err error
		_db, err = _u.ConnectDB()
		if err != nil {
			l.Fatal(err)
		}

		themeName := viper.GetString("theme")

		wh := handlers.NewWebHandlers(_db, l, _u, &TemplatesFS, _tpl, themeName)

		routes.SetupRoutes(wh)

		_appPort := viper.GetInt("app_port")

		l.Infof("Server starting on :%d", _appPort)
		l.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", _appPort), nil))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
