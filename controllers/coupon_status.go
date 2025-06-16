package controllers

import (
	"BeichenLiunx/coupon-service/coupon"
	"encoding/json"
	"net/http"
)

type CouponStatusResponse struct {
	Code   string `json:"code"`
	Status string `json:"status"`
}

func HandleGetCouponStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	code := r.URL.Query().Get("coupon_code")
	if code == "" {
		http.Error(w, "Missing coupon code", http.StatusBadRequest)
		return
	}

	status := coupon.GetCouponStatus(code)

	w.Header().Set("Content-Type", "application/json")
	response := CouponStatusResponse{
		Code:   code,
		Status: status,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
