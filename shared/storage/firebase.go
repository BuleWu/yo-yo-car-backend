package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"zavrsni/yo-yo-car/firebase"
)

func UploadProfilePicture(file multipart.File, header *multipart.FileHeader, userID string) (string, error) {
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

	objectPath := fmt.Sprintf("profile_pictures/%s-%s", userID, header.Filename)
	wc := bucket.Object(objectPath).NewWriter(ctx)
	defer wc.Close()

	wc.ContentType = header.Header.Get("Content-Type")

	if _, err := io.Copy(wc, file); err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket.Name(), objectPath)
	return url, nil
}
