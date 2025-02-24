package cmd

import (
	"fmt"
	"github.com/scalland/bitebuddy/internal/handlers"
	"github.com/scalland/bitebuddy/internal/routes"
	"github.com/scalland/bitebuddy/pkg/log"
	"github.com/scalland/bitebuddy/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/srinathgs/mysqlstore"
	"log/slog"
	"net/http"
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

		// Filesystem Store
		//sessStore := sessions.NewFilesystemStore(viper.GetString("session_store_path"), []byte(sessSecret))

		sessionStoreDBTable := viper.GetString("session_db_table")
		sessionCookiePath := viper.GetString("session_cookie_path")
		sessionCookieAge := viper.GetInt("session_cookie_validity_mins")

		var sessionSameSite http.SameSite
		sameSiteStr := viper.GetString("session_cookie_same_site")
		switch sameSiteStr {
		case "default":
			sessionSameSite = http.SameSiteDefaultMode
		case "lax":
			sessionSameSite = http.SameSiteLaxMode
		case "strict":
			sessionSameSite = http.SameSiteStrictMode
		case "none":
			sessionSameSite = http.SameSiteNoneMode
		}

		sessStore, sessStoreErr := mysqlstore.NewMySQLStoreFromConnection(_db, sessionStoreDBTable, sessionCookiePath, sessionCookieAge*60, []byte(sessSecret))
		if sessStoreErr != nil {
			l.Fatalf("%s.cmd.serve: error connecting to session store: %s", utils.APP_NAME, sessStoreErr.Error())
		}

		sessStore.Options.HttpOnly = viper.GetBool("session_cookie_http_only")
		sessStore.Options.Secure = viper.GetBool("session_cookie_secure")
		sessStore.Options.Domain = viper.GetString("session_cookie_domain")
		sessStore.Options.SameSite = sessionSameSite

		sessStore.Options.Partitioned = viper.GetBool("session_cookie_partitioned")

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
