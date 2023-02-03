package domain

import "github.com/google/uuid"

type ServiceCrudOperations[T DomainModel] interface {
	Create(m *T) error
	GetAll() ([]T, error)
	GetByID(ID uuid.UUID) (*T, error)
	Update(ID uuid.UUID, m *T) error
	Delete(ID uuid.UUID) error
}

type RepositoryCrudOperations[T DomainModel] interface {
	Save(m *T) error
	Find() ([]T, error)
	FindByID(ID uuid.UUID) (*T, error)
	Update(ID uuid.UUID, m *T) error
	Remove(ID uuid.UUID) error
}

type DomainModel interface {
	GetStringID() string
}

type Model struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

func (m Model) GetStringID() string {
	return m.ID.String()
}
