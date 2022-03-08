package serve

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"stock-ticker/pkg/stock"
)

type mockFetcher struct {
	prices []stock.Price
	err    error
}

func (f mockFetcher) Get() ([]stock.Price, error) {
	return f.prices, f.err
}

type mockRateLimiter struct {
	accepted bool
}

func (l mockRateLimiter) EnforceRateLimiting(request *http.Request) bool {
	return l.accepted
}

func TestServeHTTP(t *testing.T) {
	cases := []struct {
		name                 string
		method               string
		rateLimited          bool
		prices               []stock.Price
		fetchError           error
		expectedStatusCode   int
		expectedResponseBody string
	}{{
		name:               "POST",
		method:             "POST",
		expectedStatusCode: http.StatusMethodNotAllowed,
	}, {
		name:               "Rate limited",
		method:             "GET",
		rateLimited:        true,
		expectedStatusCode: http.StatusTooManyRequests,
	}, {
		name:                 "Fetch error",
		method:               "GET",
		fetchError:           fmt.Errorf("fetch failed"),
		expectedStatusCode:   http.StatusInternalServerError,
		expectedResponseBody: "could not fetch prices: fetch failed",
	}, {
		name:   "Successful fetch",
		method: "GET",
		prices: []stock.Price{
			{Date: "today", Close: "10"},
			{Date: "today-1", Close: "12"},
			{Date: "today-2", Close: "8"},
			{Date: "today-3", Close: "9"},
		},
		expectedStatusCode:   http.StatusOK,
		expectedResponseBody: `{"average":"9.7500","historical":[{"date":"today","close":"10"},{"date":"today-1","close":"12"},{"date":"today-2","close":"8"},{"date":"today-3","close":"9"}]}`,
	}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			handler := BuildHandler(
				mockFetcher{
					prices: tc.prices,
					err:    tc.fetchError,
				},
				mockRateLimiter{
					accepted: !tc.rateLimited,
				},
			)
			request := httptest.NewRequest(tc.method, "/", nil)
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)
			result := recorder.Result()
			defer result.Body.Close()
			assert.Equal(t, tc.expectedStatusCode, result.StatusCode, "unexpected status code")
			actualResponseBody, err := ioutil.ReadAll(result.Body)
			if !assert.NoError(t, err, "unexpected error reading response body") {
				return
			}
			assert.Equal(t, tc.expectedResponseBody, string(actualResponseBody), "unexpected response data")
		})
	}
}
