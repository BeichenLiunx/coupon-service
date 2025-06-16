package coupon

const (
	COUPON_STATUS_USED      = "USED"
	COUPON_STATUS_UNUSED    = "UNUSED"
	COUPON_STATUS_NOT_FOUND = "NOT_FOUND"
)

func GetCouponStatus(code string) string {
	gMutex.Lock()
	defer gMutex.Unlock()

	coupon, exists := gCoupons[code]
	if !exists {
		return COUPON_STATUS_NOT_FOUND
	}

	if coupon.Redeemed {
		return COUPON_STATUS_USED
	}

	return COUPON_STATUS_UNUSED
}
