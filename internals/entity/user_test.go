package entity

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserTestSuite struct {
	suite.Suite
}

func (s *UserTestSuite) TestNewUser() {
	username := "test"
	password := "Test@123@@"

	result, err := NewUser(username, password)

	s.Nil(err)
	s.NotNil(result)
	s.NotNil(result.Id)
	s.Equal(username, result.Username)
	s.NotEqual(password, result.Password)
}

func (s *UserTestSuite) TestNewUserWithInvalidUsername() {
	usernames := []string{
		"tet",
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus quis accumsan neque, in sollicitudin massa.",
	}
	password := "Test@123@@"

	for _, username := range usernames {
		result, err := NewUser(username, password)

		s.NotNil(err)
		s.Nil(result)
	}

}

func (s *UserTestSuite) TestNewUserWithInvalidPassword() {
	username := "test"
	password := []string{"12312131313", "test@@123@@", "Test@@@@@@", "Test123431938", "t1@A"}
	for _, p := range password {
		result, err := NewUser(username, p)

		s.NotNil(err)
		s.Nil(result)
	}
}

func TestUserTestSuit(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
