# Coupon Code System Design

## 1. Coupon Code Format and Generation Strategy

**Format**: `[A-Z0-9]{8}`  
- We exclude ambiguous characters (`I`, `O`, `0`, `1`) to reduce human entry errors.  
- Characters used: `A-H, J-N, P-Z, 2-9`

**Generation Strategy**:
- Pre-generate codes using a cryptographically secure random function.
- Maintain uniqueness using an in-memory map (`map[string]*Coupon`).
- Pre-fill a pool (`chan string`) with unused codes for fast access.
- Fallback: If pool is empty, generate on-demand with retry logic.
- Collision probability is extremely low (36^8 ≈ 2.8 trillion combinations).

## 2. API Overview

### Endpoints

| Method | Path               | Description                       |
|--------|--------------------|-----------------------------------|
| POST   | /generate-coupon   | Generate a new 8-char coupon      |
| POST   | /redeem-coupon     | Redeem a coupon (idempotent)      |
| GET    | /coupon-status     | Get status of a coupon by code    |

### Request/Response Examples

- **POST /generate-coupon**
  - Response: `{ "code": "AB2C4XYZ" }`

- **POST /redeem-coupon**
  - Request: `{ "code": "AB2C4XYZ" }`
  - Response: `{ "message": "Coupon redeemed successfully" }`

- **GET /coupon-status?code=AB2C4XYZ**
  - Response: `{ "code": "AB2C4XYZ", "status": "UNUSED" or "USED" or "NOT_FOUND" }`

## 3. High-Level Architecture

- **Client** → REST API Server
- **CouponManager** (in-memory store + generator)
- **RateLimiter** (per-IP protection)
- No external dependencies (standard Go libraries only)

### Scalability Considerations
- Pre-generated pool for fast access.
- Mutex locking to prevent concurrent access issues.
- Could plug into Redis or DB for distributed use in future.

## 4. Rate-Limiting and Abuse Prevention

- Per-IP rate limiting using `map[string]*RateLimitInfo`.
- Fixed 1-minute window, max 5 redemption attempts per IP.
- Background job purges stale IP entries every 2 minutes.
- Returns HTTP 429 on abuse.

## 5. Concurrency and Retry Handling

- Uses `sync.RWMutex` to guard concurrent access to coupon map.
- Redeem operation:
  - Checked under lock.
  - Ensures idempotency — re-redeeming returns success if already redeemed.
- Coupon assignment is atomic inside `gMutex.Lock()`.

## 6. Observability

- Console logs for:
  - Coupon generation path (pool vs. on-demand)
  - Redemption outcome
  - Rate-limiting triggers
- Future enhancements:
  - Add Prometheus metrics for `/generate`, `/redeem`, and rate-limited hits.

## 7. Additional Considerations

- Coupon struct supports `redeemed_at` for auditing.
- Could extend to support:
  - Expiration dates
  - Campaign metadata
  - User-specific codes
- Code is testable and designed for safe extension (e.g. swap in Redis later).

## Summary

- Simple, in-memory, secure coupon system with REST API
- Fast, safe, scalable to thousands of requests per second
- Fully written in Go with zero third-party dependencies
