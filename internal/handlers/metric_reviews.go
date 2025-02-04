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

func (o *WebHandlers) MetricReviewsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := o.db.Query("SELECT metric_review_id, review_id, metric_id, score FROM metric_reviews")
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
	o.tpl.ExecuteTemplate(w, "metric_reviews.html", mrReviews)
}

func (o *WebHandlers) MetricReviewNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		o.tpl.ExecuteTemplate(w, "metric_review_form.html", nil)
		return
	}
	reviewID, _ := strconv.ParseInt(r.FormValue("review_id"), 10, 64)
	metricID, _ := strconv.ParseInt(r.FormValue("metric_id"), 10, 64)
	score, _ := strconv.ParseFloat(r.FormValue("score"), 64)
	stmt, err := o.db.Prepare("INSERT INTO metric_reviews (review_id, metric_id, score) VALUES (?, ?, ?)")
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

func (o *WebHandlers) MetricReviewEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var mr MetricReview
		err := o.db.QueryRow("SELECT metric_review_id, review_id, metric_id, score FROM metric_reviews WHERE metric_review_id=?", id).
			Scan(&mr.ID, &mr.ReviewID, &mr.MetricID, &mr.Score)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		o.tpl.ExecuteTemplate(w, "metric_review_form.html", mr)
		return
	}
	reviewID, _ := strconv.ParseInt(r.FormValue("review_id"), 10, 64)
	metricID, _ := strconv.ParseInt(r.FormValue("metric_id"), 10, 64)
	score, _ := strconv.ParseFloat(r.FormValue("score"), 64)
	stmt, err := o.db.Prepare("UPDATE metric_reviews SET review_id=?, metric_id=?, score=? WHERE metric_review_id=?")
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

func (o *WebHandlers) MetricReviewDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := o.db.Prepare("DELETE FROM metric_reviews WHERE metric_review_id=?")
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
