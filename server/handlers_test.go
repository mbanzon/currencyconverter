package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCurrencyGet(t *testing.T) {
	var tmp currencyResponse
	r := fireReq("/currencies", http.MethodGet, nil)
	expect(t, r, http.StatusOK, true, &tmp)
}

func TestCurrencyPostKnownCurrency(t *testing.T) {
	var tmp currencyResponse
	r := fireReq("/currencies", http.MethodPost, &currencyRequest{BaseCurrency: "USD"})
	expect(t, r, http.StatusOK, true, &tmp)
}

func TestCurrencyPostInvalidRequest(t *testing.T) {
	r := fireReq("/currencies", http.MethodPost, nil)
	expect(t, r, http.StatusInternalServerError, true, nil)
}

func TestCurrencyPostUnknownCurrency(t *testing.T) {
	r := fireReq("/currencies", http.MethodPost, &currencyRequest{BaseCurrency: "FOO"})
	expect(t, r, http.StatusInternalServerError, false, nil)
}

func TestCurrencyPut(t *testing.T) {
	r := fireReq("/currencies", http.MethodPut, nil)
	expect(t, r, http.StatusBadRequest, true, nil)
}

func TestConvert(t *testing.T) {
	r := fireReq("/convert", http.MethodPost, convertRequest{
		BaseCurrency:   "DKK",
		TargetCurrency: "USD",
		Amounts:        []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	})
	var data convertResponse
	expect(t, r, http.StatusOK, true, &data)
}

func TestConvertInvalidMethod(t *testing.T) {
	r := fireReq("/convert", http.MethodPut, convertRequest{
		BaseCurrency:   "DKK",
		TargetCurrency: "USD",
		Amounts:        []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	})
	expect(t, r, http.StatusBadRequest, true, nil)
}

func TestConvertInvalidCurrency(t *testing.T) {
	r := fireReq("/convert", http.MethodPost, convertRequest{
		BaseCurrency:   "DKK",
		TargetCurrency: "INVALID",
		Amounts:        []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	})
	expect(t, r, http.StatusInternalServerError, true, nil)
}

func TestConvertInvalidRequestData(t *testing.T) {
	r := fireReq("/convert", http.MethodPost, nil)
	expect(t, r, http.StatusInternalServerError, true, nil)
}

func TestNotFoundRoute(t *testing.T) {
	r := fireReq("/notfound", http.MethodGet, nil)
	expect(t, r, http.StatusNotFound, true, nil)
}

func fireReq(endpoint, method string, data interface{}) (rec *httptest.ResponseRecorder) {
	rec = httptest.NewRecorder()
	buff := &bytes.Buffer{}
	if data != nil {
		bytes, _ := json.Marshal(&data)
		buff.Write(bytes)
	}
	req, _ := http.NewRequest(method, "", buff)
	req.RequestURI = endpoint
	server.ServeHTTP(rec, req)

	return rec
}

func expect(t *testing.T, r *httptest.ResponseRecorder, code int, hasCode bool, content interface{}) {
	if r.Code != code && hasCode {
		t.Fatal("Unexpected status:", r.Code)
	}

	if content != nil {
		data, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(data, &content)
		if err != nil {
			t.Fatal(err)
		}
	}
}
