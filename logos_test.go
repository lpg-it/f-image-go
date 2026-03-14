package fimage

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogosGetUsesExistsEndpoint(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/logos/marriott.com/exists" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"exists":true,"domain":"marriott.com","url":"https://i.f-image.com/logos/marriott.com","id":12}`))
	}))
	defer server.Close()

	client := NewClient("test-token", WithBaseURL(server.URL), WithHTTPClient(server.Client()))

	logo, err := client.Logos.Get(context.Background(), "https://www.marriott.com/path?x=1")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if logo.Domain != "marriott.com" {
		t.Fatalf("unexpected domain: %s", logo.Domain)
	}
	if logo.URL != "https://i.f-image.com/logos/marriott.com" {
		t.Fatalf("unexpected url: %s", logo.URL)
	}
	if logo.ID != 12 {
		t.Fatalf("unexpected id: %d", logo.ID)
	}
}

func TestLogosGetReturnsEmptyURLWhenMissing(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/logos/missing.com/exists" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"exists":false,"domain":"missing.com"}`))
	}))
	defer server.Close()

	client := NewClient("test-token", WithBaseURL(server.URL), WithHTTPClient(server.Client()))

	logo, err := client.Logos.Get(context.Background(), "missing.com")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if logo.Domain != "missing.com" {
		t.Fatalf("unexpected domain: %s", logo.Domain)
	}
	if logo.URL != "" {
		t.Fatalf("expected empty url, got: %s", logo.URL)
	}
}

func TestLogosGetFallsBackToLegacyEndpoint(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/logos/example.com/exists":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found"}`))
		case "/api/logos/example.com":
			_, _ = w.Write([]byte(`{"id":7,"domain":"example.com","url":"https://i.f-image.com/logos/example.com"}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient("test-token", WithBaseURL(server.URL), WithHTTPClient(server.Client()))

	logo, err := client.Logos.Get(context.Background(), "example.com")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if logo.Domain != "example.com" {
		t.Fatalf("unexpected domain: %s", logo.Domain)
	}
	if logo.URL != "https://i.f-image.com/logos/example.com" {
		t.Fatalf("unexpected url: %s", logo.URL)
	}
	if logo.ID != 7 {
		t.Fatalf("unexpected id: %d", logo.ID)
	}
}

func TestLogosGetReturnsEmptyURLWhenLegacyEndpointMisses(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/logos/example.com/exists":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found"}`))
		case "/api/logos/example.com":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Logo not found"}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient("test-token", WithBaseURL(server.URL), WithHTTPClient(server.Client()))

	logo, err := client.Logos.Get(context.Background(), "example.com")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if logo.Domain != "example.com" {
		t.Fatalf("unexpected domain: %s", logo.Domain)
	}
	if logo.URL != "" {
		t.Fatalf("expected empty url, got: %s", logo.URL)
	}
}

func TestLogosGetReturnsErrorWhenLegacyRouteIsMissing(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/logos/example.com/exists":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("404 page not found"))
		case "/api/logos/example.com":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("404 page not found"))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient("test-token", WithBaseURL(server.URL), WithHTTPClient(server.Client()))

	if _, err := client.Logos.Get(context.Background(), "example.com"); err == nil {
		t.Fatal("expected error when legacy route is missing")
	}
}
