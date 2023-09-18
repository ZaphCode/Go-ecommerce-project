package core

import (
	"testing"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/repositories/category"
	"github.com/ZaphCode/clean-arch/src/repositories/product"
	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ProductServiceSuite struct {
	suite.Suite
	service *prodService
}

func TestProductServiceSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceSuite))
}

func (s *ProductServiceSuite) SetupSuite() {
	s.T().Logf("\n-------------- init ---------------")

	prodRepo := product.NewMemoryProductRepository(
		utils.ProductExp1,
		utils.ProductExp2,
	)

	catRepo := category.NewMemoryCategoryRepository(
		utils.CategoryExp1,
		utils.CategoryExp2,
	)

	s.service = &prodService{
		prodRepo: prodRepo,
		catRepo:  catRepo,
	}
}

func (s *ProductServiceSuite) TestProductService_Create() {
	testCases := []struct {
		desc    string
		wantErr bool
		input   domain.Product
	}{
		{
			desc:    "proper work",
			wantErr: false,
			input: domain.Product{
				Category:     "clothes",
				Name:         "Blue pants",
				Description:  "incredible pants",
				Price:        2424,
				DiscountRate: 13,
				Tags:         []string{"blue", "pants", "levis"},
				Avalible:     true,
				ImagesUrl:    []string{"https://levis.com/bluepants.jpg"},
			},
		},
		{
			desc:    "category not found",
			wantErr: true,
			input: domain.Product{
				Category:     "ropa",
				Name:         "pantalones azules",
				Description:  "increibles pantalones",
				Price:        43,
				DiscountRate: 13,
				Tags:         []string{"blue", "pants", "levis"},
				Avalible:     true,
				ImagesUrl:    []string{"https://levis.com/bluepants.jpg"},
			},
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			err := s.service.Create(&tC.input)

			s.Equal(tC.wantErr, (err != nil), "expect error fail")

			if err != nil {
				s.T().Logf("\n\n Error >>> %s \n\n", err.Error())
				return
			}

			s.NotZero(tC.input.ID)

			utils.PrettyPrintTesting(s.T(), tC.input)
		})
	}
}

func (s *ProductServiceSuite) TestProductService_GetAll() {
	ps, err := s.service.GetAll()

	s.NoError(err)

	utils.PrettyPrintTesting(s.T(), ps)
}

func (s *ProductServiceSuite) TestProductService_GetByID() {
	testCases := []struct {
		desc     string
		id       uuid.UUID
		wantErr  bool
		wantProd bool
	}{
		{
			desc:     "proper work",
			id:       utils.ProductExp2.ID,
			wantErr:  false,
			wantProd: true,
		},
		{
			desc:     "not found",
			id:       uuid.New(),
			wantErr:  false,
			wantProd: false,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			got, err := s.service.GetByID(tC.id)

			s.Equal(tC.wantErr, (err != nil), "expert error fail")

			s.Equal(tC.wantProd, (got != nil), "expert product fail")

		})
	}
}

func (s *ProductServiceSuite) TestProductService_GetByCategory() {
	testCases := []struct {
		desc      string
		category  string
		wantErr   bool
		wantProds bool
	}{
		{
			desc:      "proper work",
			category:  "clothes",
			wantErr:   false,
			wantProds: true,
		},
		{
			desc:      "not found",
			category:  "ropa",
			wantErr:   false,
			wantProds: false,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			got, err := s.service.GetByCategory(tC.category)

			s.Equal(tC.wantErr, (err != nil), "expert error fail")

			s.Equal(tC.wantProds, (len(got) > 0), "expert product fail")

			utils.PrettyPrintTesting(s.T(), got)
		})
	}
}

func (s *ProductServiceSuite) TestProductService_GetLatests() {
	ps, err := s.service.GetLatestProds(1)

	s.NoError(err, "should not throw error")

	s.Len(ps, 1, "should contain only one product")

	s.Equal("Corsair void pro", ps[0].Name, "should be the last prod created")

	utils.PrettyPrintTesting(s.T(), ps)
}

func (s *ProductServiceSuite) TestProductService_Delete() {
	testCases := []struct {
		desc    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			desc:    "Not found",
			wantErr: true,
			id:      uuid.New(),
		},
		{
			desc:    "Proper work",
			wantErr: false,
			id:      utils.ProductExp1.ID,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			err := s.service.Delete(tC.id)

			s.Equal(tC.wantErr, (err != nil), "expert error fail")
		})
	}
}

