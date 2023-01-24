package user

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
)

//* Implementation

type memoryUserRepositoryImpl struct {
	mu    sync.RWMutex
	store map[uuid.UUID]domain.User
}

//* Constructor

func NewMemoryUserRepository() domain.UserRepository {
	return &memoryUserRepositoryImpl{
		mu: sync.RWMutex{},
		store: map[uuid.UUID]domain.User{
			uuid.MustParse("8ded83fe-93c8-11ed-ab0f-d8bbc1a27048"): {
				DBModel: utils.DBModel{
					ID:        uuid.MustParse("8ded83fe-93c8-11ed-ab0f-d8bbc1a27048"),
					CreatedAt: time.Now().Unix(),
					UpdatedAt: time.Now().Unix(),
				},
				CustomerID:    "",
				Username:      "John Doe",
				Email:         "john@gmail.com",
				Role:          "user",
				Password:      "kldsj-3djlvjckl-ya4ejkgrio-pdvcipo",
				VerifiedEmail: false,
				ImageUrl:      "",
				Age:           15,
			},
			uuid.MustParse("3afc3021-9395-11ed-a8b6-d8bbc1a27048"): {
				DBModel: utils.DBModel{
					ID:        uuid.MustParse("3afc3021-9395-11ed-a8b6-d8bbc1a27048"),
					CreatedAt: time.Now().Unix(),
					UpdatedAt: time.Now().Unix(),
				},
				CustomerID:    "",
				Username:      "Foo Bar",
				Email:         "foo@gmail.com",
				Role:          "user",
				Password:      "kldsj-3djlvjckl-ya4ejkgrio-pdvcipo",
				VerifiedEmail: true,
				ImageUrl:      "",
				Age:           18,
			},
		},
	}
}

//* Methods

func (r *memoryUserRepositoryImpl) Save(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[user.ID] = *user
	return nil
}

func (r *memoryUserRepositoryImpl) Find() (users []domain.User, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, v := range r.store {
		users = append(users, v)
	}
	return
}

func (r *memoryUserRepositoryImpl) FindByID(ID uuid.UUID) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.store[ID]
	if !ok {
		return nil, nil
	}
	return &user, nil
}

func (r *memoryUserRepositoryImpl) FindByField(field string, value any) (user *domain.User, err error) {
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

func (r *memoryUserRepositoryImpl) Remove(ID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.store, ID)
	return nil
}

func (r *memoryUserRepositoryImpl) Update(ID uuid.UUID, user *domain.User) error {
	return fmt.Errorf("not implemented")
}

func (r *memoryUserRepositoryImpl) UpdateField(ID uuid.UUID, field string, value any) error {
	return fmt.Errorf("not implemented")
}
