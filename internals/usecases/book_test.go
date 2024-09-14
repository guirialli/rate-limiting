package usecases

import (
	"context"
	"github.com/guirialli/rater_limit/internals/infra/database"
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

func TestBookTestSuit(t *testing.T) {
	suite.Run(t, new(BookTestSuite))
}
