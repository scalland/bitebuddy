package main

import (
	"embed"
	_ "github.com/go-sql-driver/mysql"
	"github.com/scalland/bitebuddy/cmd"
)

// -----------------------------------------------------------------
// Embed templates and static assets using Go embed

//go:embed templates/*
var templatesFS embed.FS

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
	cmd.TemplatesFS = templatesFS
	cmd.Execute()

}
