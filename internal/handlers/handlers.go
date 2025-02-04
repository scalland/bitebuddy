package handlers

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/scalland/bitebuddy/pkg/utils"
	"github.com/spf13/viper"
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

// DSN should be adjusted or read from configuration.
var DSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", viper.GetString("db_username"), viper.GetString("db_password"), viper.GetString("db_host"), viper.GetInt("db_port"), viper.GetString("db_database"))

//// User represents the users table.
//type User struct {
//	ID               int64
//	Email            string
//	MobileNumber     string
//	UserTypeID       int
//	IsActive         bool
//	CreatedAt        time.Time
//	LastLogin        time.Time
//	LastAccessedFrom string
//}

//// UsersHandler lists all users.
//func UsersHandler(w http.ResponseWriter, r *http.Request) {
//	db, err := sql.Open("mysql", DSN)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	defer db.Close()
//
//	rows, err := db.Query("SELECT user_id, email, mobile_number, user_type_id, is_active, created_at, last_login, last_accessed_from FROM users")
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	defer rows.Close()
//
//	var users []User
//	for rows.Next() {
//		var u User
//		err := rows.Scan(&u.ID, &u.Email, &u.MobileNumber, &u.UserTypeID, &u.IsActive, &u.CreatedAt, &u.LastLogin, &u.LastAccessedFrom)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		users = append(users, u)
//	}
//	tpl.ExecuteTemplate(w, "users.html", users)
//}

//// UserNewHandler handles creation of a new user.
//func UserNewHandler(w http.ResponseWriter, r *http.Request) {
//	if r.Method == http.MethodGet {
//		tpl.ExecuteTemplate(w, "user_form.html", nil)
//		return
//	}
//
//	// POST: process form submission.
//	email := r.FormValue("email")
//	mobile := r.FormValue("mobile")
//	userTypeStr := r.FormValue("user_type")
//	userType, _ := strconv.Atoi(userTypeStr)
//	isActive := r.FormValue("is_active") == "on"
//	now := time.Now()
//
//	db, err := sql.Open("mysql", DSN)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	defer db.Close()
//
//	stmt, err := db.Prepare("INSERT INTO users (email, mobile_number, user_type_id, is_active, created_at, last_login, last_accessed_from) VALUES (?, ?, ?, ?, ?, ?, ?)")
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	defer stmt.Close()
//	_, err = stmt.Exec(email, mobile, userType, isActive, now, now, "")
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	http.Redirect(w, r, "/users", http.StatusSeeOther)
//}

//// UserEditHandler handles updating a user.
//func UserEditHandler(w http.ResponseWriter, r *http.Request) {
//	db, err := sql.Open("mysql", DSN)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	defer db.Close()
//
//	idStr := r.URL.Query().Get("id")
//	id, _ := strconv.ParseInt(idStr, 10, 64)
//
//	if r.Method == http.MethodGet {
//		var u User
//		err := db.QueryRow("SELECT user_id, email, mobile_number, user_type_id, is_active, created_at, last_login, last_accessed_from FROM users WHERE user_id = ?", id).
//			Scan(&u.ID, &u.Email, &u.MobileNumber, &u.UserTypeID, &u.IsActive, &u.CreatedAt, &u.LastLogin, &u.LastAccessedFrom)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		tpl.ExecuteTemplate(w, "user_form.html", u)
//		return
//	}
//
//	// POST: update user.
//	email := r.FormValue("email")
//	mobile := r.FormValue("mobile")
//	userTypeStr := r.FormValue("user_type")
//	userType, _ := strconv.Atoi(userTypeStr)
//	isActive := r.FormValue("is_active") == "on"
//
//	stmt, err := db.Prepare("UPDATE users SET email=?, mobile_number=?, user_type_id=?, is_active=? WHERE user_id=?")
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	defer stmt.Close()
//	_, err = stmt.Exec(email, mobile, userType, isActive, id)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	http.Redirect(w, r, "/users", http.StatusSeeOther)
//}

//// UserDeleteHandler deletes a user.
//func UserDeleteHandler(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodPost {
//		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
//		return
//	}
//	db, err := sql.Open("mysql", DSN)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	defer db.Close()
//
//	idStr := r.FormValue("id")
//	id, _ := strconv.ParseInt(idStr, 10, 64)
//	stmt, err := db.Prepare("DELETE FROM users WHERE user_id=?")
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	defer stmt.Close()
//	_, err = stmt.Exec(id)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	http.Redirect(w, r, "/users", http.StatusSeeOther)
//}

// (Additional handlers for restaurants, metrics, reviews, etc. would follow a similar pattern.)
