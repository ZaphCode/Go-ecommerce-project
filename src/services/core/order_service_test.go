package core

import (
	"testing"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/address"
	"github.com/ZaphCode/clean-arch/src/repositories/order"
	"github.com/stretchr/testify/suite"
)

type OrderServiceSuite struct {
	suite.Suite
	service *orderService
}

func TestOrderServiceSuite(t *testing.T) {
	suite.Run(t, new(OrderServiceSuite))
}

func (s *OrderServiceSuite) SetupSuite() {
	s.T().Logf("\n-------------- init ---------------")

	ordRepo := order.NewMemoryOrderRepository()
	addrRepo := address.NewMemoryAddressRepository()

	s.service = &orderService{
		ordRepo:  ordRepo,
		addrRepo: addrRepo,
	}
}

//* Tests

func (s *OrderServiceSuite) TestCreateOrder() {
	ord := domain.Order{}
	s.service.Create(&ord)
}
