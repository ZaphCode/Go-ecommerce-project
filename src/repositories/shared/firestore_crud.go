package shared

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/ZaphCode/clean-arch/src/domain/shared"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FirestoreCrudRepo[T shared.IDomainModel] struct {
	Client    *firestore.Client
	CollName  string
	ModelName string
}

func (r *FirestoreCrudRepo[T]) Save(m *T) error {
	docRef := r.Client.Collection(r.CollName).Doc(r.getModelID(*m))

	_, err := docRef.Create(context.TODO(), m)

	if err != nil {
		return fmt.Errorf("creating user error: %w", err)
	}

	return nil
}

func (r *FirestoreCrudRepo[T]) Find() ([]T, error) {
	ss, err := r.Client.Collection(r.CollName).Documents(context.TODO()).GetAll()

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

func (r *FirestoreCrudRepo[T]) FindByID(ID uuid.UUID) (*T, error) {
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

func (r *FirestoreCrudRepo[T]) FindWhere(fld, cond string, val any) ([]T, error) {
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

func (r *FirestoreCrudRepo[T]) FindByField(fld string, val any) (*T, error) {
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

func (r *FirestoreCrudRepo[T]) Update(ID uuid.UUID, m *T) error {
	ref := r.Client.Collection(r.CollName).Doc(ID.String())

	updates := []firestore.Update{
		{Path: "UpdatedAt", Value: time.Now().Unix()},
	}

	v := reflect.ValueOf(m).Elem()

	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i).Name
		fieldValue := v.Field(i).Interface()

		if !v.Field(i).IsZero() {
			updates = append(updates, firestore.Update{
				Path: fieldName, Value: fieldValue,
			})
		}
	}

	_, err := ref.Update(context.Background(), updates)

	if err != nil {
		return err
	}

	return nil
}

func (r *FirestoreCrudRepo[T]) UpdateField(ID uuid.UUID, fld string, val any) error {
	ref := r.Client.Collection(r.CollName).Doc(ID.String())

	_, err := ref.Update(context.TODO(), []firestore.Update{
		{Path: fld, Value: val},
		{Path: "UpdatedAt", Value: time.Now().Unix()},
	})

	if err != nil {
		return fmt.Errorf("error updating %s: %s", r.ModelName, err)
	}

	return nil
}

func (r *FirestoreCrudRepo[T]) Remove(ID uuid.UUID) error {
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

// Helper

func (r *FirestoreCrudRepo[T]) getModelID(m T) string {
	return m.GetStringID()
}
