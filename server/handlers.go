package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// struct for the currency rates
type currencyResponse struct {
	CurrencyDate string         `json:"currency_date"`
	BaseCurrency string         `json:"base_currency"`
	Rates        []rateResponse `json:"rates"`
}

// struct for the single rates
type rateResponse struct {
	Name string  `json:"name"`
	Rate float64 `json:"rate"`
}

// struct for the currency request with a different base
type currencyRequest struct {
	BaseCurrency string `json:"base_currency"`
}

// struct for the currency convertion request
type convertRequest struct {
	BaseCurrency   string    `json:"base_currency"`
	TargetCurrency string    `json:"target_currency"`
	Amounts        []float64 `json:"amounts"`
}

// struct for the currency convertion response
type convertResponse struct {
	BaseCurrency     string    `json:"base_currency"`
	TargetCurrency   string    `json:"target_currency"`
	CurrencyDate     string    `json:"currency_date"`
	ConvertedAmounts []float64 `json:"converted_amounts"`
}

// The main serving function. This handles all requests to he server by
// delegating the requests to the other handlers.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s\n", r.Method, r.RequestURI)

	// error if there is no currencies
	if !s.hasCurrencies {
		log.Println("No currencies, returning error!")
		http.Error(w, "No currencies", http.StatusServiceUnavailable)
		return
	}

	// select the correct handler, error on unknown URI
	switch r.RequestURI {
	case "/currencies":
		s.currenciesHandler(w, r)
	case "/convert":
		s.convertHandler(w, r)
	case "/webhook":
		s.webhookHandler(w, r)
	default:
		if strings.HasPrefix(r.RequestURI, "/script?base=") {
			s.scriptHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	}
}

// Handles currency requests (/currencies)
func (s *Server) currenciesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// GET - create response with EUR base
		res, err := s.createResponse(eur)
		s.respondJson(w, res, err)
		s.currencyHits.Add(1)
	} else if r.Method == http.MethodPost {
		// POST - parse request to get base, fail on error
		var req currencyRequest
		err := s.getJsonRequest(r, &req)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// create reponse with parsed base
		res, err := s.createResponse(req.BaseCurrency)
		s.respondJson(w, res, err)
		s.currencyHits.Add(1)
	} else {
		http.Error(w, "", http.StatusBadRequest)
	}
}

// Handles the convertion requests (/convert)
func (s *Server) convertHandler(w http.ResponseWriter, r *http.Request) {
	// only handle POST, error on everything else
	if r.Method == http.MethodPost {
		// parse the convertion request
		var req convertRequest
		err := s.getJsonRequest(r, &req)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// create the convertion response
		res, err := s.createConvertResponse(req.TargetCurrency, req.BaseCurrency, req.Amounts)
		s.respondJson(w, res, err)
		s.convertHits.Add(1)
	} else {
		http.Error(w, "", http.StatusBadRequest)
	}
}

// Handles the webhook call to add webhooks
func (s *Server) webhookHandler(w http.ResponseWriter, r *http.Request) {
	// we only handle POST
	if r.Method == http.MethodPost {
		// parse the webhook
		var hook webhook
		err := s.getJsonRequest(r, &hook)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// verify the webhook data and insert
		if s.verifyWebhook(hook) {
			s.mutex.Lock()
			defer s.mutex.Unlock()

			s.webhooks[hook.Url] = hook
			s.webhookHits.Add(1)
		} else {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	} else {
		http.NotFound(w, r)
	}
}

// generic method to return JSON of v to a http.ResponseWriter, return proper
// status code if the passed error is not nil
func (s *Server) respondJson(w http.ResponseWriter, v interface{}, err error) {
	// add header for content type
	w.Header().Add("Content-Type", "application/json")

	// return internal server error if err is not nil
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// encode as JSON
	err = json.NewEncoder(w).Encode(v)

	// return error if encoding caused one
	if err != nil {
		log.Println(err)
		http.Error(w, "Error creating JSON response.", http.StatusInternalServerError)
	}
}

// parses the given http.Request into the given interface
func (s *Server) getJsonRequest(r *http.Request, v interface{}) (err error) {
	err = json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return err
	}

	return nil
}
