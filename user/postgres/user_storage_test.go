package user_postgres_test

import (
	"context"
	"database/sql"
	"github.com/rubemlrm/go-api-bootstrap/config"
	"github.com/rubemlrm/go-api-bootstrap/pkg/logger"
	"testing"
	"time"

	gooseTesting "github.com/rubemlrm/go-api-bootstrap/tests/goose"
	"github.com/rubemlrm/go-api-bootstrap/tests/testcontainers"
	user "github.com/rubemlrm/go-api-bootstrap/user"
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

func (s *UserRepositoryTestSuite) TestUserCreationWithSucess() {
	u := user.UserCreate{
		Name:     "teste2",
		Password: "changeme",
		Email:    "foo",
	}
	logger := logger.NewLogger(config.Logger{
		Level: "Debug",
	})
	repository := user_postgres.NewConnection(s.DB, logger)
	id, err := repository.Create(&u)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), id, user.ID(1))
}

func (s *UserRepositoryTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	s.Require().NoError(s.testcontainer.Terminate(ctx))
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
