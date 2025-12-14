// Example: Manage files with F-Image
//
// This example demonstrates how to list, search, delete, and move files
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

	// Example 1: List all files
	fmt.Println("=== Example 1: List files ===")
	listFiles(ctx, client)

	// Example 2: List files with pagination
	fmt.Println("\n=== Example 2: Paginated list ===")
	listFilesPaginated(ctx, client)

	// Example 3: List files in a specific album
	fmt.Println("\n=== Example 3: Files by album ===")
	listFilesByAlbum(ctx, client, 1) // Replace with your album ID

	// Example 4: Search files
	fmt.Println("\n=== Example 4: Search files ===")
	searchFiles(ctx, client, "sunset")

	// Example 5: Move file to album
	fmt.Println("\n=== Example 5: Move file to album ===")
	moveFileToAlbum(ctx, client, 1, 2) // Replace with your file ID and album ID

	// Example 6: Delete a file (move to trash)
	fmt.Println("\n=== Example 6: Delete file ===")
	deleteFile(ctx, client, 1) // Replace with your file ID

	// Example 7: Batch delete files
	fmt.Println("\n=== Example 7: Batch delete ===")
	batchDeleteFiles(ctx, client, []int64{1, 2, 3}) // Replace with your file IDs
}

// listFiles lists all files with default pagination.
func listFiles(ctx context.Context, client *fimage.Client) {
	resp, err := client.Files.List(ctx, nil)
	if err != nil {
		log.Printf("Error listing files: %v\n", err)
		return
	}

	fmt.Printf("Found %d files (showing page 1)\n", resp.Total)
	for _, file := range resp.Files {
		fmt.Printf("  [%d] %s (%dx%d, %d bytes)\n",
			file.ID, file.OriginalName, file.Width, file.Height, file.Size)
	}
}

// listFilesPaginated demonstrates pagination.
func listFilesPaginated(ctx context.Context, client *fimage.Client) {
	resp, err := client.Files.List(ctx, &fimage.ListOptions{
		Page:  1,
		Limit: 10,
	})
	if err != nil {
		log.Printf("Error listing files: %v\n", err)
		return
	}

	totalPages := (resp.Total + int64(resp.Limit) - 1) / int64(resp.Limit)
	fmt.Printf("Page %d of %d (total: %d files)\n", resp.Page, totalPages, resp.Total)

	for _, file := range resp.Files {
		fmt.Printf("  [%d] %s - %s\n", file.ID, file.OriginalName, file.URL)
	}
}

// listFilesByAlbum lists files in a specific album.
func listFilesByAlbum(ctx context.Context, client *fimage.Client, albumID int64) {
	resp, err := client.Files.List(ctx, &fimage.ListOptions{
		AlbumID: &albumID,
	})
	if err != nil {
		log.Printf("Error listing files: %v\n", err)
		return
	}

	fmt.Printf("Found %d files in album %d\n", resp.Total, albumID)
	for _, file := range resp.Files {
		fmt.Printf("  [%d] %s\n", file.ID, file.OriginalName)
	}
}

// searchFiles searches for files by name or description.
func searchFiles(ctx context.Context, client *fimage.Client, query string) {
	resp, err := client.Files.Search(ctx, &fimage.SearchOptions{
		Query: query,
		Page:  1,
		Limit: 20,
	})
	if err != nil {
		log.Printf("Error searching files: %v\n", err)
		return
	}

	fmt.Printf("Found %d files matching '%s'\n", resp.Total, query)
	for _, file := range resp.Files {
		fmt.Printf("  [%d] %s - %s\n", file.ID, file.OriginalName, file.Description)
	}
}

// moveFileToAlbum moves a file to an album.
func moveFileToAlbum(ctx context.Context, client *fimage.Client, fileID, albumID int64) {
	resp, err := client.Files.Move(ctx, fileID, &albumID)
	if err != nil {
		if fimage.IsNotFound(err) {
			log.Printf("File or album not found\n")
			return
		}
		log.Printf("Error moving file: %v\n", err)
		return
	}

	fmt.Printf("File moved successfully: %s\n", resp.Message)
}

// deleteFile moves a file to trash.
func deleteFile(ctx context.Context, client *fimage.Client, fileID int64) {
	resp, err := client.Files.Delete(ctx, fileID)
	if err != nil {
		if fimage.IsNotFound(err) {
			log.Printf("File not found\n")
			return
		}
		log.Printf("Error deleting file: %v\n", err)
		return
	}

	fmt.Printf("File deleted: %s\n", resp.Message)
	if resp.Info != "" {
		fmt.Printf("  Info: %s\n", resp.Info)
	}
}

// batchDeleteFiles deletes multiple files at once.
func batchDeleteFiles(ctx context.Context, client *fimage.Client, fileIDs []int64) {
	resp, err := client.Files.BatchDelete(ctx, fileIDs)
	if err != nil {
		log.Printf("Error batch deleting files: %v\n", err)
		return
	}

	fmt.Printf("Batch delete complete:\n")
	fmt.Printf("  Deleted: %d\n", resp.Deleted)
	fmt.Printf("  Failed: %d\n", resp.Failed)
	fmt.Printf("  Message: %s\n", resp.Message)
}
