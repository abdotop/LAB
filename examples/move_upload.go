package examples

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/abdotop/LAB/config"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// CloudinaryLab struct for examples
type CloudinaryLab struct {
	cld *cloudinary.Cloudinary
}

// MoveUploadExample demonstrates moving files from tmp to partner directory
func MoveUploadExample() {
	ctx := context.Background()

	config := config.NewCloudinaryConfig()
	cld, err := config.NewCloudinaryInstance()
	if err != nil {
		log.Fatalf("Failed to create Cloudinary instance: %v", err)
	}

	// Example: Moving a file from tmp to partner directory
	fromPublicID := "tmp/sample_image"
	toPublicID := "partner/sample_image"

	fmt.Printf("Moving upload from %s to %s...\n", fromPublicID, toPublicID)

	// In a real scenario, you would first check if the source file exists
	lab := &CloudinaryLab{cld: cld}
	
	exists, assetInfo, err := lab.CheckUploadExists(ctx, fromPublicID)
	if err != nil {
		log.Printf("Error checking if upload exists: %v", err)
		return
	}

	if !exists {
		fmt.Printf("Upload %s does not exist. You would need to upload a file first.\n", fromPublicID)
		return
	}

	fmt.Printf("✅ Found upload: %s\n", assetInfo.PublicID)

	// Move the upload
	result, err := lab.MoveUpload(ctx, fromPublicID, toPublicID)
	if err != nil {
		log.Printf("Error moving upload: %v", err)
		return
	}

	fmt.Printf("✅ Successfully moved upload to: %s\n", result.PublicID)

	// Verify the move by checking the new location
	exists, newAssetInfo, err := lab.CheckUploadExists(ctx, toPublicID)
	if err != nil {
		log.Printf("Error verifying moved upload: %v", err)
		return
	}

	if exists {
		fmt.Printf("✅ Verified: Upload now exists at %s\n", newAssetInfo.PublicID)
		
		// Generate secure URL for the partner directory
		secureURL, err := lab.GenerateSecureURL(toPublicID, time.Now().Add(1*time.Hour))
		if err != nil {
			log.Printf("Error generating secure URL: %v", err)
		} else {
			fmt.Printf("🔐 Secure URL (expires in 1 hour): %s\n", secureURL)
		}
	}
}

// CheckUploadExists checks if an upload exists
func (cl *CloudinaryLab) CheckUploadExists(ctx context.Context, publicID string) (bool, *admin.AssetResult, error) {
	result, err := cl.cld.Admin.Asset(ctx, admin.AssetParams{
		PublicID: publicID,
	})
	
	if err != nil {
		return false, nil, err
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

// GenerateSecureURL generates a secure URL with expiration
func (cl *CloudinaryLab) GenerateSecureURL(publicID string, expireAt time.Time) (string, error) {
	// Create a basic delivery URL (conceptual implementation)
	// The actual API call depends on your specific Cloudinary setup
	baseURL := fmt.Sprintf("https://res.cloudinary.com/%s/image/upload/%s", 
		cl.cld.Config.Cloud.CloudName, 
		publicID)

	// Add authentication token (conceptual implementation)
	secureURL := fmt.Sprintf("%s?auth_token=%s&expires=%d", 
		baseURL, 
		"secure-token", 
		expireAt.Unix())

	return secureURL, nil
}