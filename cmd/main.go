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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	symbol := os.Getenv("SYMBOL")
	if symbol == "" {
		log.Fatal("SYMBOL environment variable is required")
	}
	numberOfDaysAsString := os.Getenv("NDAYS")
	if numberOfDaysAsString == "" {
		log.Fatal("NDAYS environment variable is required")
	}
	numberOfDays, err := strconv.Atoi(numberOfDaysAsString)
	if err != nil || numberOfDays <= 0 {
		log.Fatal("NDAYS environment variable must be a positive integer")
	}
	apiKey := os.Getenv("APIKEY")
	if apiKey == "" {
		log.Fatal("APIKEY environment variable is required")
	}
	rateLimiter, err := serve.BuildRateLimiter(10, 10000)
	if err != nil {
		log.Fatal(err)
	}
	priceFetcher := stock.BuildFetcher(symbol, numberOfDays, apiKey)
	cachingPriceFetcher := stock.AddCaching(priceFetcher, 10*time.Minute)
	priceHandler := serve.BuildHandler(cachingPriceFetcher, rateLimiter)
	http.Handle("/", priceHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
