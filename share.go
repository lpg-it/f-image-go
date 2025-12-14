package fimage

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ShareService handles share link operations.
type ShareService struct {
	client *Client
}

// CreateShareOptions contains options for creating a share link.
type CreateShareOptions struct {
	// FileID is the ID of the file to share (for file shares).
	// Either FileID or AlbumID must be set.
	FileID *int64

	// AlbumID is the ID of the album to share (for album shares).
	// Either FileID or AlbumID must be set.
	AlbumID *int64

	// Password is an optional password for the share.
	Password string

	// ExpiresIn is the number of hours until the share expires.
	// Leave as 0 for no expiration.
	ExpiresIn int

	// MaxViews is the maximum number of views allowed.
	// Leave as 0 for unlimited views.
	MaxViews int
}

// UpdateShareOptions contains options for updating a share link.
type UpdateShareOptions struct {
	// Password sets a new password (empty string removes the password).
	Password *string

	// MaxViews sets a new view limit.
	MaxViews *int64

	// IsActive sets whether the share is active.
	IsActive *bool
}

// ShareListOptions contains options for listing share links.
type ShareListOptions struct {
	// Page is the page number (1-indexed).
	Page int

	// Limit is the number of items per page.
	Limit int
}

// List returns all share links for the authenticated user.
//
// Example:
//
//	resp, err := client.Share.List(ctx, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, share := range resp.Shares {
//	    fmt.Printf("Share: %s (views: %d)\n", share.ShareURL, share.ViewCount)
//	}
func (s *ShareService) List(ctx context.Context, opts *ShareListOptions) (*SharesListResponse, error) {
	query := url.Values{}

	if opts != nil {
		if opts.Page > 0 {
			query.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.Limit > 0 {
			query.Set("limit", strconv.Itoa(opts.Limit))
		}
	}

	var resp SharesListResponse
	if err := s.client.requestWithQuery(ctx, "/api/shares", query, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Create creates a new share link.
//
// Example:
//
//	// Share a file
//	fileID := int64(123)
//	share, err := client.Share.Create(ctx, &fimage.CreateShareOptions{
//	    FileID:    &fileID,
//	    Password:  "secret123",
//	    ExpiresIn: 24, // 24 hours
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Share URL: %s\n", share.ShareURL)
//
//	// Share an album
//	albumID := int64(456)
//	share, err = client.Share.Create(ctx, &fimage.CreateShareOptions{
//	    AlbumID:  &albumID,
//	    MaxViews: 100,
//	})
func (s *ShareService) Create(ctx context.Context, opts *CreateShareOptions) (*ShareLink, error) {
	if opts == nil || (opts.FileID == nil && opts.AlbumID == nil) {
		return nil, fmt.Errorf("either FileID or AlbumID is required")
	}

	req := struct {
		FileID    *int64 `json:"file_id,omitempty"`
		AlbumID   *int64 `json:"album_id,omitempty"`
		Password  string `json:"password,omitempty"`
		ExpiresIn int    `json:"expires_in,omitempty"`
		MaxViews  int    `json:"max_views,omitempty"`
	}{
		FileID:    opts.FileID,
		AlbumID:   opts.AlbumID,
		Password:  opts.Password,
		ExpiresIn: opts.ExpiresIn,
		MaxViews:  opts.MaxViews,
	}

	var share ShareLink
	if err := s.client.request(ctx, http.MethodPost, "/api/shares", req, &share); err != nil {
		return nil, err
	}

	return &share, nil
}

// Update updates an existing share link.
//
// Example:
//
//	isActive := false
//	share, err := client.Share.Update(ctx, 123, &fimage.UpdateShareOptions{
//	    IsActive: &isActive,
//	})
func (s *ShareService) Update(ctx context.Context, shareID int64, opts *UpdateShareOptions) (*ShareLink, error) {
	if opts == nil {
		return nil, fmt.Errorf("update options are required")
	}

	path := fmt.Sprintf("/api/shares/%d", shareID)

	req := struct {
		Password *string `json:"password,omitempty"`
		MaxViews *int64  `json:"max_views,omitempty"`
		IsActive *bool   `json:"is_active,omitempty"`
	}{
		Password: opts.Password,
		MaxViews: opts.MaxViews,
		IsActive: opts.IsActive,
	}

	var share ShareLink
	if err := s.client.request(ctx, http.MethodPut, path, req, &share); err != nil {
		return nil, err
	}

	return &share, nil
}

// Delete deletes a share link.
//
// Example:
//
//	err := client.Share.Delete(ctx, 123)
//	if err != nil {
//	    log.Fatal(err)
//	}
func (s *ShareService) Delete(ctx context.Context, shareID int64) (*MessageResponse, error) {
	path := fmt.Sprintf("/api/shares/%d", shareID)

	var resp MessageResponse
	if err := s.client.request(ctx, http.MethodDelete, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Access retrieves the content of a share link.
// This is a public endpoint that doesn't require authentication.
//
// Example:
//
//	content, err := client.Share.Access(ctx, "abc123token")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if content.RequiresPassword {
//	    // Use VerifyPassword to access
//	}
func (s *ShareService) Access(ctx context.Context, token string) (*SharedContent, error) {
	path := fmt.Sprintf("/api/s/%s", token)

	var content SharedContent
	if err := s.client.request(ctx, http.MethodGet, path, nil, &content); err != nil {
		return nil, err
	}

	return &content, nil
}

// VerifyPassword verifies the password for a password-protected share.
// This is a public endpoint that doesn't require authentication.
//
// Example:
//
//	content, err := client.Share.VerifyPassword(ctx, "abc123token", "secret123")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Access granted: %s\n", content.Type)
func (s *ShareService) VerifyPassword(ctx context.Context, token, password string) (*SharedContent, error) {
	path := fmt.Sprintf("/api/s/%s/verify", token)

	req := struct {
		Password string `json:"password"`
	}{
		Password: password,
	}

	var content SharedContent
	if err := s.client.request(ctx, http.MethodPost, path, req, &content); err != nil {
		return nil, err
	}

	return &content, nil
}

// Helper functions for creating options

// ShareFile creates share options for sharing a file.
func ShareFile(fileID int64) *CreateShareOptions {
	return &CreateShareOptions{FileID: &fileID}
}

// ShareAlbum creates share options for sharing an album.
func ShareAlbum(albumID int64) *CreateShareOptions {
	return &CreateShareOptions{AlbumID: &albumID}
}

// WithPassword adds a password to share options.
func (opts *CreateShareOptions) WithPassword(password string) *CreateShareOptions {
	opts.Password = password
	return opts
}

// WithExpiration adds an expiration time to share options.
func (opts *CreateShareOptions) WithExpiration(hours int) *CreateShareOptions {
	opts.ExpiresIn = hours
	return opts
}

// WithMaxViews adds a view limit to share options.
func (opts *CreateShareOptions) WithMaxViews(maxViews int) *CreateShareOptions {
	opts.MaxViews = maxViews
	return opts
}

// ExpiresAt returns the expiration time based on ExpiresIn hours from now.
func (opts *CreateShareOptions) ExpiresAt() *time.Time {
	if opts.ExpiresIn <= 0 {
		return nil
	}
	t := time.Now().Add(time.Duration(opts.ExpiresIn) * time.Hour)
	return &t
}
