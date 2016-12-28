package server

import (
	"fmt"
	"testing"
)

func TestCurrencyConversionNewServer(t *testing.T) {
	server, err := New()
	if err != nil {
		t.Fatal(err)
	}

	_, err = server.convert("USD", "DKK", 1)
	if err == nil {
		t.Fatal("Expected to fail!")
	}

	_, err = server.createResponse("DKK")
	if err == nil {
		t.Fatal("Expected to fail!")
	}
}

func TestCurrencyConversion(t *testing.T) {
	amount := 123.45
	converted, err := server.convert("USD", "DKK", amount)
	if err != nil {
		t.Fatal(err)
	}

	if amount == converted {
		t.Fatal("Converted value is the same")
	}

	converted2, err2 := server.convert("DKK", "USD", converted)
	if err2 != nil {
		t.Fatal(err2)
	}

	if fmt.Sprintf("%.2f", amount) != fmt.Sprintf("%.2f", converted2) {
		t.Fatal("Expected to be the same:", amount, converted2)
	}

	converted3, err3 := server.convert("DKK", "DKK", amount)
	if err3 != nil {
		t.Fatal(err3)
	}

	if converted3 != amount {
		t.Fatal("Expected to be the same:", amount, converted3)
	}
}

func TestUnknownCurrency(t *testing.T) {
	amount := 123.45
	_, err := server.convert("FOO", "DKK", amount)
	if err == nil {
		t.Fatal("Shouldn't know currency: FOO")
	}

	_, err = server.convert("DKK", "BAR", amount)
	if err == nil {
		t.Fatal("Shouldn't know currency: BAR")
	}
}

func TestConvertResponseCreation(t *testing.T) {
	amounts := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	response, err := server.createConvertResponse("DKK", "EUR", amounts)
	if err != nil {
		t.Fatal(err)
	}

	for i, a := range response.ConvertedAmounts {
		if a != float64(i)*response.ConvertedAmounts[1] {
			t.Fatal("Unexpected amount at index:", i)
		}
	}
}

func TestConvertedResponseUnknownCurrency(t *testing.T) {
	amounts := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	_, err := server.createConvertResponse("FOO", "DKK", amounts)
	if err == nil {
		t.Fatal("Currency shouldn't be known: FOO")
	}
}

func TestCreateCurrencyResponse(t *testing.T) {
	_, err := server.createResponse("USD")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateInvalidCurrencyResponse(t *testing.T) {
	_, err := server.createResponse("FOO")
	if err == nil {
		t.Fatal("Currency shouldn't be known: FOO")
	}
}
