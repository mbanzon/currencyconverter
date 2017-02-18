package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// verifies a single webhook. Looks up the base currency, parses the
// URL and attempts to call the webhook
func (s *Server) verifyWebhook(hook webhook) bool {
	if !s.hasCurrencies {
		return false
	}

	if _, hasBase := s.currencies[hook.BaseCurrency]; !hasBase {
		return false
	}

	if _, err := url.Parse(hook.Url); err != nil {
		return false
	}

	err := s.callSingleWebhook(hook)
	if err != nil {
		return false
	}

	return true
}

// calls all webhooks
func (s *Server) callWebhooks() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, hook := range s.webhooks {
		s.callSingleWebhook(hook)
	}
}

// calls a single webhook
func (s *Server) callSingleWebhook(hook webhook) (err error) {
	// creates a "response" using the base currency of the webhook
	cRes, err := s.createResponse(hook.BaseCurrency)
	if err != nil {
		log.Println("Error creating data for webhook:", err)
		return err
	}

	// encodes the "response" to be used as a request payload
	data := &bytes.Buffer{}
	err = json.NewEncoder(data).Encode(&cRes)
	if err != nil {
		log.Println("Error creating data for webhook:", err)
		return err
	}

	// creates a new request using the payload data and the webhook URL
	req, err := http.NewRequest(http.MethodPost, hook.Url, data)
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}

	// sets headers for content type and authorization with the
	// webhook secret
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", hook.Secret)

	// make the request, log return code or errors
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Webhook call error:", err)
		return err
	} else {
		log.Printf("Webhook return code: %d\n", res.StatusCode)
	}

	s.webhookTriggers.Add(1)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Webhook returned: %d", res.StatusCode)
	}

	return nil
}
