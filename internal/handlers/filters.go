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

func (wh *WebHandlers) FiltersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := wh.db.Query("SELECT filter_id, filter_type_id, filter_value FROM filters")
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
	tmpl, tmplErr := wh.ExecuteTemplate("filters", filters)
	if tmplErr != nil {
		http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
	}
	wh.WriteHTML(w, tmpl, http.StatusOK)
}

func (wh *WebHandlers) FilterNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, tmplErr := wh.ExecuteTemplate("filter_form", nil)
		if tmplErr != nil {
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
		}
		wh.WriteHTML(w, tmpl, http.StatusOK)
		return
	}
	filterTypeID, _ := strconv.ParseInt(r.FormValue("filter_type_id"), 10, 64)
	filterValue := r.FormValue("filter_value")
	stmt, err := wh.db.Prepare("INSERT INTO filters (filter_type_id, filter_value) VALUES (?, ?)")
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

func (wh *WebHandlers) FilterEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var f Filter
		err := wh.db.QueryRow("SELECT filter_id, filter_type_id, filter_value FROM filters WHERE filter_id=?", id).
			Scan(&f.ID, &f.FilterTypeID, &f.FilterValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl, tmplErr := wh.ExecuteTemplate("filter_form.html", f)
		if tmplErr != nil {
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
		}
		wh.WriteHTML(w, tmpl, http.StatusOK)
		return
	}
	filterTypeID, _ := strconv.ParseInt(r.FormValue("filter_type_id"), 10, 64)
	filterValue := r.FormValue("filter_value")
	stmt, err := wh.db.Prepare("UPDATE filters SET filter_type_id=?, filter_value=? WHERE filter_id=?")
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

func (wh *WebHandlers) FilterDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := wh.db.Prepare("DELETE FROM filters WHERE filter_id=?")
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
