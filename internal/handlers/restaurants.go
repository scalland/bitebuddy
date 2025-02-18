package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

type Restaurant struct {
	ID                int64
	Name              string
	Address           string
	Latitude          float64
	Longitude         float64
	OverallRating     float64
	PriceForTwo       float64
	ImageURL          string
	DiscountAvailable bool
	AlcoholAvailable  bool
	PortionSizeLarge  bool
}

// -----------------------------------------------------------------
// Restaurants Handlers

func (wh *WebHandlers) RestaurantsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := wh.db.Query("SELECT restaurant_id, name, address, latitude, longitude, overall_rating, price_for_two, image_url, discount_available, alcohol_available, portion_size_large FROM restaurants")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var restaurants []Restaurant
	for rows.Next() {
		var rct Restaurant
		err := rows.Scan(&rct.ID, &rct.Name, &rct.Address, &rct.Latitude, &rct.Longitude, &rct.OverallRating, &rct.PriceForTwo, &rct.ImageURL, &rct.DiscountAvailable, &rct.AlcoholAvailable, &rct.PortionSizeLarge)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		restaurants = append(restaurants, rct)
	}
	tmpl, tmplErr := wh.ExecuteTemplate("restaurants", restaurants)
	if tmplErr != nil {
		slog.Error(fmt.Sprintf("Error executing template: restaurants.html: %s", tmplErr))
		http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
	}
	wh.WriteHTML(w, tmpl, http.StatusOK)
}

func (wh *WebHandlers) RestaurantNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, tmplErr := wh.ExecuteTemplate("restaurant_form", nil)
		if tmplErr != nil {
			slog.Error(fmt.Sprintf("Error executing template: restaurant_form: %s", tmplErr))
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
		}
		wh.WriteHTML(w, tmpl, http.StatusOK)
		return
	}
	// POST
	name := r.FormValue("name")
	address := r.FormValue("address")
	lat, _ := strconv.ParseFloat(r.FormValue("latitude"), 64)
	lng, _ := strconv.ParseFloat(r.FormValue("longitude"), 64)
	overallRating, _ := strconv.ParseFloat(r.FormValue("overall_rating"), 64)
	priceForTwo, _ := strconv.ParseFloat(r.FormValue("price_for_two"), 64)
	imageURL := r.FormValue("image_url")
	discountAvailable := r.FormValue("discount_available") == "on"
	alcoholAvailable := r.FormValue("alcohol_available") == "on"
	portionSizeLarge := r.FormValue("portion_size_large") == "on"
	stmt, err := wh.db.Prepare("INSERT INTO restaurants (name, address, latitude, longitude, overall_rating, price_for_two, image_url, discount_available, alcohol_available, portion_size_large) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(name, address, lat, lng, overallRating, priceForTwo, imageURL, discountAvailable, alcoholAvailable, portionSizeLarge)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/restaurants", http.StatusSeeOther)
}

func (wh *WebHandlers) RestaurantEditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method == http.MethodGet {
		var rct Restaurant
		err := wh.db.QueryRow("SELECT restaurant_id, name, address, latitude, longitude, overall_rating, price_for_two, image_url, discount_available, alcohol_available, portion_size_large FROM restaurants WHERE restaurant_id=?", id).
			Scan(&rct.ID, &rct.Name, &rct.Address, &rct.Latitude, &rct.Longitude, &rct.OverallRating, &rct.PriceForTwo, &rct.ImageURL, &rct.DiscountAvailable, &rct.AlcoholAvailable, &rct.PortionSizeLarge)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl, tmplErr := wh.ExecuteTemplate("restaurant_form", rct)
		if tmplErr != nil {
			slog.Error(fmt.Sprintf("Error executing template: restaurant_form.%s", tmplErr))
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
		}
		wh.WriteHTML(w, tmpl, http.StatusOK)
		return
	}
	// POST update
	name := r.FormValue("name")
	address := r.FormValue("address")
	lat, _ := strconv.ParseFloat(r.FormValue("latitude"), 64)
	lng, _ := strconv.ParseFloat(r.FormValue("longitude"), 64)
	overallRating, _ := strconv.ParseFloat(r.FormValue("overall_rating"), 64)
	priceForTwo, _ := strconv.ParseFloat(r.FormValue("price_for_two"), 64)
	imageURL := r.FormValue("image_url")
	discountAvailable := r.FormValue("discount_available") == "on"
	alcoholAvailable := r.FormValue("alcohol_available") == "on"
	portionSizeLarge := r.FormValue("portion_size_large") == "on"
	stmt, err := wh.db.Prepare("UPDATE restaurants SET name=?, address=?, latitude=?, longitude=?, overall_rating=?, price_for_two=?, image_url=?, discount_available=?, alcohol_available=?, portion_size_large=? WHERE restaurant_id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(name, address, lat, lng, overallRating, priceForTwo, imageURL, discountAvailable, alcoholAvailable, portionSizeLarge, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/restaurants", http.StatusSeeOther)
}

func (wh *WebHandlers) RestaurantDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	stmt, err := wh.db.Prepare("DELETE FROM restaurants WHERE restaurant_id=?")
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
	http.Redirect(w, r, "/restaurants", http.StatusSeeOther)
}
