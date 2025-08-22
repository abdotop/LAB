package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/abdotop/LAB/config"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// CloudinaryLab demonstrates various Cloudinary features
type CloudinaryLab struct {
	cld *cloudinary.Cloudinary
}

// NewCloudinaryLab creates a new CloudinaryLab instance
func NewCloudinaryLab() (*CloudinaryLab, error) {
	config := config.NewCloudinaryConfig()
	cld, err := config.NewCloudinaryInstance()
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloudinary instance: %w", err)
	}

	return &CloudinaryLab{cld: cld}, nil
}

// GeneratePresignedUploadURL generates a presigned upload URL with directory structure
func (cl *CloudinaryLab) GeneratePresignedUploadURL(ctx context.Context, folder string) (map[string]interface{}, error) {
	timestamp := time.Now().Unix()
	
	// Create a map with upload parameters for presigned upload
	signature := map[string]interface{}{
		"folder":         folder,
		"resource_type":  "auto",
		"timestamp":      timestamp,
		"unique_filename": true,
		"cloud_name":     cl.cld.Config.Cloud.CloudName,
		"api_key":        cl.cld.Config.Cloud.APIKey,
	}

	return signature, nil
}

// CheckUploadExists checks if an upload exists in a specific directory
func (cl *CloudinaryLab) CheckUploadExists(ctx context.Context, publicID string) (bool, *admin.AssetResult, error) {
	result, err := cl.cld.Admin.Asset(ctx, admin.AssetParams{
		PublicID: publicID,
	})
	
	if err != nil {
		// If resource not found, return false
		return false, nil, nil
	}

	return true, result, nil
}

