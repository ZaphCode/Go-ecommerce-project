package shared

import (
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
)

type MemoryRepo[T domain.DomainModel] struct {
	Store *utils.SyncMap[uuid.UUID, T]
}

//* Constructor

func NewMemoryRepo[T domain.DomainModel](
	store *utils.SyncMap[uuid.UUID, T],
) *MemoryRepo[T] {
	return &MemoryRepo[T]{
		Store: store,
	}
}

func (r *MemoryRepo[T]) Save(m *T) error {
	return nil
}

func (r *MemoryRepo[T]) Find() ([]T, error) {
	return nil, nil
}

func (r *MemoryRepo[T]) FindByID(ID uuid.UUID) (*T, error) {
	return new(T), nil
}

func (r *MemoryRepo[T]) FindByField(fld string, val any) (*T, error) {

	return new(T), nil
}

func (r *MemoryRepo[T]) FindWhere(fld, cond string, val any) ([]T, error) {

	return nil, nil
}

func (r *MemoryRepo[T]) Update(ID uuid.UUID, m *T) error {

	return nil
}

func (r *MemoryRepo[T]) UpdateField(ID uuid.UUID, fld string, val any) error {
	return nil
}

func (r *MemoryRepo[T]) Remove(ID uuid.UUID) error {
	return nil
}

// Helper

func (r *MemoryRepo[T]) getModelID(m T) string {
	return m.GetStringID()
}
