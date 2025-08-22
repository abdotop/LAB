package main

import (
	"fmt"
	"os"

	"github.com/abdotop/LAB/examples"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	switch os.Args[1] {
	case "demo":
		RunMainDemo()
	case "move":
		examples.MoveUploadExample()
	case "ttl":
		examples.TemporaryTTLExample()
	case "secure":
		examples.SecureAccessExample()
	case "all":
		fmt.Println("Running all examples...\n")
		
		fmt.Println("=== Main Demo ===")
		RunMainDemo()
		
		fmt.Println("\n=== Move Upload Example ===")
		examples.MoveUploadExample()
		
		fmt.Println("\n=== Temporary TTL Example ===")
		examples.TemporaryTTLExample()
		
		fmt.Println("\n=== Secure Access Example ===")
		examples.SecureAccessExample()
	default:
		fmt.Printf("Unknown example: %s\n", os.Args[1])
		showHelp()
	}
}

func showHelp() {
	fmt.Println("Cloudinary Examples Runner")
	fmt.Println("=========================")
	fmt.Println("Usage: go run . <example>")
	fmt.Println("")
	fmt.Println("Available examples:")
	fmt.Println("  demo   - Main comprehensive demo")
	fmt.Println("  move   - Upload movement between directories")
	fmt.Println("  ttl    - Temporary uploads with TTL")
	fmt.Println("  secure - Secure directory access")
	fmt.Println("  all    - Run all examples")
	fmt.Println("")
	fmt.Println("Note: Set your Cloudinary credentials in environment variables:")
	fmt.Println("  export CLOUDINARY_CLOUD_NAME=your_cloud_name")
	fmt.Println("  export CLOUDINARY_API_KEY=your_api_key")
	fmt.Println("  export CLOUDINARY_API_SECRET=your_api_secret")
}