package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/guirialli/rater_limit/internals/entity/dtos"
	"github.com/guirialli/rater_limit/internals/infra/database"
	"github.com/guirialli/rater_limit/internals/usecases"
	"github.com/guirialli/rater_limit/test/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type SuiteAuthTest struct {
	suite.Suite
	auth    *Auth
	useCase usecases.IUser
	init    string
}

func (s *SuiteAuthTest) SetupSuite() {
	s.init = "../../../../test/database/init.sql"
}

func (s *SuiteAuthTest) SetupTest() {
	fmt.Println("---------------------------------")
	db, _ := database.NewSqlite("file::memory:?cache=shared").InitDatabaseGetConnection(s.init)
	user, _ := usecases.NewUser("12", 10, 's')
	s.useCase = user
	s.auth = NewAuth(db, user)
}

func (s *SuiteAuthTest) TearDownTest() {
	fmt.Println("---------------------------------")
}

func (s *SuiteAuthTest) TestSetup() {
	s.NotNil(s.auth)
}

func (s *SuiteAuthTest) mockUser() *dtos.RegisterForm {
	user := &dtos.RegisterForm{
		Username: "test1",
		Password: "TestA!@#1@@",
	}
	_, err := s.useCase.Register(context.Background(), s.auth.db, user)
	if err != nil {
		panic(err)
	}
	return user
}

func (s *SuiteAuthTest) TestRegister() {
	var registeredUser dtos.ResponseJson[dtos.ResponseJwt]
	body, _ := json.Marshal(mock.NewUserMock().RegisterForm())
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.auth.Register(w, req)

	s.Equal(http.StatusCreated, w.Code)
	err := json.NewDecoder(w.Body).Decode(&registeredUser)
	s.NoError(err)
	s.NotEmpty(registeredUser.Data.Token)
}

func (s *SuiteAuthTest) TestRegisterInvalid() {
	fails := []dtos.RegisterForm{{
		Username: "te", // invalid username
		Password: "Test1231@@sdas",
	}, {
		Username: "test",
		Password: "Te", // invalid password
	}}
	status := []int{http.StatusBadRequest, http.StatusBadRequest, http.StatusBadRequest}

	for i, st := range status {
		var body []byte
		if i < len(fails) {
			body, _ = json.Marshal(fails[i])
		}
		req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		s.auth.Register(w, req)

		s.Equal(st, w.Code)
	}
}

func (s *SuiteAuthTest) TestLogin() {
	var longedUser dtos.ResponseJson[dtos.ResponseJwt]
	form := s.mockUser()
	body, _ := json.Marshal(dtos.LoginForm{
		Username: form.Username,
		Password: form.Password,
	})
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.auth.Login(w, req)

	s.Equal(http.StatusOK, w.Code)
	err := json.NewDecoder(w.Body).Decode(&longedUser)
	s.NoError(err)
	s.NotEmpty(longedUser.Data.Token)
}

func (s *SuiteAuthTest) TestLoginInvalid() {
	longedUser := s.mockUser()
	fails := []dtos.LoginForm{{
		Username: longedUser.Username + "1",
		Password: longedUser.Password,
	}, {
		Username: longedUser.Username,
		Password: longedUser.Password + "1",
	}}
	status := []int{http.StatusUnauthorized, http.StatusUnauthorized, http.StatusBadRequest}

	for i, st := range status {
		var body []byte
		if i < len(fails) {
			body, _ = json.Marshal(fails[i])
		}
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		s.auth.Login(w, req)

		s.Equal(st, w.Code)
	}
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(SuiteAuthTest))
}
