package tests

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/iurikman/smartSurvey/internal/models"
	server "github.com/iurikman/smartSurvey/internal/rest"
)

func (s *IntegrationTestSuite) TestUsers() {
	company := models.Company{
		ID:   s.testCompanyID,
		Name: "testCompanyPost",
	}

	rCompany := new(models.Company)
	s.sendRequest(
		context.Background(),
		http.MethodPost,
		companiesEndpoint,
		company,
		&server.HTTPResponse{Data: &rCompany},
	)

	user1 := models.User{
		ID:       uuid.New(),
		Company:  s.testCompanyID,
		Role:     "TestRole1",
		Name:     "TestName1",
		Surname:  newString("TestSurname1"),
		Phone:    newString("1991"),
		Email:    newString("testuser1@test.org"),
		UserType: newString("Type1"),
	}

	user2 := models.User{
		ID:       uuid.New(),
		Company:  s.testCompanyID,
		Role:     "TestRole2",
		Name:     "TestName2",
		Surname:  newString("TestSurname2"),
		Phone:    newString("2991"),
		Email:    newString("testuser2@test.org"),
		UserType: newString("Type2"),
	}

	user3 := models.User{
		ID:       uuid.New(),
		Company:  s.testCompanyID,
		Role:     "TestRole3",
		Name:     "TestName3",
		Surname:  newString("TestSurname3"),
		Phone:    newString("3991"),
		Email:    newString("testuser3@test.org"),
		UserType: newString("Type3"),
	}

	user4 := models.User{
		ID:       uuid.New(),
		Company:  s.testCompanyID,
		Role:     "TestRole4",
		Name:     "",
		Surname:  newString("TestSurname4"),
		Phone:    newString("4991"),
		Email:    newString("testuser4@test.org"),
		UserType: newString("Type4"),
	}

	user5 := models.User{
		ID:       uuid.New(),
		Company:  s.testCompanyID,
		Role:     "TestRole5",
		Name:     "TestName5",
		Surname:  newString("TestSurname5"),
		Phone:    newString("5991"),
		Email:    newString(""),
		UserType: newString("Type5"),
	}

	user6 := models.User{
		ID:       uuid.New(),
		Company:  s.testCompanyID,
		Role:     "TestRole6",
		Name:     "TestName6",
		Surname:  newString("TestSurname6"),
		Phone:    newString(""),
		Email:    newString("testuser6@test.org"),
		UserType: newString("Type6"),
	}

	s.Run("users", func() {
		s.Run("POST:/users", func() {
			s.Run("201/created", func() {
				s.checkUserPost(&user1)
				s.checkUserPost(&user2)
				s.checkUserPost(&user3)
			})

			s.Run("409/duplicate user", func() {
				respUserData := new(models.User)

				resp := s.sendRequest(
					context.Background(),
					http.MethodPost,
					usersEndpoint,
					user1,
					&server.HTTPResponse{Data: &respUserData},
				)
				s.Require().Equal(http.StatusConflict, resp.StatusCode)
			})

			s.Run("400/badRequest", func() {
				resp := s.sendRequest(
					context.Background(),
					http.MethodPost,
					usersEndpoint,
					"badRequest?",
					nil,
				)
				s.Require().Equal(http.StatusBadRequest, resp.StatusCode)

				s.Run("422/StatusUnprocessableEntity/nameIsEmpty", func() {
					resp := s.sendRequest(
						context.Background(),
						http.MethodPost,
						usersEndpoint,
						user4,
						nil,
					)
					s.Require().Equal(http.StatusUnprocessableEntity, resp.StatusCode)
				})

				s.Run("422/StatusUnprocessableEntity/emailIsEmpty", func() {
					resp := s.sendRequest(
						context.Background(),
						http.MethodPost,
						usersEndpoint,
						user5,
						nil,
					)
					s.Require().Equal(http.StatusUnprocessableEntity, resp.StatusCode)
				})

				s.Run("422/StatusUnprocessableEntity/phoneIsEmpty", func() {
					resp := s.sendRequest(
						context.Background(),
						http.MethodPost,
						usersEndpoint,
						user6,
						nil,
					)
					s.Require().Equal(http.StatusUnprocessableEntity, resp.StatusCode)
				})
			})

			s.Run("GET:/users/{id}", func() {
				respUserData := new(models.User)
				resp := s.sendRequest(
					context.Background(),
					http.MethodGet,
					usersEndpoint+"/"+user3.ID.String(),
					nil,
					&server.HTTPResponse{Data: &respUserData},
				)
				s.Require().Equal(http.StatusOK, resp.StatusCode)
				s.Require().Equal(user3.ID, respUserData.ID)
				s.Require().Equal(user3.Company, respUserData.Company)
				s.Require().Equal(user3.Role, respUserData.Role)
				s.Require().Equal(user3.Name, respUserData.Name)
				s.Require().Equal(user3.Surname, respUserData.Surname)
				s.Require().Equal(user3.Phone, respUserData.Phone)
				s.Require().Equal(user3.Email, respUserData.Email)
				s.Require().Equal(user3.UserType, respUserData.UserType)
			})

			s.Run("PATCH:/users", func() {
				s.Run("404/not found", func() {
					var respUserData models.User
					resp := s.sendRequest(
						context.Background(),
						http.MethodPatch,
						usersEndpoint+"/"+uuid.New().String(),
						user3,
						&server.HTTPResponse{Data: &respUserData},
					)
					s.Require().Equal(http.StatusNotFound, resp.StatusCode)
				})
			})

			s.Run("DELETE:/users", func() {
				s.Run("404/not found", func() {
					var respUserData models.User
					resp := s.sendRequest(
						context.Background(),
						http.MethodDelete,
						usersEndpoint+"/"+uuid.New().String(),
						user3,
						&server.HTTPResponse{Data: &respUserData},
					)
					s.Require().Equal(http.StatusNotFound, resp.StatusCode)
				})
			})
		})
	})
}

func newString(value string) *string {
	return &value
}
