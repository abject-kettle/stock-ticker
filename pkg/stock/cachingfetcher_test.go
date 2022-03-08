package stock

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockFetcher struct {
	prices []Price
	err    error
	called bool
}

func (f *mockFetcher) Get() ([]Price, error) {
	f.called = true
	return f.prices, f.err
}

func TestCachingFetcherGet(t *testing.T) {
	cases := []struct {
		name                    string
		cachedPrices            []Price
		expired                 bool
		fetcherPrices           []Price
		fetcherError            error
		expectedPrices          []Price
		expectedError           error
		expectExpirationRefresh bool
	}{{
		name:                    "nothing cached",
		fetcherPrices:           []Price{{Date: "today", Close: "12.3400"}},
		expectedPrices:          []Price{{Date: "today", Close: "12.3400"}},
		expectExpirationRefresh: true,
	}, {
		name:                    "cache expired",
		cachedPrices:            []Price{{Date: "old", Close: "10.0000"}},
		expired:                 true,
		fetcherPrices:           []Price{{Date: "today", Close: "12.3400"}},
		expectedPrices:          []Price{{Date: "today", Close: "12.3400"}},
		expectExpirationRefresh: true,
	}, {
		name:                    "cache not expired",
		cachedPrices:            []Price{{Date: "old", Close: "10.0000"}},
		expired:                 false,
		fetcherPrices:           []Price{{Date: "today", Close: "12.3400"}},
		expectedPrices:          []Price{{Date: "old", Close: "10.0000"}},
		expectExpirationRefresh: false,
	}, {
		name:          "fetch error",
		cachedPrices:  []Price{{Date: "old", Close: "10.0000"}},
		expired:       true,
		fetcherError:  fmt.Errorf("could not fetch prices"),
		expectedError: fmt.Errorf("could not fetch prices"),
	}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			baseFetcher := &mockFetcher{
				prices: tc.fetcherPrices,
				err:    tc.fetcherError,
			}
			var startingExpiration time.Time
			if tc.expired {
				startingExpiration = time.Now().Add(-1 * time.Hour)
			} else {
				startingExpiration = time.Now().Add(1 * time.Hour)
			}
			const expiry = 10 * time.Minute
			f := cachingFetcher{
				fetcher:    baseFetcher,
				expiry:     expiry,
				prices:     tc.cachedPrices,
				expiration: startingExpiration,
			}
			timeBeforeGet := time.Now()
			actualPrices, actualError := f.Get()
			timeAfterGet := time.Now()
			assert.Equal(t, tc.expectedPrices, actualPrices)
			assert.Equal(t, tc.expectedError, actualError)
			if tc.expectExpirationRefresh {
				assert.LessOrEqual(t, timeBeforeGet.Add(expiry).Unix(), f.expiration.Unix())
				assert.GreaterOrEqual(t, timeAfterGet.Add(expiry).Unix(), f.expiration.Unix())
			} else {
				assert.Equal(t, startingExpiration, f.expiration)
			}
		})
	}
}
