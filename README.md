<p align="center">
  <a href="https://f-image.com" target="_blank">
    <img src="https://f-image.com/logo.png" width="120" alt="F-Image Logo">
  </a>
</p>

<h1 align="center">F-Image Go SDK</h1>

<p align="center">
  <strong>The official Go SDK for the <a href="https://f-image.com">F-Image</a> image hosting platform.</strong>
  <br />
  Fast, reliable image hosting with CDN, on-the-fly processing, and deduplication.
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/lpg-it/f-image-go"><img src="https://pkg.go.dev/badge/github.com/lpg-it/f-image-go.svg" alt="Go Reference"></a>
  <a href="https://github.com/lpg-it/f-image-go/releases"><img src="https://img.shields.io/github/v/release/lpg-it/f-image-go" alt="Release"></a>
  <a href="https://goreportcard.com/report/github.com/lpg-it/f-image-go"><img src="https://goreportcard.com/badge/github.com/lpg-it/f-image-go" alt="Go Report Card"></a>
  <a href="https://github.com/lpg-it/f-image-go/blob/main/LICENSE"><img src="https://img.shields.io/github/license/lpg-it/f-image-go" alt="License"></a>
</p>

---

## 🌟 What is F-Image?

**[F-Image](https://f-image.com)** is a professional **image hosting service** and **image CDN** designed for developers, bloggers, and businesses. It provides:

- ⚡ **Global CDN** - Lightning-fast image delivery worldwide
- 🔄 **Automatic Deduplication** - Save storage with intelligent file detection
- 📐 **On-the-fly Processing** - Auto-generate thumbnails and optimized sizes
- 🖼️ **Multi-format Support** - JPEG, PNG, GIF, WebP, BMP, and more
- 🔗 **Share Links** - Password-protected, time-limited, and view-limited sharing
- 📁 **Album Organization** - Organize images into albums
- 🏷️ **Tagging System** - Tag and categorize your images
- 🔌 **Powerful API** - Full-featured REST API for automation

**[Sign up for free](https://f-image.com/register)** to get your API token and start hosting images in minutes!

---

## 📦 Installation

```bash
go get github.com/lpg-it/f-image-go
```

**Requirements:** Go 1.21 or later

---

## 🚀 Quick Start

### Get Your API Token

1. Visit [F-Image Dashboard](https://f-image.com/dashboard/settings/api)
2. Create a new API token
3. Copy the token (starts with `fimg_live_` or `fimg_test_`)

### Upload Your First Image

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    fimage "github.com/lpg-it/f-image-go"
)

func main() {
    // Initialize client with your API token
    client := fimage.NewClient(os.Getenv("FIMAGE_API_TOKEN"))

    // Open image file
    file, err := os.Open("photo.jpg")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    // Upload image
    resp, err := client.Files.Upload(context.Background(), file, &fimage.UploadOptions{
        Filename:    "my-photo.jpg",
        Description: "A beautiful sunset",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Print the CDN URL
    fmt.Printf("Image URL: %s\n", resp.Data.URL)
    fmt.Printf("Thumbnail: %s\n", *resp.Data.ThumbnailURL)
}
```

---

## 📚 API Reference

### Client Initialization

```go
// Basic initialization
client := fimage.NewClient("your-api-token")

// With custom options
client := fimage.NewClient("your-api-token",
    fimage.WithTimeout(60*time.Second),
    fimage.WithBaseURL("https://custom-domain.example.com"),
)
```

---

### 📤 Files API

Upload, manage, and organize your images.

#### Upload Image

```go
// Upload from file
file, _ := os.Open("photo.jpg")
resp, err := client.Files.Upload(ctx, file, &fimage.UploadOptions{
    Filename:    "photo.jpg",
    Description: "My photo description",
})
fmt.Println(resp.Data.URL)    // https://i.f-image.com/images/abc123.jpg
fmt.Println(resp.Data.IsFlash) // true if deduplicated

// Upload from URL
resp, err := client.Files.UploadFromURL(ctx, "https://example.com/image.jpg")
```

#### List Files

```go
// List all files
resp, err := client.Files.List(ctx, nil)
for _, file := range resp.Files {
    fmt.Printf("[%d] %s - %s\n", file.ID, file.OriginalName, file.URL)
}

// With pagination
resp, err := client.Files.List(ctx, &fimage.ListOptions{
    Page:  2,
    Limit: 50,
})

// Filter by album
albumID := int64(123)
resp, err := client.Files.List(ctx, &fimage.ListOptions{
    AlbumID: &albumID,
})
```

#### Search Files

```go
resp, err := client.Files.Search(ctx, &fimage.SearchOptions{
    Query: "sunset beach",
    Page:  1,
    Limit: 20,
})
fmt.Printf("Found %d matching files\n", resp.Total)
```

#### Delete Files

```go
// Single file (moves to trash)
_, err := client.Files.Delete(ctx, 123)

// Batch delete
resp, err := client.Files.BatchDelete(ctx, []int64{1, 2, 3})
fmt.Printf("Deleted: %d, Failed: %d\n", resp.Deleted, resp.Failed)
```

#### Move Files to Album

```go
// Move single file
albumID := int64(456)
_, err := client.Files.Move(ctx, 123, &albumID)

// Move multiple files
_, err := client.Files.MoveMany(ctx, []int64{1, 2, 3}, &albumID)

// Remove from album
_, err := client.Files.Move(ctx, 123, nil)
```

---

### 📁 Albums API

Organize your images into albums.

#### Create Album

```go
album, err := client.Albums.Create(ctx, &fimage.CreateAlbumOptions{
    Name:        "Vacation 2024",
    Description: "Summer vacation photos",
})
fmt.Printf("Created album: %s (ID: %d)\n", album.Name, album.ID)
```

#### List Albums

```go
albums, err := client.Albums.List(ctx)
for _, album := range albums {
    fmt.Printf("%s - %d files\n", album.Name, album.FileCount)
}
```

#### Get Album

```go
album, err := client.Albums.Get(ctx, 123)
if fimage.IsNotFound(err) {
    fmt.Println("Album not found")
}
```

#### Update Album

```go
album, err := client.Albums.Update(ctx, 123, &fimage.UpdateAlbumOptions{
    Name:        "Summer Vacation 2024",
    Description: "Updated description",
})
```

#### Delete Album

```go
_, err := client.Albums.Delete(ctx, 123)
// Note: Files in the album are not deleted, they become "unalbumed"
```

---

### 🔗 Share API

Create shareable links for files and albums.

#### Create Share Link

```go
// Share a file
fileID := int64(123)
share, err := client.Share.Create(ctx, fimage.ShareFile(fileID))
fmt.Printf("Share URL: %s\n", share.ShareURL)

// Share with password protection
share, err := client.Share.Create(ctx,
    fimage.ShareFile(fileID).WithPassword("secret123"),
)

// Share with expiration (24 hours)
share, err := client.Share.Create(ctx,
    fimage.ShareFile(fileID).WithExpiration(24),
)

// Share with view limit
share, err := client.Share.Create(ctx,
    fimage.ShareFile(fileID).WithMaxViews(100),
)

// Share an album
albumID := int64(456)
share, err := client.Share.Create(ctx, fimage.ShareAlbum(albumID))
```

#### List Share Links

```go
resp, err := client.Share.List(ctx, nil)
for _, share := range resp.Shares {
    fmt.Printf("%s - %d views\n", share.ShareURL, share.ViewCount)
}
```

#### Update Share Link

```go
isActive := false
_, err := client.Share.Update(ctx, 123, &fimage.UpdateShareOptions{
    IsActive: &isActive,
})
```

#### Access Shared Content

```go
content, err := client.Share.Access(ctx, "abc123token")
if content.RequiresPassword {
    content, err = client.Share.VerifyPassword(ctx, "abc123token", "secret123")
}
```

#### Delete Share Link

```go
_, err := client.Share.Delete(ctx, 123)
```

---

### 🏷️ Tags API

Tag and categorize your images.

#### Create Tag

```go
tag, err := client.Tags.Create(ctx, &fimage.CreateTagOptions{
    Name:  "Nature",
    Color: "#4CAF50",
})
```

#### List Tags

```go
tags, err := client.Tags.List(ctx)
for _, tag := range tags {
    fmt.Printf("%s (%s) - %d files\n", tag.Name, tag.Color, tag.FileCount)
}
```

#### Tag/Untag Files

```go
// Add tag to file
_, err := client.Tags.TagFile(ctx, 123, tagID)

// Remove tag from file
_, err := client.Tags.UntagFile(ctx, 123, tagID)
```

#### Get Files by Tag

```go
resp, err := client.Tags.GetFiles(ctx, tagID, nil)
for _, file := range resp.Files {
    fmt.Println(file.OriginalName)
}
```

#### Update Tag

```go
_, err := client.Tags.Update(ctx, tagID, &fimage.UpdateTagOptions{
    Name:  "Wildlife",
    Color: "#FF9800",
})
```

#### Delete Tag

```go
_, err := client.Tags.Delete(ctx, tagID)
```

---

### 🗑️ Trash API

Manage deleted files (soft delete with 30-day retention).

#### List Trash

```go
resp, err := client.Trash.List(ctx, nil)
for _, file := range resp.Files {
    fmt.Printf("%s (deleted: %s)\n", file.OriginalName, *file.DeletedAt)
}
```

#### Restore Files

```go
// Single file
_, err := client.Trash.Restore(ctx, 123)

// Multiple files
resp, err := client.Trash.RestoreMany(ctx, []int64{1, 2, 3})
fmt.Printf("Restored: %d, Failed: %d\n", resp.Restored, resp.Failed)
```

#### Permanent Delete

```go
// Single file (CANNOT BE UNDONE!)
result, err := client.Trash.PermanentDelete(ctx, 123)

// Empty entire trash (CANNOT BE UNDONE!)
result, err := client.Trash.Empty(ctx)
fmt.Printf("Deleted: %d files\n", result.DeletedCount)
```

---

## 🛡️ Error Handling

The SDK provides typed errors for common scenarios:

```go
resp, err := client.Files.Upload(ctx, file, nil)
if err != nil {
    switch {
    case fimage.IsUnauthorized(err):
        fmt.Println("Invalid API token")
    case fimage.IsNotFound(err):
        fmt.Println("Resource not found")
    case fimage.IsQuotaExceeded(err):
        fmt.Println("Storage quota exceeded")
    case fimage.IsForbidden(err):
        fmt.Println("Access denied")
    default:
        fmt.Printf("Error: %v\n", err)
    }
}
```

---

## 📋 Response Types

### UploadResponse

```go
type UploadResponse struct {
    Success bool        `json:"success"`
    Status  int         `json:"status"`
    Data    *UploadData `json:"data"`
}

type UploadData struct {
    ID           int64   `json:"id"`
    URL          string  `json:"url"`
    MediumURL    *string `json:"medium_url,omitempty"`
    ThumbnailURL *string `json:"thumbnail_url,omitempty"`
    OriginalName string  `json:"original_name"`
    Description  string  `json:"description"`
    Size         int64   `json:"size"`
    Width        int     `json:"width"`
    Height       int     `json:"height"`
    MimeType     string  `json:"mime_type"`
    IsFlash      bool    `json:"is_flash"`
}
```

### File

```go
type File struct {
    ID           int64   `json:"id"`
    AlbumID      *int64  `json:"album_id,omitempty"`
    AlbumName    *string `json:"album_name,omitempty"`
    OriginalName string  `json:"original_name"`
    Description  string  `json:"description"`
    URL          string  `json:"url"`
    MediumURL    *string `json:"medium_url,omitempty"`
    ThumbnailURL *string `json:"thumbnail_url,omitempty"`
    Size         int64   `json:"size"`
    Width        int     `json:"width"`
    Height       int     `json:"height"`
    MimeType     string  `json:"mime_type"`
    CreatedAt    string  `json:"created_at"`
}
```

### Album

```go
type Album struct {
    ID          int64  `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    FileCount   int64  `json:"file_count"`
    CreatedAt   string `json:"created_at"`
}
```

### ShareLink

```go
type ShareLink struct {
    ID          int64      `json:"id"`
    Token       string     `json:"token"`
    ShareURL    string     `json:"share_url"`
    FileID      *int64     `json:"file_id,omitempty"`
    AlbumID     *int64     `json:"album_id,omitempty"`
    FileName    *string    `json:"file_name,omitempty"`
    AlbumName   *string    `json:"album_name,omitempty"`
    HasPassword bool       `json:"has_password"`
    ExpiresAt   *time.Time `json:"expires_at,omitempty"`
    MaxViews    *int64     `json:"max_views,omitempty"`
    ViewCount   int64      `json:"view_count"`
    IsActive    bool       `json:"is_active"`
    CreatedAt   time.Time  `json:"created_at"`
}
```

### Tag

```go
type Tag struct {
    ID        int64  `json:"id"`
    Name      string `json:"name"`
    Color     string `json:"color"`
    FileCount int64  `json:"file_count"`
}
```

---

## 💡 Examples

The `examples/` directory contains complete, runnable examples:

| Example | Description |
|---------|-------------|
| [upload](examples/upload) | Upload images from file, bytes, or URL |
| [files](examples/files) | List, search, delete, and move files |
| [albums](examples/albums) | Album CRUD operations |
| [share](examples/share) | Create and manage share links |
| [tags](examples/tags) | Tag management and file tagging |
| [trash](examples/trash) | Trash management and restoration |

Run an example:

```bash
export FIMAGE_API_TOKEN="your-api-token"
go run examples/upload/main.go
```

---

## 🔧 Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithBaseURL(url)` | Set custom API base URL | `https://f-image.com` |
| `WithTimeout(duration)` | Set HTTP client timeout | `30s` |
| `WithHTTPClient(client)` | Use custom HTTP client | Default client |
| `WithUserAgent(ua)` | Set custom User-Agent header | `f-image-go/1.0.0` |

---

## 🌐 Why Choose F-Image?

| Feature | F-Image | Imgur | Cloudinary |
|---------|---------|-------|------------|
| **Global CDN** | ✅ | ✅ | ✅ |
| **Free Tier** | ✅ Generous | ❌ Limited | ✅ Limited |
| **API Rate Limits** | ✅ High | ❌ Restrictive | ✅ Based on plan |
| **Deduplication** | ✅ Automatic | ❌ | ❌ |
| **Password-Protected Shares** | ✅ | ❌ | ❌ |
| **Time-Limited Shares** | ✅ | ❌ | ✅ |
| **Go SDK** | ✅ Official | ❌ | ✅ |
| **No Watermarks** | ✅ | ❌ | ✅ |

---

## 📖 Documentation

- **[F-Image Website](https://f-image.com)** - Learn more about F-Image
- **[API Documentation](https://f-image.com/docs)** - Full API reference
- **[Go Reference](https://pkg.go.dev/github.com/lpg-it/f-image-go)** - Go package documentation

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## ⚖️ License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<p align="center">
  <a href="https://f-image.com"><strong>🌐 F-Image Website</strong></a> •
  <a href="https://f-image.com/register"><strong>🚀 Get Started Free</strong></a> •
  <a href="https://f-image.com/docs"><strong>📖 API Docs</strong></a>
</p>

<p align="center">
  Made with ❤️ by the <a href="https://f-image.com">F-Image</a> team
</p>

---

## 🔍 SEO Keywords

Image hosting API, Go image upload SDK, image CDN Go, F-Image Go client, image hosting service API, upload images Go, image management API, Go SDK for images, API image hosting, image CDN API, cloud image storage Go, image deduplication API, thumbnail generation API, image sharing API Go, picture hosting API, photo hosting Go SDK, image upload library Go, CDN for images Go, scalable image hosting, developer image hosting API