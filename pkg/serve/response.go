package serve

import "stock-ticker/pkg/stock"

type PricesResponse struct {
	Average    string        `json:"average"`
	Historical []stock.Price `json:"historical"`
}
