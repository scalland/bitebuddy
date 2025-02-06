package handlers

import (
	"bytes"
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/scalland/bitebuddy/pkg/log"
	"github.com/scalland/bitebuddy/pkg/utils"
	"html/template"
	"net/http"
	"path/filepath"
)

type WebHandlers struct {
	db          *sql.DB
	Log         *log.Logger
	u           *utils.Utils
	templatesFS *embed.FS
	tpl         *template.Template
	themeName   string
}

func NewWebHandlers(db *sql.DB, l *log.Logger, u *utils.Utils, tFS *embed.FS, tpl *template.Template, tName string) *WebHandlers {
	return &WebHandlers{
		db:          db,
		Log:         l,
		u:           u,
		templatesFS: tFS,
		tpl:         tpl,
		themeName:   tName,
	}
}

func (wh *WebHandlers) ExecuteTemplate(templateFileNameSansExtension string, data interface{}) (bytes.Buffer, error) {
	tmpl, tmplErr := wh.generateHTML(wh.themeName, templateFileNameSansExtension, data)
	if tmplErr != nil {
		return bytes.Buffer{}, tmplErr
	}

	return tmpl, nil
}

// GenerateHTML generates the complete HTML for a given page name
func (wh *WebHandlers) generateHTML(themeName, pageName string, data interface{}) (bytes.Buffer, error) {

	wh.Log.Debugf("template theme: %s", themeName)

	baseDir := "templates/" + themeName
	wh.Log.Debugf("base template dir: %s", baseDir)

	// Define the paths to the partials and page files
	layoutPath := filepath.Join(baseDir, "partials", "layout.html")
	headerPath := filepath.Join(baseDir, "partials", "header.html")
	footerPath := filepath.Join(baseDir, "partials", "footer.html")
	pagePath := filepath.Join(baseDir, "pages", pageName+".html")

	wh.Log.Debugf("template partial (base): %s", layoutPath)
	wh.Log.Debugf("template partial (header): %s", headerPath)
	wh.Log.Debugf("template partial (footer): %s", footerPath)
	wh.Log.Debugf("template partial (page): %s", pagePath)

	// Parse the templates
	tmpl, err := template.ParseFiles(layoutPath, headerPath, footerPath, pagePath)
	if err != nil {
		wh.Log.Debugf("Error parsing template: %s", err.Error())
		return bytes.Buffer{}, fmt.Errorf("error parsing templates: %s", err.Error())
	}

	// Parse templates from the embedded filesystem.
	tmpl, err = template.ParseFS(wh.templatesFS, layoutPath, headerPath, footerPath, pagePath)
	if err != nil {
		wh.Log.Fatal(err)
	}

	// Create a buffer to store the generated HTML
	var output bytes.Buffer
	err = tmpl.Execute(&output, data)
	if err != nil {
		wh.Log.Debugf("Error parsing template: %s", err.Error())
		return bytes.Buffer{}, fmt.Errorf("error executing template: %s", err.Error())
	}

	return output, nil
}

func (wh *WebHandlers) WriteHTML(w http.ResponseWriter, data bytes.Buffer) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	write, err := w.Write(data.Bytes())
	if err != nil {
		wh.Log.Debugf("WebHandler.WriteHTML: error parsing template: %s", err.Error())
		return
	}
	wh.Log.Debugf("%d bytes written to dashboard", write)
}

func (wh *WebHandlers) GetThemeName() string {
	return wh.themeName
}

func (wh *WebHandlers) GetTemplateFS() *embed.FS {
	return wh.templatesFS
}
