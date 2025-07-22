package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	firebaseSDK "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
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

	App = app
	Client = client
}
