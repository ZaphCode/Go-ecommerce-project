package core

import (
	"fmt"
	"time"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/google/uuid"
)

type orderService struct {
	ordRepo domain.OrderRepository
	usrRepo domain.UserRepository
}

func NewOrderService(
	ordRepo domain.OrderRepository,
	usrRepo domain.UserRepository,
) domain.OrderService {
	return &orderService{
		ordRepo: ordRepo,
		usrRepo: usrRepo,
	}
}

func (s *orderService) Create(ord *domain.Order) error {
	ID, err := uuid.NewUUID()

	if err != nil {
		return fmt.Errorf("error generating uuid: %s", err)
	}

	ord.ID = ID
	ord.CreatedAt = time.Now().Unix()
	ord.UpdatedAt = time.Now().Unix()

	s.ordRepo.Save(ord)

	return nil
}

func (s *orderService) GetAll() ([]domain.Order, error) {
	return s.ordRepo.Find()
}

func (s *orderService) GetAllByUserID(ID uuid.UUID) ([]domain.Order, error) {
	return s.ordRepo.FindWhere("UserID", "==", ID)
}