func (s *ProductServiceSuite) TestProductService_Update() {
	testCases := []struct {
		desc         string
		id           uuid.UUID
		uf           domain.UpdateFields
		wantErr      bool
		validationFn func()
	}{
		{
			desc:    "prod not found",
			wantErr: true,
			id:      uuid.New(),
			uf:      domain.UpdateFields{},
		},
		{
			desc:    "nil update fields",
			wantErr: true,
			id:      utils.ProductExp2.ID,
			uf:      nil,
		},
		{
			desc:    "invalid category",
			wantErr: true,
			id:      utils.ProductExp2.ID,
			uf: domain.UpdateFields{
				"Category": "osos maduros",
				"Price":    2467,
			},
		},
		{
			desc:    "unexisting field",
			wantErr: true,
			id:      utils.ProductExp2.ID,
			uf: domain.UpdateFields{
				"Price": 2467,
				"Email": "test@test.com",
			},
		},
		{
			desc:    "invalid type",
			wantErr: true,
			id:      utils.ProductExp2.ID,
			uf: domain.UpdateFields{
				"Price": "wakawakeeheh",
				"Name":  []int{1, 2, 3},
			},
		},
		{
			desc:    "proper work",
			wantErr: false,
			id:      utils.ProductExp2.ID,
			uf: domain.UpdateFields{
				"Name":  "Corsair Void Pro masters",
				"Price": int64(7777),
			},
			validationFn: func() {
				p, _ := s.service.prodRepo.FindByID(utils.ProductExp2.ID)

				s.NotNil(p, "should exists")

				if p.Name != "Corsair Void Pro masters" {
					s.Fail("Should be updated")
				}

				if p.Price != 7777 {
					s.Fail("Should be updated")
				}
			},
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			err := s.service.Update(tC.id, tC.uf)

			s.Equal(tC.wantErr, (err != nil), "expert error fail")

			if err != nil {
				s.T().Logf("\n\n Error >>> %s \n\n", err.Error())
				return
			}

			if tC.validationFn != nil {
				tC.validationFn()
			}
		})
	}
}

func (s *ProductServiceSuite) TestProductService_SetAvalible() {
	testCases := []struct {
		desc         string
		id           uuid.UUID
		avalible     bool
		wantErr      bool
		validationFn func()
	}{
		{
			desc:    "prod not found",
			wantErr: true,
			id:      uuid.New(),
		},
		{
			desc:     "proper work",
			wantErr:  false,
			id:       utils.ProductExp2.ID,
			avalible: false,
			validationFn: func() {
				p, _ := s.service.prodRepo.FindByID(utils.ProductExp2.ID)

				s.NotNil(p, "should exists")

				if p.Avalible != false {
					s.Fail("Should be updated")
				}
			},
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			err := s.service.SetAvalible(tC.id, tC.avalible)

			s.Equal(tC.wantErr, (err != nil), "expert error fail")

			if err != nil {
				s.T().Logf("\n\n Error >>> %s \n\n", err.Error())
				return
			}

			if tC.validationFn != nil {
				tC.validationFn()
			}
		})
	}
}

func (s *ProductServiceSuite) TestProductService_CalculateTotalPrice() {
	testCases := []struct {
		desc       string
		input      []domain.OrderProduct
		wantOutput int64
		wantErr    bool
	}{
		{
			desc: "proper work: one prod",
			input: []domain.OrderProduct{
				{ID: utils.ProductExp1.ID, Quantity: 3},
			},
			wantOutput: 6192,
			wantErr:    false,
		},
		{
			desc: "proper work: two products",
			input: []domain.OrderProduct{
				{ID: utils.ProductExp1.ID, Quantity: 4},
				{ID: utils.ProductExp2.ID, Quantity: 7},
			},
			wantOutput: 46476,
			wantErr:    false,
		},
		{
			desc: "product that does not exist",
			input: []domain.OrderProduct{
				{ID: uuid.New(), Quantity: 4},
				{ID: utils.ProductExp2.ID, Quantity: 7},
			},
			wantOutput: 0,
			wantErr:    true,
		},
		{
			desc: "product that does not exist",
			input: []domain.OrderProduct{
				{ID: utils.AddrExp2.ID, Quantity: 4},
				{ID: utils.AddrExp1.ID, Quantity: 7},
			},
			wantOutput: 0,
			wantErr:    true,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			got, err := s.service.CalculateTotalPrice(tC.input)

			s.Equal(tC.wantErr, (err != nil), "bad :(")

			s.Less(tC.wantOutput-got, int64(2), "diference cannot be bigger")

			s.Greater(tC.wantOutput-got, int64(-2), "diference cannot be smaller")

			s.T().Logf("\n\n Total Price: %d\n\n", got)
		})
	}
}
