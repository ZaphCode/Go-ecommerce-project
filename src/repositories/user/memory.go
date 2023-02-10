package user

import (
	"log"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/shared"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
)

//* Implementation

type memoryUserRepo struct {
	shared.MemoryRepo[domain.User]
}

//* Constructor

func NewMemoryUserRepository(im ...domain.User) domain.UserRepository {
	store := utils.NewSyncMap[uuid.UUID, domain.User]()

	for _, m := range im {
		if err := store.Set(m.ID, m); err != nil {
			log.Fatal(err)
		}
	}

	return &memoryUserRepo{
		shared.MemoryRepo[domain.User]{
			Store: store,
		},
	}
}
