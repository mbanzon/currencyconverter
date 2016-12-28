package server

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	errorSleepTime   = 1 * time.Minute // The sleep time after an error
	successSleepTime = 1 * time.Hour   // The standard sleep time

	ecbCurrencyUrl     = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"
	currencyDateFormat = "2006-01-02" // The time format in ECB XML
	eur                = "EUR"        // The euro symbol
)

// Starts a goroutine what fetches the currencies from the ECB every hour
func (s *Server) startCurrencyUpdating() {
	log.Println("Starting currency fetching...")
	go func() {
		for {
			log.Println("Starting new currency fetch...")

			// initialize the default nap time
			napTime := successSleepTime

			if data, err := fetchCurrencyData(); err == nil {
				if ts, curr, err2 := parseCurrencyData(data); err2 == nil {
					// everything succeeded - update the currency data
					// lock while doing so
					s.mutex.Lock()
					s.hasCurrencies, s.lastUpdateTime, s.currencies = true, ts, curr
					s.mutex.Unlock()

					log.Println("Currencies updated.")

					// call the webhooks
					go s.callWebhooks()
				} else {
					// error occured - log and set smaller nap time
					log.Println("Error parsing currency data:", err)
					napTime = errorSleepTime
				}
			} else {
				// error occured - log and set smaller nap time
				log.Println("Error fetching currency data:", err)
				napTime = errorSleepTime
			}

			// nap
			log.Println("Sleeping", napTime)
			time.Sleep(napTime)
		}
	}()
}

// The currency XML data
type currencyEnvelope struct {
	Sender string `xml:"Sender>name"`
	Cube   []cube `xml:"Cube>Cube>Cube"`
}

// The time XML data
type timeEnvelope struct {
	Time timeCube `xml:"Cube>Cube"`
}

// The time holder XML data
type timeCube struct {
	Time string `xml:"time,attr"`
}

// The cube XML structure
type cube struct {
	Name string  `xml:"currency,attr"`
	Rate float64 `xml:"rate,attr"`
}

// Fetches the raw data from the ECB URL.
func fetchCurrencyData() (data []byte, err error) {
	res, err := http.Get(ecbCurrencyUrl)
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadAll(res.Body)
	res.Body.Close() // ignore error?
	return data, err
}

// Parse the raw data from the ECB. Returns the time from the XML along with a
// map of currency rates.
func parseCurrencyData(data []byte) (ts time.Time, currencies map[string]float64, err error) {
	// parse once to get the currencies, return on error
	var c currencyEnvelope
	err = xml.Unmarshal(data, &c)
	if err != nil {
		return time.Time{}, nil, err
	}

	// parse again to get the timestamp, return on error
	var t timeEnvelope
	err = xml.Unmarshal(data, &t)
	if err != nil {
		return time.Time{}, nil, err
	}

	// parse time, return on error
	ts, err = time.Parse(currencyDateFormat, t.Time.Time)
	if err != nil {
		return time.Time{}, nil, err
	}

	currencies = make(map[string]float64)

	// manually insert EUR as "1"
	currencies[eur] = 1

	// insert all rates
	for _, currency := range c.Cube {
		currencies[currency.Name] = currency.Rate
	}

	return ts, currencies, nil
}
