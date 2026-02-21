package ibge

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/VieiraGabrielAlexandre/cityscope-api/internal/contextutil"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) getJSON(ctx context.Context, path string, query url.Values, out any) error {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return fmt.Errorf("invalid base url: %w", err)
	}
	u.Path = u.Path + path
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	start := time.Now()
	requestID := contextutil.GetRequestID(ctx)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		slog.Error("ibge call failed",
			"url", u.String(),
			"error", err,
			"request_id", requestID,
		)
		return fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	latency := time.Since(start)
	slog.Info("ibge call",
		"url", u.String(),
		"status", resp.StatusCode,
		"latency", latency.String(),
		"request_id", requestID,
	)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ibge non-2xx: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}

	return nil
}
