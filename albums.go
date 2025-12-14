package fimage

import (
	"context"
	"fmt"
	"net/http"
)

// AlbumsService handles album operations.
type AlbumsService struct {
	client *Client
}

// CreateAlbumOptions contains options for creating an album.
type CreateAlbumOptions struct {
	// Name is the album name (required).
	Name string

	// Description is an optional album description.
	Description string
}

// UpdateAlbumOptions contains options for updating an album.
type UpdateAlbumOptions struct {
	// Name is the new album name (required).
	Name string

	// Description is the new album description.
	Description string
}

// List returns all albums for the authenticated user.
//
// Example:
//
//	albums, err := client.Albums.List(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, album := range albums {
//	    fmt.Printf("%s (%d files)\n", album.Name, album.FileCount)
//	}
func (s *AlbumsService) List(ctx context.Context) ([]Album, error) {
	var resp struct {
		Albums []Album `json:"albums"`
	}

	if err := s.client.request(ctx, http.MethodGet, "/api/albums", nil, &resp); err != nil {
		return nil, err
	}

	return resp.Albums, nil
}

// Get returns a specific album by ID.
//
// Example:
//
//	album, err := client.Albums.Get(ctx, 123)
//	if err != nil {
//	    if fimage.IsNotFound(err) {
//	        fmt.Println("Album not found")
//	        return
//	    }
//	    log.Fatal(err)
//	}
//	fmt.Printf("Album: %s\n", album.Name)
func (s *AlbumsService) Get(ctx context.Context, albumID int64) (*Album, error) {
	path := fmt.Sprintf("/api/albums/%d", albumID)

	var album Album
	if err := s.client.request(ctx, http.MethodGet, path, nil, &album); err != nil {
		return nil, err
	}

	return &album, nil
}

// Create creates a new album.
//
// Example:
//
//	album, err := client.Albums.Create(ctx, &fimage.CreateAlbumOptions{
//	    Name:        "Vacation Photos",
//	    Description: "Photos from our summer vacation",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Created album: %s (ID: %d)\n", album.Name, album.ID)
func (s *AlbumsService) Create(ctx context.Context, opts *CreateAlbumOptions) (*Album, error) {
	if opts == nil || opts.Name == "" {
		return nil, fmt.Errorf("album name is required")
	}

	req := struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
	}{
		Name:        opts.Name,
		Description: opts.Description,
	}

	var album Album
	if err := s.client.request(ctx, http.MethodPost, "/api/albums", req, &album); err != nil {
		return nil, err
	}

	return &album, nil
}

// Update updates an existing album.
//
// Example:
//
//	album, err := client.Albums.Update(ctx, 123, &fimage.UpdateAlbumOptions{
//	    Name:        "Summer Vacation 2024",
//	    Description: "Updated description",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Updated album: %s\n", album.Name)
func (s *AlbumsService) Update(ctx context.Context, albumID int64, opts *UpdateAlbumOptions) (*Album, error) {
	if opts == nil || opts.Name == "" {
		return nil, fmt.Errorf("album name is required")
	}

	path := fmt.Sprintf("/api/albums/%d", albumID)

	req := struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
	}{
		Name:        opts.Name,
		Description: opts.Description,
	}

	var album Album
	if err := s.client.request(ctx, http.MethodPut, path, req, &album); err != nil {
		return nil, err
	}

	return &album, nil
}

// Delete deletes an album. Files in the album are not deleted,
// they are moved to "no album".
//
// Example:
//
//	err := client.Albums.Delete(ctx, 123)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Album deleted")
func (s *AlbumsService) Delete(ctx context.Context, albumID int64) (*MessageResponse, error) {
	path := fmt.Sprintf("/api/albums/%d", albumID)

	var resp MessageResponse
	if err := s.client.request(ctx, http.MethodDelete, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
