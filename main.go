package main

import (
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/scalland/bitebuddy/internal/handlers"
	"github.com/scalland/bitebuddy/internal/routes"
	"github.com/scalland/bitebuddy/pkg/utils"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
)

// -----------------------------------------------------------------
// Embed templates and static assets using Go embed

//go:embed templates/*
var templatesFS embed.FS

//go:embed static/*
var staticFS embed.FS

var tpl *template.Template
var db *sql.DB

// -----------------------------------------------------------------
// Data structures corresponding to the database tables

type RestaurantMetric struct {
	ID           int64
	RestaurantID int64
	MetricID     int64
	AverageScore float64
}

// -----------------------------------------------------------------
// main(): initialize DB, parse templates, define routes, and start server

func main() {
	env := "local"
	u := utils.NewUtils()
	appName := utils.APP_NAME
	u.ViperReadConfig(env, appName, "app.yml")
	var err error
	db, err = u.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	// Parse templates from the embedded filesystem.
	tpl, err = template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	wh := handlers.NewWebHandlers(db, u, tpl)

	routes.SetupRoutes(staticFS, wh)

	_appPort := viper.GetInt("app_port")

	log.Printf("Server starting on :%d", _appPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", _appPort), nil))
}
