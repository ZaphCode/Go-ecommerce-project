package card

import (
	"cloud.google.com/go/firestore"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/shared"
)

//* Implementation

type firestoreCardRepo struct {
	shared.FirestoreRepo[domain.Card]
}

//* Constructor

func NewFirestoreCardRepository(
	client *firestore.Client,
	collName string,
) domain.CardRepository {
	return &firestoreCardRepo{
		shared.FirestoreRepo[domain.Card]{
			Client:    client,
			CollName:  collName,
			ModelName: "card",
		},
	}
}
