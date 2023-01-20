package user

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/ZaphCode/clean-arch/domain"
	"github.com/google/uuid"
)

//* Implementation

type mockUserRepositoryImpl struct {
	mu    sync.RWMutex
	store map[uuid.UUID]domain.User
}

//* Constructor

func NewMemoryUserRepository() domain.UserRepository {
	return &mockUserRepositoryImpl{
		mu: sync.RWMutex{},
		store: map[uuid.UUID]domain.User{
			uuid.MustParse("8ded83fe-93c8-11ed-ab0f-d8bbc1a27048"): {
				ID:            uuid.MustParse("8ded83fe-93c8-11ed-ab0f-d8bbc1a27048"),
				CustomerID:    "",
				Username:      "John Doe",
				Email:         "john@gmail.com",
				Role:          "user",
				Password:      "kldsj-3djlvjckl-ya4ejkgrio-pdvcipo",
				VerifiedEmail: false,
				ImageUrl:      "",
				Age:           15,
				CreatedAt:     time.Now().Unix(),
				UpdatedAt:     time.Now().Unix(),
			},
			uuid.MustParse("3afc3021-9395-11ed-a8b6-d8bbc1a27048"): {
				ID:            uuid.MustParse("3afc3021-9395-11ed-a8b6-d8bbc1a27048"),
				CustomerID:    "",
				Username:      "Foo Bar",
				Email:         "foo@gmail.com",
				Role:          "user",
				Password:      "kldsj-3djlvjckl-ya4ejkgrio-pdvcipo",
				VerifiedEmail: true,
				ImageUrl:      "",
				Age:           18,
				CreatedAt:     time.Now().Unix(),
				UpdatedAt:     time.Now().Unix(),
			},
		},
	}
}

//* Methods

func (r *mockUserRepositoryImpl) Save(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[user.ID] = *user
	return nil
}

func (r *mockUserRepositoryImpl) Find() (users []domain.User, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, v := range r.store {
		users = append(users, v)
	}
	return
}

func (r *mockUserRepositoryImpl) FindByID(ID uuid.UUID) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.store[ID]
	if !ok {
		return nil, nil
	}
	return &user, nil
}

func (r *mockUserRepositoryImpl) FindByField(field string, value any) (user *domain.User, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.store {
		v := reflect.ValueOf(u)
		for i := 0; i < v.NumField(); i++ {
			fieldName := v.Type().Field(i).Name
			fieldValue := v.Field(i).Interface()

			if fieldName == field && fieldValue == value {
				user = &u
				return
			}
		}
	}
	return
}

func (r *mockUserRepositoryImpl) Remove(ID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.store, ID)
	return nil
}

func (r *mockUserRepositoryImpl) Update(ID uuid.UUID, user *domain.User) error {
	return fmt.Errorf("not implemented")
}

func (r *mockUserRepositoryImpl) UpdateField(ID uuid.UUID, field string, value any) error {
	return fmt.Errorf("not implemented")
}
