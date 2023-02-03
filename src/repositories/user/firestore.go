package user

import (
	"cloud.google.com/go/firestore"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/shared"
)

//* Implementation

type firestoreUserRepo struct {
	shared.FirestoreRepo[domain.User]
}

//* Constructor

func NewFirestoreUserRepository(
	client *firestore.Client,
	collName string,
) domain.UserRepository {
	return &firestoreUserRepo{
		shared.FirestoreRepo[domain.User]{
			Client:    client,
			CollName:  collName,
			ModelName: "user",
		},
	}
}
