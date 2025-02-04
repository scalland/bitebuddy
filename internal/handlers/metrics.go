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

func (o *WebHandlers) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := o.db.Query("SELECT metric_id, metric_name, parent_metric_id, is_sub_metric, display_type_id, metric_type_id FROM metrics")
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
	o.tpl.ExecuteTemplate(w, "metrics.html", metrics)
}

func (o *WebHandlers) MetricNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		o.tpl.ExecuteTemplate(w, "metric_form.html", nil)
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
	stmt, err := o.db.Prepare("INSERT INTO metrics (metric_name, parent_metric_id, is_sub_metric, display_type_id, metric_type_id) VALUES (?, ?, ?, ?, ?)")
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

func (o *WebHandlers) MetricEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var m Metric
		err := o.db.QueryRow("SELECT metric_id, metric_name, parent_metric_id, is_sub_metric, display_type_id, metric_type_id FROM metrics WHERE metric_id=?", id).
			Scan(&m.ID, &m.MetricName, &m.ParentMetricID, &m.IsSubMetric, &m.DisplayTypeID, &m.MetricTypeID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		o.tpl.ExecuteTemplate(w, "metric_form.html", m)
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
	stmt, err := o.db.Prepare("UPDATE metrics SET metric_name=?, parent_metric_id=?, is_sub_metric=?, display_type_id=?, metric_type_id=? WHERE metric_id=?")
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

func (o *WebHandlers) MetricDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := o.db.Prepare("DELETE FROM metrics WHERE metric_id=?")
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
