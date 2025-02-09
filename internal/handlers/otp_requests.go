package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type OTPRequest struct {
	ID             int64
	UserID         int64
	OTPCode        string
	RequestedAt    time.Time
	DeliveryMethod string
	ValidTill      int64
}

// -----------------------------------------------------------------
// OTP Requests Handlers

func (wh *WebHandlers) OtpRequestsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := wh.db.Query("SELECT otp_request_id, user_id, otp_code, requested_at, delivery_method, valid_till FROM otp_requests")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var otps []OTPRequest
	for rows.Next() {
		var otpReq OTPRequest
		err := rows.Scan(&otpReq.ID, &otpReq.UserID, &otpReq.OTPCode, &otpReq.RequestedAt, &otpReq.DeliveryMethod, &otpReq.ValidTill)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		otps = append(otps, otpReq)
	}
	tmpl, tmplErr := wh.ExecuteTemplate("otp_requests", otps)
	if tmplErr != nil {
		slog.Error(fmt.Sprintf("error executing template: %s", tmplErr.Error()))
		http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
	}
	wh.WriteHTML(w, tmpl, http.StatusOK)
}

func (wh *WebHandlers) OtpRequestNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, tmplErr := wh.ExecuteTemplate("otp_request_form", nil)
		if tmplErr != nil {
			slog.Error(fmt.Sprintf("error executing template: %s", tmplErr.Error()))
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
		}
		wh.WriteHTML(w, tmpl, http.StatusOK)
		return
	}
	userID, _ := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	otpCode := r.FormValue("otp_code")
	requestedAt := time.Now()
	deliveryMethod := r.FormValue("delivery_method")
	// Note: valid_till is a generated column, so we do not insert it.
	stmt, err := wh.db.Prepare("INSERT INTO otp_requests (user_id, otp_code, requested_at, delivery_method) VALUES (?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(userID, otpCode, requestedAt, deliveryMethod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/otp_requests", http.StatusSeeOther)
}

func (wh *WebHandlers) OtpRequestEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var otpReq OTPRequest
		err := wh.db.QueryRow("SELECT otp_request_id, user_id, otp_code, requested_at, delivery_method, valid_till FROM otp_requests WHERE otp_request_id=?", id).
			Scan(&otpReq.ID, &otpReq.UserID, &otpReq.OTPCode, &otpReq.RequestedAt, &otpReq.DeliveryMethod, &otpReq.ValidTill)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl, tmplErr := wh.ExecuteTemplate("otp_request_form", wh)
		if tmplErr != nil {
			slog.Error(fmt.Sprintf("error executing template: %s", tmplErr.Error()))
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
		}
		wh.WriteHTML(w, tmpl, http.StatusOK)
		return
	}
	userID, _ := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	otpCode := r.FormValue("otp_code")
	requestedAt, _ := time.Parse("2006-01-02 15:04:05", r.FormValue("requested_at"))
	deliveryMethod := r.FormValue("delivery_method")
	stmt, err := wh.db.Prepare("UPDATE otp_requests SET user_id=?, otp_code=?, requested_at=?, delivery_method=? WHERE otp_request_id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(userID, otpCode, requestedAt, deliveryMethod, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/otp_requests", http.StatusSeeOther)
}

func (wh *WebHandlers) OtpRequestDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := wh.db.Prepare("DELETE FROM otp_requests WHERE otp_request_id=?")
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
	http.Redirect(w, r, "/otp_requests", http.StatusSeeOther)
}
