// Example: Manage tags with F-Image
//
// This example demonstrates how to create, use, and manage tags
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

	// Example 1: Create tags
	fmt.Println("=== Example 1: Create tags ===")
	natureTag, err := createTag(ctx, client, "Nature", "#4CAF50")
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	travelTag, _ := createTag(ctx, client, "Travel", "#2196F3")
	_ = travelTag

	// Example 2: List all tags
	fmt.Println("\n=== Example 2: List tags ===")
	listTags(ctx, client)

	// Example 3: Tag a file
	fmt.Println("\n=== Example 3: Tag a file ===")
	tagFile(ctx, client, 1, natureTag.ID) // Replace 1 with your file ID

	// Example 4: Get files by tag
	fmt.Println("\n=== Example 4: Get files by tag ===")
	getFilesByTag(ctx, client, natureTag.ID)

	// Example 5: Update a tag
	fmt.Println("\n=== Example 5: Update tag ===")
	updateTag(ctx, client, natureTag.ID)

	// Example 6: Untag a file
	fmt.Println("\n=== Example 6: Untag file ===")
	untagFile(ctx, client, 1, natureTag.ID) // Replace 1 with your file ID

	// Example 7: Delete a tag
	fmt.Println("\n=== Example 7: Delete tag ===")
	deleteTag(ctx, client, natureTag.ID)
}

// createTag creates a new tag with a color.
func createTag(ctx context.Context, client *fimage.Client, name, color string) (*fimage.Tag, error) {
	tag, err := client.Tags.Create(ctx, &fimage.CreateTagOptions{
		Name:  name,
		Color: color,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	fmt.Printf("Tag created:\n")
	fmt.Printf("  ID: %d\n", tag.ID)
	fmt.Printf("  Name: %s\n", tag.Name)
	fmt.Printf("  Color: %s\n", tag.Color)

	return tag, nil
}

// listTags lists all tags with their file counts.
func listTags(ctx context.Context, client *fimage.Client) {
	tags, err := client.Tags.List(ctx)
	if err != nil {
		log.Printf("Error listing tags: %v\n", err)
		return
	}

	fmt.Printf("Found %d tags:\n", len(tags))
	for _, tag := range tags {
		fmt.Printf("  [%d] %s (color: %s, files: %d)\n",
			tag.ID, tag.Name, tag.Color, tag.FileCount)
	}
}

// tagFile adds a tag to a file.
func tagFile(ctx context.Context, client *fimage.Client, fileID, tagID int64) {
	resp, err := client.Tags.TagFile(ctx, fileID, tagID)
	if err != nil {
		if fimage.IsNotFound(err) {
			log.Printf("File or tag not found\n")
			return
		}
		log.Printf("Error tagging file: %v\n", err)
		return
	}

	fmt.Printf("File tagged: %s\n", resp.Message)
}

// untagFile removes a tag from a file.
func untagFile(ctx context.Context, client *fimage.Client, fileID, tagID int64) {
	resp, err := client.Tags.UntagFile(ctx, fileID, tagID)
	if err != nil {
		log.Printf("Error untagging file: %v\n", err)
		return
	}

	fmt.Printf("File untagged: %s\n", resp.Message)
}

// getFilesByTag gets all files with a specific tag.
func getFilesByTag(ctx context.Context, client *fimage.Client, tagID int64) {
	resp, err := client.Tags.GetFiles(ctx, tagID, &fimage.TagFilesOptions{
		Page:  1,
		Limit: 20,
	})
	if err != nil {
		log.Printf("Error getting files: %v\n", err)
		return
	}

	fmt.Printf("Found %d files with this tag:\n", resp.Total)
	for _, file := range resp.Files {
		fmt.Printf("  [%d] %s\n", file.ID, file.OriginalName)
	}
}

// updateTag updates a tag's name and color.
func updateTag(ctx context.Context, client *fimage.Client, tagID int64) {
	tag, err := client.Tags.Update(ctx, tagID, &fimage.UpdateTagOptions{
		Name:  "Wildlife",
		Color: "#FF9800",
	})
	if err != nil {
		log.Printf("Error updating tag: %v\n", err)
		return
	}

	fmt.Printf("Tag updated:\n")
	fmt.Printf("  ID: %d\n", tag.ID)
	fmt.Printf("  Name: %s\n", tag.Name)
	fmt.Printf("  Color: %s\n", tag.Color)
}

// deleteTag deletes a tag (removes it from all files).
func deleteTag(ctx context.Context, client *fimage.Client, tagID int64) {
	resp, err := client.Tags.Delete(ctx, tagID)
	if err != nil {
		if fimage.IsNotFound(err) {
			log.Printf("Tag not found\n")
			return
		}
		log.Printf("Error deleting tag: %v\n", err)
		return
	}

	fmt.Printf("Tag deleted: %s\n", resp.Message)
}
