package handlers

import (
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

// -----------------------------------------------------------------
// Users Handlers

func (wh *WebHandlers) UsersHandler(w http.ResponseWriter, r *http.Request) {
	reconnectErr := wh.ReconnectDB()
	if reconnectErr != nil {
		wh.Log.Errorf("Error reconnecting to database: %s", reconnectErr.Error())
		http.Error(w, reconnectErr.Error(), http.StatusInternalServerError)
		return
	}
	rows, err := wh.db.Query("SELECT user_id, email, mobile_number, u.user_type_id AS userTypeID, ut.user_type_name, is_active, created_at, last_login, last_accessed_from FROM users AS u LEFT JOIN user_types AS ut ON u.user_type_id=ut.user_type_id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

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
	tmpl, tmplErr := wh.ExecuteTemplate("users", users)
	if tmplErr != nil {
		slog.Error(fmt.Sprintf("Error executing template: users.html: %s", tmplErr.Error()))
		http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
	}
	wh.WriteHTML(w, tmpl, http.StatusOK)
}

func (wh *WebHandlers) UserNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, tmplErr := wh.ExecuteTemplate("user_form", nil)
		if tmplErr != nil {
			slog.Error(fmt.Sprintf("Error executing template: users: %s", tmplErr.Error()))
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
		}
		wh.WriteHTML(w, tmpl, http.StatusOK)
		return
	}
	// POST
	email := r.FormValue("email")
	mobile := r.FormValue("mobile")
	userTypeStr := r.FormValue("user_type")
	userType, _ := strconv.Atoi(userTypeStr)
	isActive := r.FormValue("is_active") == "on"
	now := time.Now()
	stmt, err := wh.db.Prepare("INSERT INTO users (email, mobile_number, u.user_type_id AS userTypeID, ut.user_type_name, is_active, created_at, last_login, last_accessed_from) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(email, mobile, userType, isActive, now, now, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

type UserEditData struct {
	U             User
	UserTypesData []UserTypes
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
