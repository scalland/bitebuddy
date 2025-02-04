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
	IsActive         bool
	CreatedAt        time.Time
	LastLogin        time.Time
	LastAccessedFrom string
}

// -----------------------------------------------------------------
// Users Handlers

func (o *WebHandlers) UsersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := o.db.Query("SELECT user_id, email, mobile_number, user_type_id, is_active, created_at, last_login, last_accessed_from FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Email, &u.MobileNumber, &u.UserTypeID, &u.IsActive, &u.CreatedAt, &u.LastLogin, &u.LastAccessedFrom)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}
	tmplErr := o.tpl.ExecuteTemplate(w, "users.html", users)
	if tmplErr != nil {
		slog.Error(fmt.Sprintf("Error executing template: users.html: %s", tmplErr.Error()))
	}
}

func (o *WebHandlers) UserNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmplErr := o.tpl.ExecuteTemplate(w, "user_form.html", nil)
		if tmplErr != nil {
			slog.Error(fmt.Sprintf("Error executing template: users.html: %s", tmplErr.Error()))
		}
		return
	}
	// POST
	email := r.FormValue("email")
	mobile := r.FormValue("mobile")
	userTypeStr := r.FormValue("user_type")
	userType, _ := strconv.Atoi(userTypeStr)
	isActive := r.FormValue("is_active") == "on"
	now := time.Now()
	stmt, err := o.db.Prepare("INSERT INTO users (email, mobile_number, user_type_id, is_active, created_at, last_login, last_accessed_from) VALUES (?, ?, ?, ?, ?, ?, ?)")
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

func (o *WebHandlers) UserEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var u User
		err := o.db.QueryRow("SELECT user_id, email, mobile_number, user_type_id, is_active, created_at, last_login, last_accessed_from FROM users WHERE user_id = ?", id).
			Scan(&u.ID, &u.Email, &u.MobileNumber, &u.UserTypeID, &u.IsActive, &u.CreatedAt, &u.LastLogin, &u.LastAccessedFrom)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmplErr := o.tpl.ExecuteTemplate(w, "user_form.html", u)
		if tmplErr != nil {
			slog.Error(fmt.Sprintf("Error executing template: user_form.html: %s", tmplErr.Error()))
		}
		return
	}
	// POST update
	email := r.FormValue("email")
	mobile := r.FormValue("mobile")
	userTypeStr := r.FormValue("user_type")
	userType, _ := strconv.Atoi(userTypeStr)
	isActive := r.FormValue("is_active") == "on"
	stmt, err := o.db.Prepare("UPDATE users SET email=?, mobile_number=?, user_type_id=?, is_active=? WHERE user_id=?")
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

func (o *WebHandlers) UserDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := o.db.Prepare("DELETE FROM users WHERE user_id=?")
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
