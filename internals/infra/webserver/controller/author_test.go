package controller

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/guirialli/rater_limit/internals/entity"
	"github.com/guirialli/rater_limit/internals/entity/dtos"
	"github.com/guirialli/rater_limit/internals/infra/database"
	"github.com/guirialli/rater_limit/internals/usecases"
	"github.com/guirialli/rater_limit/test/mock"
	"github.com/stretchr/testify/suite"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

type SuiteAuthorTest struct {
	suite.Suite
	author      *Author
	useCase     usecases.IAuthor
	bookUseCase usecases.IBook
	init        string
	db          *sql.DB
}

func (s *SuiteAuthorTest) create() *entity.Author {
	author, err := s.useCase.Create(context.Background(), s.db, mock.NewAuthor().Create(nil))
	if err != nil {
		panic(err)
	}
	return author
}

func (s *SuiteAuthorTest) SetupSuite() {
	s.init = "../../../../test/database/init.sql"
}

func (s *SuiteAuthorTest) SetupTest() {
	db, _ := database.NewSqlite("file::memory:?cache=shared").InitDatabaseGetConnection(s.init)
	author := usecases.NewAuthor()
	book := usecases.NewBook()

	s.db = db
	s.bookUseCase = book
	s.useCase = author
	s.author = NewAuthor(db, author, book)
}

func (s *SuiteAuthorTest) TestSetup() {
	s.NotNil(s.db)
	s.NotNil(s.author)
	s.NotNil(s.bookUseCase)
	s.NotNil(s.author)
}

func (s *SuiteAuthorTest) TestGetAllBook() {
	length := rand.IntN(1000) + 1
	for i := 0; i < length; i++ {
		s.create()
	}
	var authors dtos.ResponseJson[[]entity.Author]
	req, _ := http.NewRequest("GET", "/authors", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.author.GetAll(w, req)

	s.Equal(w.Code, http.StatusOK)
	err := json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&authors)
	s.NoError(err)
	s.Len(authors.Data, length)
}

func (s *SuiteAuthorTest) TestGetBook() {
	req, _ := http.NewRequest("GET", "/authors/{id}", nil)
	req.Header.Set("Content-Type", "application/json")
	author := s.create()
	w := httptest.NewRecorder()
	rCtx := chi.NewRouteContext()
	rCtx.URLParams.Add("id", author.Id)

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rCtx))
	var result dtos.ResponseJson[entity.Author]
	s.author.GetById(w, req)

	s.Equal(w.Code, http.StatusOK)
	err := json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&result)
	s.NoError(err)
	s.NotNil(result.Data)
	s.Equal(result.Data.Id, author.Id)
	s.Equal(result.Data.Name, author.Name)
}

func (s *SuiteAuthorTest) TestCreate() {
	var response dtos.ResponseJson[entity.Author]
	body, _ := json.Marshal(mock.NewAuthor().Create(nil))
	req, _ := http.NewRequest("POST", "/authors", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.author.Create(w, req)

	s.Equal(w.Code, http.StatusCreated)
	err := json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&response)
	s.NoError(err)
	s.NotNil(response.Data)
	s.NotEmpty(response.Data.Id)
}

func TestAuthorSuite(t *testing.T) {
	suite.Run(t, new(SuiteAuthorTest))
}
