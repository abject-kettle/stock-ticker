package serve

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"stock-ticker/pkg/stock"
)

func BuildHandler(fetcher stock.Fetcher) http.Handler {
	return &handler{
		fetcher: fetcher,
	}
}

type handler struct {
	fetcher stock.Fetcher
}

func (h *handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !enforceRateLimiting(request) {
		response.WriteHeader(http.StatusTooManyRequests)
		return
	}
	prices, err := h.fetcher.Get()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(fmt.Sprintf("could not fetch prices: %v", err)))
	}
	priceResponse := PricesResponse{
		Average:    averagePrice(prices),
		Historical: prices,
	}
	priceResponseAsJSON, err := json.Marshal(priceResponse)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(fmt.Sprintf("could not marshal price response: %v", err)))
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
