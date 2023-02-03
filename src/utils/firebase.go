package utils

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/storage"
	"google.golang.org/api/iterator"
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

func DeleteFirestoreCollection(
	client *firestore.Client,
	collName string,
	batchSize int,
) error {
	ctx := context.Background()

	for {
		iter := client.Collection(collName).Limit(batchSize).Documents(ctx)
		numDeleted := 0

		batch := client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)

		if err != nil {
			return err
		}
	}

}
