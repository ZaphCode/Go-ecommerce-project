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
	Name: "tenis",
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
	Available:    true,
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
	Available:    true,
}

var ProductExpToDev1 = domain.Product{
	Model: domain.Model{
		ID:        uuid.MustParse("3480c083-a8ec-11ed-9883-5266b0bd59ca"),
		CreatedAt: time.Now().Add(time.Hour * 24).Unix(),
		UpdatedAt: time.Now().Add(time.Hour * 24).Unix(),
	},
	Category:     "clothes",
	Name:         "Adidas Black T-Shirt Basketball",
	Description:  "The best T-shirt in the world.",
	Price:        2599,
	DiscountRate: 30,
	ImagesUrl: []string{
		"https://titan22.com/cdn/shop/files/IR8492-A_1082x.png?v=1690430352",
		"https://titan22.com/cdn/shop/files/IR8492-B_1082x.png?v=1690430352"},
	Tags:      []string{"t-shirts", "clothes", "Adidas"},
	Available: true,
}

var ProductExpToDev2 = domain.Product{
	Model: domain.Model{
		ID:        uuid.MustParse("2229674a-00cc-4846-8f71-4b28b6e246db"),
		CreatedAt: time.Now().Add(time.Hour * 24).Unix(),
		UpdatedAt: time.Now().Add(time.Hour * 24).Unix(),
	},
	Category:     "clothes",
	Name:         "Nike Black Cup",
	Description:  "The best cup. Super comfortable.",
	Price:        1549,
	DiscountRate: 0,
	ImagesUrl: []string{
		"https://static.nike.com/a/images/t_default/84588c76-14b7-42cb-a65c-5bbb77f6699d/gorra-estructurada-con-cierre-a-presi%C3%B3n-dri-fit-rise-hR0Mq4.png",
		"https://static.nike.com/a/images/t_PDP_1728_v1/f_auto,q_auto:eco/ad1dbe01-7a5a-4862-bdd6-416e3a120703/gorra-estructurada-con-cierre-a-presi%C3%B3n-dri-fit-rise-hR0Mq4.png",
	},
	Tags:      []string{"cups", "clothes", "Nike"},
	Available: true,
}

var ProductExpToDev3 = domain.Product{
	Model: domain.Model{
		ID:        uuid.MustParse("1119674a-00cc-4846-8f71-4b28b6e246da"),
		CreatedAt: time.Now().Add(time.Hour * 24).Unix(),
		UpdatedAt: time.Now().Add(time.Hour * 24).Unix(),
	},
	Category:     "tenis",
	Name:         "Running Tenis Puma Black",
	Description:  "Very comfortable shoes for running.",
	Price:        1530,
	DiscountRate: 10,
	ImagesUrl: []string{
		"https://martimx.vtexassets.com/arquivos/ids/489205-800-800?v=637346702472670000&width=800&height=800&aspect=true",
		"https://martimx.vtexassets.com/arquivos/ids/489294-800-800?v=637346703167700000&width=800&height=800&aspect=true",
	},
	Tags:      []string{"shoes", "clothes", "puma"},
	Available: true,
}

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
