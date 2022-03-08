package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"stock-ticker/pkg/serve"
	"stock-ticker/pkg/stock"
)

func main() {
	// get the port on which to listen
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// get the ticker symbol for which to fetch prices
	symbol := os.Getenv("SYMBOL")
	if symbol == "" {
		log.Fatal("SYMBOL environment variable is required")
	}

	// get the number of days of prices to fetch
	numberOfDaysAsString := os.Getenv("NDAYS")
	if numberOfDaysAsString == "" {
		log.Fatal("NDAYS environment variable is required")
	}
	numberOfDays, err := strconv.Atoi(numberOfDaysAsString)
	if err != nil || numberOfDays <= 0 {
		log.Fatal("NDAYS environment variable must be a positive integer")
	}

	// get the API key to use when fetching prices
	apiKey := os.Getenv("APIKEY")
	if apiKey == "" {
		log.Fatal("APIKEY environment variable is required")
	}

	// build the rate limiter
	rateLimiter, err := serve.BuildRateLimiter(10, 10000)
	if err != nil {
		log.Fatal(err)
	}

	// build the price fetcher
	priceFetcher := stock.BuildFetcher(stock.SourceURL(symbol, apiKey), numberOfDays)
	cachingPriceFetcher := stock.AddCaching(priceFetcher, 10*time.Minute)

	// build the http server handler
	priceHandler := serve.BuildHandler(cachingPriceFetcher, rateLimiter)
	http.Handle("/", priceHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
