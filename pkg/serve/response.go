package serve

import "stock-ticker/pkg/stock"

// PricesResponse is the response sent to the user from the http server.
type PricesResponse struct {
	// Average is the average price of the stock over the historical data.
	Average string `json:"average"`
	// Historical is the historical closing price of the stock.
	Historical []stock.Price `json:"historical"`
}
