package domain

import "github.com/google/uuid"

type repositoryCrudOperation[T any] interface {
	Save(*T) error
	Find() ([]T, error)
	FindByID(uuid.UUID) (*T, error)
	Update(uuid.UUID, *T) error
	Remove(uuid.UUID) error
}

type serviceCrudOperation[T any] interface {
	Create(*T) error
	GetAll() ([]T, error)
	GetByID(uuid.UUID) (*T, error)
	Update(uuid.UUID, *T) error
	Delete(uuid.UUID) error
}
