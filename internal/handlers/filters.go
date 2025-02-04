package handlers

import (
	"net/http"
	"strconv"
)

type Filter struct {
	ID           int64
	FilterTypeID int64
	FilterValue  string
}

// -----------------------------------------------------------------
// Filters Handlers

func (o *WebHandlers) FiltersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := o.db.Query("SELECT filter_id, filter_type_id, filter_value FROM filters")
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
	o.tpl.ExecuteTemplate(w, "filters.html", filters)
}

func (o *WebHandlers) FilterNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		o.tpl.ExecuteTemplate(w, "filter_form.html", nil)
		return
	}
	filterTypeID, _ := strconv.ParseInt(r.FormValue("filter_type_id"), 10, 64)
	filterValue := r.FormValue("filter_value")
	stmt, err := o.db.Prepare("INSERT INTO filters (filter_type_id, filter_value) VALUES (?, ?)")
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

func (o *WebHandlers) FilterEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var f Filter
		err := o.db.QueryRow("SELECT filter_id, filter_type_id, filter_value FROM filters WHERE filter_id=?", id).
			Scan(&f.ID, &f.FilterTypeID, &f.FilterValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		o.tpl.ExecuteTemplate(w, "filter_form.html", f)
		return
	}
	filterTypeID, _ := strconv.ParseInt(r.FormValue("filter_type_id"), 10, 64)
	filterValue := r.FormValue("filter_value")
	stmt, err := o.db.Prepare("UPDATE filters SET filter_type_id=?, filter_value=? WHERE filter_id=?")
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

func (o *WebHandlers) FilterDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := o.db.Prepare("DELETE FROM filters WHERE filter_id=?")
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
