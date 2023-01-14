package utils

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/storage"
)

func GetFirestoreClient(app *firebase.App) *firestore.Client {
	client, err := app.Firestore(context.Background())

	if err != nil || client == nil {
		log.Fatal("Error getting firestore client", err)
	}

	return client
}

func GetStorageClient(app *firebase.App) *storage.Client {
	client, err := app.Storage(context.Background())

	if err != nil || client == nil {
		log.Fatal("Error getting storage client", err)
	}

	return client
}
