package product

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

type firestoreProductRepo struct {
	shared.FirestoreCrudRepo[domain.Product]
}

//* Constructor

func NewFirestoreProductRepository(
	client *firestore.Client,
	collName string,
) domain.ProductRepository {
	return &firestoreProductRepo{
		FirestoreCrudRepo: shared.FirestoreCrudRepo[domain.Product]{
			Client:   client,
			CollName: collName,
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
		return nil, fmt.Errorf("documents.GetAll(): %w", err)
	}

	ps := make([]domain.Product, len(ss))

	for i, s := range ss {
		var p domain.Product

		if err := s.DataTo(&p); err != nil {
			return nil, fmt.Errorf("snapshot.DataTo(): %w", err)
		}

		ps[i] = p
	}

	return ps, nil
}

func (r *firestoreProductRepo) FindWhere(field string, cond string, val any) ([]domain.Product, error) {
	ss, err := r.Client.
		Collection(r.CollName).
		Where(field, cond, val).
		Documents(context.TODO()).
		GetAll()

	if err != nil {
		return nil, fmt.Errorf("documents.GetAll(): %w", err)
	}

	ps := make([]domain.Product, len(ss))

	for i, s := range ss {
		var p domain.Product

		if err := s.DataTo(&p); err != nil {
			return nil, fmt.Errorf("snapshot.DataTo(): %w", err)
		}

		ps[i] = p
	}

	return ps, nil
}

func (r *firestoreProductRepo) UpdateField(ID uuid.UUID, field string, val any) error {
	docRef := r.Client.Collection(r.CollName).Doc(ID.String())

	_, err := docRef.Update(context.TODO(), []firestore.Update{
		{Path: field, Value: val},
		{Path: "UpdatedAt", Value: time.Now().Unix()},
	})

	if err != nil {
		return err
	}

	return nil
}
