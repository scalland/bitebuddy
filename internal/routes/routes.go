package routes

import (
	"github.com/scalland/bitebuddy/internal/handlers"
	"io/fs"
	"net/http"
)

//// SetupRoutes creates and returns a mux.Router with all your endpoints.
//func SetupRoutes() *mux.Router {
//	router := mux.NewRouter()
//
//	//// Serve static assets (using Go embed; adjust if needed)
//	//// (Assuming your static assets are embedded in internal/handlers or another package.)
//	//router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(handlers.StaticFS))))
//
//	// Serve static assets from the embedded FS.
//	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(assets.FS))))
//
//	// Dashboard
//	router.HandleFunc("/", handlers.DashboardHandler).Methods("GET")
//
//	// Users endpoints
//	router.HandleFunc("/users", handlers.UsersHandler).Methods("GET")
//	router.HandleFunc("/users/new", handlers.UserNewHandler).Methods("GET", "POST")
//	router.HandleFunc("/users/edit", handlers.UserEditHandler).Methods("GET", "POST")
//	router.HandleFunc("/users/delete", handlers.UserDeleteHandler).Methods("POST")
//
//	// (Repeat for Restaurants, Metrics, Reviews, Metric Reviews, Filter Types, Filters, OTP Requests)
//	// For example:
//	router.HandleFunc("/restaurants", handlers.RestaurantsHandler).Methods("GET")
//	router.HandleFunc("/restaurants/new", handlers.RestaurantNewHandler).Methods("GET", "POST")
//	router.HandleFunc("/restaurants/edit", handlers.RestaurantEditHandler).Methods("GET", "POST")
//	router.HandleFunc("/restaurants/delete", handlers.RestaurantDeleteHandler).Methods("POST")
//
//	// ... and so on for other endpoints.
//
//	return router
//}

func SetupRoutes(staticFS fs.FS, wh *handlers.WebHandlers) {

	//// Serve static files from the embedded static directory.
	//staticServer := http.FileServer(http.FS(staticFS))
	//http.Handle("/static/", http.StripPrefix("/static/", staticServer))

	// Serve static files.
	http.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Define routes:
	http.HandleFunc("/", wh.DashboardHandler)

	// Users CRUD
	http.HandleFunc("/users", wh.UsersHandler)
	http.HandleFunc("/users/new", wh.UserNewHandler)
	http.HandleFunc("/users/edit", wh.UserEditHandler)
	http.HandleFunc("/users/delete", wh.UserDeleteHandler)

	// Restaurants CRUD
	http.HandleFunc("/restaurants", wh.RestaurantsHandler)
	http.HandleFunc("/restaurants/new", wh.RestaurantNewHandler)
	http.HandleFunc("/restaurants/edit", wh.RestaurantEditHandler)
	http.HandleFunc("/restaurants/delete", wh.RestaurantDeleteHandler)

	// Metrics CRUD
	http.HandleFunc("/metrics", wh.MetricsHandler)
	http.HandleFunc("/metrics/new", wh.MetricNewHandler)
	http.HandleFunc("/metrics/edit", wh.MetricEditHandler)
	http.HandleFunc("/metrics/delete", wh.MetricDeleteHandler)

	// Reviews CRUD
	http.HandleFunc("/reviews", wh.ReviewsHandler)
	http.HandleFunc("/reviews/new", wh.ReviewNewHandler)
	http.HandleFunc("/reviews/edit", wh.ReviewEditHandler)
	http.HandleFunc("/reviews/delete", wh.ReviewDeleteHandler)

	// Metric Reviews CRUD
	http.HandleFunc("/metric_reviews", wh.MetricReviewsHandler)
	http.HandleFunc("/metric_reviews/new", wh.MetricReviewNewHandler)
	http.HandleFunc("/metric_reviews/edit", wh.MetricReviewEditHandler)
	http.HandleFunc("/metric_reviews/delete", wh.MetricReviewDeleteHandler)

	// Filter Types CRUD
	http.HandleFunc("/filter_types", wh.FilterTypesHandler)
	http.HandleFunc("/filter_types/new", wh.FilterTypeNewHandler)
	http.HandleFunc("/filter_types/edit", wh.FilterTypeEditHandler)
	http.HandleFunc("/filter_types/delete", wh.FilterTypeDeleteHandler)

	// Filters CRUD
	http.HandleFunc("/filters", wh.FiltersHandler)
	http.HandleFunc("/filters/new", wh.FilterNewHandler)
	http.HandleFunc("/filters/edit", wh.FilterEditHandler)
	http.HandleFunc("/filters/delete", wh.FiltersHandler)

	// OTP Requests CRUD
	http.HandleFunc("/otp_requests", wh.OtpRequestsHandler)
	http.HandleFunc("/otp_requests/new", wh.OtpRequestNewHandler)
	http.HandleFunc("/otp_requests/edit", wh.OtpRequestEditHandler)
	http.HandleFunc("/otp_requests/delete", wh.OtpRequestDeleteHandler)
}
