package stock

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Fetcher interface {
	Get() ([]Price, error)
}

func SourceURL(symbol string, apiKey string) string {
	sourceURL := url.URL{}
	sourceURL.Scheme = "https"
	sourceURL.Host = "www.alphavantage.co"
	sourceURL.Path = "query"
	query := url.Values{}
	query.Set("apikey", apiKey)
	query.Set("function", "TIME_SERIES_DAILY")
	query.Set("symbol", symbol)
	query.Set("outputsize", "compact")
	query.Set("datatype", "csv")
	sourceURL.RawQuery = query.Encode()
	return sourceURL.String()
}

func BuildFetcher(sourceURL string, numberOfDays int) Fetcher {
	return &fetcher{
		sourceURL:    sourceURL,
		numberOfDays: numberOfDays,
	}
}

type fetcher struct {
	sourceURL    string
	numberOfDays int
}

func (f *fetcher) Get() ([]Price, error) {
	resp, err := http.Get(f.sourceURL)
	if err != nil {
		return nil, fmt.Errorf("could not get prices: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get prices due to %s", http.StatusText(resp.StatusCode))
	}
	prices := make([]Price, 0, f.numberOfDays)
	scanner := bufio.NewScanner(resp.Body)
	if !scanner.Scan() {
		return nil, fmt.Errorf("response did not include a header")
	}
	columnNames := scanner.Text()
	timestampColumn, closeColumn := -1, -1
	for i, name := range strings.Split(columnNames, ",") {
		switch name {
		case "timestamp":
			timestampColumn = i
		case "close":
			closeColumn = i
		}
	}
	if timestampColumn < 0 {
		return nil, fmt.Errorf("response did not include a timestamp column")
	}
	if closeColumn < 0 {
		return nil, fmt.Errorf("response did not include a close column")
	}
	for len(prices) < f.numberOfDays && scanner.Scan() {
		line := scanner.Text()
		lineParts := strings.Split(line, ",")
		if len(lineParts) < timestampColumn+1 || len(lineParts) < closeColumn+1 {
			return nil, fmt.Errorf("malformed line in response")
		}
		prices = append(prices, Price{
			Date:  lineParts[timestampColumn],
			Close: lineParts[closeColumn],
		})
	}
	return prices, nil
}
