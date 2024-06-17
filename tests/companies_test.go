package tests

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/iurikman/smartSurvey/internal/models"
	server "github.com/iurikman/smartSurvey/internal/rest"
)

func (s *IntegrationTestSuite) TestCompanies() {
	rCompany := new(models.Company)
	company1 := models.Company{
		ID:   s.testCompanyID,
		Name: "testCompanyPost1",
	}
	company2 := models.Company{
		ID:   uuid.New(),
		Name: "testCompanyPost2",
	}
	company3 := models.Company{
		ID:   uuid.New(),
		Name: "",
	}
	s.Run("companies", func() {
		s.Run("POST:/companies", func() {
			s.Run("201/created", func() {
				resp := s.sendRequest(
					context.Background(),
					http.MethodPost,
					companiesEndpoint,
					company1,
					&server.HTTPResponse{Data: &rCompany},
				)
				s.Require().Equal(http.StatusCreated, resp.StatusCode)
				s.testCompanyID = rCompany.ID
			})

			s.Run("400/badRequest/companyNameIsEmpty", func() {
				resp := s.sendRequest(
					context.Background(),
					http.MethodPost,
					companiesEndpoint,
					company3,
					nil,
				)
				s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
			})

			s.Run("409/duplicate company", func() {
				resp := s.sendRequest(
					context.Background(),
					http.MethodPost,
					companiesEndpoint,
					company1,
					&server.HTTPResponse{Data: &rCompany},
				)
				s.Require().Equal(http.StatusConflict, resp.StatusCode)
			})
		})

		s.Run("PATCH:/companies", func() {
			s.Run("404/userNotFound", func() {
				resp := s.sendRequest(
					context.Background(),
					http.MethodPatch,
					companiesEndpoint+"/"+uuid.New().String(),
					company2,
					&server.HTTPResponse{Data: &rCompany},
				)
				s.Require().Equal(http.StatusNotFound, resp.StatusCode)
			})

			s.Run("400/badRequest/companyNameIsEmpty", func() {
				resp := s.sendRequest(
					context.Background(),
					http.MethodPost,
					usersEndpoint,
					company3,
					nil,
				)
				s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
			})
		})
	})
}
