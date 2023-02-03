package domain

import (
	"github.com/google/uuid"
)

//* Model

type Category struct {
	Model
	Name string `json:"name"`
}

//* Service

type CategoryService interface {
	Create(c *Category) error
	GetByID(ID uuid.UUID) (*Category, error)
	GetByName(n string) (*Category, error)
	GetAll() ([]Category, error)
	Delete(ID uuid.UUID) error
}

//* Repository

type CategoryRepository interface {
	Find() ([]Category, error)
	FindByID(ID uuid.UUID) (*Category, error)
	FindByField(f string, v any) (*Category, error)
	Save(c *Category) error
	Remove(ID uuid.UUID) error
}
