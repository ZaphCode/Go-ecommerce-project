package utils

import "errors"

//* OAuth providers

const (
	GoogleProvider  = "google"
	DiscordProvider = "discord"
	GithubProvider  = "github"
)

func GetOAuthProviders() []string {
	return []string{GoogleProvider, DiscordProvider, GithubProvider}
}

//* User roles

const (
	UserRole      = "user"
	ModeratorRole = "moderator"
	AdminRole     = "admin"
)

func GetUserRoles() []string {
	return []string{UserRole, ModeratorRole, AdminRole}
}

//* Firestore collection names

const (
	UserColl  = "users"
	ProdColl  = "products"
	AddrColl  = "addresses"
	OrderColl = "orders"
	CardColl  = "cards"
	CategColl = "categories"
)

//* Errors

var (
	ErrNotFound = errors.New("resourse not found")
)
