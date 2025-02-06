package handlers

import (
	"net/http"
)

// -----------------------------------------------------------------
// Dashboard handler

func (wh *WebHandlers) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	data, err := wh.ExecuteTemplate("dashboard", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	wh.WriteHTML(w, data)
}
