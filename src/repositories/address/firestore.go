package address

import (
	"cloud.google.com/go/firestore"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/shared"
)

//* Implementation

type firestoreAddressRepo struct {
	shared.FirestoreCrudRepo[domain.Address]
}

//* Constructor

func NewFirestoreAddressRepository(
	client *firestore.Client,
	collName string,
) domain.AddressRepository {
	return &firestoreAddressRepo{
		shared.FirestoreCrudRepo[domain.Address]{
			Client:    client,
			CollName:  collName,
			ModelName: "address",
		},
	}
}
