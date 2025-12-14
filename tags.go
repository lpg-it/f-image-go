package fimage

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// TagsService handles tag operations.
type TagsService struct {
	client *Client
}

// CreateTagOptions contains options for creating a tag.
type CreateTagOptions struct {
	// Name is the tag name (required).
	Name string

	// Color is the tag color in hex format (e.g., "#FF5733").
	Color string
}

// UpdateTagOptions contains options for updating a tag.
type UpdateTagOptions struct {
	// Name is the new tag name.
	Name string

	// Color is the new tag color in hex format.
	Color string
}

// TagFilesOptions contains options for listing files by tag.
type TagFilesOptions struct {
	// Page is the page number (1-indexed).
	Page int

	// Limit is the number of items per page.
	Limit int
}

// List returns all tags for the authenticated user.
//
// Example:
//
//	tags, err := client.Tags.List(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, tag := range tags {
//	    fmt.Printf("%s (%d files)\n", tag.Name, tag.FileCount)
//	}
func (s *TagsService) List(ctx context.Context) ([]Tag, error) {
	var tags []Tag
	if err := s.client.request(ctx, http.MethodGet, "/api/tags", nil, &tags); err != nil {
		return nil, err
	}

	return tags, nil
}

// Create creates a new tag.
//
// Example:
//
//	tag, err := client.Tags.Create(ctx, &fimage.CreateTagOptions{
//	    Name:  "Nature",
//	    Color: "#4CAF50",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Created tag: %s (ID: %d)\n", tag.Name, tag.ID)
func (s *TagsService) Create(ctx context.Context, opts *CreateTagOptions) (*Tag, error) {
	if opts == nil || opts.Name == "" {
		return nil, fmt.Errorf("tag name is required")
	}

	req := struct {
		Name  string `json:"name"`
		Color string `json:"color,omitempty"`
	}{
		Name:  opts.Name,
		Color: opts.Color,
	}

	var tag Tag
	if err := s.client.request(ctx, http.MethodPost, "/api/tags", req, &tag); err != nil {
		return nil, err
	}

	return &tag, nil
}

// Update updates an existing tag.
//
// Example:
//
//	tag, err := client.Tags.Update(ctx, 123, &fimage.UpdateTagOptions{
//	    Name:  "Wildlife",
//	    Color: "#2196F3",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Updated tag: %s\n", tag.Name)
func (s *TagsService) Update(ctx context.Context, tagID int64, opts *UpdateTagOptions) (*Tag, error) {
	if opts == nil {
		return nil, fmt.Errorf("update options are required")
	}

	path := fmt.Sprintf("/api/tags/%d", tagID)

	req := struct {
		Name  string `json:"name,omitempty"`
		Color string `json:"color,omitempty"`
	}{
		Name:  opts.Name,
		Color: opts.Color,
	}

	var tag Tag
	if err := s.client.request(ctx, http.MethodPut, path, req, &tag); err != nil {
		return nil, err
	}

	return &tag, nil
}

// Delete deletes a tag. The tag is removed from all files.
//
// Example:
//
//	err := client.Tags.Delete(ctx, 123)
//	if err != nil {
//	    log.Fatal(err)
//	}
func (s *TagsService) Delete(ctx context.Context, tagID int64) (*MessageResponse, error) {
	path := fmt.Sprintf("/api/tags/%d", tagID)

	var resp MessageResponse
	if err := s.client.request(ctx, http.MethodDelete, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// TagFile adds a tag to a file.
//
// Example:
//
//	err := client.Tags.TagFile(ctx, 456, 123) // Add tag 123 to file 456
//	if err != nil {
//	    log.Fatal(err)
//	}
func (s *TagsService) TagFile(ctx context.Context, fileID, tagID int64) (*MessageResponse, error) {
	req := struct {
		FileID int64 `json:"file_id"`
		TagID  int64 `json:"tag_id"`
	}{
		FileID: fileID,
		TagID:  tagID,
	}

	var resp MessageResponse
	if err := s.client.request(ctx, http.MethodPost, "/api/tags/file", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// UntagFile removes a tag from a file.
//
// Example:
//
//	err := client.Tags.UntagFile(ctx, 456, 123) // Remove tag 123 from file 456
//	if err != nil {
//	    log.Fatal(err)
//	}
func (s *TagsService) UntagFile(ctx context.Context, fileID, tagID int64) (*MessageResponse, error) {
	req := struct {
		FileID int64 `json:"file_id"`
		TagID  int64 `json:"tag_id"`
	}{
		FileID: fileID,
		TagID:  tagID,
	}

	var resp MessageResponse
	if err := s.client.request(ctx, http.MethodDelete, "/api/tags/file", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetFiles returns all files with a specific tag.
//
// Example:
//
//	resp, err := client.Tags.GetFiles(ctx, 123, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, file := range resp.Files {
//	    fmt.Println(file.OriginalName)
//	}
func (s *TagsService) GetFiles(ctx context.Context, tagID int64, opts *TagFilesOptions) (*FilesListResponse, error) {
	path := fmt.Sprintf("/api/tags/%d/files", tagID)

	query := url.Values{}
	if opts != nil {
		if opts.Page > 0 {
			query.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.Limit > 0 {
			query.Set("limit", strconv.Itoa(opts.Limit))
		}
	}

	if len(query) > 0 {
		path = path + "?" + query.Encode()
	}

	var resp FilesListResponse
	if err := s.client.request(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
