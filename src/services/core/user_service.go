package core

import (
	"fmt"
	"time"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// TODO: Add card and addr repo to user service

type userService struct {
	usrRepo domain.UserRepository
	//cardRepo domain.CardRepository
	//addrRepo domain.AddressRepository
}

func NewUserService(usrRepo domain.UserRepository) domain.UserService {
	return &userService{usrRepo: usrRepo}
}

func (s *userService) Create(user *domain.User) error {
	eusr, err := s.usrRepo.FindByField("Email", user.Email)

	if err != nil {
		return fmt.Errorf("internal server error: %s", err)
	}

	if eusr != nil {
		return fmt.Errorf("email taken")
	}

	ID, err := uuid.NewUUID()

	if err != nil {
		return fmt.Errorf("uuid generation error: %s", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("error hashing password: %s", err)
	}

	user.ID = ID
	user.Password = string(hash)
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()

	if user.Age == 0 {
		user.Age = 18
	}

	if user.Role == "" {
		user.Role = utils.UserRole
	}

	if user.ImageUrl == "" {
		user.ImageUrl = fmt.Sprintf("https://api.dicebear.com/5.x/bottts-neutral/svg?seed=%d", user.CreatedAt)
	}

	if err := s.usrRepo.Save(user); err != nil {
		return err
	}

	user.Password = ""

	return nil
}

func (s *userService) GetAll() ([]domain.User, error) {
	return s.usrRepo.Find()
}

func (s *userService) GetByID(ID uuid.UUID) (*domain.User, error) {
	return s.usrRepo.FindByID(ID)
}

func (s *userService) GetByEmail(email string) (*domain.User, error) {
	return s.usrRepo.FindByField("Email", email)
}

func (s *userService) VerifyEmail(ID uuid.UUID) error {
	return s.usrRepo.UpdateField(ID, "VerifiedEmail", true)
}

func (s *userService) UpdatePassword(ID uuid.UUID, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("error hasing password %w", err)
	}

	return s.usrRepo.UpdateField(ID, "Password", string(hash))
}

func (s *userService) Update(ID uuid.UUID, uf domain.UpdateFields) error {
	delete(uf, "Email")
	delete(uf, "Model")
	return s.usrRepo.Update(ID, uf)
}

func (s *userService) Delete(ID uuid.UUID) error {
	return s.usrRepo.Remove(ID)
}
