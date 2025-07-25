package firebase

import (
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/iam"
	"cloud.google.com/go/iam/apiv1/iampb"
	"cloud.google.com/go/storage"
	"context"
	firebaseSDK "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"io"
	"log"
	"os"
)

var (
	App    *firebaseSDK.App
	Client *firestore.Client
)

func InitFirebase(ctx context.Context, credentialsFile string, projectID string, storageBucket string) {
	conf := &firebaseSDK.Config{
		ProjectID:     projectID,
		StorageBucket: storageBucket,
	}

	opt := option.WithCredentialsFile(credentialsFile)

	app, err := firebaseSDK.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	if err := setBucketPublicIAM(os.Stdout, storageBucket); err != nil {
		log.Fatalf("Failed to set bucket public: %v", err)
	}

	App = app
	Client = client
}

// setBucketPublicIAM makes all objects in a bucket publicly readable.
func setBucketPublicIAM(w io.Writer, bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	policy, err := client.Bucket(bucketName).IAM().V3().Policy(ctx)
	if err != nil {
		return fmt.Errorf("Bucket(%q).IAM().V3().Policy: %w", bucketName, err)
	}
	role := "roles/storage.objectViewer"
	policy.Bindings = append(policy.Bindings, &iampb.Binding{
		Role:    role,
		Members: []string{iam.AllUsers},
	})
	if err := client.Bucket(bucketName).IAM().V3().SetPolicy(ctx, policy); err != nil {
		return fmt.Errorf("Bucket(%q).IAM().SetPolicy: %w", bucketName, err)
	}
	fmt.Fprintf(w, "Bucket %v is now publicly readable\n", bucketName)
	return nil
}
