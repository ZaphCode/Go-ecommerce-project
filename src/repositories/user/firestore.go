package user

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/shared"
	"github.com/google/uuid"
)

//* Implementation

type firestoreUserRepo struct {
	shared.FirestoreCrudRepo[domain.User]
}

//* Constructor

func NewFirestoreUserRepository(
	client *firestore.Client,
	collName string,
) domain.UserRepository {
	return &firestoreUserRepo{
		FirestoreCrudRepo: shared.FirestoreCrudRepo[domain.User]{
			Client:   client,
			CollName: collName,
		},
	}
}

//* Methods

func (r *firestoreUserRepo) FindByField(field string, data any) (*domain.User, error) {
	snapshots, err := r.Client.
		Collection(r.CollName).
		Where(field, "==", data).
		Limit(1).
		Documents(context.TODO()).
		GetAll()

	if err != nil {
		return nil, fmt.Errorf("documents.GetAll(): %w", err)
	}

	if len(snapshots) <= 0 {
		return nil, nil
	}

	var user domain.User

	for _, snapshot := range snapshots {
		if err := snapshot.DataTo(&user); err != nil {
			return nil, fmt.Errorf("snapshot.DataTo(): %w", err)
		} else {
			break
		}
	}

	return &user, nil
}

func (r *firestoreUserRepo) UpdateField(ID uuid.UUID, field string, data any) error {
	docRef := r.Client.Collection(r.CollName).Doc(ID.String())

	_, err := docRef.Update(context.TODO(), []firestore.Update{
		{Path: field, Value: data},
		{Path: "UpdatedAt", Value: time.Now().Unix()},
	})

	if err != nil {
		return err
	}

	return nil
}
