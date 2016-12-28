package server

import (
	"fmt"
)

// Takes a string identifying a currency and returns a container with
// the known rates relative to the given base.
func (s *Server) createResponse(base string) (r *currencyResponse, err error) {
	// return if we don't have any currencies (job might still be fetching)
	if !s.hasCurrencies {
		return nil, fmt.Errorf("Currencies not fetched")
	}

	// return if the base given is unknown
	baserate, found := s.currencies[base]
	if !found {
		return nil, fmt.Errorf("Unknown currency: %s", base)
	}

	// create the response container
	response := currencyResponse{}
	response.BaseCurrency = base
	response.CurrencyDate = s.lastUpdateTime.Format("2006-01-02")

	// fill the converted rates
	for name, rate := range s.currencies {
		relativeRate := rate / baserate
		r := rateResponse{
			Name: name,
			Rate: relativeRate,
		}

		response.Rates = append(response.Rates, r)
	}

	return &response, nil
}

// Takes a base currency and a target currency and converts a slice of amounts
// from one to another.
func (s *Server) createConvertResponse(to, from string, amounts []float64) (r *convertResponse, err error) {
	// create the response package
	response := convertResponse{}
	response.BaseCurrency = from
	response.TargetCurrency = to
	response.CurrencyDate = s.lastUpdateTime.Format("2006-01-02")

	// convert the amounts, one at a time
	for _, amount := range amounts {
		converted, err := s.convert(to, from, amount)
		// return the whole function if the single conversion fails!
		if err != nil {
			return nil, err
		}

		// append to converted amounts if conversion succeeded
		response.ConvertedAmounts = append(response.ConvertedAmounts, converted)
	}

	return &response, nil
}

// Converts a single amount from one currency to another
func (s *Server) convert(to, from string, amount float64) (result float64, err error) {
	// error if base currency is not known
	baserate, found := s.currencies[from]
	if !found {
		return 0.0, fmt.Errorf("Unknown currency: %s", from)
	}

	// error if target currency is not known
	targetrate, found := s.currencies[to]
	if !found {
		return 0.0, fmt.Errorf("Unknown currency: %s", to)
	}

	// convert!
	result = amount / baserate * targetrate
	return result, nil
}
