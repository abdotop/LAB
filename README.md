# Cloudinary Features Discovery Lab

This repository demonstrates various Cloudinary features and capabilities using the Go SDK. It showcases advanced file management, security, and automation features.

## 🚀 Features Implemented

### 1. Presigned Upload URLs with Directory Structure
- Generate secure, presigned upload URLs for specific directories (e.g., `tmp`)
- Automatic file organization and unique filename generation
- Time-limited upload signatures for enhanced security

### 2. Upload Management
- Check for existence of uploads in specific directories
- Move uploads between directories (e.g., `tmp` → `partner`)
- List all uploads in a folder with pagination
- Delete uploads with proper cleanup

### 3. Secure Directory Access
- Make directories accessible only with generated keys
- IP-based access restrictions (configurable)
- Role-based access control with different permission levels
- Time-bound access with automatic expiration

### 4. Access Key Management
- Generate long-term API keys for persistent access
- Create temporary access keys with TTL (Time To Live)
- Key rotation capabilities for enhanced security
- Different access levels (read-only, full-access, admin)

### 5. Temporary Data with TTL
- Set automatic expiration for uploads (e.g., 5 minutes)
- Tag-based cleanup system for temporary files
- Context-based metadata for lifecycle management
- Automated cleanup strategies

### 6. Single-Use Upload Signatures
- Generate temporary, one-time upload signatures
- Prevent signature reuse for enhanced security
- Custom expiration times for different use cases
- Integration with client-side upload workflows

## 🛠️ Setup

### Prerequisites
- Go 1.19 or later
- Cloudinary account with API credentials

### Installation

1. Clone the repository:
```bash
git clone https://github.com/abdotop/LAB.git
cd LAB
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
export CLOUDINARY_CLOUD_NAME=your_cloud_name
export CLOUDINARY_API_KEY=your_api_key
export CLOUDINARY_API_SECRET=your_api_secret
```

### Running the Examples

#### Main Demo
```bash
go run main.go
```

#### Specific Examples
```bash
# Move upload between directories
go run examples/move_upload.go

# Temporary uploads with TTL
go run examples/temporary_ttl.go

# Secure directory access
go run examples/secure_access.go
```

## 📖 Usage Examples

### Generate Presigned Upload URL

```go
lab, _ := NewCloudinaryLab()
presignedURL, err := lab.GeneratePresignedUploadURL(ctx, "tmp")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Upload URL parameters: %+v\n", presignedURL)
```

### Check Upload Existence

```go
exists, assetInfo, err := lab.CheckUploadExists(ctx, "tmp/my_file")
if exists {
    fmt.Printf("File found: %s (%d bytes)\n", assetInfo.PublicID, assetInfo.Bytes)
}
```

### Move Upload Between Directories

```go
result, err := lab.MoveUpload(ctx, "tmp/my_file", "partner/my_file")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("File moved to: %s\n", result.PublicID)
```

### Generate Secure URL with Expiration

```go
secureURL, err := lab.GenerateSecureURL("partner/secure_file", time.Now().Add(1*time.Hour))
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Secure URL (expires in 1 hour): %s\n", secureURL)
```

### Set Temporary Data Lifetime

```go
result, err := lab.SetTempDataLifetime(ctx, "tmp/temp_file", 5) // 5 minutes
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Temporary file expires in 5 minutes: %s\n", result.PublicID)
```

### Delete Upload

```go
result, err := lab.DeleteUpload(ctx, "tmp/unwanted_file", "image")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("File deleted: %s\n", result.Result)
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `CLOUDINARY_CLOUD_NAME` | Your Cloudinary cloud name | Yes |
| `CLOUDINARY_API_KEY` | Your Cloudinary API key | Yes |
| `CLOUDINARY_API_SECRET` | Your Cloudinary API secret | Yes |

### Directory Structure

```
LAB/
├── config/
│   └── cloudinary.go      # Cloudinary configuration
├── examples/
│   ├── move_upload.go     # Upload movement example
│   ├── temporary_ttl.go   # Temporary files with TTL
│   └── secure_access.go   # Secure directory access
├── main.go                # Main demonstration
├── go.mod
├── go.sum
└── README.md
```

## 🔐 Security Features

### Authentication Methods
- **API Key Authentication**: Standard Cloudinary API authentication
- **Signed URLs**: Time-limited URLs with cryptographic signatures
- **Access Tokens**: Temporary tokens for specific resources
- **IP Restrictions**: Whitelist-based access control

### Access Control Levels
- **Public**: Open access for public resources
- **Authenticated**: Requires valid authentication
- **Restricted**: IP-based or key-based restrictions
- **Private**: Full access control with custom permissions

## 🚨 Best Practices

1. **Never hardcode API credentials** - Use environment variables
2. **Implement key rotation** - Regularly update access keys
3. **Set appropriate TTL** - Don't make temporary URLs too long-lived
4. **Monitor usage** - Track API calls and storage usage
5. **Clean up temporary files** - Implement automated cleanup for temp directories
6. **Validate uploads** - Check file types and sizes before processing
7. **Use HTTPS** - Always use secure connections for API calls

## 🧪 Testing

The examples in this repository are designed to work with any valid Cloudinary account. Make sure to:

1. Set up your Cloudinary credentials
2. Test with small files first
3. Monitor your Cloudinary dashboard for uploads
4. Clean up test files to avoid unnecessary storage costs

## 📚 Additional Resources

- [Cloudinary Go SDK Documentation](https://cloudinary.com/documentation/go_integration)
- [Cloudinary Upload API](https://cloudinary.com/documentation/upload_images)
- [Cloudinary Admin API](https://cloudinary.com/documentation/admin_api)
- [Security and Access Control](https://cloudinary.com/documentation/control_access_to_media)

## 🤝 Contributing

Feel free to submit issues, fork the repository, and create pull requests for any improvements.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.