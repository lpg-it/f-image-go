package fimage

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// TrashService handles trash operations.
type TrashService struct {
	client *Client
}

// TrashListOptions contains options for listing trash items.
type TrashListOptions struct {
	// Page is the page number (1-indexed).
	Page int

	// Limit is the number of items per page.
	Limit int
}

// List returns all files in the trash.
//
// Example:
//
//	resp, err := client.Trash.List(ctx, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, file := range resp.Files {
//	    fmt.Printf("%s (deleted: %s)\n", file.OriginalName, *file.DeletedAt)
//	}
func (s *TrashService) List(ctx context.Context, opts *TrashListOptions) (*TrashListResponse, error) {
	query := url.Values{}

	if opts != nil {
		if opts.Page > 0 {
			query.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.Limit > 0 {
			query.Set("limit", strconv.Itoa(opts.Limit))
		}
	}

	var resp TrashListResponse
	if err := s.client.requestWithQuery(ctx, "/api/trash", query, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Restore restores a single file from trash.
//
// Example:
//
//	resp, err := client.Trash.Restore(ctx, 123)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(resp.Message)
func (s *TrashService) Restore(ctx context.Context, fileID int64) (*RestoreResponse, error) {
	path := fmt.Sprintf("/api/trash/%d/restore", fileID)

	var resp RestoreResponse
	if err := s.client.request(ctx, http.MethodPost, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// RestoreMany restores multiple files from trash.
//
// Example:
//
//	resp, err := client.Trash.RestoreMany(ctx, []int64{1, 2, 3})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Restored: %d, Failed: %d\n", resp.Restored, resp.Failed)
func (s *TrashService) RestoreMany(ctx context.Context, fileIDs []int64) (*RestoreResponse, error) {
	req := struct {
		FileIDs []int64 `json:"file_ids"`
	}{
		FileIDs: fileIDs,
	}

	var resp RestoreResponse
	if err := s.client.request(ctx, http.MethodPost, "/api/trash/restore", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// PermanentDelete permanently deletes a file from trash.
// This action cannot be undone.
//
// Example:
//
//	result, err := client.Trash.PermanentDelete(ctx, 123)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if result.Success {
//	    fmt.Println("File permanently deleted")
//	} else {
//	    fmt.Printf("Failed: %s\n", result.Message)
//	}
func (s *TrashService) PermanentDelete(ctx context.Context, fileID int64) (*DeleteResult, error) {
	path := fmt.Sprintf("/api/trash/%d", fileID)

	var result DeleteResult
	if err := s.client.request(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Empty permanently deletes all files from trash.
// This action cannot be undone.
//
// Example:
//
//	result, err := client.Trash.Empty(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Deleted: %d files\n", result.DeletedCount)
//	if result.FailedCount > 0 {
//	    fmt.Printf("Failed: %d files (may have active share links)\n", result.FailedCount)
//	}
func (s *TrashService) Empty(ctx context.Context) (*DeleteResult, error) {
	var result DeleteResult
	if err := s.client.request(ctx, http.MethodDelete, "/api/trash/empty", nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
