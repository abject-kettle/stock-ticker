package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

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
	priceFetcher := stock.BuildFetcher(symbol, numberOfDays, apiKey)
	priceHandler := serve.BuildHandler(priceFetcher)
	http.Handle("/", priceHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
