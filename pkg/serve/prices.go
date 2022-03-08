package serve

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"stock-ticker/pkg/stock"
)

// BuildHandler builds an http server handler.
func BuildHandler(fetcher stock.Fetcher, rateLimiter RateLimiter) http.Handler {
	return &handler{
		fetcher:     fetcher,
		rateLimiter: rateLimiter,
	}
}

type handler struct {
	fetcher     stock.Fetcher
	rateLimiter RateLimiter
}

// ServeHTTP implements http.HandlerFunc.
// The response from the server is a JSON document containing the historical closing prices for the desired stock and
// the average of the closing prices across that history.
//
// Example response:
// {"average":"295.9420","historical":[{"date":"2022-03-04","close":"289.8600"},{"date":"2022-03-03","close":"295.9200"},{"date":"2022-03-02","close":"300.1900"},{"date":"2022-03-01","close":"294.9500"},{"date":"2022-02-28","close":"298.7900"}]}
func (h *handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !h.rateLimiter.EnforceRateLimiting(request) {
		response.WriteHeader(http.StatusTooManyRequests)
		return
	}
	prices, err := h.fetcher.Get()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(fmt.Sprintf("could not fetch prices: %v", err)))
		return
	}
	priceResponse := PricesResponse{
		Average:    averagePrice(prices),
		Historical: prices,
	}
	priceResponseAsJSON, err := json.Marshal(priceResponse)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(fmt.Sprintf("could not marshal price response: %v", err)))
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write(priceResponseAsJSON)
}

func averagePrice(prices []stock.Price) string {
	sum := 0.
	count := 0
	for _, p := range prices {
		price, err := strconv.ParseFloat(p.Close, 64)
		if err != nil {
			continue
		}
		sum += price
		count++
	}
	return fmt.Sprintf("%.4f", sum/float64(count))
}
