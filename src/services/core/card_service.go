package core

import (
	"fmt"
	"time"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/google/uuid"
)

type cardService struct {
	cardRepo domain.CardRepository
	usrRepo  domain.UserRepository
}

func NewCardService(
	cardRepo domain.CardRepository,
	ursRepo domain.UserRepository,
) domain.CardService {
	return &cardService{
		cardRepo: cardRepo,
		usrRepo:  ursRepo,
	}
}

func (s *cardService) Create(card *domain.Card) error {
	ID, err := uuid.NewUUID()

	if err != nil {
		return fmt.Errorf("error generating uuid: %s", err)
	}

	usr, err := s.usrRepo.FindByID(card.UserID)

	if err != nil || usr == nil {
		return fmt.Errorf("error getting card owner")
	}

	card.ID = ID
	card.CreatedAt = time.Now().Unix()
	card.UpdatedAt = time.Now().Unix()

	s.cardRepo.Save(card)

	return nil
}

func (s *cardService) GetAll() ([]domain.Card, error) {
	return s.cardRepo.Find()
}

func (s *cardService) GetByID(ID uuid.UUID) (*domain.Card, error) {
	return s.cardRepo.FindByID(ID)
}

func (s *cardService) Update(ID uuid.UUID, uf domain.UpdateFields) error {
	delete(uf, "UserID")
	delete(uf, "Model")
	return s.cardRepo.Update(ID, uf)
}

func (s *cardService) Delete(ID uuid.UUID) error {
	return s.cardRepo.Remove(ID)
}

func (s *cardService) GetAllByUserID(ID uuid.UUID) ([]domain.Card, error) {
	usr, err := s.usrRepo.FindByID(ID)

	if err != nil || usr == nil {
		return nil, fmt.Errorf("error getting user")
	}

	return s.cardRepo.FindWhere("UserID", "==", ID)
}
