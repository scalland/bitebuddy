package cmd

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/scalland/bitebuddy/pkg/utils"
	"html/template"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "bitebuddy",
		Short: "BiteBuddy is the Food/Restaurant discovery dashboard",
		Long:  `BiteBuddy is a fully responsive, mobile‑friendly web dashboard written in Go for managing your food/restaurant data.`,
	}
	_u          *utils.Utils
	TemplatesFS embed.FS
	_tpl        *template.Template
	_db         *sql.DB
)

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initializer)
}

func initializer() {
	_u = utils.NewUtils()
	env := "local"
	appName := utils.APP_NAME
	_u.ViperReadConfig(env, appName, "app.yml")
}
