package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"

	"github.com/iurikman/smartSurvey/internal/models"

	"github.com/iurikman/smartSurvey/internal/config"
	server "github.com/iurikman/smartSurvey/internal/rest"
	"github.com/iurikman/smartSurvey/internal/service"
	"github.com/iurikman/smartSurvey/internal/store"
	_ "github.com/jackc/pgx/v5/stdlib"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/suite"
)

const (
	bindAddr          = "http://localhost:8080/api"
	usersEndpoint     = "/users"
	companiesEndpoint = "/companies"
)

type IntegrationTestSuite struct {
	suite.Suite
	ctx           context.Context
	cfg           config.Config
	pgStore       *store.Postgres
	serviceLayer  service.Service
	serverOne     server.Server
	authToken     string
	testCompanyID uuid.UUID
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	ctx := context.Background()

	cfg := config.New()

	pgStore, err := store.New(ctx, store.Config{
		PGUser:     cfg.PGUser,
		PGPassword: cfg.PGPassword,
		PGHost:     cfg.PGHost,
		PGPort:     cfg.PGPort,
		PGDatabase: cfg.PGDatabase,
	})
	s.Require().NoError(err)

	serviceLayer := service.New(pgStore)

	serverOne := server.NewServer(
		server.Config{BindAddress: cfg.BindAddress},
		serviceLayer)

	err = pgStore.Migrate(migrate.Up)
	s.Require().NoError(err)

	err = pgStore.Truncate(ctx, "users", "companies")
	s.Require().NoError(err)

	go func() {
		err = serverOne.Start(ctx)
		s.Require().NoError(err)
	}()
}

func (s *IntegrationTestSuite) sendRequest(ctx context.Context, method, endpoint string, body interface{}, dest interface{}) *http.Response {
	s.T().Helper()

	reqBody, err := json.Marshal(body)
	s.Require().NoError(err)

	req, err := http.NewRequestWithContext(ctx, method, bindAddr+endpoint, bytes.NewReader(reqBody))
	s.Require().NoError(err)

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)

	defer func() {
		err = resp.Body.Close()
		s.Require().NoError(err)
	}()

	if dest != nil {
		err = json.NewDecoder(resp.Body).Decode(&dest)
		s.Require().NoError(err)
	}

	return resp
}

func (s *IntegrationTestSuite) checkUserPost(user *models.User) {
	respUserData := new(models.User)

	resp := s.sendRequest(
		context.Background(),
		http.MethodPost,
		usersEndpoint,
		user,
		&server.HTTPResponse{Data: &respUserData},
	)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)
	s.Require().Equal(user.Company, respUserData.Company)
	s.Require().Equal(user.Role, respUserData.Role)
	s.Require().Equal(user.Name, respUserData.Name)
	s.Require().Equal(user.Surname, respUserData.Surname)
	s.Require().Equal(user.Phone, respUserData.Phone)
	s.Require().Equal(user.Email, respUserData.Email)
	s.Require().Equal(user.UserType, respUserData.UserType)
	s.Require().NotZero(respUserData.ID)
	user.ID = respUserData.ID
}

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

				s.Run("400/badRequest/nameIsEmpty", func() {
					resp := s.sendRequest(
						context.Background(),
						http.MethodPost,
						usersEndpoint,
						user4,
						nil,
					)
					s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
				})

				s.Run("400/badRequest/emailIsEmpty", func() {
					resp := s.sendRequest(
						context.Background(),
						http.MethodPost,
						usersEndpoint,
						user5,
						nil,
					)
					s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
				})

				s.Run("400/badRequest/phoneIsEmpty", func() {
					resp := s.sendRequest(
						context.Background(),
						http.MethodPost,
						usersEndpoint,
						user6,
						nil,
					)
					s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
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
