package handlers

import (
	"database/sql"
	"fmt"
	"github.com/scalland/bitebuddy/pkg/utils"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

type UserTypes struct {
	UserTypeID   int64
	UserType     string
	UserTypeName string
}

func (wh *WebHandlers) UserTypeNameToString(userTypeValue string) string {
	// 1. Strip leading and trailing underscores
	trimmedValue := strings.Trim(userTypeValue, "_")

	// 2. Split the string into words by underscore
	words := strings.Split(trimmedValue, "_")

	var capitalizedWords []string
	for _, word := range words {
		if len(word) > 0 { // Handle empty words if any (though unlikely after trim and split)
			// Capitalize the first letter and keep the rest lowercase
			runes := []rune(word)                // Convert string to rune slice to handle Unicode correctly
			runes[0] = unicode.ToUpper(runes[0]) // Capitalize the first rune
			capitalizedWords = append(capitalizedWords, string(runes))
		}
	}

	// 3. Join the capitalized words with spaces (or you can adjust as needed)
	return strings.Join(capitalizedWords, " ")
}

type UserTypesHandlerTemplateData struct {
	IsLoggedIn      bool
	IsLoggedInAdmin bool
	Errors          []string
	UserTypes       []UserTypes
}

// UserTypesHandler - This function handles the web requests for /user_types
func (wh *WebHandlers) UserTypesHandler(w http.ResponseWriter, r *http.Request) {
	userTypes, userTypesErr := wh.GetUserTypes()
	if userTypesErr != nil {
		wh.Log.Errorf("Error reconnecting to database: %s", userTypesErr.Error())
		http.Error(w, userTypesErr.Error(), http.StatusInternalServerError)
		return
	}

	templateData := UserTypesHandlerTemplateData{
		IsLoggedIn:      wh.isLoggedIn,
		IsLoggedInAdmin: wh.isAdmin,
		Errors:          nil,
		UserTypes:       userTypes,
	}

	tmpl, tmplErr := wh.ExecuteTemplate("user_types", templateData)
	if tmplErr != nil {
		slog.Error(fmt.Sprintf("Error executing template: users.html: %s", tmplErr.Error()))
		http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
	}
	wh.WriteHTML(w, tmpl, http.StatusOK)
}

func (wh *WebHandlers) GetUserTypes() ([]UserTypes, error) {
	tUserTypes := UserTypes{
		UserTypeID: 0,
		UserType:   "",
	}
	reconnectErr := wh.ReconnectDB()
	if reconnectErr != nil {
		wh.Log.Errorf("Error reconnecting to database: %s", reconnectErr.Error())
		return []UserTypes{tUserTypes}, reconnectErr
	}
	rows, err := wh.db.Query("SELECT user_type_id, user_type_name FROM user_types ORDER BY user_type_id ASC")
	if err != nil {
		return []UserTypes{tUserTypes}, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			wh.Log.Errorf("handlers.WebHandlers.GetUserTypes: error closing rows: %s", err.Error())
		}
	}(rows)

	var userTypes []UserTypes
	for rows.Next() {
		var u UserTypes
		err := rows.Scan(&u.UserTypeID, &u.UserType)
		if err != nil {
			return []UserTypes{tUserTypes}, err
		}
		u.UserTypeName = wh.UserTypeNameToString(u.UserType)
		userTypes = append(userTypes, u)
	}

	return userTypes, nil
}

type UserTypesNewHandlerTemplateData struct {
	IsLoggedIn      bool
	IsLoggedInAdmin bool
	Errors          []string
	UserTypes       UserTypes
}

func (wh *WebHandlers) UserTypesNewHandler(w http.ResponseWriter, r *http.Request) {
	userTypes := UserTypes{
		UserTypeID: 0,
		UserType:   "",
	}
	if r.Method == http.MethodGet {
		templateData := UserTypesNewHandlerTemplateData{
			IsLoggedIn:      wh.isLoggedIn,
			IsLoggedInAdmin: wh.isAdmin,
			Errors:          nil,
			UserTypes:       userTypes,
		}
		tmpl, tmplErr := wh.ExecuteTemplate("user_types_form", templateData)
		if tmplErr != nil {
			slog.Error(fmt.Sprintf("wh.UserTypesNewHandler: Error executing template: users: %s", tmplErr.Error()))
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
		}
		wh.WriteHTML(w, tmpl, http.StatusOK)
		return
	}
	// POST
	userTypeName := r.FormValue("usertypename")
	stmt, err := wh.db.Prepare("INSERT INTO user_types (user_type_name) VALUES (?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
		if err != nil {
			wh.Log.Errorf("handlers.WebHandlers.UserTypesNewHandler: error closing statement: %s", err.Error())
		}
	}(stmt)

	_, err = stmt.Exec(userTypeName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/user_types", http.StatusSeeOther)
}

type UserTypesEditHandlerTemplateData UserTypesNewHandlerTemplateData

func (wh *WebHandlers) UserTypesEditHandler(w http.ResponseWriter, r *http.Request) {
	wh.Log.Debugf("%s.handlers.UserTypesEditHandler: %s %s", utils.APP_NAME, r.Method, r.URL.Path)
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var u UserTypes
		err := wh.db.QueryRow("SELECT user_type_id, user_type_name FROM user_types WHERE user_type_id = ?", id).
			Scan(&u.UserTypeID, &u.UserTypeName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		templateData := UserTypesNewHandlerTemplateData{
			IsLoggedIn:      wh.isLoggedIn,
			IsLoggedInAdmin: wh.isAdmin,
			Errors:          nil,
			UserTypes:       u,
		}
		tmpl, tmplErr := wh.ExecuteTemplate("user_types_form", templateData)
		if tmplErr != nil {
			wh.Log.Errorf(fmt.Sprintf("Error executing template: user_form.html: %s", tmplErr.Error()))
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
		}
		wh.WriteHTML(w, tmpl, http.StatusOK)
		return
	}
	// POST update
	email := r.FormValue("usertypename")
	stmt, err := wh.db.Prepare("UPDATE user_types SET user_type_name=? WHERE user_type_id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
		if err != nil {
			wh.Log.Errorf("handlers.WebHandlers.UserTypesEditHandler: error closing statement: %s", err.Error())
		}
	}(stmt)
	_, err = stmt.Exec(email, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/user_types", http.StatusSeeOther)
}

func (wh *WebHandlers) UserTypesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	deleteRestrictedUserTypes := viper.GetString("delete_restricted_user_types")
	idStr := r.FormValue("id")

	if strings.Contains(deleteRestrictedUserTypes, idStr) {
		http.Error(w, fmt.Sprintf("User Type '%s' is configured as restricted from deletetion", idStr), http.StatusConflict)
		return
	}
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := wh.db.Prepare("DELETE FROM user_types WHERE user_type_id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
		if err != nil {
			wh.Log.Errorf("handlers.WebHandlers.UserTypesDeleteHandler: error closing statement: %s", err.Error())
		}
	}(stmt)

	_, err = stmt.Exec(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/user_types", http.StatusSeeOther)
}
