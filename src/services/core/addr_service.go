package core

import (
	"fmt"
	"time"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/google/uuid"
)

type addressService struct {
	addrRepo domain.AddressRepository
	usrRepo  domain.UserRepository
}

func NewAddressService(
	addrRepo domain.AddressRepository,
	usrRepo domain.UserRepository,
) domain.AddressService {
	return &addressService{
		addrRepo: addrRepo,
		usrRepo:  usrRepo,
	}
}

func (s *addressService) Create(addr *domain.Address) error {
	ID, err := uuid.NewUUID()

	if err != nil {
		return fmt.Errorf("error generating uuid: %s", err)
	}

	usr, err := s.usrRepo.FindByID(addr.UserID)

	if err != nil || usr == nil {
		return fmt.Errorf("error getting user. %s", err)
	}

	addr.ID = ID
	addr.CreatedAt = time.Now().Unix()
	addr.UpdatedAt = time.Now().Unix()

	s.addrRepo.Save(addr)

	return nil
}

func (s *addressService) GetAll() ([]domain.Address, error) {
	return s.addrRepo.Find()
}

func (s *addressService) GetByID(ID uuid.UUID) (*domain.Address, error) {
	return s.addrRepo.FindByID(ID)
}

func (s *addressService) Update(ID uuid.UUID, uf domain.UpdateFields) error {
	delete(uf, "UserID")
	delete(uf, "Model")
	return s.addrRepo.Update(ID, uf)
}

func (s *addressService) Delete(ID uuid.UUID) error {
	return s.addrRepo.Remove(ID)
}

func (s *addressService) GetAllByUserID(ID uuid.UUID) ([]domain.Address, error) {
	usr, err := s.usrRepo.FindByID(ID)

	if err != nil || usr == nil {
		return nil, fmt.Errorf("error getting user")
	}

	return s.addrRepo.FindWhere("UserID", "==", ID)
}
