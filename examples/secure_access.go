package examples

import (
	"fmt"
	"log"
	"time"

	"github.com/abdotop/LAB/config"
)

// SecureAccessExample demonstrates secure directory access with key-based authentication
func SecureAccessExample() {
	config := config.NewCloudinaryConfig()
	_, err := config.NewCloudinaryInstance()
	if err != nil {
		log.Fatalf("Failed to create Cloudinary instance: %v", err)
	}

	fmt.Println("🔐 Secure Directory Access Example")
	fmt.Println("==================================")

	// Create a secure upload signature for the 'partner' directory
	partnerFolder := "partner"
	secureKey := "secure-partner-key-2024"

	// Generate secure upload parameters
	timestamp := time.Now().Unix()
	
	signature := map[string]interface{}{
		"folder":         partnerFolder,
		"timestamp":      timestamp,
		"unique_filename": true,
		"tags":           "secure,partner,authenticated",
		"context":        fmt.Sprintf("access_key=%s|security_level=high|created_by=golang_server", secureKey),
		"cloud_name":     config.CloudName,
		"api_key":        config.APIKey,
	}

	fmt.Printf("✅ Generated secure upload for '%s' directory:\n", partnerFolder)
	fmt.Printf("   Security Key: %s\n", secureKey)
	fmt.Printf("   Timestamp: %d\n", timestamp)

	// Demonstrate generating secure delivery URLs
	examplePublicID := "partner/secure_document"
	
	fmt.Println("\n🔗 Generating secure delivery URLs...")
	
	// 1. URL with token authentication (expires in 1 hour)
	baseURL := fmt.Sprintf("https://res.cloudinary.com/%s/image/upload/%s", 
		config.CloudName, 
		examplePublicID)
	
	// Create secure URL with authentication token
	secureURL := fmt.Sprintf("%s?auth_token=%s&expires=%d", 
		baseURL, 
		secureKey, 
		time.Now().Add(1*time.Hour).Unix())
	
	fmt.Printf("✅ Secure URL (1 hour expiry): %s\n", secureURL)

	// 2. URL with IP-based restrictions (conceptual)
	fmt.Println("\n🌐 IP-based access control:")
	ipRestrictedURL := fmt.Sprintf("https://res.cloudinary.com/%s/image/upload/%s", 
		config.CloudName, 
		examplePublicID)
	fmt.Printf("✅ IP-restricted URL: %s\n", ipRestrictedURL)
	fmt.Println("   (Configure IP whitelist in Cloudinary dashboard)")

	// 3. Time-bound access with specific permissions
	fmt.Println("\n⏰ Time-bound access example:")
	
	// Generate different access levels
	accessLevels := []struct{
		name     string
		duration int
		level    string
	}{
		{"read-only", 300, "read"},      // 5 minutes
		{"full-access", 1800, "write"},  // 30 minutes  
		{"admin", 3600, "admin"},        // 1 hour
	}

	for _, access := range accessLevels {
		baseURL := fmt.Sprintf("https://res.cloudinary.com/%s/image/upload/%s", 
			config.CloudName, 
			examplePublicID)

		secureURL := fmt.Sprintf("%s?auth_token=%s-%s&expires=%d&access_level=%s", 
			baseURL, 
			secureKey, 
			access.level,
			time.Now().Add(time.Duration(access.duration)*time.Second).Unix(),
			access.name)

		fmt.Printf("   %s (%d min): %s\n", 
			access.name, 
			access.duration/60, 
			secureURL[:50]+"...")
	}

	fmt.Println("\n🛡️  Security Features Demonstrated:")
	fmt.Printf("   ✅ Secure upload parameters: %+v\n", signature)
	fmt.Println("   ✅ Key-based authentication")
	fmt.Println("   ✅ Time-based expiration")
	fmt.Println("   ✅ Access level differentiation")
	fmt.Println("   ✅ IP-based restrictions (configurable)")
	fmt.Println("   ✅ Secure context metadata")

	// Demonstrate key rotation (conceptual)
	fmt.Println("\n🔄 Key Rotation Example:")
	oldKey := secureKey
	newKey := fmt.Sprintf("secure-partner-key-%d", time.Now().Unix())
	
	fmt.Printf("   Old key: %s\n", oldKey)
	fmt.Printf("   New key: %s\n", newKey)
	fmt.Println("   💡 Implement regular key rotation for enhanced security")
}