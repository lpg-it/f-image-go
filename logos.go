package fimage

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// LogosService handles domain logo lookups.
type LogosService struct {
	client *Client
}

// Get returns the public logo URL for a domain when it exists.
//
// The returned Logo always includes the normalized domain. When no logo exists,
// the returned Logo has an empty URL and no error.
func (s *LogosService) Get(ctx context.Context, domain string) (*Logo, error) {
	normalizedDomain := normalizeLogoLookupDomain(domain)
	if normalizedDomain == "" {
		return nil, fmt.Errorf("domain is required")
	}

	path := fmt.Sprintf("/api/logos/%s/exists", url.PathEscape(normalizedDomain))

	var resp struct {
		Exists bool   `json:"exists"`
		Domain string `json:"domain"`
		URL    string `json:"url"`
		ID     int64  `json:"id"`
	}
	if err := s.client.request(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}

	logo := &Logo{
		ID:     resp.ID,
		Domain: normalizedDomain,
		URL:    resp.URL,
		Exists: resp.Exists,
	}
	if resp.Domain != "" {
		logo.Domain = resp.Domain
	}
	if !resp.Exists {
		logo.ID = 0
		logo.URL = ""
		logo.Exists = false
	}

	return logo, nil
}

// Resolve returns a stored logo or fetches and saves one for the domain.
//
// The input may be a bare domain or a URL. If no logo can be found, the
// returned Logo has an empty URL and Exists=false instead of an error.
func (s *LogosService) Resolve(ctx context.Context, domain string) (*Logo, error) {
	normalizedDomain := normalizeLogoLookupDomain(domain)
	if normalizedDomain == "" {
		return nil, fmt.Errorf("domain is required")
	}

	var resp struct {
		Exists   bool   `json:"exists"`
		Domain   string `json:"domain"`
		URL      string `json:"url"`
		ID       int64  `json:"id"`
		Source   string `json:"source"`
		Provider string `json:"provider"`
	}
	if err := s.client.request(ctx, http.MethodPost, "/api/logos/resolve", map[string]string{
		"domain": domain,
	}, &resp); err != nil {
		return nil, err
	}

	logo := &Logo{
		ID:       resp.ID,
		Domain:   normalizedDomain,
		URL:      resp.URL,
		Exists:   resp.Exists,
		Source:   resp.Source,
		Provider: resp.Provider,
	}
	if resp.Domain != "" {
		logo.Domain = resp.Domain
	}
	if !resp.Exists {
		logo.ID = 0
		logo.URL = ""
		logo.Exists = false
		logo.Source = ""
		logo.Provider = ""
	}

	return logo, nil
}

func normalizeLogoLookupDomain(domain string) string {
	cleaned := strings.TrimSpace(strings.ToLower(domain))
	if cleaned == "" {
		return ""
	}

	candidate := cleaned
	if !strings.Contains(candidate, "://") {
		candidate = "https://" + candidate
	}

	parsed, err := url.Parse(candidate)
	if err == nil {
		if host := parsed.Hostname(); host != "" {
			cleaned = host
		}
	}

	cleaned = strings.TrimPrefix(cleaned, "www.")
	cleaned = strings.TrimSuffix(cleaned, ".")

	return cleaned
}
