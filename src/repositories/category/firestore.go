package category

import (
	"cloud.google.com/go/firestore"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/google/uuid"
)

//* Implementation

type firestoreCategoryRepo struct {
	Client   *firestore.Client
	CollName string
}

//* Constructor

func NewFirestoreCategoryRepository(
	client *firestore.Client,
	collName string,
) domain.CategoryRepository {
	return &firestoreCategoryRepo{}
}

func (r *firestoreCategoryRepo) Save(c *domain.Category) error

func (r *firestoreCategoryRepo) FindByField(f string, v any) (*domain.Category, error)

func (r *firestoreCategoryRepo) Remove(ID uuid.UUID) error
