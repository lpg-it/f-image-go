package fimage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the default F-Image API base URL.
	DefaultBaseURL = "https://f-image.com"

	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 30 * time.Second

	// Version is the current SDK version.
	Version = "1.0.0"
)

// Client is the F-Image API client.
type Client struct {
	// BaseURL is the base URL for API requests.
	BaseURL string

	// HTTPClient is the HTTP client used for API requests.
	HTTPClient *http.Client

	// apiToken is the API token for authentication.
	apiToken string

	// userAgent is the User-Agent header value.
	userAgent string

	// Services
	Files  *FilesService
	Albums *AlbumsService
	Share  *ShareService
	Tags   *TagsService
	Trash  *TrashService
}

// ClientOption is a function that configures the Client.
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the client.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.BaseURL = strings.TrimSuffix(baseURL, "/")
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.HTTPClient.Timeout = timeout
	}
}

// WithUserAgent sets a custom User-Agent header.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// NewClient creates a new F-Image API client.
//
// The apiToken is required and can be obtained from your F-Image dashboard
// at https://f-image.com/dashboard/settings/api
//
// Example:
//
//	client := fimage.NewClient("fimg_live_your_token_here")
//
//	// With options
//	client := fimage.NewClient("fimg_live_your_token_here",
//	    fimage.WithTimeout(60*time.Second),
//	    fimage.WithBaseURL("https://custom-api.example.com"),
//	)
func NewClient(apiToken string, opts ...ClientOption) *Client {
	c := &Client{
		BaseURL: DefaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		apiToken:  apiToken,
		userAgent: fmt.Sprintf("f-image-go/%s", Version),
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	// Initialize services
	c.Files = &FilesService{client: c}
	c.Albums = &AlbumsService{client: c}
	c.Share = &ShareService{client: c}
	c.Tags = &TagsService{client: c}
	c.Trash = &TrashService{client: c}

	return c
}

// request performs an HTTP request and decodes the response.
func (c *Client) request(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	// Build URL
	reqURL := c.BaseURL + path

	// Prepare request body
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("User-Agent", c.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseAPIError(resp.StatusCode, respBody)
	}

	// Decode response
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// requestWithQuery performs an HTTP GET request with query parameters.
func (c *Client) requestWithQuery(ctx context.Context, path string, query url.Values, result interface{}) error {
	if len(query) > 0 {
		path = path + "?" + query.Encode()
	}
	return c.request(ctx, http.MethodGet, path, nil, result)
}

// uploadMultipart performs a multipart file upload.
func (c *Client) uploadMultipart(ctx context.Context, path string, reader io.Reader, filename string, fields map[string]string) ([]byte, error) {
	// Create multipart writer
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add file field
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, reader); err != nil {
		return nil, fmt.Errorf("failed to copy file data: %w", err)
	}

	// Add other fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return nil, fmt.Errorf("failed to write field %s: %w", key, err)
		}
	}

	// Close writer
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Build URL
	reqURL := c.BaseURL + path

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseAPIError(resp.StatusCode, respBody)
	}

	return respBody, nil
}

// parseAPIError parses an API error response.
func parseAPIError(statusCode int, body []byte) error {
	var errResp struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &errResp); err != nil {
		return &APIError{
			StatusCode: statusCode,
			Message:    string(body),
		}
	}

	msg := errResp.Error
	if msg == "" {
		msg = errResp.Message
	}
	if msg == "" {
		msg = http.StatusText(statusCode)
	}

	return &APIError{
		StatusCode: statusCode,
		Message:    msg,
	}
}
