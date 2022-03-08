package stock

// Price is the historical price of a stock on a given day.
type Price struct {
	// Date of the price.
	Date string `json:"date"`
	// Closing price of the stock on the day.
	Close string `json:"close"`
}
