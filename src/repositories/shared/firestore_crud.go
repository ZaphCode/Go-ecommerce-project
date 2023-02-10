package shared

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//* Firestore generic repo

type FirestoreRepo[T domain.DomainModel] struct {
	Client    *firestore.Client
	CollName  string
	ModelName string
}

//* Constructor

func NewFirestoreRepo[T domain.DomainModel](
	client *firestore.Client,
	collName string,
	modelName string,
) *FirestoreRepo[T] {
	return &FirestoreRepo[T]{
		Client:    client,
		CollName:  collName,
		ModelName: modelName,
	}
}

func (r *FirestoreRepo[T]) Save(m *T) error {
	if m == nil {
		return fmt.Errorf("cannot accept nil value")
	}

	docRef := r.Client.Collection(r.CollName).Doc((*m).GetStringID())

	_, err := docRef.Create(context.TODO(), m)

	if err != nil {
		return fmt.Errorf("creating user error: %w", err)
	}

	return nil
}

func (r *FirestoreRepo[T]) Find() ([]T, error) {
	ss, err := r.Client.Collection(r.CollName).Documents(context.TODO()).GetAll()

	if err != nil {
		return nil, fmt.Errorf("error fetching %s: %w", r.CollName, err)
	}

	if len(ss) == 0 {
		return []T{}, nil
	}

	ms := make([]T, len(ss))

	for i, s := range ss {
		var m T

		if err := s.DataTo(&m); err != nil {
			return nil, fmt.Errorf("error parsing a %s: %w", r.ModelName, err)
		}

		ms[i] = m
	}

	return ms, nil
}

func (r *FirestoreRepo[T]) FindByID(ID uuid.UUID) (*T, error) {
	ref := r.Client.Collection(r.CollName).Doc(ID.String())

	s, err := ref.Get(context.TODO())

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("docRef.Get(): %s", err.Error())
	}

	var m T

	if err := s.DataTo(&m); err != nil {
		return nil, fmt.Errorf("snapshot.DataTo(): %s", err.Error())
	}

	return &m, nil
}

func (r *FirestoreRepo[T]) FindByField(fld string, val any) (*T, error) {
	fv, err := utils.GetStructField(new(T), fld)

	if err != nil {
		return nil, err
	}

	if !utils.IsSameType(fv, val) {
		return nil, fmt.Errorf("invalid type. the %s field only acepts %T values", fld, fv)
	}

	ss, err := r.Client.
		Collection(r.CollName).
		Where(fld, "==", val).
		Limit(1).
		Documents(context.TODO()).
		GetAll()

	if err != nil {
		return nil, fmt.Errorf("error getting %s documents: %s", r.ModelName, err)
	}

	if len(ss) <= 0 {
		return nil, nil
	}

	var m T

	for _, s := range ss {
		if err := s.DataTo(&m); err != nil {
			return nil, fmt.Errorf("snapshot.DataTo(): %w", err)
		} else {
			break
		}
	}

	return &m, nil
}

func (r *FirestoreRepo[T]) FindWhere(fld, cond string, val any) ([]T, error) {
	if _, err := utils.GetStructField(new(T), fld); err != nil {
		return nil, err
	}

	ss, err := r.Client.
		Collection(r.CollName).
		Where(fld, cond, val).
		Documents(context.TODO()).
		GetAll()

	if err != nil {
		return nil, fmt.Errorf("documents.GetAll(): %w", err)
	}

	ms := make([]T, len(ss))

	for i, s := range ss {
		var m T

		if err := s.DataTo(&m); err != nil {
			return nil, fmt.Errorf("snapshot.DataTo(): %w", err)
		}

		ms[i] = m
	}

	return ms, nil
}

func (r *FirestoreRepo[T]) Update(ID uuid.UUID, uf domain.UpdateFields) error {
	if uf == nil {
		return fmt.Errorf("cannot accept nil value")
	}

	ref := r.Client.Collection(r.CollName).Doc(ID.String())

	updates := []firestore.Update{
		{Path: "UpdatedAt", Value: time.Now().Unix()},
	}

	for fldName, fldVal := range uf {
		updates = append(updates, firestore.Update{
			Path: fldName, Value: fldVal,
		})
	}

	_, err := ref.Update(context.Background(), updates)

	if err != nil {
		return err
	}

	return nil
}

func (r *FirestoreRepo[T]) UpdateField(ID uuid.UUID, fld string, val any) error {
	fv, err := utils.GetStructField(new(T), fld)

	if err != nil {
		return err
	}

	if !utils.IsSameType(fv, val) {
		return fmt.Errorf("invalid type. the %s field only acepts %T values", fld, fv)
	}

	ref := r.Client.Collection(r.CollName).Doc(ID.String())

	_, err = ref.Update(context.TODO(), []firestore.Update{
		{Path: fld, Value: val},
		{Path: "UpdatedAt", Value: time.Now().Unix()},
	})

	if err != nil {
		return fmt.Errorf("error updating %s: %s", r.ModelName, err)
	}

	return nil
}

func (r *FirestoreRepo[T]) Remove(ID uuid.UUID) error {
	m, err := r.FindByID(ID)

	if err != nil {
		return fmt.Errorf("error looking for %s: %s", r.ModelName, err.Error())
	}

	if m == nil {
		return fmt.Errorf("%s not found", r.ModelName)
	}

	ref := r.Client.Collection(r.CollName).Doc(ID.String())

	_, err = ref.Delete(context.TODO())

	if err != nil {
		return fmt.Errorf("error deleting %s: %s", r.ModelName, err)
	}

	return nil
}

func (r *FirestoreRepo[T]) clear() error {
	return utils.DeleteFirestoreCollection(r.Client, r.CollName, 5)
}
