package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"zavrsni/yo-yo-car/firebase"
)

func UploadProfilePicture(file multipart.File, header *multipart.FileHeader, userID string) (string, error) {
	// Validate file size
	const maxFileSize = 5 * 1024 * 1024 // 5 MB
	if header.Size > maxFileSize {
		return "", fmt.Errorf("file size exceeds the maximum limit of 5 MB")
	}

	// Sanitize filename
	safeFilename := SanitizeFilename(header.Filename)
	ctx := context.Background()
	app := firebase.App

	client, err := app.Storage(ctx)
	if err != nil {
		return "", err
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return "", err
	}

	objectPath := fmt.Sprintf("profile_pictures/%s-%s", userID, safeFilename)
	wc := bucket.Object(objectPath).NewWriter(ctx)
	defer wc.Close()

	// Define an allowlist of safe MIME types
	allowedContentTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}

	// Get the Content-Type from the header
	contentType := header.Header.Get("Content-Type")

	// Validate the Content-Type
	if !allowedContentTypes[contentType] {
		return "", fmt.Errorf("invalid content type: %s", contentType)
	}

	wc.ContentType = contentType
	if _, err := io.Copy(wc, file); err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket.BucketName(), objectPath)
	return url, nil
}
