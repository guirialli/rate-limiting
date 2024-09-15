package usecases

import (
	"context"
	"github.com/guirialli/rater_limit/internals/infra/database"
	"github.com/guirialli/rater_limit/internals/vos"
	"github.com/guirialli/rater_limit/test/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type BookTestSuite struct {
	suite.Suite
	useCase  *Book
	db       *database.Sqlite
	fileInit string
}

func (s *BookTestSuite) SetupTest() {
	db := database.NewSqlite("test.db")
	s.db = db
	s.useCase = NewBook()
	s.fileInit = "../../test/database/init.sql"
}

func (s *BookTestSuite) TestCreateBook() {
	db, err := s.db.InitDatabaseGetConnection(s.fileInit)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx := context.Background()
	var description *string
	bookCreate := mock.NewBookMock().MockCreate(description)

	book, err := s.useCase.Create(ctx, db, bookCreate)

	s.Nil(err)
	s.NotNil(book)

	s.Equal(book.Title, bookCreate.Title)
	s.Equal(book.Pages, bookCreate.Pages)
	s.Equal(book.Description, bookCreate.Description)
	s.Equal(book.Author, bookCreate.Author)
}

func (s *BookTestSuite) TestFindAll() {
	db, err := s.db.InitDatabaseGetConnection(s.fileInit)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx := context.Background()
	book, err := s.useCase.FindAll(ctx, db)
	s.Nil(err)
	s.NotNil(book)

}

func (s *BookTestSuite) TestFindById() {
	db, err := s.db.InitDatabaseGetConnection(s.fileInit)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	description := "test"
	bookCreate := mock.NewBookMock().MockCreate(&description)
	book, _ := s.useCase.Create(ctx, db, bookCreate)

	result, err := s.useCase.FindById(ctx, db, book.Id)
	s.Nil(err)
	s.NotNil(result)

	s.Equal(book.Id, result.Id)
	s.Equal(book.Title, result.Title)
	s.Equal(book.Pages, result.Pages)
	s.Equal(*book.Description, *result.Description)
	s.Equal(book.Author, result.Author)
}

func (s *BookTestSuite) TestUpdateBook() {
	db, err := s.db.InitDatabaseGetConnection(s.fileInit)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	description := "test"
	bookCreate := mock.NewBookMock().MockCreate(&description)
	book, _ := s.useCase.Create(ctx, db, bookCreate)

	description2 := "test2"
	result, err := s.useCase.Update(ctx, db, book.Id, &vos.BookUpdate{
		Title:       "test2",
		Pages:       0,
		Author:      "test2",
		Description: &description2,
	})
	s.Nil(err)
	s.NotNil(result)

	s.Equal(book.Id, result.Id)
	s.NotEqual(description, *result.Description)
	s.NotEqual(book.Pages, result.Pages)
	s.NotEqual(book.Title, result.Title)
	s.NotEqual(book.Author, result.Author)
}

func (s *BookTestSuite) TestDeleteBook() {
	db, err := s.db.InitDatabaseGetConnection(s.fileInit)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	description := "test"
	bookCreate := mock.NewBookMock().MockCreate(&description)
	book, _ := s.useCase.Create(ctx, db, bookCreate)

	err = s.useCase.Delete(ctx, db, book.Id)
	s.Nil(err)
	_, err = s.useCase.FindById(ctx, db, book.Id)
	s.NotNil(err)
}

func TestBookTestSuit(t *testing.T) {
	suite.Run(t, new(BookTestSuite))
}
