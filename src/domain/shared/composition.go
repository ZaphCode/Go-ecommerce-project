package shared

import "github.com/google/uuid"

type ServiceCrudOperations[T IDomainModel] interface {
	Create(m *T) error
	GetAll() ([]T, error)
	GetByID(ID uuid.UUID) (*T, error)
	Update(ID uuid.UUID, m *T) error
	Delete(ID uuid.UUID) error
}

type RepositoryCrudOperations[T IDomainModel] interface {
	Save(m *T) error
	Find() ([]T, error)
	FindByID(ID uuid.UUID) (*T, error)
	Update(ID uuid.UUID, m *T) error
	Remove(ID uuid.UUID) error
}

type DomainModel struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

func (m DomainModel) GetStringID() string {
	return m.ID.String()
}

type IDomainModel interface {
	GetStringID() string
}
