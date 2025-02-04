package main

import (
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/scalland/bitebuddy/pkg/utils"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
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

type Restaurant struct {
	ID                int64
	Name              string
	Address           string
	Latitude          float64
	Longitude         float64
	OverallRating     float64
	PriceForTwo       float64
	ImageURL          string
	DiscountAvailable bool
	AlcoholAvailable  bool
	PortionSizeLarge  bool
}

type Metric struct {
	ID             int64
	MetricName     string
	ParentMetricID sql.NullInt64
	IsSubMetric    bool
	DisplayTypeID  int
	MetricTypeID   int
}

type RestaurantMetric struct {
	ID           int64
	RestaurantID int64
	MetricID     int64
	AverageScore float64
}

type Review struct {
	ID           int64
	RestaurantID int64
	UserID       int64
	OverallScore float64
	ReviewText   string
	CreatedAt    time.Time
}

type MetricReview struct {
	ID       int64
	ReviewID int64
	MetricID int64
	Score    float64
}

type FilterType struct {
	ID             int64
	FilterTypeName string
}

type Filter struct {
	ID           int64
	FilterTypeID int64
	FilterValue  string
}

type OTPRequest struct {
	ID             int64
	UserID         int64
	OTPCode        string
	RequestedAt    time.Time
	DeliveryMethod string
	ValidTill      int64
}

// -----------------------------------------------------------------
// main(): initialize DB, parse templates, define routes, and start server

func main() {
	env := "local"
	u := utils.NewUtils()
	appName := utils.APP_NAME
	u.ViperReadConfig(env, appName, "app.yml")
	var err error

	// Connect to the database (adjust DSN as needed)
	// DSN format: username:password@tcp(host:port)/dbname?parseTime=true
	db, err = sql.Open(viper.GetString("db_driver"), fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", viper.GetString("db_username"), viper.GetString("db_password"), viper.GetString("db_host"), viper.GetInt("db_port"), viper.GetString("db_database")))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Parse templates from the embedded filesystem.
	tpl, err = template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// Serve static files.
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Define routes:
	http.HandleFunc("/", dashboardHandler)

	// Users CRUD
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/new", userNewHandler)
	http.HandleFunc("/users/edit", userEditHandler)
	http.HandleFunc("/users/delete", userDeleteHandler)

	// Restaurants CRUD
	http.HandleFunc("/restaurants", restaurantsHandler)
	http.HandleFunc("/restaurants/new", restaurantNewHandler)
	http.HandleFunc("/restaurants/edit", restaurantEditHandler)
	http.HandleFunc("/restaurants/delete", restaurantDeleteHandler)

	// Metrics CRUD
	http.HandleFunc("/metrics", metricsHandler)
	http.HandleFunc("/metrics/new", metricNewHandler)
	http.HandleFunc("/metrics/edit", metricEditHandler)
	http.HandleFunc("/metrics/delete", metricDeleteHandler)

	// Reviews CRUD
	http.HandleFunc("/reviews", reviewsHandler)
	http.HandleFunc("/reviews/new", reviewNewHandler)
	http.HandleFunc("/reviews/edit", reviewEditHandler)
	http.HandleFunc("/reviews/delete", reviewDeleteHandler)

	// Metric Reviews CRUD
	http.HandleFunc("/metric_reviews", metricReviewsHandler)
	http.HandleFunc("/metric_reviews/new", metricReviewNewHandler)
	http.HandleFunc("/metric_reviews/edit", metricReviewEditHandler)
	http.HandleFunc("/metric_reviews/delete", metricReviewDeleteHandler)

	// Filter Types CRUD
	http.HandleFunc("/filter_types", filterTypesHandler)
	http.HandleFunc("/filter_types/new", filterTypeNewHandler)
	http.HandleFunc("/filter_types/edit", filterTypeEditHandler)
	http.HandleFunc("/filter_types/delete", filterTypeDeleteHandler)

	// Filters CRUD
	http.HandleFunc("/filters", filtersHandler)
	http.HandleFunc("/filters/new", filterNewHandler)
	http.HandleFunc("/filters/edit", filterEditHandler)
	http.HandleFunc("/filters/delete", filterDeleteHandler)

	// OTP Requests CRUD
	http.HandleFunc("/otp_requests", otpRequestsHandler)
	http.HandleFunc("/otp_requests/new", otpRequestNewHandler)
	http.HandleFunc("/otp_requests/edit", otpRequestEditHandler)
	http.HandleFunc("/otp_requests/delete", otpRequestDeleteHandler)

	_appPort := viper.GetInt("app_port")

	log.Printf("Server starting on :%d", _appPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", _appPort), nil))
}

// -----------------------------------------------------------------
// Dashboard handler

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "dashboard.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// -----------------------------------------------------------------
// Users Handlers

func usersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT user_id, email, mobile_number, user_type_id, is_active, created_at, last_login, last_accessed_from FROM users")
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
	tpl.ExecuteTemplate(w, "users.html", users)
}

func userNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "user_form.html", nil)
		return
	}
	// POST
	email := r.FormValue("email")
	mobile := r.FormValue("mobile")
	userTypeStr := r.FormValue("user_type")
	userType, _ := strconv.Atoi(userTypeStr)
	isActive := r.FormValue("is_active") == "on"
	now := time.Now()
	stmt, err := db.Prepare("INSERT INTO users (email, mobile_number, user_type_id, is_active, created_at, last_login, last_accessed_from) VALUES (?, ?, ?, ?, ?, ?, ?)")
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

func userEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var u User
		err := db.QueryRow("SELECT user_id, email, mobile_number, user_type_id, is_active, created_at, last_login, last_accessed_from FROM users WHERE user_id = ?", id).
			Scan(&u.ID, &u.Email, &u.MobileNumber, &u.UserTypeID, &u.IsActive, &u.CreatedAt, &u.LastLogin, &u.LastAccessedFrom)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tpl.ExecuteTemplate(w, "user_form.html", u)
		return
	}
	// POST update
	email := r.FormValue("email")
	mobile := r.FormValue("mobile")
	userTypeStr := r.FormValue("user_type")
	userType, _ := strconv.Atoi(userTypeStr)
	isActive := r.FormValue("is_active") == "on"
	stmt, err := db.Prepare("UPDATE users SET email=?, mobile_number=?, user_type_id=?, is_active=? WHERE user_id=?")
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

func userDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := db.Prepare("DELETE FROM users WHERE user_id=?")
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

// -----------------------------------------------------------------
// Restaurants Handlers

func restaurantsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT restaurant_id, name, address, latitude, longitude, overall_rating, price_for_two, image_url, discount_available, alcohol_available, portion_size_large FROM restaurants")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var restaurants []Restaurant
	for rows.Next() {
		var rct Restaurant
		err := rows.Scan(&rct.ID, &rct.Name, &rct.Address, &rct.Latitude, &rct.Longitude, &rct.OverallRating, &rct.PriceForTwo, &rct.ImageURL, &rct.DiscountAvailable, &rct.AlcoholAvailable, &rct.PortionSizeLarge)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		restaurants = append(restaurants, rct)
	}
	tpl.ExecuteTemplate(w, "restaurants.html", restaurants)
}

func restaurantNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "restaurant_form.html", nil)
		return
	}
	// POST
	name := r.FormValue("name")
	address := r.FormValue("address")
	lat, _ := strconv.ParseFloat(r.FormValue("latitude"), 64)
	lng, _ := strconv.ParseFloat(r.FormValue("longitude"), 64)
	overallRating, _ := strconv.ParseFloat(r.FormValue("overall_rating"), 64)
	priceForTwo, _ := strconv.ParseFloat(r.FormValue("price_for_two"), 64)
	imageURL := r.FormValue("image_url")
	discountAvailable := r.FormValue("discount_available") == "on"
	alcoholAvailable := r.FormValue("alcohol_available") == "on"
	portionSizeLarge := r.FormValue("portion_size_large") == "on"
	stmt, err := db.Prepare("INSERT INTO restaurants (name, address, latitude, longitude, overall_rating, price_for_two, image_url, discount_available, alcohol_available, portion_size_large) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(name, address, lat, lng, overallRating, priceForTwo, imageURL, discountAvailable, alcoholAvailable, portionSizeLarge)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/restaurants", http.StatusSeeOther)
}

func restaurantEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var rct Restaurant
		err := db.QueryRow("SELECT restaurant_id, name, address, latitude, longitude, overall_rating, price_for_two, image_url, discount_available, alcohol_available, portion_size_large FROM restaurants WHERE restaurant_id=?", id).
			Scan(&rct.ID, &rct.Name, &rct.Address, &rct.Latitude, &rct.Longitude, &rct.OverallRating, &rct.PriceForTwo, &rct.ImageURL, &rct.DiscountAvailable, &rct.AlcoholAvailable, &rct.PortionSizeLarge)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tpl.ExecuteTemplate(w, "restaurant_form.html", rct)
		return
	}
	// POST update
	name := r.FormValue("name")
	address := r.FormValue("address")
	lat, _ := strconv.ParseFloat(r.FormValue("latitude"), 64)
	lng, _ := strconv.ParseFloat(r.FormValue("longitude"), 64)
	overallRating, _ := strconv.ParseFloat(r.FormValue("overall_rating"), 64)
	priceForTwo, _ := strconv.ParseFloat(r.FormValue("price_for_two"), 64)
	imageURL := r.FormValue("image_url")
	discountAvailable := r.FormValue("discount_available") == "on"
	alcoholAvailable := r.FormValue("alcohol_available") == "on"
	portionSizeLarge := r.FormValue("portion_size_large") == "on"
	stmt, err := db.Prepare("UPDATE restaurants SET name=?, address=?, latitude=?, longitude=?, overall_rating=?, price_for_two=?, image_url=?, discount_available=?, alcohol_available=?, portion_size_large=? WHERE restaurant_id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(name, address, lat, lng, overallRating, priceForTwo, imageURL, discountAvailable, alcoholAvailable, portionSizeLarge, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/restaurants", http.StatusSeeOther)
}

func restaurantDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := db.Prepare("DELETE FROM restaurants WHERE restaurant_id=?")
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
	http.Redirect(w, r, "/restaurants", http.StatusSeeOther)
}

// -----------------------------------------------------------------
// Metrics Handlers

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT metric_id, metric_name, parent_metric_id, is_sub_metric, display_type_id, metric_type_id FROM metrics")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var metrics []Metric
	for rows.Next() {
		var m Metric
		err := rows.Scan(&m.ID, &m.MetricName, &m.ParentMetricID, &m.IsSubMetric, &m.DisplayTypeID, &m.MetricTypeID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		metrics = append(metrics, m)
	}
	tpl.ExecuteTemplate(w, "metrics.html", metrics)
}

func metricNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "metric_form.html", nil)
		return
	}
	metricName := r.FormValue("metric_name")
	parentStr := r.FormValue("parent_metric_id")
	var parentID sql.NullInt64
	if parentStr != "" {
		idVal, _ := strconv.ParseInt(parentStr, 10, 64)
		parentID = sql.NullInt64{Int64: idVal, Valid: true}
	}
	isSubMetric := r.FormValue("is_sub_metric") == "on"
	displayTypeID, _ := strconv.Atoi(r.FormValue("display_type_id"))
	metricTypeID, _ := strconv.Atoi(r.FormValue("metric_type_id"))
	stmt, err := db.Prepare("INSERT INTO metrics (metric_name, parent_metric_id, is_sub_metric, display_type_id, metric_type_id) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(metricName, parentID, isSubMetric, displayTypeID, metricTypeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/metrics", http.StatusSeeOther)
}

func metricEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var m Metric
		err := db.QueryRow("SELECT metric_id, metric_name, parent_metric_id, is_sub_metric, display_type_id, metric_type_id FROM metrics WHERE metric_id=?", id).
			Scan(&m.ID, &m.MetricName, &m.ParentMetricID, &m.IsSubMetric, &m.DisplayTypeID, &m.MetricTypeID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tpl.ExecuteTemplate(w, "metric_form.html", m)
		return
	}
	metricName := r.FormValue("metric_name")
	parentStr := r.FormValue("parent_metric_id")
	var parentID sql.NullInt64
	if parentStr != "" {
		idVal, _ := strconv.ParseInt(parentStr, 10, 64)
		parentID = sql.NullInt64{Int64: idVal, Valid: true}
	}
	isSubMetric := r.FormValue("is_sub_metric") == "on"
	displayTypeID, _ := strconv.Atoi(r.FormValue("display_type_id"))
	metricTypeID, _ := strconv.Atoi(r.FormValue("metric_type_id"))
	stmt, err := db.Prepare("UPDATE metrics SET metric_name=?, parent_metric_id=?, is_sub_metric=?, display_type_id=?, metric_type_id=? WHERE metric_id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(metricName, parentID, isSubMetric, displayTypeID, metricTypeID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/metrics", http.StatusSeeOther)
}

func metricDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := db.Prepare("DELETE FROM metrics WHERE metric_id=?")
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
	http.Redirect(w, r, "/metrics", http.StatusSeeOther)
}

// -----------------------------------------------------------------
// Reviews Handlers

func reviewsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT review_id, restaurant_id, user_id, overall_score, review_text, created_at FROM reviews")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var reviews []Review
	for rows.Next() {
		var rev Review
		err := rows.Scan(&rev.ID, &rev.RestaurantID, &rev.UserID, &rev.OverallScore, &rev.ReviewText, &rev.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reviews = append(reviews, rev)
	}
	tpl.ExecuteTemplate(w, "reviews.html", reviews)
}

func reviewNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "review_form.html", nil)
		return
	}
	restaurantID, _ := strconv.ParseInt(r.FormValue("restaurant_id"), 10, 64)
	userID, _ := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	overallScore, _ := strconv.ParseFloat(r.FormValue("overall_score"), 64)
	reviewText := r.FormValue("review_text")
	now := time.Now()
	stmt, err := db.Prepare("INSERT INTO reviews (restaurant_id, user_id, overall_score, review_text, created_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(restaurantID, userID, overallScore, reviewText, now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/reviews", http.StatusSeeOther)
}

func reviewEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var rev Review
		err := db.QueryRow("SELECT review_id, restaurant_id, user_id, overall_score, review_text, created_at FROM reviews WHERE review_id=?", id).
			Scan(&rev.ID, &rev.RestaurantID, &rev.UserID, &rev.OverallScore, &rev.ReviewText, &rev.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tpl.ExecuteTemplate(w, "review_form.html", rev)
		return
	}
	restaurantID, _ := strconv.ParseInt(r.FormValue("restaurant_id"), 10, 64)
	userID, _ := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	overallScore, _ := strconv.ParseFloat(r.FormValue("overall_score"), 64)
	reviewText := r.FormValue("review_text")
	stmt, err := db.Prepare("UPDATE reviews SET restaurant_id=?, user_id=?, overall_score=?, review_text=? WHERE review_id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(restaurantID, userID, overallScore, reviewText, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/reviews", http.StatusSeeOther)
}

func reviewDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := db.Prepare("DELETE FROM reviews WHERE review_id=?")
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
	http.Redirect(w, r, "/reviews", http.StatusSeeOther)
}

// -----------------------------------------------------------------
// Metric Reviews Handlers

func metricReviewsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT metric_review_id, review_id, metric_id, score FROM metric_reviews")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var mrReviews []MetricReview
	for rows.Next() {
		var mr MetricReview
		err := rows.Scan(&mr.ID, &mr.ReviewID, &mr.MetricID, &mr.Score)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mrReviews = append(mrReviews, mr)
	}
	tpl.ExecuteTemplate(w, "metric_reviews.html", mrReviews)
}

func metricReviewNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "metric_review_form.html", nil)
		return
	}
	reviewID, _ := strconv.ParseInt(r.FormValue("review_id"), 10, 64)
	metricID, _ := strconv.ParseInt(r.FormValue("metric_id"), 10, 64)
	score, _ := strconv.ParseFloat(r.FormValue("score"), 64)
	stmt, err := db.Prepare("INSERT INTO metric_reviews (review_id, metric_id, score) VALUES (?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(reviewID, metricID, score)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/metric_reviews", http.StatusSeeOther)
}

func metricReviewEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var mr MetricReview
		err := db.QueryRow("SELECT metric_review_id, review_id, metric_id, score FROM metric_reviews WHERE metric_review_id=?", id).
			Scan(&mr.ID, &mr.ReviewID, &mr.MetricID, &mr.Score)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tpl.ExecuteTemplate(w, "metric_review_form.html", mr)
		return
	}
	reviewID, _ := strconv.ParseInt(r.FormValue("review_id"), 10, 64)
	metricID, _ := strconv.ParseInt(r.FormValue("metric_id"), 10, 64)
	score, _ := strconv.ParseFloat(r.FormValue("score"), 64)
	stmt, err := db.Prepare("UPDATE metric_reviews SET review_id=?, metric_id=?, score=? WHERE metric_review_id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(reviewID, metricID, score, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/metric_reviews", http.StatusSeeOther)
}

func metricReviewDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := db.Prepare("DELETE FROM metric_reviews WHERE metric_review_id=?")
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
	http.Redirect(w, r, "/metric_reviews", http.StatusSeeOther)
}

// -----------------------------------------------------------------
// Filter Types Handlers

func filterTypesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT filter_type_id, filter_type_name FROM filter_types")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var fTypes []FilterType
	for rows.Next() {
		var ft FilterType
		err := rows.Scan(&ft.ID, &ft.FilterTypeName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fTypes = append(fTypes, ft)
	}
	tpl.ExecuteTemplate(w, "filter_types.html", fTypes)
}

func filterTypeNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "filter_type_form.html", nil)
		return
	}
	filterTypeName := r.FormValue("filter_type_name")
	stmt, err := db.Prepare("INSERT INTO filter_types (filter_type_name) VALUES (?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(filterTypeName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/filter_types", http.StatusSeeOther)
}

func filterTypeEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var ft FilterType
		err := db.QueryRow("SELECT filter_type_id, filter_type_name FROM filter_types WHERE filter_type_id=?", id).
			Scan(&ft.ID, &ft.FilterTypeName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tpl.ExecuteTemplate(w, "filter_type_form.html", ft)
		return
	}
	filterTypeName := r.FormValue("filter_type_name")
	stmt, err := db.Prepare("UPDATE filter_types SET filter_type_name=? WHERE filter_type_id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(filterTypeName, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/filter_types", http.StatusSeeOther)
}

func filterTypeDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := db.Prepare("DELETE FROM filter_types WHERE filter_type_id=?")
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
	http.Redirect(w, r, "/filter_types", http.StatusSeeOther)
}

// -----------------------------------------------------------------
// Filters Handlers

func filtersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT filter_id, filter_type_id, filter_value FROM filters")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var filters []Filter
	for rows.Next() {
		var f Filter
		err := rows.Scan(&f.ID, &f.FilterTypeID, &f.FilterValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		filters = append(filters, f)
	}
	tpl.ExecuteTemplate(w, "filters.html", filters)
}

func filterNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "filter_form.html", nil)
		return
	}
	filterTypeID, _ := strconv.ParseInt(r.FormValue("filter_type_id"), 10, 64)
	filterValue := r.FormValue("filter_value")
	stmt, err := db.Prepare("INSERT INTO filters (filter_type_id, filter_value) VALUES (?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(filterTypeID, filterValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/filters", http.StatusSeeOther)
}

func filterEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var f Filter
		err := db.QueryRow("SELECT filter_id, filter_type_id, filter_value FROM filters WHERE filter_id=?", id).
			Scan(&f.ID, &f.FilterTypeID, &f.FilterValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tpl.ExecuteTemplate(w, "filter_form.html", f)
		return
	}
	filterTypeID, _ := strconv.ParseInt(r.FormValue("filter_type_id"), 10, 64)
	filterValue := r.FormValue("filter_value")
	stmt, err := db.Prepare("UPDATE filters SET filter_type_id=?, filter_value=? WHERE filter_id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(filterTypeID, filterValue, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/filters", http.StatusSeeOther)
}

func filterDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := db.Prepare("DELETE FROM filters WHERE filter_id=?")
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
	http.Redirect(w, r, "/filters", http.StatusSeeOther)
}

// -----------------------------------------------------------------
// OTP Requests Handlers

func otpRequestsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT otp_request_id, user_id, otp_code, requested_at, delivery_method, valid_till FROM otp_requests")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var otps []OTPRequest
	for rows.Next() {
		var o OTPRequest
		err := rows.Scan(&o.ID, &o.UserID, &o.OTPCode, &o.RequestedAt, &o.DeliveryMethod, &o.ValidTill)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		otps = append(otps, o)
	}
	tpl.ExecuteTemplate(w, "otp_requests.html", otps)
}

func otpRequestNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "otp_request_form.html", nil)
		return
	}
	userID, _ := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	otpCode := r.FormValue("otp_code")
	requestedAt := time.Now()
	deliveryMethod := r.FormValue("delivery_method")
	// Note: valid_till is a generated column, so we do not insert it.
	stmt, err := db.Prepare("INSERT INTO otp_requests (user_id, otp_code, requested_at, delivery_method) VALUES (?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(userID, otpCode, requestedAt, deliveryMethod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/otp_requests", http.StatusSeeOther)
}

func otpRequestEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var o OTPRequest
		err := db.QueryRow("SELECT otp_request_id, user_id, otp_code, requested_at, delivery_method, valid_till FROM otp_requests WHERE otp_request_id=?", id).
			Scan(&o.ID, &o.UserID, &o.OTPCode, &o.RequestedAt, &o.DeliveryMethod, &o.ValidTill)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tpl.ExecuteTemplate(w, "otp_request_form.html", o)
		return
	}
	userID, _ := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	otpCode := r.FormValue("otp_code")
	requestedAt, _ := time.Parse("2006-01-02 15:04:05", r.FormValue("requested_at"))
	deliveryMethod := r.FormValue("delivery_method")
	stmt, err := db.Prepare("UPDATE otp_requests SET user_id=?, otp_code=?, requested_at=?, delivery_method=? WHERE otp_request_id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(userID, otpCode, requestedAt, deliveryMethod, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/otp_requests", http.StatusSeeOther)
}

func otpRequestDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := db.Prepare("DELETE FROM otp_requests WHERE otp_request_id=?")
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
	http.Redirect(w, r, "/otp_requests", http.StatusSeeOther)
}
