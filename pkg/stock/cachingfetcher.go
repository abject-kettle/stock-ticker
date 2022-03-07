package stock

import "time"

func AddCaching(fetcher Fetcher, expiry time.Duration) Fetcher {
	return &cachingFetcher{
		fetcher: fetcher,
		expiry:  expiry,
	}
}

type cachingFetcher struct {
	fetcher    Fetcher
	expiry     time.Duration
	prices     []Price
	expiration time.Time
}

func (f *cachingFetcher) Get() ([]Price, error) {
	if len(f.prices) > 0 || time.Until(f.expiration) > 0 {
		return f.prices, nil
	}
	prices, err := f.fetcher.Get()
	if err != nil {
		return nil, err
	}
	f.prices = prices
	f.expiration = time.Now().Add(f.expiry)
	return prices, nil
}
