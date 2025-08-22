package config

import (
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

// CloudinaryConfig holds Cloudinary configuration
type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
}

// NewCloudinaryConfig creates a new Cloudinary configuration
func NewCloudinaryConfig() *CloudinaryConfig {
	return &CloudinaryConfig{
		CloudName: getEnvOrDefault("CLOUDINARY_CLOUD_NAME", "demo"),
		APIKey:    getEnvOrDefault("CLOUDINARY_API_KEY", ""),
		APISecret: getEnvOrDefault("CLOUDINARY_API_SECRET", ""),
	}
}

// NewCloudinaryInstance creates a new Cloudinary instance
func (c *CloudinaryConfig) NewCloudinaryInstance() (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(c.CloudName, c.APIKey, c.APISecret)
	if err != nil {
		return nil, err
	}
	return cld, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}