package handlers

import (
	"net/http"
	"strconv"
)

type FilterType struct {
	ID             int64
	FilterTypeName string
}

// -----------------------------------------------------------------
// Filter Types Handlers

func (wh *WebHandlers) FilterTypesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := wh.db.Query("SELECT filter_type_id, filter_type_name FROM filter_types")
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
	tmpl, tmplErr := wh.ExecuteTemplate("filter_types", fTypes)
	if tmplErr != nil {
		http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
	}
	wh.WriteHTML(w, tmpl, http.StatusOK)
}

func (wh *WebHandlers) FilterTypeNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, tmplErr := wh.ExecuteTemplate("filter_type_form", nil)
		if tmplErr != nil {
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
		}
		wh.WriteHTML(w, tmpl, http.StatusOK)
		return
	}
	filterTypeName := r.FormValue("filter_type_name")
	stmt, err := wh.db.Prepare("INSERT INTO filter_types (filter_type_name) VALUES (?)")
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

func (wh *WebHandlers) FilterTypeEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var ft FilterType
		err := wh.db.QueryRow("SELECT filter_type_id, filter_type_name FROM filter_types WHERE filter_type_id=?", id).
			Scan(&ft.ID, &ft.FilterTypeName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl, tmplErr := wh.ExecuteTemplate("filter_type_form", ft)
		if tmplErr != nil {
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
		}
		wh.WriteHTML(w, tmpl, http.StatusOK)
		return
	}
	filterTypeName := r.FormValue("filter_type_name")
	stmt, err := wh.db.Prepare("UPDATE filter_types SET filter_type_name=? WHERE filter_type_id=?")
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

func (wh *WebHandlers) FilterTypeDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := wh.db.Prepare("DELETE FROM filter_types WHERE filter_type_id=?")
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
