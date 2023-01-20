package core

import (
	"fmt"
	"time"

	"github.com/ZaphCode/clean-arch/domain"
	"github.com/ZaphCode/clean-arch/infrastructure/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(user *domain.User) error {
	ID, err := uuid.NewUUID()

	if err != nil {
		return fmt.Errorf("uuid.NewUUID: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("bycrypt.GenPass(): %w", err)
	}

	user.ID = ID
	user.Password = string(hash)
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()
	user.Role = domain.UserRole
	user.ImageUrl = fmt.Sprintf("https://api.dicebear.com/5.x/bottts-neutral/svg?seed=%d", user.CreatedAt)
	user.VerifiedEmail = false

	if err := s.repo.Save(user); err != nil {
		return err
	}

	user.Password = ""

	return nil
}

func (s *userService) CreateFromOAuth(user *domain.User) error {
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
		return err
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

func (s *userService) GetByEmail(email string) (*domain.User, error) {
	return s.repo.FindByField("Email", email)
}

func (s *userService) VerifyEmail(ID uuid.UUID) error {
	return s.repo.UpdateField(ID, "VerifiedEmail", true)
}

func (s *userService) UpdatePassword(ID uuid.UUID, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("bycrypt.GenPass(): %w", err)
	}

	return s.repo.UpdateField(ID, "Password", string(hash))
}

func (s *userService) Update(ID uuid.UUID, user *domain.User) error {
	if user.Role != "" && !utils.ItemInSlice(user.Role, domain.GetUserRoles()) {
		return fmt.Errorf("invalid role")
	}
	return s.repo.Update(ID, &domain.User{
		Username:   user.Username,
		CustomerID: user.CustomerID,
		ImageUrl:   user.ImageUrl,
		Age:        user.Age,
		Role:       user.Role,
	})
}

func (s *userService) Delete(ID uuid.UUID) error {
	return s.repo.Remove(ID)
}
