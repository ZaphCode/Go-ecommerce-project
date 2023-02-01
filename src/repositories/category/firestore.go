package category

import (
	"cloud.google.com/go/firestore"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/shared"
)

//* Implementation

type firestoreCategoryRepo struct {
	shared.FirestoreCrudRepo[domain.Category]
}

//* Constructor

func NewFirestoreCategoryRepository(
	client *firestore.Client,
	collName string,
) domain.CategoryRepository {
	return &firestoreCategoryRepo{
		FirestoreCrudRepo: shared.FirestoreCrudRepo[domain.Category]{
			Client:    client,
			CollName:  collName,
			ModelName: "category",
		},
	}
}
