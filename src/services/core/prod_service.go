package core

import (
	"fmt"
	"time"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/google/uuid"
)

type prodService struct {
	repo domain.ProductRepository
}

func NewProductService(repo domain.ProductRepository) domain.ProductService {
	return &prodService{repo: repo}
}

func (s *prodService) Create(prod *domain.Product) error {
	ID, err := uuid.NewUUID()

	if err != nil {
		return fmt.Errorf("error generating uuid: %s", err)
	}

	prod.ID = ID
	prod.CreatedAt = time.Now().Unix()
	prod.UpdatedAt = time.Now().Unix()

	s.repo.Save(prod)

	return nil
}

func (s *prodService) GetAll() ([]domain.Product, error) {
	return s.repo.Find()
}

func (s *prodService) GetByID(ID uuid.UUID) (*domain.Product, error) {
	return s.repo.FindByID(ID)
}

func (s *prodService) GetLatestProds(lim int) ([]domain.Product, error) {
	return s.repo.FindOrderBy("CreatedAt", "ASC")
}

func (s *prodService) GetByTags(tags ...string) ([]domain.Product, error) {
	return s.repo.FindWhere("Tags", "array-contains-any", tags)
}

func (s *prodService) GetByCategory(c string) ([]domain.Product, error) {
	return s.repo.FindWhere("Category", "==", c)
}

func (s *prodService) Update(ID uuid.UUID, prod *domain.Product) error {
	p, err := s.repo.FindByID(ID)

	if err != nil || p == nil {
		return fmt.Errorf("error getting product")
	}

	return s.repo.Update(ID, &domain.Product{
		Price:       prod.Price,
		Category:    prod.Category,
		Description: prod.Description,
		Name:        prod.Name,
		Tags:        prod.Tags,
		ImagesUrl:   prod.ImagesUrl,
	})
}

func (s *prodService) SetAvalible(ID uuid.UUID, avl bool) error {
	p, err := s.repo.FindByID(ID)

	if err != nil {
		return err
	}

	if p == nil {
		return fmt.Errorf("product not found")
	}

	return s.repo.UpdateField(ID, "Avalible", avl)
}

func (s *prodService) Delete(ID uuid.UUID) error {
	p, err := s.repo.FindByID(ID)

	if err != nil {
		return err
	}

	if p == nil {
		return fmt.Errorf("product not found")
	}

	return s.repo.Remove(ID)
}
