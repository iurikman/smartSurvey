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
	bindAddr          = "http://localhost:8080/api/v1"
	usersEndpoint     = "/users"
	companiesEndpoint = "/companies"
)

type IntegrationTestSuite struct {
	suite.Suite
	cancel        context.CancelFunc
	pgStore       *store.Postgres
	service       *service.Service
	server        *server.Server
	authToken     string
	testCompanyID uuid.UUID
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	cfg := config.New()

	pgStore, err := store.New(ctx, store.Config{
		PGUser:     cfg.PGUser,
		PGPassword: cfg.PGPassword,
		PGHost:     cfg.PGHost,
		PGPort:     cfg.PGPort,
		PGDatabase: cfg.PGDatabase,
	})
	s.Require().NoError(err)

	s.pgStore = pgStore

	s.service = service.New(pgStore)

	s.server = server.NewServer(
		server.Config{BindAddress: cfg.BindAddress},
		s.service)

	err = pgStore.Migrate(migrate.Up)
	s.Require().NoError(err)

	err = pgStore.Truncate(ctx, "users", "companies")
	s.Require().NoError(err)

	go func() {
		err = s.server.Start(ctx)
		s.Require().NoError(err)
	}()
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.cancel()
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
