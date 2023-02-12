package shared

import (
	"errors"
	"fmt"
	"reflect"

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
	if m == nil {
		return fmt.Errorf("you cant instert nil value")
	}

	id, err := uuid.Parse((*m).GetStringID())

	if err != nil {
		return err
	}

	return r.Store.Set(id, *m)
}

func (r *MemoryRepo[T]) Find() ([]T, error) {
	return r.Store.GetAll()
}

func (r *MemoryRepo[T]) FindByID(ID uuid.UUID) (*T, error) {
	m, err := r.Store.Get(ID)

	if err != nil {

		if errors.Is(err, utils.ErrNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &m, nil
}

func (r *MemoryRepo[T]) FindByField(fld string, val any) (*T, error) {
	mdls, err := r.Store.GetAll()

	if err != nil {
		return nil, err
	}

	var model *T

	for _, mdl := range mdls {
		fv, err := utils.GetStructField(mdl, fld)

		if err != nil {
			return nil, err
		}

		if !utils.IsSameType(fv, val) {
			return nil, fmt.Errorf("the value of the field %q is %T and you send a %T", fld, fv, val)
		}

		if fv == val {
			model = &mdl
			break
		}
	}

	return model, nil
}

func (r *MemoryRepo[T]) FindWhere(fld, cond string, val any) ([]T, error) {
	mdls, err := r.Store.GetAll()

	if err != nil {
		return nil, err
	}

	ms := []T{}

	for _, mdl := range mdls {
		fv, err := utils.GetStructField(mdl, fld)

		if err != nil {
			return nil, err
		}

		var statement bool

		switch cond {
		case "==":
			statement = reflect.DeepEqual(fv, val)
		case "!=":
			statement = !reflect.DeepEqual(fv, val)
		default:
			return nil, fmt.Errorf("invalid condition")
		}

		if statement {
			ms = append(ms, mdl)
		}
	}

	return ms, nil
}

func (r *MemoryRepo[T]) Update(ID uuid.UUID, uf domain.UpdateFields) error {
	if uf == nil {
		return fmt.Errorf("you cant insert nil")
	}

	m, err := r.FindByID(ID)

	if err != nil {
		return err
	}

	if m == nil {
		return fmt.Errorf("data not found")
	}

	if err := utils.UpdateStructFields(m, uf); err != nil {
		return err
	}

	return r.Store.Update(ID, *m)
}

func (r *MemoryRepo[T]) UpdateField(ID uuid.UUID, fld string, val any) error {
	m, err := r.FindByID(ID)

	if err != nil {
		return err
	}

	if m == nil {
		return fmt.Errorf("data not found")
	}

	if err := utils.SetStructField(m, fld, val); err != nil {
		return err
	}

	return r.Store.Update(ID, *m)
}

func (r *MemoryRepo[T]) Remove(ID uuid.UUID) error {
	return r.Store.Remove(ID)
}

func (r *MemoryRepo[T]) clear() error {
	r.Store.Clear()
	return nil
}
