package coupon

import (
	"BeichenLiunx/coupon-service/ratelimiter"
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Coupon struct {
	Code       string
	Redeemed   bool
	RedeemedAt *time.Time
}

const (
	MAX_RETRYS = 50
)

var (
	gCoupons  = make(map[string]*Coupon)
	gCodePool = make(chan string, 1000) // Pre-allocate space for 1000 codes
	gMutex    = sync.RWMutex{}
	gCharset  = []rune("ABCDEFGHJKLMNPQRSTUVWXYZ23456789")
)

func GenerateCouponCode() (string, error) {
	var code string

	// Try to get a code from the pool first
	select {
	case code = <-gCodePool:
		fmt.Println("Using pre-generated coupon code:", code)

	default:
		retry := 0
		for {
			if retry >= MAX_RETRYS {
				return "", ErrGenerateMaxRetryExceeded
			}

			code = generateCouponCode()
			gMutex.RLock()
			_, exists := gCoupons[code]
			gMutex.RUnlock()
			if !exists {
				fmt.Println("Generated new coupon code:", code)
				break
			}
			retry++
		}
	}

	fmt.Println("Final coupon code to be stored:", code)

	gMutex.Lock()
	gCoupons[code] = &Coupon{Code: code}
	gMutex.Unlock()

	fmt.Println("Final coupon code to be store222d:", code)

	return code, nil
}

func generateCouponCode() string {
	b := make([]rune, 8)
	for i := range b {
		b[i] = gCharset[rand.Intn(len(gCharset))]
	}
	return string(b)
}

func Init(ctx context.Context) error {
	gRateLimiter = ratelimiter.NewRateLimiter(5, ctx)
	go backgroundCouponGenerator(ctx)
	return nil
}

func backgroundCouponGenerator(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			code := generateCouponCode()
			gMutex.RLock()
			_, exists := gCoupons[code]
			gMutex.RUnlock()
			if !exists {
				select {
				case gCodePool <- code:

				default:
				}
			}
		}
	}
}
