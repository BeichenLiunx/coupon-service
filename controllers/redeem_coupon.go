package controllers

import (
	"BeichenLiunx/coupon-service/coupon"
	"encoding/json"
	"net/http"
	"strings"
)

func HandleRedeemCoupon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	code := strings.ToUpper(req.Code)

	err := coupon.RedeemCouponCode(code, r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Coupon redeemed successfully",
	})
}
