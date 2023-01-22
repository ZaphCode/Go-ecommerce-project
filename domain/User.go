package domain

import "github.com/google/uuid"

// * Model

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

//* Service

type UserService interface {
	serviceCrudOperation[User]
	GetByEmail(email string) (*User, error)
	VerifyEmail(uuid.UUID) error
	UpdatePassword(uuid.UUID, string) error
}

//* Repository

type UserRepository interface {
	repositoryCrudOperation[User]
	FindByField(string, any) (*User, error)
	UpdateField(uuid.UUID, string, any) error
}

//* Utils

const (
	UserRole      = "user"
	ModeratorRole = "moderator"
	AdminRole     = "admin"
)

func GetUserRoles() []string {
	return []string{UserRole, ModeratorRole, AdminRole}
}
