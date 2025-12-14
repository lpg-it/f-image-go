// Example: Manage trash with F-Image
//
// This example demonstrates how to view, restore, and permanently delete
// files in the trash using the F-Image Go SDK.
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

	// Example 1: List trash items
	fmt.Println("=== Example 1: List trash ===")
	listTrash(ctx, client)

	// Example 2: Restore a single file
	fmt.Println("\n=== Example 2: Restore file ===")
	restoreFile(ctx, client, 1) // Replace with your file ID

	// Example 3: Restore multiple files
	fmt.Println("\n=== Example 3: Restore multiple files ===")
	restoreMultiple(ctx, client, []int64{1, 2, 3}) // Replace with your file IDs

	// Example 4: Permanently delete a file
	fmt.Println("\n=== Example 4: Permanent delete ===")
	permanentDelete(ctx, client, 1) // Replace with your file ID

	// Example 5: Empty trash (DANGEROUS!)
	fmt.Println("\n=== Example 5: Empty trash ===")
	emptyTrash(ctx, client)
}

// listTrash lists all files in the trash.
func listTrash(ctx context.Context, client *fimage.Client) {
	resp, err := client.Trash.List(ctx, nil)
	if err != nil {
		log.Printf("Error listing trash: %v\n", err)
		return
	}

	if resp.Total == 0 {
		fmt.Println("Trash is empty")
		return
	}

	fmt.Printf("Found %d files in trash:\n", resp.Total)
	for _, file := range resp.Files {
		fmt.Printf("  [%d] %s\n", file.ID, file.OriginalName)
		if file.DeletedAt != nil {
			fmt.Printf("       Deleted: %s\n", *file.DeletedAt)
		}
		fmt.Printf("       Size: %d bytes\n", file.Size)
	}
}

// listTrashPaginated demonstrates pagination for trash items.
func listTrashPaginated(ctx context.Context, client *fimage.Client) {
	resp, err := client.Trash.List(ctx, &fimage.TrashListOptions{
		Page:  1,
		Limit: 10,
	})
	if err != nil {
		log.Printf("Error listing trash: %v\n", err)
		return
	}

	totalPages := (resp.Total + int64(resp.Limit) - 1) / int64(resp.Limit)
	fmt.Printf("Page %d of %d (total: %d files)\n", resp.Page, totalPages, resp.Total)

	for _, file := range resp.Files {
		fmt.Printf("  [%d] %s\n", file.ID, file.OriginalName)
	}
}

// restoreFile restores a single file from trash.
func restoreFile(ctx context.Context, client *fimage.Client, fileID int64) {
	resp, err := client.Trash.Restore(ctx, fileID)
	if err != nil {
		if fimage.IsNotFound(err) {
			log.Printf("File not found in trash\n")
			return
		}
		log.Printf("Error restoring file: %v\n", err)
		return
	}

	fmt.Printf("File restored: %s\n", resp.Message)
}

// restoreMultiple restores multiple files from trash.
func restoreMultiple(ctx context.Context, client *fimage.Client, fileIDs []int64) {
	resp, err := client.Trash.RestoreMany(ctx, fileIDs)
	if err != nil {
		log.Printf("Error restoring files: %v\n", err)
		return
	}

	fmt.Printf("Restore complete:\n")
	fmt.Printf("  Restored: %d\n", resp.Restored)
	fmt.Printf("  Failed: %d\n", resp.Failed)
}

// permanentDelete permanently deletes a file from trash.
// WARNING: This action cannot be undone!
func permanentDelete(ctx context.Context, client *fimage.Client, fileID int64) {
	result, err := client.Trash.PermanentDelete(ctx, fileID)
	if err != nil {
		if fimage.IsNotFound(err) {
			log.Printf("File not found in trash\n")
			return
		}
		log.Printf("Error deleting file: %v\n", err)
		return
	}

	if result.Success {
		fmt.Println("File permanently deleted")
	} else {
		fmt.Printf("Delete failed: %s\n", result.Message)

		// Check for share links blocking deletion
		if len(result.FailedDeletions) > 0 {
			for _, failed := range result.FailedDeletions {
				fmt.Printf("  File %s: %s\n", failed.FileName, failed.Reason)
				if len(failed.ShareLinks) > 0 {
					fmt.Printf("    Active share links: %d\n", len(failed.ShareLinks))
				}
			}
		}
	}
}

// emptyTrash permanently deletes all files from trash.
// WARNING: This action cannot be undone!
func emptyTrash(ctx context.Context, client *fimage.Client) {
	// Safety confirmation (in real apps, use proper confirmation)
	fmt.Println("WARNING: This will permanently delete all files in trash!")
	fmt.Println("In a real application, ask for confirmation before proceeding.")
	fmt.Println("Skipping for safety...")
	return

	// Uncomment to actually empty trash:
	/*
		result, err := client.Trash.Empty(ctx)
		if err != nil {
			log.Printf("Error emptying trash: %v\n", err)
			return
		}

		fmt.Printf("Trash emptied:\n")
		fmt.Printf("  Deleted: %d files\n", result.DeletedCount)
		fmt.Printf("  Failed: %d files\n", result.FailedCount)

		if len(result.FailedDeletions) > 0 {
			fmt.Println("\nSome files could not be deleted:")
			for _, failed := range result.FailedDeletions {
				fmt.Printf("  - %s: %s\n", failed.FileName, failed.Reason)
			}
		}
	*/
}
