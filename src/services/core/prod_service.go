package core

import (
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/google/uuid"
)

type prodService struct {
	repo domain.ProductRepository
}

func NewProductService(repo domain.ProductRepository) domain.ProductService {
	return &prodService{repo: repo}
}

func (s *prodService) Create(prod *domain.Product) error
func (s *prodService) GetAll() ([]domain.Product, error)
func (s *prodService) GetByID(ID uuid.UUID) (*domain.Product, error)
func (s *prodService) Update(ID uuid.UUID, prod *domain.Product) error
func (s *prodService) Delete(ID uuid.UUID) error
func (s *prodService) GetLatestProds(lim int) ([]domain.Product, error)
func (s *prodService) GetByTags(tags ...string) ([]domain.Product, error)
