package user

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/ZaphCode/clean-arch/domain"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//* Implementation

type firestoreUserRepositoryImpl struct {
	client   *firestore.Client
	collName string
}

//* Constructor

func NewFirestoreUserRepository(
	client *firestore.Client,
	collName string,
) domain.UserRepository {
	return &firestoreUserRepositoryImpl{
		client:   client,
		collName: collName,
	}
}

//* Methods

func (r *firestoreUserRepositoryImpl) Save(user *domain.User) error {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	docRef := r.client.Collection(r.collName).Doc(user.ID.String())

	existingUser, err := r.FindByField("Email", user.Email)

	if err != nil {
		return fmt.Errorf("internal server error: %w", err)
	}

	if existingUser != nil {
		return fmt.Errorf("email taken")
	}

	_, err = docRef.Create(ctx, user)

	if err != nil {
		return fmt.Errorf("creating user error: %w", err)
	}

	return nil
}

func (r *firestoreUserRepositoryImpl) Find() ([]domain.User, error) {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	snapshots, err := r.client.Collection(r.collName).Documents(ctx).GetAll()

	if err != nil {
		return nil, fmt.Errorf("documents.GetAll(): %w", err)
	}

	users := make([]domain.User, len(snapshots))

	for i, snapshot := range snapshots {
		var user domain.User

		if err := snapshot.DataTo(&user); err != nil {
			return nil, fmt.Errorf("snapshot.DataTo(): %w", err)
		}

		users[i] = user
	}

	return users, nil
}

func (r *firestoreUserRepositoryImpl) FindByID(ID uuid.UUID) (*domain.User, error) {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	docRef := r.client.Collection(r.collName).Doc(ID.String())

	snapshot, err := docRef.Get(ctx)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("docRef.Get(): %s", err.Error())
	}

	var user domain.User

	if err := snapshot.DataTo(&user); err != nil {
		return nil, fmt.Errorf("snapshot.DataTo(): %s", err.Error())
	}

	return &user, nil
}

func (r *firestoreUserRepositoryImpl) FindByField(field string, data any) (*domain.User, error) {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	snapshots, err := r.client.
		Collection(r.collName).
		Where(field, "==", data).
		Limit(1).
		Documents(ctx).
		GetAll()

	if err != nil {
		return nil, fmt.Errorf("documents.GetAll(): %w", err)
	}

	if len(snapshots) <= 0 {
		return nil, nil
	}

	var user domain.User

	for _, snapshot := range snapshots {
		if err := snapshot.DataTo(&user); err != nil {
			return nil, fmt.Errorf("snapshot.DataTo(): %w", err)
		} else {
			break
		}
	}

	return &user, nil
}

func (r *firestoreUserRepositoryImpl) Remove(ID uuid.UUID) error {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	docRef := r.client.Collection(r.collName).Doc(ID.String())

	_, err := docRef.Delete(ctx)

	if err != nil {
		return fmt.Errorf("docRef.Get(): %s", err.Error())
	}

	return nil
}

func (r *firestoreUserRepositoryImpl) Update(ID uuid.UUID, user *domain.User) error {
	docRef := r.client.Collection(r.collName).Doc(ID.String())

	updates := []firestore.Update{
		{Path: "UpdatedAt", Value: time.Now().Unix()},
	}

	v := reflect.ValueOf(user).Elem()

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

func (r *firestoreUserRepositoryImpl) UpdateField(ID uuid.UUID, field string, data any) error {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	docRef := r.client.Collection(r.collName).Doc(ID.String())

	_, err := docRef.Update(ctx, []firestore.Update{
		{Path: field, Value: data},
		{Path: "UpdatedAt", Value: time.Now().Unix()},
	})

	if err != nil {
		return err
	}

	return nil
}
