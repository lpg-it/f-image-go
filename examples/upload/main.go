// Example: Upload images to F-Image
//
// This example demonstrates how to upload images using the F-Image Go SDK.
// It covers both file uploads and URL uploads.
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
	"strings"
	"time"

	fimage "github.com/lpg-it/f-image-go"
)

func main() {
	// Get API token from environment variable
	apiToken := os.Getenv("FIMAGE_API_TOKEN")
	if apiToken == "" {
		log.Fatal("FIMAGE_API_TOKEN environment variable is required")
	}

	// Create a new client
	client := fimage.NewClient(apiToken,
		fimage.WithTimeout(60*time.Second), // Increase timeout for large uploads
	)

	ctx := context.Background()

	// Example 1: Upload a file from disk
	fmt.Println("=== Example 1: Upload from file ===")
	if err := uploadFromFile(ctx, client); err != nil {
		log.Printf("Error: %v\n", err)
	}

	// Example 2: Upload from an io.Reader (e.g., bytes buffer)
	fmt.Println("\n=== Example 2: Upload from bytes ===")
	if err := uploadFromBytes(ctx, client); err != nil {
		log.Printf("Error: %v\n", err)
	}

	// Example 3: Upload from a public URL
	fmt.Println("\n=== Example 3: Upload from URL ===")
	if err := uploadFromURL(ctx, client); err != nil {
		log.Printf("Error: %v\n", err)
	}
}

// uploadFromFile demonstrates uploading a file from disk.
func uploadFromFile(ctx context.Context, client *fimage.Client) error {
	// Open the file
	file, err := os.Open("example.jpg")
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Upload with options
	resp, err := client.Files.Upload(ctx, file, &fimage.UploadOptions{
		Filename:    "my-photo.jpg",
		Description: "A beautiful landscape photo",
	})
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}

	// Print the result
	fmt.Printf("Upload successful!\n")
	fmt.Printf("  ID: %d\n", resp.Data.ID)
	fmt.Printf("  URL: %s\n", resp.Data.URL)
	fmt.Printf("  Size: %d bytes\n", resp.Data.Size)
	fmt.Printf("  Dimensions: %dx%d\n", resp.Data.Width, resp.Data.Height)
	fmt.Printf("  MIME Type: %s\n", resp.Data.MimeType)
	fmt.Printf("  Is Flash Upload: %v\n", resp.Data.IsFlash)

	if resp.Data.ThumbnailURL != nil {
		fmt.Printf("  Thumbnail: %s\n", *resp.Data.ThumbnailURL)
	}
	if resp.Data.MediumURL != nil {
		fmt.Printf("  Medium: %s\n", *resp.Data.MediumURL)
	}

	return nil
}

// uploadFromBytes demonstrates uploading from a bytes buffer.
func uploadFromBytes(ctx context.Context, client *fimage.Client) error {
	// Create some sample image data (in real usage, this would be actual image data)
	imageData := []byte{0xFF, 0xD8, 0xFF, 0xE0} // JPEG magic bytes (minimal example)

	// Convert to io.Reader
	reader := strings.NewReader(string(imageData))

	// Upload
	resp, err := client.Files.Upload(ctx, reader, &fimage.UploadOptions{
		Filename:    "generated-image.jpg",
		Description: "Programmatically generated image",
	})
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}

	fmt.Printf("Upload successful! URL: %s\n", resp.Data.URL)
	return nil
}

// uploadFromURL demonstrates uploading from a public URL.
func uploadFromURL(ctx context.Context, client *fimage.Client) error {
	// Upload from a public image URL
	imageURL := "https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=1200"

	resp, err := client.Files.UploadFromURL(ctx, imageURL)
	if err != nil {
		return fmt.Errorf("upload from URL failed: %w", err)
	}

	fmt.Printf("Upload from URL successful!\n")
	fmt.Printf("  ID: %d\n", resp.Data.ID)
	fmt.Printf("  URL: %s\n", resp.Data.URL)
	fmt.Printf("  Original Name: %s\n", resp.Data.OriginalName)
	fmt.Printf("  Dimensions: %dx%d\n", resp.Data.Width, resp.Data.Height)

	return nil
}
