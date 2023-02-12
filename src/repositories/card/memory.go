package card

import (
	"log"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/shared"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
)

//* Implementation

type memoryCardRepo struct {
	shared.MemoryRepo[domain.Card]
}

//* Constructor

func NewMemoryCardRepository(im ...domain.Card) domain.CardRepository {
	store := utils.NewSyncMap[uuid.UUID, domain.Card]()

	for _, m := range im {
		if err := store.Set(m.ID, m); err != nil {
			log.Fatal(err)
		}
	}

	return &memoryCardRepo{
		shared.MemoryRepo[domain.Card]{
			Store: store,
		},
	}
}
