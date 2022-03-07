package serve

import (
	"fmt"
	"net/http"

	lru "github.com/hashicorp/golang-lru"
	"k8s.io/client-go/util/flowcontrol"
)

var (
	overallRateLimiter    flowcontrol.PassiveRateLimiter
	userRateLimitersCache *lru.Cache
)

func InitializeRateLimiters() error {
	overallRateLimiter = flowcontrol.NewTokenBucketPassiveRateLimiter(10000, 20000)
	var err error
	userRateLimitersCache, err = lru.New(1000)
	if err != nil {
		return fmt.Errorf("failed to create cache for user rate limiters: %w", err)
	}
	return nil
}

func enforceRateLimiting(request *http.Request) bool {
	overallAccepted := overallRateLimiter.TryAccept()
	userRateLimiter, ok := userRateLimitersCache.Get(request.RemoteAddr)
	if !ok {
		userRateLimiter = flowcontrol.NewTokenBucketPassiveRateLimiter(10, 20)
		userRateLimitersCache.Add(request.RemoteAddr, userRateLimiter)
	}
	userAccepted := userRateLimiter.(flowcontrol.PassiveRateLimiter).TryAccept()
	return overallAccepted && userAccepted
}
