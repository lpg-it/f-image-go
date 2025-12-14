package fimage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// FilesService handles file operations.
type FilesService struct {
	client *Client
}

// UploadOptions contains options for uploading a file.
type UploadOptions struct {
	// Filename is the name to use for the uploaded file.
	// If empty, a default name will be used.
	Filename string

	// Description is an optional description for the file.
	Description string

	// AlbumID is the optional album to add the file to.
	AlbumID *int64
}

// Upload uploads an image file.
//
// Example:
//
//	file, _ := os.Open("photo.jpg")
//	defer file.Close()
//
//	resp, err := client.Files.Upload(ctx, file, &fimage.UploadOptions{
//	    Filename:    "my-photo.jpg",
//	    Description: "A beautiful sunset",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Uploaded: %s\n", resp.Data.URL)
func (s *FilesService) Upload(ctx context.Context, reader io.Reader, opts *UploadOptions) (*UploadResponse, error) {
	if opts == nil {
		opts = &UploadOptions{}
	}

	filename := opts.Filename
	if filename == "" {
		filename = "image.jpg"
	}

	fields := make(map[string]string)
	if opts.Description != "" {
		fields["description"] = opts.Description
	}

	respBody, err := s.client.uploadMultipart(ctx, "/api/files/upload", reader, filename, fields)
	if err != nil {
		return nil, err
	}

	var resp UploadResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &resp, nil
}

// UploadFromURLOptions contains options for uploading from a URL.
type UploadFromURLOptions struct {
	// URL is the URL to download and upload from.
	URL string
}

// UploadFromURL uploads an image from a public URL.
//
// Example:
//
//	resp, err := client.Files.UploadFromURL(ctx, "https://example.com/image.jpg")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Uploaded: %s\n", resp.Data.URL)
func (s *FilesService) UploadFromURL(ctx context.Context, imageURL string) (*UploadResponse, error) {
	req := struct {
		URL string `json:"url"`
	}{
		URL: imageURL,
	}

	var resp UploadResponse
	if err := s.client.request(ctx, http.MethodPost, "/api/files/upload_from_url", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ListOptions contains options for listing files.
type ListOptions struct {
	// Page is the page number (1-indexed).
	Page int

	// Limit is the number of items per page (max 100).
	Limit int

	// AlbumID filters files by album. Use 0 for files without an album.
	AlbumID *int64
}

// List returns a paginated list of files.
//
// Example:
//
//	// Get first page of files
//	resp, err := client.Files.List(ctx, nil)
//
//	// Get files from a specific album
//	albumID := int64(123)
//	resp, err := client.Files.List(ctx, &fimage.ListOptions{
//	    AlbumID: &albumID,
//	    Page:    1,
//	    Limit:   50,
//	})
func (s *FilesService) List(ctx context.Context, opts *ListOptions) (*FilesListResponse, error) {
	query := url.Values{}

	if opts != nil {
		if opts.Page > 0 {
			query.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.Limit > 0 {
			query.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.AlbumID != nil {
			query.Set("album_id", strconv.FormatInt(*opts.AlbumID, 10))
		}
	}

	var resp FilesListResponse
	if err := s.client.requestWithQuery(ctx, "/api/files", query, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// SearchOptions contains options for searching files.
type SearchOptions struct {
	// Query is the search query string.
	Query string

	// Page is the page number (1-indexed).
	Page int

	// Limit is the number of items per page (max 100).
	Limit int
}

// Search searches for files by filename or description.
//
// Example:
//
//	resp, err := client.Files.Search(ctx, &fimage.SearchOptions{
//	    Query: "sunset",
//	    Page:  1,
//	    Limit: 20,
//	})
//	for _, file := range resp.Files {
//	    fmt.Println(file.OriginalName)
//	}
func (s *FilesService) Search(ctx context.Context, opts *SearchOptions) (*FilesListResponse, error) {
	if opts == nil || opts.Query == "" {
		return nil, fmt.Errorf("search query is required")
	}

	query := url.Values{}
	query.Set("q", opts.Query)

	if opts.Page > 0 {
		query.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.Limit > 0 {
		query.Set("limit", strconv.Itoa(opts.Limit))
	}

	var resp FilesListResponse
	if err := s.client.requestWithQuery(ctx, "/api/files/search", query, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Delete moves a file to trash (soft delete).
//
// Example:
//
//	err := client.Files.Delete(ctx, 123)
//	if err != nil {
//	    log.Fatal(err)
//	}
func (s *FilesService) Delete(ctx context.Context, fileID int64) (*MessageResponse, error) {
	path := fmt.Sprintf("/api/files/%d", fileID)

	var resp MessageResponse
	if err := s.client.request(ctx, http.MethodDelete, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// BatchDelete moves multiple files to trash.
//
// Example:
//
//	resp, err := client.Files.BatchDelete(ctx, []int64{1, 2, 3})
//	fmt.Printf("Deleted: %d, Failed: %d\n", resp.Deleted, resp.Failed)
func (s *FilesService) BatchDelete(ctx context.Context, fileIDs []int64) (*BatchDeleteResponse, error) {
	req := struct {
		FileIDs []int64 `json:"file_ids"`
	}{
		FileIDs: fileIDs,
	}

	var resp BatchDeleteResponse
	if err := s.client.request(ctx, http.MethodPost, "/api/files/batch-delete", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Move moves a single file to an album.
// Set albumID to nil to remove the file from its current album.
//
// Example:
//
//	// Move to album
//	albumID := int64(123)
//	err := client.Files.Move(ctx, 456, &albumID)
//
//	// Remove from album
//	err = client.Files.Move(ctx, 456, nil)
func (s *FilesService) Move(ctx context.Context, fileID int64, albumID *int64) (*MessageResponse, error) {
	path := fmt.Sprintf("/api/files/%d/move", fileID)

	query := url.Values{}
	if albumID != nil {
		query.Set("album_id", strconv.FormatInt(*albumID, 10))
	}

	if len(query) > 0 {
		path = path + "?" + query.Encode()
	}

	var resp MessageResponse
	if err := s.client.request(ctx, http.MethodPut, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// MoveMany moves multiple files to an album.
// Set albumID to nil to remove the files from their current album.
//
// Example:
//
//	albumID := int64(123)
//	err := client.Files.MoveMany(ctx, []int64{1, 2, 3}, &albumID)
func (s *FilesService) MoveMany(ctx context.Context, fileIDs []int64, albumID *int64) (*MessageResponse, error) {
	req := struct {
		FileIDs []int64 `json:"file_ids"`
		AlbumID *int64  `json:"album_id,omitempty"`
	}{
		FileIDs: fileIDs,
		AlbumID: albumID,
	}

	var resp MessageResponse
	if err := s.client.request(ctx, http.MethodPut, "/api/files/move", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
