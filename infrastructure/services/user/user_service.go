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

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("bycrypt.GenPass: %w", err)
	}

	user.ID = ID
	user.Password = string(hash)
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()
	user.Role = domain.UserRole
	user.ImageUrl = fmt.Sprintf("https://api.dicebear.com/5.x/bottts-neutral/svg?seed=%d", user.CreatedAt)
	user.VerifiedEmail = false

	if err := s.repo.Save(user); err != nil {
		return fmt.Errorf("repo.Save(): %w", err)
	}

	user.Password = ""

	return nil
}

func (s *userService) CreateFromGoogleUser(user *domain.User) error {
	ID, err := uuid.NewUUID()

	if err != nil {
		return fmt.Errorf("uuid.NewUUID() err: %w", err)
	}

	random, err := uuid.NewRandom()

	if err != nil {
		return fmt.Errorf("uuid.NewRandom() err: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(random.String()), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("bycrypt.GenPass() err: %w", err)
	}

	user.ID = ID
	user.Password = string(hash)
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()
	user.Role = domain.UserRole
	user.Age = 18

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
