package utils

import "github.com/google/uuid"

type ServiceCrudOperations[T any] interface {
	Create(m *T) error
	GetAll() ([]T, error)
	GetByID(ID uuid.UUID) (*T, error)
	Update(ID uuid.UUID, m *T) error
	Delete(ID uuid.UUID) error
}

type RepositoryCrudOperations[T any] interface {
	Save(m *T) error
	Find() ([]T, error)
	FindByID(ID uuid.UUID) (*T, error)
	Update(ID uuid.UUID, m *T) error
	Remove(ID uuid.UUID) error
}

type DBModel struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}
