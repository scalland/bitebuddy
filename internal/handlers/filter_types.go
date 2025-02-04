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

func (o *WebHandlers) FilterTypesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := o.db.Query("SELECT filter_type_id, filter_type_name FROM filter_types")
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
	o.tpl.ExecuteTemplate(w, "filter_types.html", fTypes)
}

func (o *WebHandlers) FilterTypeNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		o.tpl.ExecuteTemplate(w, "filter_type_form.html", nil)
		return
	}
	filterTypeName := r.FormValue("filter_type_name")
	stmt, err := o.db.Prepare("INSERT INTO filter_types (filter_type_name) VALUES (?)")
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

func (o *WebHandlers) FilterTypeEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var ft FilterType
		err := o.db.QueryRow("SELECT filter_type_id, filter_type_name FROM filter_types WHERE filter_type_id=?", id).
			Scan(&ft.ID, &ft.FilterTypeName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		o.tpl.ExecuteTemplate(w, "filter_type_form.html", ft)
		return
	}
	filterTypeName := r.FormValue("filter_type_name")
	stmt, err := o.db.Prepare("UPDATE filter_types SET filter_type_name=? WHERE filter_type_id=?")
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

func (o *WebHandlers) FilterTypeDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := o.db.Prepare("DELETE FROM filter_types WHERE filter_type_id=?")
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
