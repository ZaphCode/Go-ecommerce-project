package domain

import "github.com/ZaphCode/clean-arch/src/utils"

//* Model

type PaymentMethod struct {
	utils.DBModel
}

//* Service

type PaymentMethodService interface {
	utils.ServiceCrudOperations[PaymentMethod]
}

//* Repository

type PaymentMethodRepository interface {
	utils.RepositoryCrudOperations[PaymentMethod]
}
