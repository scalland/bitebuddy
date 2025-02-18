package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/scalland/bitebuddy/pkg/log"
	"github.com/scalland/bitebuddy/pkg/utils"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
)

type WebHandlers struct {
	adminUserTypeID int
	db              *sql.DB
	Log             *log.Logger
	u               *utils.Utils
	templatesFS     *embed.FS
	tpl             *template.Template
	themeName       string
	store           *sessions.FilesystemStore
	sessionName     string
}

func NewWebHandlers(db *sql.DB, l *log.Logger, u *utils.Utils, tFS *embed.FS, sessionStore *sessions.FilesystemStore, tpl *template.Template, tName, sName string, adminUserTypeID int) *WebHandlers {
	return &WebHandlers{
		adminUserTypeID: adminUserTypeID,
		db:              db,
		Log:             l,
		u:               u,
		templatesFS:     tFS,
		store:           sessionStore,
		tpl:             tpl,
		themeName:       tName,
		sessionName:     sName,
	}
}

// GetSession - GetSession returns the values returned by gorilla/sessions.CookieStore.Get() function. It uses the sessionName which has
// been set in cmd/serve.go for this instance of the WebHandlers. For details on gorilla function, go to: https://pkg.go.dev/github.com/gorilla/sessions@v1.4.0#CookieStore.Get
func (wh *WebHandlers) GetSession(r *http.Request) (*sessions.Session, error) {
	return wh.store.Get(r, wh.sessionName)
}

func (wh *WebHandlers) IsLoggedIn(r *http.Request, w http.ResponseWriter) bool {
	session, err := wh.store.Get(r, wh.sessionName)
	if err != nil {
		wh.Log.Debugf("handlers.WebHandlers.IsLoggedIn: error getting session: %s", err.Error())
		return false
	}
	userID, okUserID := session.Values["user_id"].(int64)
	userTypeID, okUserType := session.Values["user_type_id"].(int)
	isLoggedIn, okIsLoggedIn := session.Values["is_logged_in"].(bool)
	wh.Log.Debugf("handlers.WebHandlers.IsLoggedIn: userID = %d, userTypeID = %d, isLoggedIn = %t", userID, userTypeID, isLoggedIn)
	var dbUserID, dbUserTypeID int
	if okIsLoggedIn && okUserType && okUserID && isLoggedIn == true {
		// user is logged-in according to session. Let us check if they exist in the DB
		wh.Log.Debugf("handlers.WebHandlers.IsLoggedIn: SQLEXEC: SELECT user_id, user_type_id FROM users WHERE user_id = %d AND user_type_id = %d", userID, userTypeID)
		err = wh.db.QueryRow("SELECT user_id, user_type_id FROM users WHERE user_id = ? AND user_type_id = ?", userID, userTypeID).Scan(&dbUserID, &dbUserTypeID)
		if err != nil {
			wh.Log.Debugf("handlers.WebHandlers.IsLoggedIn: error validating user in DB: %s", err.Error())
			return false
		}
		return true // user validated from both session and DB
	}
	wh.Log.Debugf("handlers.WebHandlers.IsLoggedIn: user could not be verified from session. Won't check in DB")
	wh.Log.Debugf("handler.WebHandlers.IsLoggedIn: checking if the session values are set or not and adding accordingly")
	if !okUserID {
		session.Values["user_id"] = 0
	}
	if !okUserType {
		session.Values["user_type_id"] = 0
	}
	if !okIsLoggedIn {
		session.Values["is_logged_in"] = false
	}
	err = session.Save(r, w)
	if err != nil {
		wh.Log.Debugf("handlers.WebHandlers.IsLoggedIn: error saving session: %s", err.Error())
	}
	return false // user could not be validated from session. DB was not checked
}

