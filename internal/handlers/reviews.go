package handlers

import (
	"net/http"
	"strconv"
	"time"
)

type Review struct {
	ID           int64
	RestaurantID int64
	UserID       int64
	OverallScore float64
	ReviewText   string
	CreatedAt    time.Time
}

// -----------------------------------------------------------------
// Reviews Handlers

func (wh *WebHandlers) ReviewsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := wh.db.Query("SELECT review_id, restaurant_id, user_id, overall_score, review_text, created_at FROM reviews")
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
	tmpl, tmplErr := wh.ExecuteTemplate("reviews", reviews)
	if tmplErr != nil {
		http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
	}
	wh.WriteHTML(w, tmpl)
}

func (wh *WebHandlers) ReviewNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, tmplErr := wh.ExecuteTemplate("review_form", nil)
		if tmplErr != nil {
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
		}
		wh.WriteHTML(w, tmpl)
		return
	}
	restaurantID, _ := strconv.ParseInt(r.FormValue("restaurant_id"), 10, 64)
	userID, _ := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	overallScore, _ := strconv.ParseFloat(r.FormValue("overall_score"), 64)
	reviewText := r.FormValue("review_text")
	now := time.Now()
	stmt, err := wh.db.Prepare("INSERT INTO reviews (restaurant_id, user_id, overall_score, review_text, created_at) VALUES (?, ?, ?, ?, ?)")
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

func (wh *WebHandlers) ReviewEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var rev Review
		err := wh.db.QueryRow("SELECT review_id, restaurant_id, user_id, overall_score, review_text, created_at FROM reviews WHERE review_id=?", id).
			Scan(&rev.ID, &rev.RestaurantID, &rev.UserID, &rev.OverallScore, &rev.ReviewText, &rev.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl, tmplErr := wh.ExecuteTemplate("review_form", rev)
		if tmplErr != nil {
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
		}
		wh.WriteHTML(w, tmpl)
		return
	}
	restaurantID, _ := strconv.ParseInt(r.FormValue("restaurant_id"), 10, 64)
	userID, _ := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	overallScore, _ := strconv.ParseFloat(r.FormValue("overall_score"), 64)
	reviewText := r.FormValue("review_text")
	stmt, err := wh.db.Prepare("UPDATE reviews SET restaurant_id=?, user_id=?, overall_score=?, review_text=? WHERE review_id=?")
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

func (wh *WebHandlers) ReviewDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := wh.db.Prepare("DELETE FROM reviews WHERE review_id=?")
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
