// Package fimage provides the official Go SDK for the F-Image API.
//
// F-Image (https://f-image.com) is a professional image hosting platform
// with features like deduplication, on-the-fly image processing, multi-size
// thumbnails, and a powerful API.
//
// # Quick Start
//
//	client := fimage.NewClient("your-api-token")
//
//	// Upload an image
//	file, _ := os.Open("photo.jpg")
//	resp, _ := client.Files.Upload(ctx, file, &fimage.UploadOptions{
//	    Filename:    "photo.jpg",
//	    Description: "My photo",
//	})
//	fmt.Println(resp.Data.URL)
//
// # Authentication
//
// All API requests require an API token. You can obtain one from your
// F-Image dashboard at https://f-image.com/dashboard/settings/api
//
// # Available Services
//
//   - Files: Upload, list, search, and manage images
//   - Albums: Organize images into albums
//   - Share: Create and manage share links
//   - Tags: Tag and categorize images
//   - Trash: Manage deleted files
package fimage
