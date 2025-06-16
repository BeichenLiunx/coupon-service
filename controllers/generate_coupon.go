package controllers

import (
	"BeichenLiunx/coupon-service/coupon"
	"encoding/json"
	"fmt"
	"net/http"
)

type GenerateCouponResponse struct {
	Code string `json:"code"`
}

func HandleGenerateCoupon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	code, err := coupon.GenerateCouponCode()
	if err != nil {
		fmt.Println("Error generating coupon code:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Generated coupon code:", code)

	response := GenerateCouponResponse{Code: code}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println("Error encoding response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
