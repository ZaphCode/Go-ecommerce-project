package core

import (
	"fmt"
	"time"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/google/uuid"
)

type categoryService struct {
	catRepo  domain.CategoryRepository
	prodRepo domain.ProductRepository
}

func NewCategoryService(
	catRepo domain.CategoryRepository,
	prodRepo domain.ProductRepository,
) domain.CategoryService {
	return &categoryService{
		catRepo:  catRepo,
		prodRepo: prodRepo,
	}
}

func (s *categoryService) Create(c *domain.Category) error {
	cat, err := s.GetByName(c.Name)

	if err != nil {
		return err
	}

	if cat != nil {
		return fmt.Errorf("that category already exists")
	}

	ID, err := uuid.NewUUID()

	if err != nil {
		return fmt.Errorf("error generating uuid: %s", err)
	}

	c.ID = ID
	c.CreatedAt = time.Now().Unix()
	c.UpdatedAt = time.Now().Unix()

	return s.catRepo.Save(c)
}

func (s *categoryService) GetAll() ([]domain.Category, error) {
	return s.catRepo.Find()
}

func (s *categoryService) GetByID(ID uuid.UUID) (*domain.Category, error) {
	return s.catRepo.FindByID(ID)
}

func (s *categoryService) GetByName(n string) (*domain.Category, error) {
	cat, err := s.catRepo.FindByField("Name", n)

	if err != nil {
		return nil, err
	}

	return cat, nil
}

func (s *categoryService) Delete(ID uuid.UUID) error {
	c, err := s.catRepo.FindByID(ID)

	if err != nil {
		return err
	}

	if c == nil {
		return fmt.Errorf("category not found")
	}

	ps, err := s.prodRepo.FindWhere("Category", "==", c.Name)

	if err != nil {
		return err
	}

	if len(ps) > 0 {
		return fmt.Errorf("cannot remove %q category because it has products related", c.Name)
	}

	return s.catRepo.Remove(ID)
}
