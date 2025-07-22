package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
)

var App *firebase.App
var Firestore *firestore.Client

func InitFirebase(ctx context.Context, credentialsFile string, projectID string) {
	conf := &firebase.Config{ProjectID: projectID}
	opt := option.WithCredentialsFile(credentialsFile)

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
	}
	App = app

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing Firestore client: %v", err)
	}
	Firestore = client
}
