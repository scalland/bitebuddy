package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
)

type Metric struct {
	ID             int64
	MetricName     string
	ParentMetricID sql.NullInt64
	IsSubMetric    bool
	DisplayTypeID  int
	MetricTypeID   int
}

// -----------------------------------------------------------------
// Metrics Handlers

func (wh *WebHandlers) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := wh.db.Query("SELECT metric_id, metric_name, parent_metric_id, is_sub_metric, display_type_id, metric_type_id FROM metrics")
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
	tmpl, tmplErr := wh.ExecuteTemplate("metrics", metrics)
	if tmplErr != nil {
		http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
	}
	wh.WriteHTML(w, tmpl, http.StatusOK)
}

func (wh *WebHandlers) MetricNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, tmplErr := wh.ExecuteTemplate("metric_form.html", nil)
		if tmplErr != nil {
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
		}
		wh.WriteHTML(w, tmpl, http.StatusOK)
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
	stmt, err := wh.db.Prepare("INSERT INTO metrics (metric_name, parent_metric_id, is_sub_metric, display_type_id, metric_type_id) VALUES (?, ?, ?, ?, ?)")
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

func (wh *WebHandlers) MetricEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var m Metric
		err := wh.db.QueryRow("SELECT metric_id, metric_name, parent_metric_id, is_sub_metric, display_type_id, metric_type_id FROM metrics WHERE metric_id=?", id).
			Scan(&m.ID, &m.MetricName, &m.ParentMetricID, &m.IsSubMetric, &m.DisplayTypeID, &m.MetricTypeID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl, tmplErr := wh.ExecuteTemplate("metric_form", m)
		if tmplErr != nil {
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
		}
		wh.WriteHTML(w, tmpl, http.StatusOK)
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
	stmt, err := wh.db.Prepare("UPDATE metrics SET metric_name=?, parent_metric_id=?, is_sub_metric=?, display_type_id=?, metric_type_id=? WHERE metric_id=?")
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

func (wh *WebHandlers) MetricDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := wh.db.Prepare("DELETE FROM metrics WHERE metric_id=?")
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
