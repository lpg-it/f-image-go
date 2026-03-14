package fimage

import (
	"context"
	"mime"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUploadLogoOrGetURLReturnsExistingLogoWithoutUpload(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/logos/marriott.com/exists":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"exists":true,"domain":"marriott.com","url":"https://i.f-image.com/logos/marriott.com","id":12}`))
		case "/api/files/upload":
			t.Fatal("upload endpoint should not be called when logo already exists")
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient("test-token", WithBaseURL(server.URL), WithHTTPClient(server.Client()))

	logo, err := client.Files.UploadLogoOrGetURL(context.Background(), nil, &UploadOptions{
		Domain: "https://www.marriott.com/path?x=1",
	})
	if err != nil {
		t.Fatalf("UploadLogoOrGetURL returned error: %v", err)
	}
	if logo.Domain != "marriott.com" {
		t.Fatalf("unexpected domain: %s", logo.Domain)
	}
	if logo.URL != "https://i.f-image.com/logos/marriott.com" {
		t.Fatalf("unexpected url: %s", logo.URL)
	}
}

func TestUploadLogoOrGetURLUploadsWhenMissing(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/logos/marriott.com/exists":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"exists":false,"domain":"marriott.com"}`))
		case "/api/files/upload":
			if r.Method != http.MethodPost {
				t.Fatalf("unexpected method: %s", r.Method)
			}
			if got := r.URL.Query().Get("type"); got != "logo" {
				t.Fatalf("unexpected upload type query: %q", got)
			}
			if got := r.URL.Query().Get("domain"); got != "marriott.com" {
				t.Fatalf("unexpected domain query: %q", got)
			}

			mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
			if err != nil {
				t.Fatalf("failed to parse content type: %v", err)
			}
			if mediaType != "multipart/form-data" {
				t.Fatalf("unexpected content type: %s", mediaType)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"success":true,"status":200,"data":{"id":9,"url":"https://i.f-image.com/logos/marriott.com","upload_type":"logo","domain":"marriott.com","mime_type":"image/png"}}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient("test-token", WithBaseURL(server.URL), WithHTTPClient(server.Client()))

	logo, err := client.Files.UploadLogoOrGetURL(context.Background(), strings.NewReader("fake-image"), &UploadOptions{
		Filename: "logo.png",
		Domain:   "marriott.com",
	})
	if err != nil {
		t.Fatalf("UploadLogoOrGetURL returned error: %v", err)
	}
	if logo.Domain != "marriott.com" {
		t.Fatalf("unexpected domain: %s", logo.Domain)
	}
	if logo.URL != "https://i.f-image.com/logos/marriott.com" {
		t.Fatalf("unexpected url: %s", logo.URL)
	}
	if logo.ID != 9 {
		t.Fatalf("unexpected id: %d", logo.ID)
	}
}

func TestUploadLogoOrGetURLReturnsConflictURLAsSuccess(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/logos/marriott.com/exists":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"exists":false,"domain":"marriott.com"}`))
		case "/api/files/upload":
			query := r.URL.Query()
			if query.Get("type") != "logo" {
				t.Fatalf("unexpected upload type query: %q", query.Get("type"))
			}
			if query.Get("domain") != "marriott.com" {
				t.Fatalf("unexpected domain query: %q", query.Get("domain"))
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(`{"error":"logo already exists for domain","url":"https://i.f-image.com/logos/marriott.com","domain":"marriott.com","exists":true,"force_update_required":true}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient("test-token", WithBaseURL(server.URL), WithHTTPClient(server.Client()))

	logo, err := client.Files.UploadLogoOrGetURL(context.Background(), strings.NewReader("fake-image"), &UploadOptions{
		Filename: "logo.png",
		Domain:   "marriott.com",
	})
	if err != nil {
		t.Fatalf("UploadLogoOrGetURL returned error: %v", err)
	}
	if logo.Domain != "marriott.com" {
		t.Fatalf("unexpected domain: %s", logo.Domain)
	}
	if logo.URL != "https://i.f-image.com/logos/marriott.com" {
		t.Fatalf("unexpected url: %s", logo.URL)
	}
}
