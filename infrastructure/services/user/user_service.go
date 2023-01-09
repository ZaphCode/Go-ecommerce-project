package user

import (
	"fmt"
	"time"

	"github.com/ZaphCode/clean-arch/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

type userService struct {
	repo domain.UserRepository
}

func (s *userService) Create(user *domain.User) error {
	ID, err := uuid.NewUUID()

	if err != nil {
		return fmt.Errorf("uuid.NewUUID: %w", err)
	}

	user.ID = ID

	user.CreatedAt = time.Now().Unix()

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("bycrypt.GenPass: %w", err)
	}

	user.Password = string(hash)

	user.IsAdmin = false

	if err := s.repo.Save(user); err != nil {
		return fmt.Errorf("repo.Save(): %w", err)
	}

	user.Password = ""

	return nil
}

func (s *userService) GetAll() ([]domain.User, error) {
	return s.repo.Find()
}

func (s *userService) GetByID(ID uuid.UUID) (*domain.User, error) {
	return s.repo.FindByID(ID)
}

func (s *userService) Delete(ID uuid.UUID) error {
	return s.repo.Remove(ID)
}
