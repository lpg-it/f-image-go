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

func TestLogosGetReturnsErrorWhenExistsRouteMissing(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/logos/example.com/exists" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found"}`))
	}))
	defer server.Close()

	client := NewClient("test-token", WithBaseURL(server.URL), WithHTTPClient(server.Client()))

	if _, err := client.Logos.Get(context.Background(), "example.com"); err == nil {
		t.Fatal("expected error when exists route is missing")
	}
}

func TestLogosResolveFetchesWhenMissing(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/api/logos/resolve" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"exists":true,"domain":"stripe.com","url":"https://i.f-image.com/logos/stripe.com","id":21,"source":"fetched","provider":"hunter"}`))
	}))
	defer server.Close()

	client := NewClient("test-token", WithBaseURL(server.URL), WithHTTPClient(server.Client()))

	logo, err := client.Logos.Resolve(context.Background(), "https://www.stripe.com/payments")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if logo.Domain != "stripe.com" {
		t.Fatalf("unexpected domain: %s", logo.Domain)
	}
	if logo.URL != "https://i.f-image.com/logos/stripe.com" {
		t.Fatalf("unexpected url: %s", logo.URL)
	}
	if logo.Source != "fetched" || logo.Provider != "hunter" {
		t.Fatalf("unexpected source metadata: %#v", logo)
	}
}

func TestLogosResolveReturnsEmptyWhenMissing(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"exists":false,"domain":"missing.com"}`))
	}))
	defer server.Close()

	client := NewClient("test-token", WithBaseURL(server.URL), WithHTTPClient(server.Client()))

	logo, err := client.Logos.Resolve(context.Background(), "missing.com")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if logo.URL != "" || logo.Exists {
		t.Fatalf("expected empty miss, got %#v", logo)
	}
}
