package product

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/shared"
)

//* Implementation

type firestoreProductRepo struct {
	shared.FirestoreRepo[domain.Product]
}

//* Constructor

func NewFirestoreProductRepository(
	client *firestore.Client,
	collName string,
) domain.ProductRepository {
	return &firestoreProductRepo{
		FirestoreRepo: shared.FirestoreRepo[domain.Product]{
			Client:    client,
			CollName:  collName,
			ModelName: "product",
		},
	}
}

func (r *firestoreProductRepo) FindOrderBy(field string, ord string) ([]domain.Product, error) {
	var d firestore.Direction

	switch ord {
	case "ASC":
		d = firestore.Asc
	case "DESC":
		d = firestore.Desc
	default:
		return nil, fmt.Errorf("invalid order method. use 'ASC' or 'DESC'")
	}

	ss, err := r.Client.Collection(r.CollName).OrderBy(field, d).Documents(context.TODO()).GetAll()

	if err != nil {
		return nil, fmt.Errorf("error getting documents: %s", err)
	}

	ps := make([]domain.Product, len(ss))

	for i, s := range ss {
		var p domain.Product

		if err := s.DataTo(&p); err != nil {
			return nil, fmt.Errorf("error parsing document %s", err)
		}

		ps[i] = p
	}

	return ps, nil
}
