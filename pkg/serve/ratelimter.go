package serve

import (
	"fmt"
	"net/http"

	lru "github.com/hashicorp/golang-lru"
	"k8s.io/client-go/util/flowcontrol"
)

type RateLimiter interface {
	EnforceRateLimiting(request *http.Request) bool
}

func BuildRateLimiter(userQPS int, globalQPS int) (RateLimiter, error) {
	overallRateLimiter := flowcontrol.NewTokenBucketPassiveRateLimiter(float32(globalQPS), 2*globalQPS)
	userRateLimitersCache, err := lru.New(1000)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache for user rate limiters: %w", err)
	}
	return &rateLimiter{
		overallRateLimiter:    overallRateLimiter,
		userRateLimitersCache: userRateLimitersCache,
		userQPS:               userQPS,
	}, nil
}

type rateLimiter struct {
	overallRateLimiter    flowcontrol.PassiveRateLimiter
	userRateLimitersCache *lru.Cache
	userQPS               int
}

func (l *rateLimiter) EnforceRateLimiting(request *http.Request) bool {
	overallAccepted := l.overallRateLimiter.TryAccept()
	userRateLimiter, ok := l.userRateLimitersCache.Get(request.RemoteAddr)
	if !ok {
		userRateLimiter = flowcontrol.NewTokenBucketPassiveRateLimiter(float32(l.userQPS), 2*l.userQPS)
		l.userRateLimitersCache.Add(request.RemoteAddr, userRateLimiter)
	}
	userAccepted := userRateLimiter.(flowcontrol.PassiveRateLimiter).TryAccept()
	return overallAccepted && userAccepted
}
