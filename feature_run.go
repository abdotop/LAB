package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

func uploadForDemo(lab *CloudinaryLab, ctx context.Context, publicID string, folder string, tags []string) (string, error) {
	tempFile, err := os.Create("temp_demo_file.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	_, _ = tempFile.WriteString("This is a demo file.")
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	uploadResult, err := lab.cld.Upload.Upload(ctx, tempFile.Name(), uploader.UploadParams{
		PublicID: publicID,
		Folder:   folder,
		Tags:     tags,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	return uploadResult.PublicID, nil
}

func runFeatureDemo() {
	if os.Getenv("CLOUDINARY_URL") == "" {
		log.Fatal("CLOUDINARY_URL is not set. Please set it to run this demo.")
	}

	ctx := context.Background()
	lab, err := NewCloudinaryLab()
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary lab: %v", err)
	}

	fmt.Println("🚀 Starting New Feature Demonstration Script")
	fmt.Println("============================================")

	// 1. Create a Folder
	folderName := "demo_folder_" + uuid.New().String()
	fmt.Printf("\n1. Creating folder: '%s'...\n", folderName)
	createFolderResult, err := lab.CreateFolder(ctx, folderName)
	if err != nil || !createFolderResult.Success {
		log.Fatalf("❌ Failed to create folder: %v", err)
	}
	fmt.Printf("✅ Folder '%s' created successfully.\n", createFolderResult.Path)

	// 2. Generate a Presigned Upload URL
	fmt.Println("\n2. Generating a temporary, single-use presigned URL...")
	presignedParams, err := lab.GeneratePresignedUploadURL(ctx, folderName, 10*time.Minute)
	if err != nil {
		log.Fatalf("❌ Failed to generate presigned URL: %v", err)
	}
	fmt.Println("✅ Generated presigned upload parameters:")
	fmt.Printf("   - API Key: %s\n", presignedParams["api_key"])
	fmt.Printf("   - Signature: %s...\n", presignedParams["signature"][:10])
	fmt.Printf("   - Timestamp: %s\n", presignedParams["timestamp"])
	fmt.Printf("   - Public ID: %s\n", presignedParams["public_id"])
	fmt.Println("   (Note: A client would now use these params to upload a file directly to Cloudinary)")

	// 3. Upload a file to test other features
	demoTag := "demo_tag_" + uuid.New().String()
	uploadPublicID := "demo_asset_" + uuid.New().String()
	fmt.Printf("\n3. Uploading a test asset with tag '%s'...\n", demoTag)
	uploadedPublicID, err := uploadForDemo(lab, ctx, uploadPublicID, "", []string{demoTag})
	if err != nil {
		log.Fatalf("❌ Failed to upload test asset: %v", err)
	}
	fmt.Printf("✅ Test asset uploaded. Public ID: %s\n", uploadedPublicID)
	defer lab.DeleteUpload(ctx, uploadedPublicID, "raw") // Defer cleanup

	// 4. List uploads by that tag
	fmt.Printf("\n4. Listing assets with tag '%s'...\n", demoTag)
	assets, err := lab.ListUploadsByTag(ctx, demoTag, api.File, 10)
	if err != nil {
		log.Fatalf("❌ Failed to list assets by tag: %v", err)
	}
	if len(assets.Assets) != 1 || assets.Assets[0].PublicID != uploadedPublicID {
		log.Fatalf("❌ Verification failed: Did not find the correct asset by tag.")
	}
	fmt.Printf("✅ Found asset '%s' by tag successfully.\n", assets.Assets[0].PublicID)

	// 5. Move the asset
	newFolderNameForMove := "moved_assets_" + uuid.New().String()
	newPublicID := newFolderNameForMove + "/" + uploadPublicID
	fmt.Printf("\n5. Moving asset '%s' to folder '%s'...\n", uploadedPublicID, newFolderNameForMove)
	// 5a. Rename
	renameResult, err := lab.MoveUpload(ctx, uploadedPublicID, newPublicID, "raw")
	if err != nil {
		log.Fatalf("❌ Failed to rename asset: %v", err)
	}
	movedPublicID := renameResult.PublicID
	fmt.Printf("✅ Asset renamed to '%s'.\n", movedPublicID)
	defer lab.DeleteUpload(ctx, movedPublicID, "raw") // Defer cleanup for moved asset
	// 5b. Update folder
	_, err = lab.UpdateAssetFolder(ctx, movedPublicID, newFolderNameForMove, api.File)
	if err != nil {
		log.Fatalf("❌ Failed to update asset folder: %v", err)
	}
	fmt.Printf("✅ Asset folder updated to '%s'.\n", newFolderNameForMove)

	// 6. Check for its existence in the new location
	fmt.Printf("\n6. Verifying existence of moved asset '%s'...\n", movedPublicID)
	exists, _, err := lab.CheckUploadExists(ctx, movedPublicID, api.File)
	if err != nil || !exists {
		log.Fatalf("❌ Failed to find asset after move: %v", err)
	}
	fmt.Println("✅ Asset found in new location.")

	fmt.Println("\n============================================")
	fmt.Println("🎉 All new features demonstrated successfully!")
}
