package handlers

import (
	"net/http"
)

// -----------------------------------------------------------------
// Dashboard handler

type DashboardHandlerData struct {
	IsLoggedIn      bool
	IsLoggedInAdmin bool
	Errors          []string
}

func (wh *WebHandlers) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	wh.Log.Debugf("inside dashboard handler")
	data, err := wh.ExecuteTemplate("dashboard", DashboardHandlerData{IsLoggedIn: wh.IsLoggedIn(r), IsLoggedInAdmin: wh.IsLoggedInAdmin(r), Errors: nil})
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	wh.WriteHTML(w, data, http.StatusOK)
}
