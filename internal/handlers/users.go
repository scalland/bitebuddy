package handlers

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	ID               int64
	Email            string
	MobileNumber     string
	UserTypeID       int
	UserTypeName     string
	IsActive         bool
	CreatedAt        time.Time
	LastLogin        time.Time
	LastAccessedFrom string
}

type UsersHandlerTemplateData struct {
	IsLoggedIn      bool
	IsLoggedInAdmin bool
	Errors          []string
	Users           []User
}

// -----------------------------------------------------------------
// Users Handlers

func (wh *WebHandlers) UsersHandler(w http.ResponseWriter, r *http.Request) {
	reconnectErr := wh.ReconnectDB()
	if reconnectErr != nil {
		wh.Log.Errorf("handlers.WebHandlers.UsersHandler: error reconnecting to database: %s", reconnectErr.Error())
		http.Error(w, reconnectErr.Error(), http.StatusInternalServerError)
		return
	}
	rows, err := wh.db.Query("SELECT user_id, email, mobile_number, u.user_type_id AS userTypeID, ut.user_type_name, is_active, created_at, last_login, last_accessed_from FROM users AS u LEFT JOIN user_types AS ut ON u.user_type_id=ut.user_type_id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			wh.Log.Errorf("handlers.WebHandlers.UsersHandler: error closing rows: %s", err.Error())
		}
	}(rows)

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Email, &u.MobileNumber, &u.UserTypeID, &u.UserTypeName, &u.IsActive, &u.CreatedAt, &u.LastLogin, &u.LastAccessedFrom)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u.UserTypeName = wh.UserTypeNameToString(u.UserTypeName)
		users = append(users, u)
	}

	templateData := UsersHandlerTemplateData{
		IsLoggedIn:      wh.isLoggedIn,
		IsLoggedInAdmin: wh.isAdmin,
		Errors:          []string{},
		Users:           users,
	}

	tmpl, tmplErr := wh.ExecuteTemplate("users", templateData)
	if tmplErr != nil {
		slog.Error(fmt.Sprintf("Error executing template: users.html: %s", tmplErr.Error()))
		http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
	}
	wh.WriteHTML(w, tmpl, http.StatusOK)
}

type UserEditData struct {
	IsLoggedIn      bool
	IsLoggedInAdmin bool
	Errors          []string
	U               User
	UserTypesData   []UserTypes
}

func (wh *WebHandlers) UserNewHandler(w http.ResponseWriter, r *http.Request) {
	templateData := UserEditData{
		IsLoggedIn:      wh.isLoggedIn,
		IsLoggedInAdmin: wh.isAdmin,
		Errors:          []string{},
		U: User{
			ID:               0,
			Email:            "",
			MobileNumber:     "",
			UserTypeID:       0,
			UserTypeName:     "",
			IsActive:         false,
			CreatedAt:        time.Time{},
			LastLogin:        time.Time{},
			LastAccessedFrom: "",
		},
		UserTypesData: nil,
	}
	var err error
	templateData.UserTypesData, err = wh.GetUserTypes()
	if err != nil {
		wh.Log.Errorf("handlers.WebHandlers.UserNewHandler: error getting user types data: %s", err.Error())
		templateData.Errors = append(templateData.Errors, fmt.Sprintf("error getting user types data: %s", err.Error()))
	}

	if r.Method == http.MethodGet {
		tmpl, tmplErr := wh.ExecuteTemplate("user_form", templateData)
		if tmplErr != nil {
			slog.Error(fmt.Sprintf("handlers.WebHandlers.UserNewHandler: error executing template: users: %s", tmplErr.Error()))
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
			return
		}
		wh.WriteHTML(w, tmpl, http.StatusOK)
		return
	}
	// POST
	templateData.U.Email = r.FormValue("email")
	templateData.U.MobileNumber = r.FormValue("mobile")
	userTypeStr := r.FormValue("user_type")
	templateData.U.UserTypeID = wh.u.Atoi(userTypeStr)
	templateData.U.IsActive = r.FormValue("is_active") == "on"
	lastAccessedFrom := "0.0.0.0"
	now := time.Now()
	stmt, err := wh.db.Prepare("INSERT INTO users(email, mobile_number, user_type_id, is_active, created_at, last_login, last_accessed_from) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		tErr := fmt.Sprintf("error creating user: %s", err.Error())
		templateData.Errors = append(templateData.Errors, tErr)
		wh.Log.Errorf("handlers.WebHandlers.UserNewHandler: %s", tErr)
		tmpl, tmplErr := wh.ExecuteTemplate("user_form", templateData)
		if tmplErr != nil {
			slog.Error(fmt.Sprintf("handlers.WebHandlers.UserNewHandler: error executing template: users: %s", tmplErr.Error()))
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
			return
		}
		wh.WriteHTML(w, tmpl, http.StatusOK)
		return
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			wh.Log.Errorf("handlers.WebHandlers.UserNewHandler: error closing rows: %s", err.Error())
		}
	}(stmt)
	_, err = stmt.Exec(templateData.U.Email, templateData.U.MobileNumber, templateData.U.UserTypeID, templateData.U.IsActive, now, now, lastAccessedFrom)
	if err != nil {
		tErr := fmt.Sprintf("error creating user: %s", err.Error())
		templateData.Errors = append(templateData.Errors, tErr)
		wh.Log.Errorf("handlers.WebHandlers.UserNewHandler: %s", tErr)
		tmpl, tmplErr := wh.ExecuteTemplate("user_form", templateData)
		if tmplErr != nil {
			slog.Error(fmt.Sprintf("handlers.WebHandlers.UserNewHandler: error executing template: users: %s", tmplErr.Error()))
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
			return
		}
		wh.WriteHTML(w, tmpl, http.StatusOK)
		return
	}
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (wh *WebHandlers) UserEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var (
			ued    UserEditData
			utdErr error
			u      = ued.U
		)
		err := wh.db.QueryRow("SELECT user_id, email, mobile_number, u.user_type_id, ut.user_type_name, is_active, created_at, last_login, last_accessed_from FROM users AS u LEFT JOIN user_types AS ut ON u.user_type_id=ut.user_type_id WHERE user_id = ?", id).
			Scan(&u.ID, &u.Email, &u.MobileNumber, &u.UserTypeID, &u.UserTypeName, &u.IsActive, &u.CreatedAt, &u.LastLogin, &u.LastAccessedFrom)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ued.U = u
		u.UserTypeName = wh.UserTypeNameToString(u.UserTypeName)
		ued.UserTypesData, utdErr = wh.GetUserTypes()
		if utdErr != nil {
			http.Error(w, utdErr.Error(), http.StatusInternalServerError)
		}
		tmpl, tmplErr := wh.ExecuteTemplate("user_form", ued)
		if tmplErr != nil {
			slog.Error(fmt.Sprintf("Error executing template: user_form.html: %s", tmplErr.Error()))
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
		}
		wh.WriteHTML(w, tmpl, http.StatusOK)
		return
	}
	// POST update
	email := r.FormValue("email")
	mobile := r.FormValue("mobile")
	userTypeStr := r.FormValue("user_type")
	userType, _ := strconv.Atoi(userTypeStr)
	isActive := r.FormValue("is_active") == "on"
	stmt, err := wh.db.Prepare("UPDATE users SET email=?, mobile_number=?, user_type_id=?, is_active=? WHERE user_id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(email, mobile, userType, isActive, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (wh *WebHandlers) UserDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := wh.db.Prepare("DELETE FROM users WHERE user_id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
