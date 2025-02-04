package handlers

import "net/http"

// -----------------------------------------------------------------
// Dashboard handler

func (o *WebHandlers) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	err := o.tpl.ExecuteTemplate(w, "dashboard.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
