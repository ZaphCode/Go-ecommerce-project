package core

import (
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/google/uuid"
)

type categoryService struct {
	repo domain.CategoryRepository
}

func NewCategoryService(repo domain.CategoryRepository) domain.CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(c *domain.Category) error
func (s *categoryService) GetAll() ([]domain.Category, error)
func (s *categoryService) GetByName(n string) (*domain.Category, error)
func (s *categoryService) Delete(ID uuid.UUID) error