func (wh *WebHandlers) IsLoggedInAdmin(r *http.Request, w http.ResponseWriter) bool {
	session, err := wh.store.Get(r, wh.sessionName)
	if err != nil {
		wh.Log.Debugf("handlers.WebHandlers.IsLoggedIn: error getting session: %s", err.Error())
		return false
	}
	userID, okUserID := session.Values["user_id"].(int64)
	userTypeID, okUserType := session.Values["user_type_id"].(int)
	isLoggedIn, okIsLoggedIn := session.Values["is_logged_in"].(bool)
	wh.Log.Debugf("handlers.WebHandlers.IsLoggedInAdmin: userID = %d, userTypeID = %d, isLoggedIn = %t", userID, userTypeID, isLoggedIn)
	var dbUserID, dbUserTypeID int
	if okIsLoggedIn && okUserType && okUserID && isLoggedIn == true && userTypeID == wh.adminUserTypeID {
		// user is logged-in according to session and are an admin user. Let us check if they exist in the DB
		wh.Log.Debugf("handlers.WebHandlers.IsLoggedInAdmin: SQLEXEC: SELECT user_id, user_type_id FROM users WHERE user_id = %d AND user_type_id = %d", userID, userTypeID)
		err := wh.db.QueryRow("SELECT user_id, user_type_id FROM users WHERE user_id = ? AND user_type_id = ?", userID, userTypeID).Scan(&dbUserID, &dbUserTypeID)
		if err != nil {
			wh.Log.Debugf("handlers.WebHandlers.IsLoggedInAdmin: error validating user in DB: %s", err.Error())
			return false
		}
		return true // user validated from both session and DB
	}
	wh.Log.Debugf("handlers.WebHandlers.IsLoggedInAdmin: user could not be verified from session. Won't check in DB")
	wh.Log.Debugf("handler.WebHandlers.IsLoggedInAdmin: checking if the session values are set or not and adding accordingly")
	if !okUserID {
		session.Values["user_id"] = 0
	}
	if !okUserType {
		session.Values["user_type_id"] = 0
	}
	if !okIsLoggedIn {
		session.Values["is_logged_in"] = false
	}
	err = session.Save(r, w)
	if err != nil {
		wh.Log.Debugf("handlers.WebHandlers.IsLoggedInAdmin: error saving session: %s", err.Error())
	}
	return false // user could not be validated from session. DB was not checked
}

func (wh *WebHandlers) ReconnectDB() error {
	connErr := wh.db.PingContext(context.TODO())
	if connErr != nil {
		wh.Log.Errorf("handlers.WebHandlers.ReconnectDB: db is disconnected (%s)", connErr.Error())
		err := wh.db.Close()
		if err != nil {
			wh.Log.Errorf("handlers.WebHandlers.ReconnectDB: error closing database (%s)", err.Error())
		}
		wh.db = nil
		wh.db, err = wh.u.ConnectDB()
		if err != nil {
			wh.Log.Errorf("handlers.WebHandlers.ReconnectDB: error connecting to database (%s)", err.Error())
			return err
		}
		wh.Log.Debugf("handlers.WebHandlers.ReconnectDB: database is connected")
		return nil
	} else {
		wh.Log.Infof("handlers.WebHandlers.ReconnectDB: db connection is alive. won't reconnect")
		return nil
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

	baseDir := filepath.Join("templates", themeName)
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

func (wh *WebHandlers) WriteHTML(w http.ResponseWriter, data bytes.Buffer, httpStatus int) {
	switch httpStatus < 100 || httpStatus > 600 {
	case true:
		httpStatus = http.StatusOK
	}
	w.WriteHeader(httpStatus)
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

func (wh *WebHandlers) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request details here (e.g., using a logging library)
		wh.Log.Infof("Request: %s %s from %s", r.Method, r.URL, r.RemoteAddr)

		// Call the next handler
		next.ServeHTTP(w, r)
	})

}

// ServeStaticFilesWithContentType -Custom handler to wrap http.FileServer and set correct Content-Type for CSS
func (wh *WebHandlers) ServeStaticFilesWithContentType(fsys fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(fsys))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the requested file path ends with ".css"
		if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		}

		if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		}

		fileServer.ServeHTTP(w, r) // Call the original FileServer to handle the file serving
	})
}
