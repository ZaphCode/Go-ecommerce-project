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

// func (r *firestoreCategoryRepo) Remove(ID uuid.UUID) error {
// 	c, err := r.FindByID(ID)

// 	if err != nil {
// 		return fmt.Errorf("error looking for %s: %s", r.ModelName, err.Error())
// 	}

// 	if c == nil {
// 		return fmt.Errorf("%s not found", r.ModelName)
// 	}

// 	ss, err := r.Client.Collection(utils.ProdColl).Documents(context.TODO()).GetAll()

// 	if err != nil {
// 		return fmt.Errorf("error getting products with that category")
// 	}

// 	if len(ss) > 0 {
// 		return fmt.Errorf("that category has %d products", len(ss))
// 	}

// 	ref := r.Client.Collection(r.CollName).Doc(ID.String())

// 	_, err = ref.Delete(context.TODO())

// 	if err != nil {
// 		return fmt.Errorf("error deleting %s: %s", r.ModelName, err)
// 	}

// 	return nil
// }
