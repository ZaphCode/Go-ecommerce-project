package product

import (
	"fmt"
	"log"
	"sort"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/shared"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
)

//* Implementation

type memoryProductRepo struct {
	shared.MemoryRepo[domain.Product]
}

//* Constructor

func NewMemoryProductRepository(im ...domain.Product) domain.ProductRepository {
	store := utils.NewSyncMap[uuid.UUID, domain.Product]()

	for _, m := range im {
		if err := store.Set(m.ID, m); err != nil {
			log.Fatal(err)
		}
	}

	return &memoryProductRepo{
		shared.MemoryRepo[domain.Product]{
			Store: store,
		},
	}
}

func (r *memoryProductRepo) FindOrderBy(field string, ord string) ([]domain.Product, error) {
	ps, err := r.Store.GetAll()

	if err != nil {
		return nil, err
	}

	switch ord {
	case "ASC":
		sort.Slice(ps, func(i, j int) bool {
			return ps[i].CreatedAt < ps[j].CreatedAt
		})
	case "DESC":
		sort.Slice(ps, func(i, j int) bool {
			return ps[i].CreatedAt > ps[j].CreatedAt
		})
	default:
		return nil, fmt.Errorf("invalid order method. use 'ASC' or 'DESC'")
	}

	return ps, nil
}
