package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/scalland/bitebuddy/internal/handlers"
	"io/fs"
	"net/http"
)

func SetupRoutes(wh *handlers.WebHandlers) *mux.Router {

	router := mux.NewRouter()
	router.Use(wh.LoggerMiddleware)

	staticFilesDirPath := fmt.Sprintf("templates/%s/static", wh.GetThemeName())
	wh.Log.Infof("Serving static files from %s", staticFilesDirPath)
	staticFilesDirFS, staticFilesDirFSErr := fs.Sub(wh.GetTemplateFS(), staticFilesDirPath)
	if staticFilesDirFSErr != nil {
		wh.Log.Errorf("Error getting a sub-directory as a filesystem from webDir: %s", staticFilesDirFSErr.Error())
	}

	wh.Log.Infof("Serving static files from %s", staticFilesDirFS)
	// Serve static assets from your embedded FS.
	//router.PathPrefix("/static/*").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(staticFilesDirFS))))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", wh.ServeStaticFilesWithContentType(staticFilesDirFS)))

	// Public routes.
	router.HandleFunc("/login", wh.LoginHandler).Methods("GET", "POST")
	router.HandleFunc("/logout", wh.LogoutHandler).Methods("GET")

	// Protected routes.
	// Wrap them with RequireAuth middleware.

	router.Handle("/", wh.RequireAuth(http.HandlerFunc(wh.DashboardHandler))).Methods("GET")

	// User Types CRUD
	router.Handle("/user_types", wh.RequireAuth(http.HandlerFunc(wh.UserTypesHandler))).Methods("GET")
	// If user management routes should be restricted to admins:
	router.Handle("/user_types/new", wh.RequireAdmin(http.HandlerFunc(wh.UserTypesNewHandler))).Methods("GET", "POST")
	router.Handle("/user_types/edit", wh.RequireAdmin(http.HandlerFunc(wh.UserTypesEditHandler))).Methods("GET", "POST")
	router.Handle("/user_types/delete", wh.RequireAdmin(http.HandlerFunc(wh.UserTypesDeleteHandler))).Methods("POST")

	// Users CRUD
	router.Handle("/users", wh.RequireAuth(http.HandlerFunc(wh.UsersHandler))).Methods("GET")
	// If user management routes should be restricted to admins:
	router.Handle("/users/new", wh.RequireAdmin(http.HandlerFunc(wh.UserNewHandler))).Methods("GET", "POST")
	router.Handle("/users/edit", wh.RequireAdmin(http.HandlerFunc(wh.UserEditHandler))).Methods("GET", "POST")
	router.Handle("/users/delete", wh.RequireAdmin(http.HandlerFunc(wh.UserDeleteHandler))).Methods("POST")

	// Restaurants CRUD
	router.Handle("/restaurants", wh.RequireAuth(http.HandlerFunc(wh.RestaurantsHandler))).Methods("GET")
	// Admin-only routes
	router.Handle("/restaurants/new", wh.RequireAdmin(http.HandlerFunc(wh.RestaurantNewHandler))).Methods("GET", "POST")
	router.Handle("/restaurants/edit", wh.RequireAdmin(http.HandlerFunc(wh.RestaurantEditHandler))).Methods("GET", "POST")
	router.Handle("/restaurants/delete", wh.RequireAdmin(http.HandlerFunc(wh.RestaurantDeleteHandler))).Methods("POST")

	// Metrics CRUD
	router.Handle("/metrics", wh.RequireAuth(http.HandlerFunc(wh.MetricsHandler))).Methods("GET")
	// Admin-only routes
	router.Handle("/metrics/new", wh.RequireAdmin(http.HandlerFunc(wh.MetricNewHandler))).Methods("GET", "POST")
	router.Handle("/metrics/edit", wh.RequireAdmin(http.HandlerFunc(wh.MetricEditHandler))).Methods("GET", "POST")
	router.Handle("/metrics/delete", wh.RequireAdmin(http.HandlerFunc(wh.MetricDeleteHandler))).Methods("POST")

	// Reviews CRUD
	router.Handle("/reviews", wh.RequireAuth(http.HandlerFunc(wh.ReviewsHandler))).Methods("GET")
	// Admin-only routes
	router.Handle("/reviews/new", wh.RequireAdmin(http.HandlerFunc(wh.ReviewNewHandler))).Methods("GET", "POST")
	router.Handle("/reviews/edit", wh.RequireAdmin(http.HandlerFunc(wh.ReviewEditHandler))).Methods("GET", "POST")
	router.Handle("/reviews/delete", wh.RequireAdmin(http.HandlerFunc(wh.ReviewDeleteHandler))).Methods("POST")

	// Metric Reviews CRUD
	router.Handle("/metric_reviews", wh.RequireAuth(http.HandlerFunc(wh.MetricReviewsHandler))).Methods("GET")
	// Admin-only routes
	router.Handle("/metric_reviews/new", wh.RequireAdmin(http.HandlerFunc(wh.MetricReviewNewHandler))).Methods("GET", "POST")
	router.Handle("/metric_reviews/edit", wh.RequireAdmin(http.HandlerFunc(wh.MetricReviewEditHandler))).Methods("GET", "POST")
	router.Handle("/metric_reviews/delete", wh.RequireAdmin(http.HandlerFunc(wh.MetricReviewDeleteHandler))).Methods("POST")

	// Metric Reviews CRUD
	router.Handle("/filter_types", wh.RequireAuth(http.HandlerFunc(wh.FilterTypesHandler))).Methods("GET")
	// Admin-only routes
	router.Handle("/filter_types/new", wh.RequireAdmin(http.HandlerFunc(wh.FilterTypeNewHandler))).Methods("GET", "POST")
	router.Handle("/filter_types/edit", wh.RequireAdmin(http.HandlerFunc(wh.FilterTypeEditHandler))).Methods("GET", "POST")
	router.Handle("/filter_types/delete", wh.RequireAdmin(http.HandlerFunc(wh.FilterTypeDeleteHandler))).Methods("POST")

	// Filters CRUD
	router.Handle("/filters", wh.RequireAuth(http.HandlerFunc(wh.FiltersHandler))).Methods("GET")
	// Admin-only routes
	router.Handle("/filters/new", wh.RequireAdmin(http.HandlerFunc(wh.FilterNewHandler))).Methods("GET", "POST")
	router.Handle("/filters/edit", wh.RequireAdmin(http.HandlerFunc(wh.FilterEditHandler))).Methods("GET", "POST")
	router.Handle("/filters/delete", wh.RequireAdmin(http.HandlerFunc(wh.FiltersHandler))).Methods("POST")

	// Filters CRUD
	router.Handle("/otp_requests", wh.RequireAuth(http.HandlerFunc(wh.OtpRequestsHandler))).Methods("GET")
	// Admin-only routes
	router.Handle("/otp_requests/new", wh.RequireAdmin(http.HandlerFunc(wh.OtpRequestNewHandler))).Methods("GET", "POST")
	router.Handle("/otp_requests/edit", wh.RequireAdmin(http.HandlerFunc(wh.OtpRequestEditHandler))).Methods("GET", "POST")
	router.Handle("/otp_requests/delete", wh.RequireAdmin(http.HandlerFunc(wh.OtpRequestDeleteHandler))).Methods("POST")

	// Logout handler
	router.HandleFunc("/logout", wh.LogoutHandler).Methods("GET")

	// ALWAYS KEEP THIS HANDLER AS LAST otherwise it would override others see. https://stackoverflow.com/a/56937571/6670698
	// Re-define the default NotFound handler
	router.NotFoundHandler = router.NewRoute().HandlerFunc(http.NotFound).GetHandler()

	return router
}
