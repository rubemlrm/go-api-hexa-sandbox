//go:build integration

package user_postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	gooseTesting "github.com/rubemlrm/go-api-bootstrap/tests/goose"
	"github.com/rubemlrm/go-api-bootstrap/tests/testcontainers"
	"github.com/rubemlrm/go-api-bootstrap/user"
	storage "github.com/rubemlrm/go-api-bootstrap/user/postgres"

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

	fakerData := user.User{}
	err := faker.FakeData(&fakerData)
	s.Require().NoError(err)

	repository := storage.NewConnection(s.DB)
	id, err := repository.Create(&fakerData)

	assert.Equal(s.T(), 2, id)
	assert.NoError(s.T(), err)
}

func (s *UserRepositoryTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	s.Require().NoError(s.testcontainer.Terminate(ctx))
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
