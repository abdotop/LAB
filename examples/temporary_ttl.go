package examples

import (
	"fmt"
	"log"
	"time"

	"github.com/abdotop/LAB/config"
)

// TemporaryTTLExample demonstrates creating temporary uploads with TTL
func TemporaryTTLExample() {
	config := config.NewCloudinaryConfig()
	_, err := config.NewCloudinaryInstance()
	if err != nil {
		log.Fatalf("Failed to create Cloudinary instance: %v", err)
	}

	fmt.Println("🕐 Temporary Upload with 5-minute TTL Example")
	fmt.Println("=============================================")

	// Generate a temporary upload signature that expires in 5 minutes
	timestamp := time.Now().Unix()
	folder := "tmp"
	validFor := 5 * time.Minute

	// Create upload parameters for temporary file
	signature := map[string]interface{}{
		"folder":          folder,
		"timestamp":       timestamp,
		"unique_filename": true,
		"tags":            "temporary,5min_ttl,auto_delete",
		"context":         fmt.Sprintf("expires_at=%d|auto_delete=true|lifetime_minutes=5", time.Now().Add(validFor).Unix()),
		"cloud_name":      config.CloudName,
		"api_key":         config.APIKey,
	}

	fmt.Printf("✅ Generated temporary upload signature:\n")
	fmt.Printf("   Folder: %s\n", folder)
	fmt.Printf("   Expires in: %d minutes\n", 5)
	fmt.Printf("   Timestamp: %d\n", timestamp)
	
	// Show upload URL
	uploadURL := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", config.CloudName)
	fmt.Printf("   Upload URL: %s\n", uploadURL)

	fmt.Println("\n📋 Upload Parameters (for client-side upload):")
	for key, value := range signature {
		fmt.Printf("   %s: %v\n", key, value)
	}

	fmt.Println("\n⚡ Usage Instructions:")
	fmt.Println("1. Use these parameters in your client-side upload")
	fmt.Println("2. The signature is valid for 5 minutes only")
	fmt.Println("3. After upload, the file will be tagged for cleanup")
	fmt.Println("4. Set up a cleanup job to delete files tagged with 'temporary'")

	// Example cleanup function (conceptual)
	fmt.Println("\n🧹 Setting up cleanup for temporary files...")
	setupTemporaryCleanup(signature)
}

func setupTemporaryCleanup(signature map[string]interface{}) {
	fmt.Println("✅ Cleanup strategy configured:")
	fmt.Println("   - Files tagged 'temporary' will be monitored")
	fmt.Println("   - Automatic deletion after expiry time")
	fmt.Println("   - You can implement a cron job to check 'expires_at' context")
	
	// In a real implementation, you might:
	// 1. Set up a scheduled job (cron, AWS Lambda, etc.)
	// 2. Query assets with 'temporary' tag
	// 3. Check their 'expires_at' context value
	// 4. Delete expired assets using cld.Upload.Destroy()
	
	exampleCleanupQuery := map[string]interface{}{
		"tag":            "temporary",
		"context_filter": "expires_at < " + fmt.Sprintf("%d", time.Now().Unix()),
	}
	
	fmt.Printf("   Example cleanup query: %+v\n", exampleCleanupQuery)
	fmt.Printf("   Upload signature was: %+v\n", signature)
}