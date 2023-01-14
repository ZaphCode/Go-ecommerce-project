package domain

import "github.com/google/uuid"

const (
	UserRole      = "user"
	ModeratorRole = "moderator"
	AdminRole     = "admin"
)

func GetUserRoles() []string {
	return []string{UserRole, ModeratorRole, AdminRole}
}

type User struct {
	ID            uuid.UUID `json:"id"`
	CustomerID    string    `json:"customer_id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Role          string    `json:"role"`
	Password      string    `json:"-"`
	VerifiedEmail bool      `json:"verified_email"`
	ImageUrl      string    `json:"image_url"`
	Age           uint16    `json:"age"`
	CreatedAt     int64     `json:"created_at"`
	UpdatedAt     int64     `json:"updated_at"`
}

type UserService interface {
	Create(*User) error
	CreateFromOAuth(*User) error
	GetAll() ([]User, error)
	GetByID(uuid.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	VerifyEmail(uuid.UUID) error
	UpdatePassword(uuid.UUID, string) error
	Update(uuid.UUID, *User) error
	Delete(uuid.UUID) error
}

type UserRepository interface {
	Save(*User) error
	Find() ([]User, error)
	FindByID(uuid.UUID) (*User, error)
	FindByField(string, any) (*User, error)
	Remove(uuid.UUID) error
	Update(uuid.UUID, *User) error
	UpdateField(uuid.UUID, string, any) error
}
