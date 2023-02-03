package order

import (
	"cloud.google.com/go/firestore"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/shared"
)

//* Implementation

type firestoreOrderRepo struct {
	shared.FirestoreRepo[domain.Order]
}

//* Constructor

func NewFirestoreOrderRepository(
	client *firestore.Client,
	collName string,
) domain.OrderRepository {
	return &firestoreOrderRepo{
		shared.FirestoreRepo[domain.Order]{
			Client:    client,
			CollName:  collName,
			ModelName: "order",
		},
	}
}
