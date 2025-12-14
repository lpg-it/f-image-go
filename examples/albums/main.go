// Example: Manage albums with F-Image
//
// This example demonstrates how to create, list, update, and delete albums
// using the F-Image Go SDK.
//
// Usage:
//
//	export FIMAGE_API_TOKEN="your-api-token"
//	go run main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	fimage "github.com/lpg-it/f-image-go"
)

func main() {
	// Get API token from environment variable
	apiToken := os.Getenv("FIMAGE_API_TOKEN")
	if apiToken == "" {
		log.Fatal("FIMAGE_API_TOKEN environment variable is required")
	}

	// Create a new client
	client := fimage.NewClient(apiToken)

	ctx := context.Background()

	// Example 1: Create an album
	fmt.Println("=== Example 1: Create album ===")
	album, err := createAlbum(ctx, client)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	// Example 2: List all albums
	fmt.Println("\n=== Example 2: List albums ===")
	listAlbums(ctx, client)

	// Example 3: Get a specific album
	fmt.Println("\n=== Example 3: Get album ===")
	getAlbum(ctx, client, album.ID)

	// Example 4: Update an album
	fmt.Println("\n=== Example 4: Update album ===")
	updateAlbum(ctx, client, album.ID)

	// Example 5: Delete an album
	fmt.Println("\n=== Example 5: Delete album ===")
	deleteAlbum(ctx, client, album.ID)
}

// createAlbum creates a new album.
func createAlbum(ctx context.Context, client *fimage.Client) (*fimage.Album, error) {
	album, err := client.Albums.Create(ctx, &fimage.CreateAlbumOptions{
		Name:        "Vacation Photos 2024",
		Description: "Photos from our summer vacation",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create album: %w", err)
	}

	fmt.Printf("Album created successfully!\n")
	fmt.Printf("  ID: %d\n", album.ID)
	fmt.Printf("  Name: %s\n", album.Name)
	fmt.Printf("  Description: %s\n", album.Description)
	fmt.Printf("  Created At: %s\n", album.CreatedAt)

	return album, nil
}

// listAlbums lists all albums.
func listAlbums(ctx context.Context, client *fimage.Client) {
	albums, err := client.Albums.List(ctx)
	if err != nil {
		log.Printf("Error listing albums: %v\n", err)
		return
	}

	fmt.Printf("Found %d albums:\n", len(albums))
	for _, album := range albums {
		fmt.Printf("  [%d] %s (%d files)\n", album.ID, album.Name, album.FileCount)
		if album.Description != "" {
			fmt.Printf("       Description: %s\n", album.Description)
		}
	}
}

// getAlbum gets a specific album by ID.
func getAlbum(ctx context.Context, client *fimage.Client, albumID int64) {
	album, err := client.Albums.Get(ctx, albumID)
	if err != nil {
		if fimage.IsNotFound(err) {
			log.Printf("Album not found\n")
			return
		}
		log.Printf("Error getting album: %v\n", err)
		return
	}

	fmt.Printf("Album details:\n")
	fmt.Printf("  ID: %d\n", album.ID)
	fmt.Printf("  Name: %s\n", album.Name)
	fmt.Printf("  Description: %s\n", album.Description)
	fmt.Printf("  File Count: %d\n", album.FileCount)
	fmt.Printf("  Created At: %s\n", album.CreatedAt)
}

// updateAlbum updates an album.
func updateAlbum(ctx context.Context, client *fimage.Client, albumID int64) {
	album, err := client.Albums.Update(ctx, albumID, &fimage.UpdateAlbumOptions{
		Name:        "Summer Vacation 2024",
		Description: "Updated: Photos from our amazing summer vacation",
	})
	if err != nil {
		log.Printf("Error updating album: %v\n", err)
		return
	}

	fmt.Printf("Album updated successfully!\n")
	fmt.Printf("  ID: %d\n", album.ID)
	fmt.Printf("  Name: %s\n", album.Name)
	fmt.Printf("  Description: %s\n", album.Description)
}

// deleteAlbum deletes an album.
func deleteAlbum(ctx context.Context, client *fimage.Client, albumID int64) {
	resp, err := client.Albums.Delete(ctx, albumID)
	if err != nil {
		if fimage.IsNotFound(err) {
			log.Printf("Album not found\n")
			return
		}
		log.Printf("Error deleting album: %v\n", err)
		return
	}

	fmt.Printf("Album deleted: %s\n", resp.Message)
}
