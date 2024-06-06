package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/rubemlrm/go-api-bootstrap/config"
	"github.com/rubemlrm/go-api-bootstrap/pkg/logger"
	"github.com/rubemlrm/go-api-bootstrap/user/factories"

	gooseTesting "github.com/rubemlrm/go-api-bootstrap/tests/goose"
	"github.com/rubemlrm/go-api-bootstrap/tests/testcontainers"
	"github.com/rubemlrm/go-api-bootstrap/user"
	user_postgres "github.com/rubemlrm/go-api-bootstrap/user/postgres"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	testcontainer *testcontainers.TestContainer
	suite.Suite
	DB *sql.DB
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	psqlContainer, err := testcontainers.StartContainer(ctx)
	s.Require().NoError(err)

	s.testcontainer = psqlContainer

	err = gooseTesting.RunMigrations(s.testcontainer.DSN)
	s.Require().NoError(err)

	s.DB, err = sql.Open("postgres", psqlContainer.DSN)
	s.Require().NoError(err)
}

func (s *UserRepositoryTestSuite) TestUserGet() {
	s.T().Run("get user with success", func(t *testing.T) {
		lg := logger.NewLogger(config.Logger{
			Level: "Debug",
		})
		uu := factories.GenerateUsers(10)
		ux := uu[0]
		err := factories.GenerateUsersOnDB(s.DB, uu)
		assert.NoError(t, err)

		repository := user_postgres.NewConnection(s.DB, lg)
		u, err := repository.Get(ux.ID)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), u.Name, ux.Name)
	})
}

func (s *UserRepositoryTestSuite) TestUserCreation() {
	uf := factories.UserFactory{}
	tests := []struct {
		name          string
		user          user.UserCreate
		expectedError error
	}{
		{
			name:          "create user with success",
			user:          *uf.CreateUserCreate(),
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			lg := logger.NewLogger(config.Logger{
				Level: "Debug",
			})
			repository := user_postgres.NewConnection(s.DB, lg)
			id, err := repository.Create(&tt.user)
			if tt.expectedError != nil {
				assert.NoError(s.T(), err)
			}
			assert.Equal(s.T(), id, user.ID(1))
		})
	}
}

func (s *UserRepositoryTestSuite) TestUserList() {
	tests := []struct {
		name          string
		requiredSeed  bool
		expectedTotal int
		expectedError error
	}{
		{
			name:          "Get empty user list without errors",
			requiredSeed:  false,
			expectedTotal: 0,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			lg := logger.NewLogger(config.Logger{
				Level: "Debug",
			})
			repository := user_postgres.NewConnection(s.DB, lg)
			uu, err := repository.All()

			if tt.requiredSeed == true {
				fu := factories.GenerateUsers(tt.expectedTotal)
				err = factories.GenerateUsersOnDB(s.DB, fu)
				assert.NoError(t, err)
			}

			if tt.expectedError != nil {
				assert.Error(s.T(), tt.expectedError, err)
			}
			assert.Equal(s.T(), tt.expectedTotal, len(*uu))
		})
	}
}

func (s *UserRepositoryTestSuite) TearDownSuite() {
	err := gooseTesting.RollbackMigrations(s.testcontainer.DSN)
	if err != nil {
		panic(err)
	}
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	s.Require().NoError(s.testcontainer.Terminate(ctx))
}

func (s *UserRepositoryTestSuite) TearDownTest() {
	_, err := s.DB.Exec("TRUNCATE table users")
	if err != nil {
		panic(err)
	}
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
