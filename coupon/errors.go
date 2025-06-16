package coupon

import "errors"

var (
	ErrGenerateMaxRetryExceeded = errors.New("generate coupon max retry exceeded")
	ErrCouponNotFound           = errors.New("coupon not found")
	ErrCouponAlreadyRedeemed    = errors.New("coupon already redeemed")
	ErrRateLimited              = errors.New("rate limit exceeded, please try again later")
)
