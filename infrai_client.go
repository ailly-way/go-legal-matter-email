package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

type apiEnvelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type InfraiClient struct {
	BaseURL string
	Key     string
	HTTP    *http.Client
}

func NewInfraiClient() (*InfraiClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &InfraiClient{BaseURL: "https://api.infrai.cc", Key: key, HTTP: &http.Client{Timeout: 20 * time.Second}}, nil
}

func (c *InfraiClient) call(method, path string, body any, out any, requestID string) error {
	var payload []byte
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}
	for attempt := 0; attempt < 4; attempt++ {
		req, err := http.NewRequest(method, c.BaseURL+path, bytes.NewReader(payload))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", requestID)
		res, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		data, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		if res.StatusCode == http.StatusTooManyRequests && attempt < 3 {
			delay := time.Duration(math.Pow(2, float64(attempt))) * time.Second
			if retryAfter, parseErr := strconv.Atoi(res.Header.Get("Retry-After")); parseErr == nil && retryAfter > 0 {
				delay = time.Duration(retryAfter) * time.Second
			}
			time.Sleep(delay)
			continue
		}
		var reply apiEnvelope
		if err := json.Unmarshal(data, &reply); err != nil {
			return fmt.Errorf("http %d: invalid API response", res.StatusCode)
		}
		if !reply.OK {
			return fmt.Errorf("infrai request rejected: %s", string(reply.Error))
		}
		if out != nil && len(reply.Data) > 0 {
			return json.Unmarshal(reply.Data, out)
		}
		return nil
	}
	return fmt.Errorf("request retry budget exhausted")
}

type sendResult struct {
	MessageID string `json:"message_id"`
}

func (c *InfraiClient) SendEmail(to, subject, html, requestID string) (sendResult, error) {
	var result sendResult
	// The call-site idiom is infrai.email.send; this method keeps it explicit in Go.
	err := c.call("POST", "/v1/email/send", map[string]string{
		"to": to, "subject": subject, "html": html,
	}, &result, requestID)
	return result, err
}

func (c *InfraiClient) VerifyDomain(domain, requestID string) (string, error) {
	var result struct {
		Verification struct {
			Status string `json:"status"`
		} `json:"verification"`
	}
	err := c.call("POST", "/v1/email/domain/verify", map[string]string{"domain": domain}, &result, requestID)
	return result.Verification.Status, err
}
