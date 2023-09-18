package core

import (
	"fmt"
	"time"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
)

type orderService struct {
	ordRepo  domain.OrderRepository
	addrRepo domain.AddressRepository
}

func NewOrderService(
	ordRepo domain.OrderRepository,
	addrRepo domain.AddressRepository,
) domain.OrderService {
	return &orderService{
		ordRepo:  ordRepo,
		addrRepo: addrRepo,
	}
}

func (s *orderService) Create(ord *domain.Order) error {
	ID, err := uuid.NewUUID()

	if err != nil {
		return fmt.Errorf("error generating uuid: %s", err)
	}

	if addr, err := s.addrRepo.FindByID(ord.AddressID); err != nil || addr == nil {
		return fmt.Errorf("invalid address id")
	}

	ord.ID = ID
	ord.Status = utils.StatusPending
	ord.CreatedAt = time.Now().Unix()
	ord.UpdatedAt = time.Now().Unix()

	s.ordRepo.Save(ord)

	return nil
}

func (s *orderService) GetAll() ([]domain.Order, error) {
	return s.ordRepo.Find()
}

func (s *orderService) UpdateStatus(ID uuid.UUID, status string) error {
	return s.ordRepo.UpdateField(ID, "Status", status)
}

func (s *orderService) SetPaidStatus(ID uuid.UUID, paid bool) error {
	return s.ordRepo.UpdateField(ID, "Paid", paid)
}

func (s *orderService) GetAllByUserID(ID uuid.UUID) ([]domain.Order, error) {
	return s.ordRepo.FindWhere("UserID", "==", ID)
}

func (s *orderService) Delete(ID uuid.UUID) error {
	return s.ordRepo.Remove(ID)
}