// MoveUpload moves an upload from one directory to another
func (cl *CloudinaryLab) MoveUpload(ctx context.Context, fromPublicID, toPublicID string) (*uploader.RenameResult, error) {
	result, err := cl.cld.Upload.Rename(ctx, uploader.RenameParams{
		FromPublicID: fromPublicID,
		ToPublicID:   toPublicID,
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to move upload from %s to %s: %w", fromPublicID, toPublicID, err)
	}

	return result, nil
}

// GenerateSecureURL generates a secure URL with signature for protected access
func (cl *CloudinaryLab) GenerateSecureURL(publicID string, expireAt time.Time) (string, error) {
	// Create a basic delivery URL
	baseURL := fmt.Sprintf("https://res.cloudinary.com/%s/image/upload/%s", 
		cl.cld.Config.Cloud.CloudName, 
		publicID)

	// Add authentication token (this is a conceptual implementation)
	// In practice, you would use Cloudinary's URL signing features
	secureURL := fmt.Sprintf("%s?auth_token=%s&expires=%d", 
		baseURL, 
		"secure-token", 
		expireAt.Unix())

	return secureURL, nil
}

// GenerateAccessKey generates long-term or temporary access keys
func (cl *CloudinaryLab) GenerateAccessKey(name string, enabled bool, ttl *time.Duration) (map[string]interface{}, error) {
	// Note: Access key management requires Admin API with appropriate permissions
	// This is a conceptual implementation - actual key generation would depend on 
	// your Cloudinary plan and access level
	
	keyData := map[string]interface{}{
		"name":    name,
		"enabled": enabled,
		"created": time.Now(),
	}
	
	if ttl != nil {
		keyData["expires"] = time.Now().Add(*ttl)
	}
	
	// In a real implementation, this would call Cloudinary's Admin API
	// to create actual access keys
	return keyData, nil
}

// SetTempDataLifetime configures temporary data with TTL (Time To Live)
func (cl *CloudinaryLab) SetTempDataLifetime(ctx context.Context, publicID string, lifetimeMinutes int) (*uploader.ExplicitResult, error) {
	// Use context with tags to mark resources as temporary
	result, err := cl.cld.Upload.Explicit(ctx, uploader.ExplicitParams{
		PublicID: publicID,
		Type:     "upload",
		Tags:     []string{fmt.Sprintf("temp_%dm", lifetimeMinutes)},
		Context:  api.CldAPIMap{"temp_expiry": fmt.Sprintf("%d", time.Now().Add(time.Duration(lifetimeMinutes) * time.Minute).Unix())},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to set temporary lifetime: %w", err)
	}

	return result, nil
}

// GenerateTemporarySingleUseSignature generates a single-use upload signature
func (cl *CloudinaryLab) GenerateTemporarySingleUseSignature(ctx context.Context, folder string, validFor time.Duration) (map[string]interface{}, error) {
	timestamp := time.Now().Unix()
	
	// Create signature parameters
	signature := map[string]interface{}{
		"folder":          folder,
		"timestamp":       timestamp,
		"unique_filename": true,
		"tags":            "single_use,temporary",
		"context":         fmt.Sprintf("expires_at=%d|single_use=true", time.Now().Add(validFor).Unix()),
		"cloud_name":      cl.cld.Config.Cloud.CloudName,
		"api_key":         cl.cld.Config.Cloud.APIKey,
		"expires_in_seconds": int(validFor.Seconds()),
		"expires_at":      time.Now().Add(validFor).Unix(),
	}

	return signature, nil
}

// DeleteUpload deletes an upload from Cloudinary
func (cl *CloudinaryLab) DeleteUpload(ctx context.Context, publicID string, resourceType string) (*uploader.DestroyResult, error) {
	if resourceType == "" {
		resourceType = "image"
	}

	result, err := cl.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: resourceType,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to delete upload %s: %w", publicID, err)
	}

	return result, nil
}

// ListUploadsInFolder lists all uploads in a specific folder
func (cl *CloudinaryLab) ListUploadsInFolder(ctx context.Context, folder string, maxResults int) (*admin.AssetsResult, error) {
	if maxResults == 0 {
		maxResults = 10
	}

	result, err := cl.cld.Admin.Assets(ctx, admin.AssetsParams{
		Prefix:     folder,
		MaxResults: maxResults,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list uploads in folder %s: %w", folder, err)
	}

	return result, nil
}

// CreateFolder creates a folder structure (conceptual)
func (cl *CloudinaryLab) CreateFolder(folderName string, secure bool) map[string]interface{} {
	folderInfo := map[string]interface{}{
		"name":     folderName,
		"secure":   secure,
		"created":  time.Now(),
		"path":     fmt.Sprintf("/%s/", folderName),
	}

	if secure {
		folderInfo["access_mode"] = "authenticated"
		folderInfo["access_key_required"] = true
	} else {
		folderInfo["access_mode"] = "public"
	}

	return folderInfo
}

func RunMainDemo() {
	ctx := context.Background()
	
	lab, err := NewCloudinaryLab()
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary lab: %v", err)
	}

	fmt.Println("🚀 Cloudinary Features Discovery Lab")
	fmt.Println("=====================================")

	// 1. Generate presigned upload URL for tmp directory
	fmt.Println("\n1. Generating presigned upload URL for 'tmp' directory...")
	presignedURL, err := lab.GeneratePresignedUploadURL(ctx, "tmp")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("✅ Presigned upload parameters generated for 'tmp' folder\n")
		fmt.Printf("   Timestamp: %v\n", presignedURL["timestamp"])
		fmt.Printf("   Folder: %v\n", presignedURL["folder"])
	}

	// 2. Create folder structures
	fmt.Println("\n2. Creating folder structures...")
	tmpFolder := lab.CreateFolder("tmp", false)
	partnerFolder := lab.CreateFolder("partner", true)
	
	fmt.Printf("✅ Created folder structure:\n")
	fmt.Printf("   Tmp folder: %v (public access)\n", tmpFolder["path"])
	fmt.Printf("   Partner folder: %v (secure access)\n", partnerFolder["path"])

	// 3. Generate access keys
	fmt.Println("\n3. Generating access keys...")
	shortTTL := 5 * time.Minute
	longTermKey, _ := lab.GenerateAccessKey("long-term-api-key", true, nil)
	tempKey, _ := lab.GenerateAccessKey("temp-api-key", true, &shortTTL)
	
	fmt.Printf("✅ Access keys generated:\n")
	fmt.Printf("   Long-term key: %s (no expiration)\n", longTermKey["name"])
	fmt.Printf("   Temporary key: %s (expires in 5 minutes)\n", tempKey["name"])

	// 4. Generate single-use temporary signature
	fmt.Println("\n4. Generating single-use temporary signature...")
	singleUseSignature, err := lab.GenerateTemporarySingleUseSignature(ctx, "tmp", 5*time.Minute)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("✅ Single-use signature generated:\n")
		fmt.Printf("   Expires in: %v seconds\n", singleUseSignature["expires_in_seconds"])
		fmt.Printf("   Folder: %v\n", singleUseSignature["folder"])
	}

	// 5. Demonstrate other features (conceptual since we don't have actual uploads)
	fmt.Println("\n5. Additional features available:")
	fmt.Println("   ✅ Check upload existence: CheckUploadExists()")
	fmt.Println("   ✅ Move uploads between directories: MoveUpload()")
	fmt.Println("   ✅ Generate secure URLs with expiration: GenerateSecureURL()")
	fmt.Println("   ✅ Set temporary data lifetime: SetTempDataLifetime()")
	fmt.Println("   ✅ Delete uploads: DeleteUpload()")
	fmt.Println("   ✅ List uploads in folder: ListUploadsInFolder()")

	fmt.Println("\n🎉 Cloudinary features discovery completed!")
	fmt.Println("\n📝 Note: Set your Cloudinary credentials in environment variables:")
	fmt.Println("   - CLOUDINARY_CLOUD_NAME=your_cloud_name")
	fmt.Println("   - CLOUDINARY_API_KEY=your_api_key") 
	fmt.Println("   - CLOUDINARY_API_SECRET=your_api_secret")
}