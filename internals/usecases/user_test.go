package usecases

import (
	"context"
	"github.com/guirialli/rater_limit/internals/infra/database"
	"github.com/guirialli/rater_limit/test/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserTestSuite struct {
	suite.Suite
	userCase *User
	db       *database.Sqlite
	fileInit string
}

func (s *UserTestSuite) SetupTest() {
	userCase, err := NewUser("12", 1, 's')
	if err != nil {
		panic(err)
	}
	s.userCase = userCase
	s.db = database.NewSqlite("file::memory:?cache=shared")
	s.fileInit = "../../test/database/init.sql"
}

func (s *UserTestSuite) TestNewUser() {
	runes := []rune{'s', 'h', 'd'}
	for _, u := range runes {
		user, err := NewUser("12", 1, u)
		s.NoError(err)
		s.NotNil(user)
	}
}

func (s *UserTestSuite) TestNewUserInvalidUnitTime() {
	user, err := NewUser("12", 1, 'a')
	s.Error(err)
	s.Nil(user)
}

func (s *UserTestSuite) TestNewTokenAuth() {
	auth := s.userCase.NewTokenAuth()
	s.NotNil(auth)
}

func (s *UserTestSuite) TestLogin() {
	db, _ := s.db.InitDatabaseGetConnection(s.fileInit)
	ctx := context.Background()
	s.userCase.Register(ctx, db, mock.NewUserMock().RegisterForm())

	result, err := s.userCase.Login(ctx, db, mock.NewUserMock().LoginForm())

	s.NoError(err)
	s.NotNil(result)
	s.NotEmpty(result)
}

func (s *UserTestSuite) TestRegister() {
	db, _ := s.db.InitDatabaseGetConnection(s.fileInit)
	form := mock.NewUserMock().RegisterForm()

	result, err := s.userCase.Register(context.Background(), db, form)

	s.Nil(err)
	s.NotNil(result)
	s.NotEmpty(result)
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
