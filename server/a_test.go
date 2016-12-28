package server

import (
	"testing"
	"time"
)

func TestServerCreation(t *testing.T) {
	if server == nil {
		t.Fatal("Server not initialized!")
	}

	if runError != nil {
		t.Fatal("Error starting the server!")
	}

	for loops := 0; !server.hasCurrencies && loops < 10; loops++ {
		t.Log("Sleep:", loops+1)
		time.Sleep(time.Duration(500*(loops+1)) * time.Millisecond)
	}

	if !server.hasCurrencies {
		t.Fatal("Currencies not loaded!")
	}
}
