package main

import (
	"BeichenLiunx/coupon-service/controllers"
	"BeichenLiunx/coupon-service/coupon"
	"context"
	"log"
	"net/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	coupon.Init(ctx)

	http.HandleFunc("/generate-coupon", controllers.HandleGenerateCoupon)
	http.HandleFunc("/redeem-coupon", controllers.HandleRedeemCoupon)
	http.HandleFunc("/coupon-status", controllers.HandleGetCouponStatus)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
