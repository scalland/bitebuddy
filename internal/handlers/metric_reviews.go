package handlers

import (
	"net/http"
	"strconv"
)

type MetricReview struct {
	ID       int64
	ReviewID int64
	MetricID int64
	Score    float64
}

// -----------------------------------------------------------------
// Metric Reviews Handlers

func (wh *WebHandlers) MetricReviewsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := wh.db.Query("SELECT metric_review_id, review_id, metric_id, score FROM metric_reviews")
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
	tmpl, tmplErr := wh.ExecuteTemplate("metric_reviews", mrReviews)
	if tmplErr != nil {
		http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
	}
	wh.WriteHTML(w, tmpl)
}

func (wh *WebHandlers) MetricReviewNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, tmplErr := wh.ExecuteTemplate("metric_review_form", nil)
		if tmplErr != nil {
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
		}
		wh.WriteHTML(w, tmpl)
		return
	}
	reviewID, _ := strconv.ParseInt(r.FormValue("review_id"), 10, 64)
	metricID, _ := strconv.ParseInt(r.FormValue("metric_id"), 10, 64)
	score, _ := strconv.ParseFloat(r.FormValue("score"), 64)
	stmt, err := wh.db.Prepare("INSERT INTO metric_reviews (review_id, metric_id, score) VALUES (?, ?, ?)")
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

func (wh *WebHandlers) MetricReviewEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var mr MetricReview
		err := wh.db.QueryRow("SELECT metric_review_id, review_id, metric_id, score FROM metric_reviews WHERE metric_review_id=?", id).
			Scan(&mr.ID, &mr.ReviewID, &mr.MetricID, &mr.Score)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl, tmplErr := wh.ExecuteTemplate("metric_review_form", mr)
		if tmplErr != nil {
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
		}
		wh.WriteHTML(w, tmpl)
		return
	}
	reviewID, _ := strconv.ParseInt(r.FormValue("review_id"), 10, 64)
	metricID, _ := strconv.ParseInt(r.FormValue("metric_id"), 10, 64)
	score, _ := strconv.ParseFloat(r.FormValue("score"), 64)
	stmt, err := wh.db.Prepare("UPDATE metric_reviews SET review_id=?, metric_id=?, score=? WHERE metric_review_id=?")
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

func (wh *WebHandlers) MetricReviewDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := wh.db.Prepare("DELETE FROM metric_reviews WHERE metric_review_id=?")
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
