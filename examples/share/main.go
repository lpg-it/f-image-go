// Example: Share files and albums with F-Image
//
// This example demonstrates how to create, manage, and access share links
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

	// Example 1: Create a simple share link
	fmt.Println("=== Example 1: Simple share ===")
	share, err := createSimpleShare(ctx, client, 1) // Replace with your file ID
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Share URL: %s\n", share.ShareURL)
	}

	// Example 2: Create a password-protected share
	fmt.Println("\n=== Example 2: Password-protected share ===")
	createPasswordShare(ctx, client, 1) // Replace with your file ID

	// Example 3: Create a time-limited share
	fmt.Println("\n=== Example 3: Time-limited share ===")
	createExpiringShare(ctx, client, 1) // Replace with your file ID

	// Example 4: Create a share with view limit
	fmt.Println("\n=== Example 4: View-limited share ===")
	createViewLimitedShare(ctx, client, 1) // Replace with your file ID

	// Example 5: Share an album
	fmt.Println("\n=== Example 5: Share album ===")
	shareAlbum(ctx, client, 1) // Replace with your album ID

	// Example 6: List all shares
	fmt.Println("\n=== Example 6: List shares ===")
	listShares(ctx, client)

	// Example 7: Update a share
	if share != nil {
		fmt.Println("\n=== Example 7: Update share ===")
		updateShare(ctx, client, share.ID)
	}

	// Example 8: Access a shared item
	if share != nil {
		fmt.Println("\n=== Example 8: Access shared content ===")
		accessShare(ctx, client, share.Token)
	}

	// Example 9: Delete a share
	if share != nil {
		fmt.Println("\n=== Example 9: Delete share ===")
		deleteShare(ctx, client, share.ID)
	}
}

// createSimpleShare creates a basic share link for a file.
func createSimpleShare(ctx context.Context, client *fimage.Client, fileID int64) (*fimage.ShareLink, error) {
	share, err := client.Share.Create(ctx, fimage.ShareFile(fileID))
	if err != nil {
		return nil, fmt.Errorf("failed to create share: %w", err)
	}

	fmt.Printf("Share created:\n")
	fmt.Printf("  ID: %d\n", share.ID)
	fmt.Printf("  Token: %s\n", share.Token)
	fmt.Printf("  URL: %s\n", share.ShareURL)
	fmt.Printf("  Is Active: %v\n", share.IsActive)

	return share, nil
}

// createPasswordShare creates a password-protected share link.
func createPasswordShare(ctx context.Context, client *fimage.Client, fileID int64) {
	share, err := client.Share.Create(ctx,
		fimage.ShareFile(fileID).WithPassword("secret123"),
	)
	if err != nil {
		log.Printf("Error creating password share: %v\n", err)
		return
	}

	fmt.Printf("Password-protected share created:\n")
	fmt.Printf("  URL: %s\n", share.ShareURL)
	fmt.Printf("  Has Password: %v\n", share.HasPassword)
}

// createExpiringShare creates a share link that expires after 24 hours.
func createExpiringShare(ctx context.Context, client *fimage.Client, fileID int64) {
	share, err := client.Share.Create(ctx,
		fimage.ShareFile(fileID).WithExpiration(24), // 24 hours
	)
	if err != nil {
		log.Printf("Error creating expiring share: %v\n", err)
		return
	}

	fmt.Printf("Time-limited share created:\n")
	fmt.Printf("  URL: %s\n", share.ShareURL)
	if share.ExpiresAt != nil {
		fmt.Printf("  Expires At: %s\n", share.ExpiresAt.Format("2006-01-02 15:04:05"))
	}
}

// createViewLimitedShare creates a share link with a view limit.
func createViewLimitedShare(ctx context.Context, client *fimage.Client, fileID int64) {
	share, err := client.Share.Create(ctx,
		fimage.ShareFile(fileID).WithMaxViews(10),
	)
	if err != nil {
		log.Printf("Error creating view-limited share: %v\n", err)
		return
	}

	fmt.Printf("View-limited share created:\n")
	fmt.Printf("  URL: %s\n", share.ShareURL)
	if share.MaxViews != nil {
		fmt.Printf("  Max Views: %d\n", *share.MaxViews)
	}
	fmt.Printf("  Current Views: %d\n", share.ViewCount)
}

// shareAlbum creates a share link for an entire album.
func shareAlbum(ctx context.Context, client *fimage.Client, albumID int64) {
	share, err := client.Share.Create(ctx, fimage.ShareAlbum(albumID))
	if err != nil {
		log.Printf("Error sharing album: %v\n", err)
		return
	}

	fmt.Printf("Album share created:\n")
	fmt.Printf("  URL: %s\n", share.ShareURL)
	if share.AlbumName != nil {
		fmt.Printf("  Album: %s\n", *share.AlbumName)
	}
}

// listShares lists all share links.
func listShares(ctx context.Context, client *fimage.Client) {
	resp, err := client.Share.List(ctx, nil)
	if err != nil {
		log.Printf("Error listing shares: %v\n", err)
		return
	}

	fmt.Printf("Found %d shares:\n", resp.Total)
	for _, share := range resp.Shares {
		name := ""
		if share.FileName != nil {
			name = *share.FileName
		} else if share.AlbumName != nil {
			name = *share.AlbumName + " (album)"
		}

		status := "active"
		if !share.IsActive {
			status = "inactive"
		}

		fmt.Printf("  [%d] %s - %s (%s, %d views)\n",
			share.ID, name, share.ShareURL, status, share.ViewCount)
	}
}

// updateShare updates a share link.
func updateShare(ctx context.Context, client *fimage.Client, shareID int64) {
	isActive := false
	share, err := client.Share.Update(ctx, shareID, &fimage.UpdateShareOptions{
		IsActive: &isActive,
	})
	if err != nil {
		log.Printf("Error updating share: %v\n", err)
		return
	}

	fmt.Printf("Share updated:\n")
	fmt.Printf("  Is Active: %v\n", share.IsActive)
}

// accessShare accesses a shared item.
func accessShare(ctx context.Context, client *fimage.Client, token string) {
	content, err := client.Share.Access(ctx, token)
	if err != nil {
		log.Printf("Error accessing share: %v\n", err)
		return
	}

	if content.RequiresPassword {
		fmt.Println("This share requires a password. Use VerifyPassword().")
		return
	}

	fmt.Printf("Shared content:\n")
	fmt.Printf("  Type: %s\n", content.Type)

	if content.File != nil {
		fmt.Printf("  File: %s\n", content.File.OriginalName)
		fmt.Printf("  URL: %s\n", content.File.URL)
	}

	if content.Album != nil {
		fmt.Printf("  Album: %s\n", content.Album.Name)
		fmt.Printf("  Files: %d\n", len(content.Files))
	}
}

// deleteShare deletes a share link.
func deleteShare(ctx context.Context, client *fimage.Client, shareID int64) {
	resp, err := client.Share.Delete(ctx, shareID)
	if err != nil {
		log.Printf("Error deleting share: %v\n", err)
		return
	}

	fmt.Printf("Share deleted: %s\n", resp.Message)
}
