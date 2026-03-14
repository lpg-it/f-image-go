package fimage

import (
	"context"
	"errors"
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

	logo, err := s.getViaExistsEndpoint(ctx, normalizedDomain)
	if err == nil {
		return logo, nil
	}
	if !IsNotFound(err) {
		return nil, err
	}

	return s.getViaLegacyEndpoint(ctx, normalizedDomain)
}

func (s *LogosService) getViaExistsEndpoint(ctx context.Context, domain string) (*Logo, error) {
	path := fmt.Sprintf("/api/logos/%s/exists", url.PathEscape(domain))

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
		Domain: domain,
		URL:    resp.URL,
	}
	if resp.Domain != "" {
		logo.Domain = resp.Domain
	}
	if !resp.Exists {
		logo.ID = 0
		logo.URL = ""
	}

	return logo, nil
}

func (s *LogosService) getViaLegacyEndpoint(ctx context.Context, domain string) (*Logo, error) {
	path := fmt.Sprintf("/api/logos/%s", url.PathEscape(domain))

	var logo Logo
	if err := s.client.request(ctx, http.MethodGet, path, nil, &logo); err != nil {
		if IsNotFound(err) {
			if isRouteNotFound(err) {
				return nil, err
			}
			return &Logo{Domain: domain}, nil
		}
		return nil, err
	}

	if logo.Domain == "" {
		logo.Domain = domain
	}

	return &logo, nil
}

func isRouteNotFound(err error) bool {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return false
	}

	message := strings.TrimSpace(strings.ToLower(apiErr.Message))

	return message == "404 page not found" || message == strings.ToLower(http.StatusText(http.StatusNotFound))
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
