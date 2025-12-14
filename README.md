<p align="center">
  <a href="https://f-image.com" target="_blank">
    <img src="https://i.f-image.com/your-logo.png" width="120" alt="F-Image Logo">
  </a>
</p>

<h1 align="center">F-Image Go SDK</h1>

<p align="center">
  <strong>The official Go client for the <a href="https://f-image.com">F-Image API</a>.</strong>
  <br />
  High-performance image hosting and asset management, now in your Go applications.
</p>

<p align="center">
    <a href="https://pkg.go.dev/github.com/your-github-username/f-image-go"><img src="https://pkg.go.dev/badge/github.com/your-github-username/f-image-go.svg" alt="Go Reference"></a>
    <a href="https://github.com/your-github-username/f-image-go/releases"><img src="https://img.shields.io/github/v/release/your-github-username/f-image-go" alt="Release"></a>
    <a href="https://github.com/your-github-username/f-image-go/actions/workflows/test.yml"><img src="https://github.com/your-github-username/f-image-go/actions/workflows/test.yml/badge.svg" alt="Tests"></a>
    <a href="https://goreportcard.com/report/github.com/your-github-username/f-image-go"><img src="https://goreportcard.com/badge/github.com/your-github-username/f-image-go" alt="Go Report Card"></a>
</p>

---

**[f-image.com](https://f-image.com)** provides a fast, reliable, and feature-rich platform for your image hosting needs, including deduplication, on-the-fly image processing, and a powerful API. This Go SDK makes integrating F-Image into your Go projects effortless.

**Don't have an account yet?** [**Sign up for free at f-image.com**](https://f-image.com/register) to get your API Token.

## ✨ Features

- **Fluent API**: Clean, modern, and easy-to-use client.
- **Multipart Uploads**: Efficiently upload files from disk or memory (`io.Reader`).
- **Asset Management**: List, delete, and manage your images and albums.
- **Typed Structs**: All API responses are strongly typed for safety.
- **Error Handling**: Clear, descriptive errors for easier debugging.
- **Context Aware**: All API calls support `context.Context` for cancellation and timeouts.

## 🚀 Installation

```bash
go get github.com/your-github-username/f-image-go
```

## 📚 Quick Start: Upload an Image

First, get your API Token from your **[F-Image Dashboard](https://f-image.com/dashboard/settings/api)**.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	
	fimage "github.com/your-github-username/f-image-go"
)

func main() {
	// Initialize the client with your API Token
	// It's recommended to load this from an environment variable.
	client := fimage.NewClient("fimg_live_YOUR_API_TOKEN_HERE")
	
	// Open the file you want to upload
	file, err := os.Open("path/to/your/image.jpg")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()
	
	// Prepare the upload options
	opts := &fimage.UploadOptions{
		Filename:    "my-cat.jpg",
		Description: "A beautiful cat enjoying the sun.",
		AlbumID:     123, // Optional
	}
	
	// Upload the file
	resp, err := client.Upload(context.Background(), file, opts)
	if err != nil {
		log.Fatalf("Upload failed: %v", err)
	}
	
	// Success!
	fmt.Println("Image uploaded successfully!")
	fmt.Printf("URL: %s\n", resp.Data.URL)
	fmt.Printf("Thumbnail URL: %s\n", resp.Data.ThumbnailURL)
}
```

## 📖 Documentation

For detailed information on all available methods and parameters, please see the **[Go Reference](https://pkg.go.dev/github.com/your-github-username/f-image-go)**.

For full API endpoint details, visit the **[F-Image API Docs](https://f-image.com/docs)**.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1.  Fork the repository.
2.  Create your feature branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the branch (`git push origin feature/AmazingFeature`).
5.  Open a Pull Request.

## ⚖️ License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.