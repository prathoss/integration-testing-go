package picture

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/prathoss/integration_testing/domain"
	"github.com/prathoss/integration_testing/logging"
)

type Client struct {
	baseUri string
	client  *http.Client
}

func NewClient(uri string) *Client {
	return &Client{
		baseUri: uri,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (c Client) GetPicturesByAuthor(
	ctx context.Context,
	authorID uint,
) ([]domain.Picture, error) {
	// Construct URL with query parameter
	reqURL := fmt.Sprintf("%s/api/v1/pictures?author=%d", c.baseUri, authorID)

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.ErrorContext(ctx, "error closing response body", logging.Err(err))
		}
	}()

	// Check for non-200 status
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("reading error response: %w", err)
		}
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	// Decode response
	var pictures []domain.Picture
	if err := json.NewDecoder(resp.Body).Decode(&pictures); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return pictures, nil
}
