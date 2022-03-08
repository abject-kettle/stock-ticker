package stock

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetcherGet(t *testing.T) {
	cases := []struct {
		name               string
		numberOfDays       int
		responseStatusCode int
		responseBody       string
		expectedPrices     []Price
		expectedError      error
	}{{
		name:               "query error",
		numberOfDays:       5,
		responseStatusCode: http.StatusInternalServerError,
		expectedError:      fmt.Errorf("could not get prices due to Internal Server Error"),
	}, {
		name:               "missing response body",
		numberOfDays:       5,
		responseStatusCode: http.StatusOK,
		expectedError:      fmt.Errorf("response did not include a header"),
	}, {
		name:               "missing timestamp column",
		numberOfDays:       5,
		responseStatusCode: http.StatusOK,
		responseBody: `open,high,low,close,volume
294.2900,295.6600,287.1650,289.8600,32369655
302.8900,303.1300,294.0500,295.9200,27314469
295.3600,301.4700,293.6980,300.1900,31873007`,
		expectedError: fmt.Errorf("response did not include a timestamp column"),
	}, {
		name:               "missing close column",
		numberOfDays:       5,
		responseStatusCode: http.StatusOK,
		responseBody: `timestamp,open,high,low,volume
2022-03-04,294.2900,295.6600,287.1650,32369655
2022-03-03,302.8900,303.1300,294.0500,27314469
2022-03-02,295.3600,301.4700,293.6980,31873007`,
		expectedError: fmt.Errorf("response did not include a close column"),
	}, {
		name:               "malformed line",
		numberOfDays:       5,
		responseStatusCode: http.StatusOK,
		responseBody: `timestamp,open,high,low,close,volume
2022-03-04,294.2900,295.6600,287.1650,289.8600,32369655
2022-03-03,302.8900
2022-03-02,295.3600,301.4700,293.6980,300.1900,31873007`,
		expectedError: fmt.Errorf("malformed line in response"),
	}, {
		name:               "more days in response",
		numberOfDays:       5,
		responseStatusCode: http.StatusOK,
		responseBody: `timestamp,open,high,low,close,volume
2022-03-04,294.2900,295.6600,287.1650,289.8600,32369655
2022-03-03,302.8900,303.1300,294.0500,295.9200,27314469
2022-03-02,295.3600,301.4700,293.6980,300.1900,31873007
2022-03-01,296.4000,299.9700,292.1500,294.9500,31217778
2022-02-28,294.3100,299.1400,293.0000,298.7900,34627457
2022-02-25,295.1400,297.6300,291.6550,297.3100,32546721
2022-02-24,272.5100,295.1600,271.5200,294.5900,56989686
2022-02-23,290.1800,291.7000,280.1000,280.2700,37811167
2022-02-22,285.0000,291.5400,284.5000,287.7200,41569319`,
		expectedPrices: []Price{
			{Date: "2022-03-04", Close: "289.8600"},
			{Date: "2022-03-03", Close: "295.9200"},
			{Date: "2022-03-02", Close: "300.1900"},
			{Date: "2022-03-01", Close: "294.9500"},
			{Date: "2022-02-28", Close: "298.7900"},
		},
	}, {
		name:               "fewer days in response",
		numberOfDays:       5,
		responseStatusCode: http.StatusOK,
		responseBody: `timestamp,open,high,low,close,volume
2022-03-04,294.2900,295.6600,287.1650,289.8600,32369655
2022-03-03,302.8900,303.1300,294.0500,295.9200,27314469
2022-03-02,295.3600,301.4700,293.6980,300.1900,31873007`,
		expectedPrices: []Price{
			{Date: "2022-03-04", Close: "289.8600"},
			{Date: "2022-03-03", Close: "295.9200"},
			{Date: "2022-03-02", Close: "300.1900"},
		},
	}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				response.WriteHeader(tc.responseStatusCode)
				response.Write([]byte(tc.responseBody))
			}))
			defer server.Close()
			f := BuildFetcher(server.URL, tc.numberOfDays)
			actualPrices, actualError := f.Get()
			assert.Equal(t, tc.expectedPrices, actualPrices)
			assert.Equal(t, tc.expectedError, actualError)
		})
	}
}
