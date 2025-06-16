package coupon

import (
	"BeichenLiunx/coupon-service/ratelimiter"
	"time"
)

var (
	gRateLimiter *ratelimiter.RateLimiter
)

func RedeemCouponCode(code, ip string) error {
	if limited := gRateLimiter.IsRateLimited(ip); limited {
		return ErrRateLimited
	}

	gMutex.Lock()
	defer gMutex.Unlock()

	coupon, exists := gCoupons[code]
	if !exists {
		return ErrCouponNotFound
	}

	if coupon.Redeemed {
		return ErrCouponAlreadyRedeemed
	}

	now := time.Now()
	coupon.Redeemed = true
	coupon.RedeemedAt = &now

	return nil
}
