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

func (s *SuiteAuthorTest) createBook(authorId string) *entity.Book {
	book, err := s.bookUseCase.Create(context.Background(), s.db, mock.NewBookMock().CreateWithAuthor(authorId, nil))
	if err != nil {
		panic(err)
	}
	return book
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
	s.author = NewAuthor(db, author, book, NewUtils())
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
	var authors dtos.ResponseJson[[]dtos.ResponseAuthor]
	req, _ := http.NewRequest("GET", "/authors", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.author.GetAll(w, req)

	s.Equal(http.StatusOK, w.Code)
	err := json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&authors)
	s.NoError(err)
	s.Len(authors.Data, length)
}

func (s *SuiteAuthorTest) TestGetAllWithBooks() {
	author := s.create()
	book := s.createBook(author.Id)
	var result dtos.ResponseJson[[]dtos.AuthorWithBooks]

	req, _ := http.NewRequest("GET", "/authors/books", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.author.GetAllWithBooks(w, req)

	s.Equal(http.StatusOK, w.Code)
	err := json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&result)
	s.NoError(err)
	for _, ab := range result.Data {
		s.Equal(author.Id, *ab.Author.Id)
		s.Equal(book.Id, ab.Books[0].Id)
	}
}

func (s *SuiteAuthorTest) TestGetById() {
	author := s.create()
	status := []int{http.StatusOK, http.StatusBadRequest, http.StatusNotFound}
	ids := []string{author.Id, "", "123"}

	for i, id := range ids {
		req, _ := http.NewRequest("GET", "/authors/{id}", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		rCtx := chi.NewRouteContext()
		rCtx.URLParams.Add("id", id)

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rCtx))
		var result dtos.ResponseJson[dtos.ResponseAuthor]
		s.author.GetById(w, req)

		s.Equal(status[i], w.Code)
		if status[i] >= 200 && status[i] < 300 {
			err := json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&result)
			s.NoError(err)
			s.NotNil(result.Data)
			s.Equal(author.Id, *result.Data.Id)
			s.Equal(author.Name, *result.Data.Name)
		}

	}

}

func (s *SuiteAuthorTest) TestGetByIdWithBooks() {
	author := s.create()
	book := s.createBook(author.Id)
	status := []int{http.StatusOK, http.StatusBadRequest, http.StatusNotFound}
	ids := []string{author.Id, "", "123"}

	for i, id := range ids {
		req, _ := http.NewRequest("GET", "/authors/{id}/books/", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		rCtx := chi.NewRouteContext()
		rCtx.URLParams.Add("id", id)

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rCtx))
		var result dtos.ResponseJson[dtos.AuthorWithBooks]
		s.author.GetByIdWithBooks(w, req)

		s.Equal(status[i], w.Code)
		if status[i] >= 200 && status[i] < 300 {
			err := json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&result)
			s.NoError(err)
			s.NotNil(result.Data)
			s.Equal(author.Id, *result.Data.Author.Id)
			s.Equal(book.Id, result.Data.Books[0].Id)
		}
	}
}

func (s *SuiteAuthorTest) TestCreate() {
	status := []int{http.StatusCreated, http.StatusBadRequest}
	var response dtos.ResponseJson[dtos.ResponseAuthor]
	for i, st := range status {
		var body []byte
		if i == 0 {
			body, _ = json.Marshal(mock.NewAuthor().CreateBody(nil))
		}

		req, _ := http.NewRequest("POST", "/authors", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		s.author.Create(w, req)

		s.Equal(st, w.Code)
		if st >= 200 && st < 300 {
			err := json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&response)
			s.NoError(err)
			s.NotNil(response.Data)
			s.NotEmpty(response.Data.Id)
		}

	}

}

func (s *SuiteAuthorTest) TestUpdate() {
	author := s.create()
	var response dtos.ResponseJson[dtos.ResponseAuthor]
	newName := author.Name + "1"
	newBirthDay := author.Birthday.Add(1000).Format("2006-01-02")
	newDescription := "bla"
	update := dtos.AuthorBody{
		Name:        &newName,
		Birthday:    &newBirthDay,
		Description: &newDescription,
	}
	status := []int{http.StatusOK, http.StatusBadRequest, http.StatusNotFound, http.StatusBadRequest}
	bodyS := []interface{}{update, update, update, ""}
	ids := []string{author.Id, "", "123", author.Id}

	for i, id := range ids {
		body, _ := json.Marshal(bodyS[i])
		req, _ := http.NewRequest("PUT", "/authors/{id}", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		rCtx := chi.NewRouteContext()
		rCtx.URLParams.Add("id", id)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rCtx))

		s.author.Update(w, req)

		s.Equal(status[i], w.Code)
		if status[i] >= 200 && status[i] < 300 {
			err := json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&response)
			s.NoError(err)
			s.NotNil(response.Data)
			s.Equal(author.Id, *response.Data.Id)
			s.NotEqual(author.Name, response.Data.Name)
		}

	}

}

func (s *SuiteAuthorTest) TestPatch() {
	author := s.create()
	var response dtos.ResponseJson[dtos.ResponseAuthor]
	newName := author.Name + "2"
	newBirthDay := author.Birthday.Add(1000).Format("2006-01-02")
	newDescription := "bla"
	update := dtos.AuthorBody{
		Name:        &newName,
		Birthday:    &newBirthDay,
		Description: &newDescription,
	}
	status := []int{http.StatusOK, http.StatusBadRequest, http.StatusNotFound, http.StatusBadRequest}
	bodyS := []interface{}{update, update, update, ""}
	ids := []string{author.Id, "", "123", author.Id}

	for i, id := range ids {
		body, _ := json.Marshal(bodyS[i])
		req, _ := http.NewRequest("PATCH", "/authors/{id}", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		rCtx := chi.NewRouteContext()
		rCtx.URLParams.Add("id", id)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rCtx))

		s.author.Patch(w, req)

		s.Equal(status[i], w.Code)
		if status[i] >= 200 && status[i] < 300 {
			err := json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&response)
			s.NoError(err)
			s.NotNil(response.Data)
			s.Equal(author.Id, *response.Data.Id)
			s.NotEqual(author.Name, response.Data.Name)
		}

	}

}

func (s *SuiteAuthorTest) TestDelete() {
	author := s.create()
	ids := []string{author.Id, "", "123"}
	status := []int{http.StatusNoContent, http.StatusBadRequest, http.StatusNotFound}
	for i, id := range ids {
		req, _ := http.NewRequest("DELETE", "/authors/{id}", http.NoBody)
		req.Header.Set("Content-Type", "application/json")
		rCtx := chi.NewRouteContext()
		rCtx.URLParams.Add("id", id)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rCtx))
		w := httptest.NewRecorder()

		s.author.Delete(w, req)

		s.Equal(status[i], w.Code)
	}
}

func TestAuthorSuite(t *testing.T) {
	suite.Run(t, new(SuiteAuthorTest))
}
