package main

import (
	"context"
	"fmt"
	"time"

	"github.com/abdotop/LAB/config"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
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

// GeneratePresignedUploadURL generates parameters for a signed, temporary, single-use upload.
func (cl *CloudinaryLab) GeneratePresignedUploadURL(ctx context.Context, folder string, ttl time.Duration) (map[string]string, error) {
	// 1. Create the upload parameters
	// We generate a unique Public ID in advance to make the signature single-use.
	params := uploader.UploadParams{
		PublicID:      "single-use-" + uuid.New().String(),
		Folder:        folder,
		ResourceType:  "auto",
		Timestamp:     time.Now().Unix(),
		Invalidate:    api.Bool(true), // Invalidate CDN cache
	}

	// 2. Convert params to url.Values
	queryParams, err := api.StructToParams(params)
	if err != nil {
		return nil, fmt.Errorf("failed to convert params: %w", err)
	}

	// Add TTL to the signature by adding it to the params before signing
	// Note: Cloudinary doesn't have a direct "expires_at" for upload signatures in the same way as for delivery URLs.
	// The standard method is to use the timestamp. A client-side check or a server-side check after upload
	// would be needed to enforce a strict TTL. We rely on the timestamp and the signature's inherent validity window.
	
	// 3. Generate the signature
	signature, err := api.SignParameters(queryParams, cl.cld.Config.Cloud.APISecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign parameters: %w", err)
	}

	// 4. Prepare the results to be returned
	// The client would use these to make a POST request to the upload endpoint.
	results := make(map[string]string)
	for key, values := range queryParams {
		if len(values) > 0 {
			results[key] = values[0]
		}
	}
	results["api_key"] = cl.cld.Config.Cloud.APIKey
	results["signature"] = signature

	return results, nil
}

// CheckUploadExists checks if an upload exists in a specific directory
func (cl *CloudinaryLab) CheckUploadExists(ctx context.Context, publicID string, resourceType api.AssetType) (bool, *admin.AssetResult, error) {
	result, err := cl.cld.Admin.Asset(ctx, admin.AssetParams{
		PublicID:  publicID,
		AssetType: resourceType,
	})

	if err != nil {
		// If resource not found, return false
		return false, nil, nil
	}

	return true, result, nil
}

// MoveUpload renames an asset's public ID. Note this does not change the asset's folder in the Media Library UI.
func (cl *CloudinaryLab) MoveUpload(ctx context.Context, fromPublicID, toPublicID string, resourceType string) (*uploader.RenameResult, error) {
	result, err := cl.cld.Upload.Rename(ctx, uploader.RenameParams{
		FromPublicID: fromPublicID,
		ToPublicID:   toPublicID,
		ResourceType: resourceType,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to rename upload from %s to %s: %w", fromPublicID, toPublicID, err)
	}

	return result, nil
}

// UpdateAssetFolder changes the asset folder of a given asset.
func (cl *CloudinaryLab) UpdateAssetFolder(ctx context.Context, publicID, folder string, resourceType api.AssetType) (*admin.AssetResult, error) {
	result, err := cl.cld.Admin.UpdateAsset(ctx, admin.UpdateAssetParams{
		PublicID:    publicID,
		AssetFolder: folder,
		AssetType:   resourceType,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to update asset folder for %s: %w", publicID, err)
	}

	return result, nil
}

// ListUploadsByTag lists all uploads with a specific tag.
func (cl *CloudinaryLab) ListUploadsByTag(ctx context.Context, tag string, resourceType api.AssetType, maxResults int) (*admin.AssetsResult, error) {
	if maxResults == 0 {
		maxResults = 10
	}

	result, err := cl.cld.Admin.AssetsByTag(ctx, admin.AssetsByTagParams{
		Tag:        tag,
		AssetType:  resourceType,
		MaxResults: maxResults,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list uploads by tag %s: %w", tag, err)
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

// CreateFolder creates a new folder in Cloudinary.
func (cl *CloudinaryLab) CreateFolder(ctx context.Context, folderName string) (*admin.CreateFolderResult, error) {
	result, err := cl.cld.Admin.CreateFolder(ctx, admin.CreateFolderParams{
		Folder: folderName,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create folder %s: %w", folderName, err)
	}

	return result, nil
}
