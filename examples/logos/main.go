// Example: Look up domain logos with F-Image
//
// This example demonstrates how to resolve a domain logo URL. F-Image returns
// a stored logo when present, otherwise it fetches and saves one.
//
// Usage:
//
//	export FIMAGE_API_TOKEN="your-api-token"
//	go run main.go marriott.com
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	fimage "github.com/lpg-it/f-image-go"
)

func main() {
	apiToken := os.Getenv("FIMAGE_API_TOKEN")
	if apiToken == "" {
		log.Fatal("FIMAGE_API_TOKEN environment variable is required")
	}

	domain := "marriott.com"
	if len(os.Args) > 1 {
		domain = os.Args[1]
	}

	client := fimage.NewClient(apiToken)

	logo, err := client.Logos.Resolve(context.Background(), domain)
	if err != nil {
		log.Fatalf("logo resolve failed: %v", err)
	}

	fmt.Printf("Domain: %s\n", logo.Domain)
	if logo.URL == "" {
		fmt.Println("Logo not found")
		return
	}

	fmt.Printf("Logo URL: %s\n", logo.URL)
	if logo.Source != "" {
		fmt.Printf("Source: %s\n", logo.Source)
	}
	if logo.Provider != "" {
		fmt.Printf("Provider: %s\n", logo.Provider)
	}
}
