package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/rubemlrm/go-api-bootstrap/config"
	"github.com/rubemlrm/go-api-bootstrap/pkg/logger"
	"github.com/rubemlrm/go-api-bootstrap/user/factories"

	gooseTesting "github.com/rubemlrm/go-api-bootstrap/tests/goose"
	"github.com/rubemlrm/go-api-bootstrap/tests/testcontainers"
	"github.com/rubemlrm/go-api-bootstrap/user"
	userpostgres "github.com/rubemlrm/go-api-bootstrap/user/postgres"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	testcontainer *testcontainers.PostgresTestContainer
	suite.Suite
	DB *sql.DB
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	psqlContainer, err := testcontainers.StartPostgresContainer(ctx)
	s.Require().NoError(err)

	s.testcontainer = psqlContainer

	err = gooseTesting.RunMigrations(s.testcontainer.DSN)
	s.Require().NoError(err)

	s.DB, err = sql.Open("postgres", psqlContainer.DSN)
	s.Require().NoError(err)
}

func (s *UserRepositoryTestSuite) TestUserCreation() {
	uf := factories.UserFactory{}
	tests := []struct {
		name             string
		user             user.UserCreate
		expectedError    bool
		want             string
		connectionString *sql.DB
	}{
		{
			name:             "Create user with success",
			user:             *uf.CreateUserCreate(),
			expectedError:    false,
			connectionString: s.DB,
		},
		{
			name:             "Simulate error on user creation",
			expectedError:    true,
			user:             *uf.CreateUserCreate(),
			connectionString: s.generateWrongSQLConnection(),
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			lg := logger.NewLogger(config.Logger{
				Level: "Debug",
			})
			ctx := context.Background()
			repository := userpostgres.NewConnection(tt.connectionString, lg)
			id, err := repository.Create(ctx, &tt.user)
			if tt.expectedError {
				assert.Error(s.T(), err)
			} else {
				assert.Equal(s.T(), id, user.ID(1))
				assert.NoError(s.T(), err)
			}
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
		{
			name:          "Get user list without errors",
			requiredSeed:  true,
			expectedTotal: 5,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			lg := logger.NewLogger(config.Logger{
				Level: "Debug",
			})
			ctx := context.Background()
			repository := userpostgres.NewConnection(s.DB, lg)

			if tt.requiredSeed == true {
				fu := factories.GenerateUsers(tt.expectedTotal)
				err := factories.GenerateUsersOnDB(s.DB, fu)
				assert.NoError(t, err)
			}

			uu, err := repository.All(ctx)

			if tt.expectedError != nil {
				assert.Error(s.T(), tt.expectedError, err)
			}
			assert.Equal(s.T(), tt.expectedTotal, len(*uu))
		})
	}
}

func (s *UserRepositoryTestSuite) TestUserGet() {
	tests := []struct {
		name          string
		requiredSeed  bool
		expectedError bool
		wantError     error
		fakeUserID    bool
	}{
		{
			name:          "Get user with success",
			requiredSeed:  true,
			expectedError: false,
			fakeUserID:    false,
		},
		{
			name:          "User not found",
			requiredSeed:  true,
			expectedError: true,
			wantError:     fmt.Errorf("not found result"),
			fakeUserID:    true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			userID := user.ID(9999)
			lg := logger.NewLogger(config.Logger{
				Level: "Debug",
			})
			ctx := context.Background()
			repository := userpostgres.NewConnection(s.DB, lg)
			var uu []user.User
			if tt.requiredSeed == true {
				fu := factories.GenerateUsers(10)
				err := factories.GenerateUsersOnDB(s.DB, fu)
				assert.NoError(t, err)
				uu = fu
				if !tt.fakeUserID {
					userID = uu[0].ID
				}
			}

			u, err := repository.Get(ctx, userID)

			if tt.expectedError {
				assert.Error(s.T(), tt.wantError, err)
				assert.Nil(s.T(), u)
				return
			}
			assert.Equal(s.T(), u.ID, uu[0].ID)
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

func (s *UserRepositoryTestSuite) generateWrongSQLConnection() *sql.DB {
	db, _ := sql.Open("postgres", "wrong connection string")
	return db
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
