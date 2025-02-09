package cmd

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/scalland/bitebuddy/internal/handlers"
	"github.com/scalland/bitebuddy/pkg/log"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"

	"github.com/scalland/bitebuddy/internal/routes"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the BiteBuddy web server",
	Run: func(cmd *cobra.Command, args []string) {
		loggerOpts := &slog.HandlerOptions{
			AddSource:   false,
			Level:       log.LevelTrace,
			ReplaceAttr: log.CustomLogLevel,
		}
		l := log.New(loggerOpts)
		var err error
		_db, err = _u.ConnectDB()
		if err != nil {
			l.Fatal(err)
		}

		sessName := viper.GetString("session_name")
		sessSecret := viper.GetString("session_secret")
		//var sessStore interface{}
		//sessStoreType := viper.GetString("session_store_type")
		//switch sessStoreType {
		//case "cookie":
		//	sessStore = sessions.NewCookieStore([]byte(sessSecret))
		//case "filesystem":
		//	sessStore = sessions.NewFilesystemStore(viper.GetString("session_store_path"), []byte(sessSecret))
		//}

		// Initialize a session store (in production, use a secure key)
		//var sessStore = sessions.NewCookieStore([]byte(viper.GetString("session_secret")))

		sessStore := sessions.NewFilesystemStore(viper.GetString("session_store_path"), []byte(sessSecret))

		themeName := viper.GetString("theme")

		adminUserTypeID := viper.GetInt("admin_user_type_id")

		wh := handlers.NewWebHandlers(_db, l, _u, &TemplatesFS, sessStore, _tpl, themeName, sessName, adminUserTypeID)

		_appPort := viper.GetInt("app_port")

		l.Infof("Server starting on :%d", _appPort)
		l.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", _appPort), routes.SetupRoutes(wh)))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
