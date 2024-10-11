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

type SuiteBookTest struct {
	suite.Suite
	book          *Book
	useCase       usecases.IBook
	authorUseCase usecases.IAuthor
	init          string
	db            *sql.DB
}

func (s *SuiteBookTest) SetupSuite() {
	s.init = "../../../../test/database/init.sql"
}

func (s *SuiteBookTest) SetupTest() {
	db, _ := database.NewSqlite("file::memory:?cache=shared").InitDatabaseGetConnection(s.init)
	author := usecases.NewAuthor()
	book := usecases.NewBook()

	s.db = db
	s.useCase = book
	s.authorUseCase = author
	s.book = NewBook(s.db, book, author)
}

func (s *SuiteBookTest) TestSetup() {
	s.NotNil(s.db)
	s.NotNil(s.authorUseCase)
	s.NotNil(s.useCase)
	s.NotNil(s.book)
}

func (s *SuiteBookTest) isSuccessReq(st int) bool {
	return st >= 200 && st < 300
}

func (s *SuiteBookTest) create() (*entity.Book, *entity.Author) {
	ctx := context.Background()
	author, err := s.authorUseCase.Create(ctx, s.db, mock.NewAuthor().Create(nil))
	if err != nil {
		panic(err)
	}
	book, err := s.useCase.Create(ctx, s.db, mock.NewBookMock().CreateWithAuthor(author.Id, nil))
	if err != nil {
		panic(err)
	}

	return book, author
}

func (s *SuiteBookTest) TestGetAll() {
	length := rand.IntN(1000) + 1
	for i := 0; i < length; i++ {
		s.create()
	}
	var response dtos.ResponseJson[[]entity.Book]
	req, _ := http.NewRequest("GET", "/books", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.book.GetAll(w, req)

	s.Equal(http.StatusOK, w.Code)
	err := json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&response)
	s.NoError(err)
	s.Equal(w.Code, response.Status)
	s.Len(response.Data, length)
}

func (s *SuiteBookTest) TestGetAllWithAuthor() {
	book, author := s.create()
	var result dtos.ResponseJson[[]dtos.BookWithAuthor]

	req, _ := http.NewRequest("GET", "/books/author", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.book.GetAllWithAuthors(w, req)

	s.Equal(http.StatusOK, w.Code)
	err := json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&result)
	s.NoError(err)
	for _, ab := range result.Data {
		s.Equal(author.Id, ab.Author.Id)
		s.Equal(book.Id, ab.Book.Id)
	}
}

func (s *SuiteBookTest) TestGetById() {
	book, _ := s.create()
	status := []int{http.StatusOK, http.StatusBadRequest, http.StatusNotFound}
	ids := []string{book.Id, "", "123"}
	for i, id := range ids {
		req, _ := http.NewRequest("GET", "/books/{id}", nil)
		w := httptest.NewRecorder()
		rCtx := chi.NewRouteContext()
		rCtx.URLParams.Add("id", id)

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rCtx))
		s.book.GetById(w, req)
		s.Equal(status[i], w.Code)

		if s.isSuccessReq(status[i]) {
			var response dtos.ResponseJson[entity.Book]
			err := json.NewDecoder(w.Body).Decode(&response)
			s.NoError(err)
			s.NotNil(response.Data)
			s.Equal(response.Data.Id, book.Id)
		}

	}

}

func (s *SuiteBookTest) TestCreate() {
	status := []int{http.StatusCreated, http.StatusBadRequest}
	var response dtos.ResponseJson[entity.Author]
	for i, st := range status {
		var body []byte
		if i == 0 {
			body, _ = json.Marshal(mock.NewAuthor().Create(nil))
		}

		req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		s.book.Create(w, req)

		s.Equal(st, w.Code)
		if s.isSuccessReq(st) {
			err := json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&response)
			s.NoError(err)
			s.NotNil(response.Data)
			s.NotEmpty(response.Data.Id)
		}
	}
}

func (s *SuiteBookTest) TestUpdate() {
	book, _ := s.create()
	var response dtos.ResponseJson[entity.Book]
	newDescription := "bla"
	update := dtos.BookUpdate{
		Title:       book.Title + "1",
		Pages:       book.Pages + 1,
		Description: &newDescription,
	}
	status := []int{http.StatusOK, http.StatusBadRequest, http.StatusNotFound, http.StatusBadRequest}
	bodyS := []interface{}{update, update, update, ""}
	ids := []string{book.Id, "", "123", book.Id}

	for i, id := range ids {
		body, _ := json.Marshal(bodyS[i])
		req, _ := http.NewRequest("PUT", "/books/{id}", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		rCtx := chi.NewRouteContext()
		rCtx.URLParams.Add("id", id)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rCtx))

		s.book.Update(w, req)

		s.Equal(status[i], w.Code)
		if s.isSuccessReq(status[i]) {
			err := json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&response)
			s.NoError(err)
			s.NotNil(response.Data)
			s.Equal(book.Id, response.Data.Id)
			s.NotEqual(book.Title, response.Data.Title)
			s.NotEqual(book.Pages, response.Data.Pages)
		}
	}
}

func (s *SuiteBookTest) TestPatch() {
	book, _ := s.create()
	var response dtos.ResponseJson[entity.Book]
	newTitle := book.Title + "1"
	newDescription := "bla"
	newPages := book.Pages + 1
	update := dtos.BookPatch{
		Title:       &newTitle,
		Pages:       &newPages,
		Description: &newDescription,
	}
	status := []int{http.StatusOK, http.StatusBadRequest, http.StatusNotFound, http.StatusBadRequest}
	bodyS := []interface{}{update, update, update, ""}
	ids := []string{book.Id, "", "123", book.Id}

	for i, id := range ids {
		body, _ := json.Marshal(bodyS[i])
		req, _ := http.NewRequest("PATCH", "/books/{id}", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		rCtx := chi.NewRouteContext()
		rCtx.URLParams.Add("id", id)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rCtx))

		s.book.Path(w, req)

		s.Equal(status[i], w.Code)
		if s.isSuccessReq(status[i]) {
			err := json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&response)
			s.NoError(err)
			s.NotNil(response.Data)
			s.Equal(book.Id, response.Data.Id)
			s.NotEqual(book.Title, response.Data.Title)
			s.NotEqual(book.Pages, response.Data.Pages)
		}
	}
}

func (s *SuiteBookTest) TestDelete() {
	book, _ := s.create()
	ids := []string{book.Id, "", "123"}
	status := []int{http.StatusNoContent, http.StatusBadRequest, http.StatusNotFound}
	for i, id := range ids {
		req, _ := http.NewRequest("DELETE", "/books/{id}", http.NoBody)
		req.Header.Set("Content-Type", "application/json")
		rCtx := chi.NewRouteContext()
		rCtx.URLParams.Add("id", id)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rCtx))
		w := httptest.NewRecorder()

		s.book.Delete(w, req)

		s.Equal(status[i], w.Code)
	}
}

func TestBookSuite(t *testing.T) {
	suite.Run(t, new(SuiteBookTest))
}
