package card

import (
	"cloud.google.com/go/firestore"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/shared"
)

//* Implementation

type firestoreCardRepo struct {
	shared.FirestoreCrudRepo[domain.Card]
}

//* Constructor

func NewFirestoreCardRepository(
	client *firestore.Client,
	collName string,
) domain.CardRepository {
	return &firestoreCardRepo{
		shared.FirestoreCrudRepo[domain.Card]{
			Client:    client,
			CollName:  collName,
			ModelName: "card",
		},
	}
}
