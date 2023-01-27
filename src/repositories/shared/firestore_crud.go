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
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	snapshots, err := r.Client.Collection(r.CollName).Documents(ctx).GetAll()

	if err != nil {
		return nil, fmt.Errorf("documents.GetAll(): %w", err)
	}

	models := make([]T, len(snapshots))

	for i, snapshot := range snapshots {
		var model T

		if err := snapshot.DataTo(&model); err != nil {
			return nil, fmt.Errorf("snapshot.DataTo(): %w", err)
		}

		models[i] = model
	}

	return models, nil
}

func (r *FirestoreCrudRepo[T]) FindByID(ID uuid.UUID) (*T, error) {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	docRef := r.Client.Collection(r.CollName).Doc(ID.String())

	snapshot, err := docRef.Get(ctx)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("docRef.Get(): %s", err.Error())
	}

	var model T

	if err := snapshot.DataTo(&model); err != nil {
		return nil, fmt.Errorf("snapshot.DataTo(): %s", err.Error())
	}

	return &model, nil
}

func (r *FirestoreCrudRepo[T]) Update(ID uuid.UUID, m *T) error {
	docRef := r.Client.Collection(r.CollName).Doc(ID.String())

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

	_, err := docRef.Update(context.Background(), updates)

	if err != nil {
		return err
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

	docRef := r.Client.Collection(r.CollName).Doc(ID.String())

	_, err = docRef.Delete(context.TODO())

	if err != nil {
		return fmt.Errorf("docRef.Get(): %s", err.Error())
	}

	return nil
}

// Helper
func (r *FirestoreCrudRepo[T]) getModelID(m T) string {
	return m.GetStringID()
}
