package domain

import (
	"github.com/google/uuid"
)

// --------------------------------------------------------------

type ServiceCrudOperations[T DomainModel] interface {
	Create(m *T) error
	GetAll() ([]T, error)
	GetByID(ID uuid.UUID) (*T, error)
	Update(ID uuid.UUID, uf UpdateFields) error
	Delete(ID uuid.UUID) error
}

type RepositoryCrudOperations[T DomainModel] interface {
	Save(m *T) error
	Find() ([]T, error)
	FindByID(ID uuid.UUID) (*T, error)
	Update(ID uuid.UUID, uf UpdateFields) error
	Remove(ID uuid.UUID) error
}

// ---------------------------------------------------------------

type DomainModel interface {
	User | Address | Category | Product | Order | Card | ExampleModel

	GetStringID() string
	GetCreatedDate() int64
}

type Model struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

func (m Model) GetStringID() string {
	return m.ID.String()
}

func (m Model) GetCreatedDate() int64 {
	return m.CreatedAt
}

type UpdateFields map[string]interface{}

type ExampleModel struct {
	Model
	Name  string   `json:"name"`
	Tags  []string `json:"tags"`
	Check bool     `json:"check"`
	Num   int      `json:"num"`
	Float float64  `json:"float"`
}
