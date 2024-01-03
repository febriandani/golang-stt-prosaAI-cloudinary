package utils

// Import Cloudinary and other necessary libraries
//===================
import (
	"context"
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/spf13/viper"
)

func CredentialsCloudinary() (*cloudinary.Cloudinary, error) {
	// Add your Cloudinary credentials, set configuration parameter
	// Secure=true to return "https" URLs, and create a context
	//===================
	cld, err := cloudinary.New()
	if err != nil {
		return nil, err
	}
	cld.Config.URL.Secure = true
	cld.Config.Cloud.APIKey = viper.GetString("CLOUDINARY.CLOUD_API_KEY")
	cld.Config.Cloud.APISecret = viper.GetString("CLOUDINARY.CLOUD_API_SECRET")
	cld.Config.Cloud.CloudName = viper.GetString("CLOUDINARY.CLOUD_CLOUDNAME")
	cld.Config.URL.Domain = viper.GetString("CLOUDINARY.URL_DOMAIN")

	return cld, nil
}

func UploadCloudinary(cld *cloudinary.Cloudinary, ctx context.Context, filename string) (string, error) {

	// Upload the image.
	// Set the asset's public ID and allow overwriting the asset with new versions
	resp, err := cld.Upload.Upload(
		ctx,
		filename,
		uploader.UploadParams{
			ResourceType: "video",
			Folder:       "stt",
		})

	if err != nil {
		log.Println("error ", err)
		return "", err
	}

	// Log the delivery URL
	log.Println("****2. Upload an image****\nDelivery URL:", resp.SecureURL)
	return resp.SecureURL, nil
}
