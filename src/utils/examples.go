package utils

import (
	"time"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/google/uuid"
)

//* Users

var UserAdmin = domain.User{
	Model: domain.Model{
		ID:        uuid.MustParse("e44ef83a-a1c7-11ed-a865-7e82d40d4740"),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	},
	CustomerID:    "",
	Username:      "Zaphkiel",
	Email:         "zaph@fapi.com",
	Password:      "$2a$10$D/cBZACDWfS4r910QhyIhucC/IKKD.4ilevC44j2CozjW0fscNBaG",
	Role:          AdminRole,
	VerifiedEmail: true,
	ImageUrl:      "https://i.etsystatic.com/15149849/r/il/16852c/2126930485/il_570xN.2126930485_oub4.jpg",
	Age:           19,
}

var UserExp1 = domain.User{
	Model: domain.Model{
		ID:        uuid.MustParse("6ac802ef-4c9e-4c03-8271-7abb13c5318b"),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	},
	CustomerID:    "cus_NokAmwAreSjg1Y",
	Username:      "John Doe",
	Email:         "john@testing.com",
	Password:      "$2a$10$0W5T1cZiUrBc8xJCViVh2.BYE7oMXuE8ogfXjDjbu4KRg4GiQqgeS",
	Role:          UserRole,
	VerifiedEmail: true,
	ImageUrl:      "https://i.etsystatic.com/15149849/r/il/16852c/2126930485/il_570xN.2126930485_oub4.jpg",
	Age:           20,
}

var UserExp2 = domain.User{
	Model: domain.Model{
		ID:        uuid.MustParse("4657dbe5-a4ab-4cf0-94fc-58122b3a71d4"),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	},
	CustomerID:    "",
	Username:      "Foo Bar",
	Email:         "foo@testing.com",
	Password:      "$2a$10$0W5T1cZiUrBc8xJCViVh2.BYE7oMXuE8ogfXjDjbu4KRg4GiQqgeS",
	Role:          ModeratorRole,
	VerifiedEmail: true,
	ImageUrl:      "https://i.etsystatic.com/15149849/r/il/16852c/2126930485/il_570xN.2126930485_oub4.jpg",
	Age:           23,
}

//* Categories

var CategoryExp1 = domain.Category{
	Model: domain.Model{
		ID:        uuid.MustParse("05f2fad6-e109-410c-aab4-792ac524414b"),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	},
	Name: "headsets",
}

var CategoryExp2 = domain.Category{
	Model: domain.Model{
		ID:        uuid.MustParse("5837ba7e-30e5-4ced-b5a6-d437bff7deee"),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	},
	Name: "clothes",
}

var CategoryExp3 = domain.Category{
	Model: domain.Model{
		ID:        uuid.MustParse("be17f397-d069-403f-ae0c-8f6c65b12415"),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	},
	Name: "mouses",
}

//* Products

var ProductExp1 = domain.Product{
	Model: domain.Model{
		ID:        uuid.MustParse("9964ac5b-ba6d-4b46-84c3-70889137674d"),
		CreatedAt: time.Now().Add(time.Hour * 4).Unix(),
		UpdatedAt: time.Now().Add(time.Hour * 4).Unix(),
	},
	Category:     "clothes",
	Name:         "Black T-shirt",
	Description:  "the best black t-shirt.",
	Price:        2400,
	DiscountRate: 14,
	ImagesUrl:    []string{"https://parspng.com/wp-content/uploads/2022/07/Tshirtpng.parspng.com_.png"},
	Tags:         []string{"clothes", "t-shirt", "black"},
	Avalible:     true,
}

var ProductExp2 = domain.Product{
	Model: domain.Model{
		ID:        uuid.MustParse("3582d8ea-44c8-4bcd-a08b-c773783d5493"),
		CreatedAt: time.Now().Add(time.Hour * 24).Unix(),
		UpdatedAt: time.Now().Add(time.Hour * 24).Unix(),
	},
	Category:     "headsets",
	Name:         "Corsair void pro",
	Description:  "the best headset.",
	Price:        6000,
	DiscountRate: 9,
	ImagesUrl:    []string{"https://http2.mlstatic.com/D_NQ_NP_798698-MLA41021638035_032020-O.jpg"},
	Tags:         []string{"headsets", "corsair", "technology"},
	Avalible:     true,
}

//* Cards

// var CardExp1 = payment.Card{
// 	// Model: domain.Model{
// 	// 	ID:        uuid.MustParse("e67ab973-8a23-4fab-8724-58410a82fde4"),
// 	// 	CreatedAt: time.Now().Unix(),
// 	// 	UpdatedAt: time.Now().Unix(),
// 	// },
// 	Country:   "USA",
// 	Name:      "main card",
// 	ExpMonth:  12,
// 	ExpYear:   2024,
// 	Brand:     "visa",
// 	Last4:     "4534",
// 	PaymentID: "cd_fasdfjkalkfjklvlakje",
// }

// var CardExp2 = payment.Card{
// 	// Model: domain.Model{
// 	// 	ID:        uuid.MustParse("b601782a-9fc2-4a27-9a50-3cb55133d2f6"),
// 	// 	CreatedAt: time.Now().Unix(),
// 	// 	UpdatedAt: time.Now().Unix(),
// 	// },
// 	Country:   "USA",
// 	Name:      "digital card",
// 	ExpMonth:  3,
// 	ExpYear:   2025,
// 	Brand:     "visa",
// 	Last4:     "6660",
// 	PaymentID: "cd_jfladsflkajcvlajldjeP",
// }

//* Addresses

var AddrExp1 = domain.Address{
	Model: domain.Model{
		ID:        uuid.MustParse("f54e460a-0a6d-4462-b177-da36b524c0cb"),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	},
	UserID:     UserExp1.ID,
	Country:    "USA",
	Name:       "Word address",
	City:       "Washintong City",
	PostalCode: "21515",
	Line1:      "Polaco street",
	Line2:      "Sterling 21",
	State:      "Washintong",
}

var AddrExp2 = domain.Address{
	Model: domain.Model{
		ID:        uuid.MustParse("d10d5c0d-875b-41e4-a205-14f9eb3ac4a9"),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	},
	UserID:     UserExp1.ID,
	Country:    "Mexico",
	Name:       "Cangrejos home",
	City:       "Cabo San Lucas",
	PostalCode: "23473",
	Line1:      "Calle Ballena",
	Line2:      "Calle Pargo",
	State:      "Baja California Sur",
}
