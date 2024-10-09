package customerrors

import "errors"

var (
	SkuNotFound        = errors.New("sku not found")
	InvalidSkuId       = errors.New("invalid sku")
	InvalidUserId      = errors.New("invalid user id")
	InvalidCount       = errors.New("invalid count")
	NotEnoughStock     = errors.New("not enough stock")
	RateLimitSetupFail = errors.New("rateLimit and burstLimit must be greater than zero")
)
