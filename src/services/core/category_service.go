package core

import (
	"fmt"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/google/uuid"
)

type categoryService struct {
	repo domain.CategoryRepository
}

func NewCategoryService(repo domain.CategoryRepository) domain.CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(c *domain.Category) error {
	cat, err := s.GetByName(c.Name)

	if err != nil {
		return err
	}

	if cat != nil {
		return fmt.Errorf("that category already exists")
	}

	return s.repo.Save(c)
}

func (s *categoryService) GetAll() ([]domain.Category, error) {
	return s.repo.Find()
}

func (s *categoryService) GetByName(n string) (*domain.Category, error) {
	cat, err := s.repo.FindByField("Name", n)

	if err != nil {
		return nil, err
	}

	return cat, nil
}

func (s *categoryService) Delete(ID uuid.UUID) error {
	c, err := s.repo.FindByID(ID)

	if err != nil {
		return err
	}

	if c == nil {
		return fmt.Errorf("category not found")
	}

	return s.repo.Remove(ID)
}
