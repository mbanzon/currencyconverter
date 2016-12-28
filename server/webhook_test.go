package server

import (
	"net/http"
	"testing"
)

const (
	webhookServerAddr = "127.0.0.1:54189"
)

var (
	hookServer *webhookServer
)

type webhookServer struct {
	secret string
	calls  int
}

func (s *webhookServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != s.secret {
		http.Error(w, "", http.StatusForbidden)
		return
	}

	s.calls++
}

func init() {
	hookServer = &webhookServer{secret: "verysecret"}
	go http.ListenAndServe(webhookServerAddr, hookServer)
}

func TestWebhookRegister(t *testing.T) {
	count := hookServer.calls
	r := fireReq("/webhook", http.MethodPost, webhook{
		BaseCurrency: "DKK",
		Secret:       "verysecret",
		Url:          "http://" + webhookServerAddr,
	})
	expect(t, r, http.StatusOK, true, nil)
	if count+1 != hookServer.calls {
		t.Fatal("Expected one call after registration")
	}
}

func TestWebhookCalling(t *testing.T) {
	count := hookServer.calls
	server.callWebhooks()
	if count+1 != hookServer.calls {
		t.Fatal("Expected one call extra")
	}
}

func TestWebhookWrongSecret(t *testing.T) {
	count := hookServer.calls
	r := fireReq("/webhook", http.MethodPost, webhook{
		BaseCurrency: "DKK",
		Secret:       "wrongsecret",
		Url:          "http://" + webhookServerAddr,
	})
	expect(t, r, http.StatusInternalServerError, true, nil)
	if count != hookServer.calls {
		t.Fatal("Expected no call after registration")
	}
}

func TestWebhookWrongBase(t *testing.T) {
	r := fireReq("/webhook", http.MethodPost, webhook{
		BaseCurrency: "INVALID",
		Secret:       "verysecret",
		Url:          "http://" + webhookServerAddr,
	})
	expect(t, r, http.StatusInternalServerError, true, nil)
}

func TestWebhookGet(t *testing.T) {
	r := fireReq("/webhook", http.MethodGet, nil)
	expect(t, r, http.StatusNotFound, true, nil)
}
