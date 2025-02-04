package handlers

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/scalland/bitebuddy/pkg/utils"
	"html/template"
)

type WebHandlers struct {
	db  *sql.DB
	u   *utils.Utils
	tpl *template.Template
}

func NewWebHandlers(db *sql.DB, u *utils.Utils, tpl *template.Template) *WebHandlers {
	return &WebHandlers{db: db, u: u, tpl: tpl}
}
