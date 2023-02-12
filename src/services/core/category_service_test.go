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

type CategoryServiceSuite struct {
	suite.Suite
	service *categoryService
}

func TestCategoryServiceSuite(t *testing.T) {
	suite.Run(t, new(CategoryServiceSuite))
}

func (s *CategoryServiceSuite) SetupSuite() {
	s.T().Logf("\n-------------- init ---------------")

	prodRepo := product.NewMemoryProductRepository(
		utils.ProductExp2, //headsets
	)

	catRepo := category.NewMemoryCategoryRepository(
		utils.CategoryExp1, // headsets
		utils.CategoryExp2,
	)

	s.service = &categoryService{
		catRepo:  catRepo,
		prodRepo: prodRepo,
	}
}

func (s *CategoryServiceSuite) TestCategoryService_Create() {
	testCases := []struct {
		desc    string
		input   domain.Category
		wantErr bool
	}{
		{
			desc: "error: already exist",
			input: domain.Category{
				Name: "clothes",
			},
			wantErr: true,
		},
		{
			desc: "proper work",
			input: domain.Category{
				Name: "keyboards",
			},
			wantErr: false,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			err := s.service.Create(&tC.input)

			if err != nil {
				s.True(tC.wantErr, "error should be expected")
				s.T().Logf("\n\n Error >>> %s \n\n", err.Error())
				return
			}

			s.NotZero(tC.input.ID, "should not be zero")

			utils.PrettyPrintTesting(s.T(), tC.input)
		})
	}
}

func (s *CategoryServiceSuite) TestCategoryService_Delete() {
	testCases := []struct {
		desc    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			desc:    "error: not found",
			id:      uuid.New(),
			wantErr: true,
		},
		{
			desc:    "error: category with products",
			id:      utils.CategoryExp1.ID,
			wantErr: true,
		},
		{
			desc:    "proper work",
			id:      utils.CategoryExp2.ID,
			wantErr: false,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			err := s.service.Delete(tC.id)

			if err != nil {
				s.T().Logf("\n\n Error >>> %v \n\n", err)
			}

			s.Equal(tC.wantErr, (err != nil), "expect error fail")
		})
	}
}

func (s *CategoryServiceSuite) TestCategoryService_GetAll() {
	cs, err := s.service.GetAll()

	s.NoError(err, "should not throw error")

	s.Greater(len(cs), 0, "should contain at least one category")

	utils.PrettyPrintTesting(s.T(), cs)
}

func (s *CategoryServiceSuite) TestCategoryService_GetByID() {
	testCases := []struct {
		desc     string
		id       uuid.UUID
		wantErr  bool
		wantProd bool
	}{
		{
			desc:     "proper work",
			id:       utils.CategoryExp1.ID,
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

			s.Equal(tC.wantErr, (err != nil), "expect error fail")

			s.Equal(tC.wantProd, (got != nil), "expect category fail")
		})
	}
}

func (s *CategoryServiceSuite) TestCategoryService_GetByName() {
	testCases := []struct {
		desc     string
		name     string
		wantErr  bool
		wantProd bool
	}{
		{
			desc:     "proper work",
			name:     "headsets",
			wantErr:  false,
			wantProd: true,
		},
		{
			desc:     "not found",
			name:     "frutas",
			wantErr:  false,
			wantProd: false,
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			got, err := s.service.GetByName(tC.name)

			s.Equal(tC.wantErr, (err != nil), "expect error fail")

			s.Equal(tC.wantProd, (got != nil), "expect category fail")
		})
	}
}
