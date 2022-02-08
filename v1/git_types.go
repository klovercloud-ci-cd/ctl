package v1

import "time"

// GitWebhook contains github web hook data
type GitWebhook struct {
	Type   string   `json:"type"`
	ID     int      `json:"id"`
	Active bool     `json:"active"`
	Events []string `json:"events"`
	Config struct {
		URL         string `json:"url"`
		InsecureSsl string `json:"insecure_ssl"`
		ContentType string `json:"content_type"`
	} `json:"config"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedAt     time.Time `json:"created_at"`
	URL           string    `json:"url"`
	TestURL       string    `json:"test_url"`
	PingURL       string    `json:"ping_url"`
	DeliveriesURL string    `json:"deliveries_url"`
}