package usecases

import (
	"context"
	"github.com/guirialli/rater_limit/internals/entity/dtos"
	"github.com/guirialli/rater_limit/internals/infra/database"
	"github.com/guirialli/rater_limit/test/mock"
	"github.com/stretchr/testify/suite"
	"math/rand/v2"
	"testing"
)

type BookTestSuite struct {
	suite.Suite
	useCase       *Book
	db            *database.Sqlite
	authorUseCase *Author
	fileInit      string
}

func (s *BookTestSuite) SetupTest() {
	db := database.NewSqlite("file::memory:?cache=shared")
	s.db = db
	s.useCase = NewBook()
	s.authorUseCase = NewAuthor()
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
	bookCreate := mock.NewBookMock().Create(description)

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
	length := rand.IntN(1000) + 1

	ctx := context.Background()
	for i := 0; i < length; i++ {
		_, err = s.useCase.Create(ctx, db, mock.NewBookMock().Create(nil))
	}
	if err != nil {
		panic(err)
	}
	book, err := s.useCase.FindAll(ctx, db)
	s.Nil(err)
	s.NotNil(book)
	s.Len(book, length)

}

func (s *BookTestSuite) TestFindAllWithAuthor() {
	db, _ := s.db.InitDatabaseGetConnection(s.fileInit)
	defer db.Close()
	ctx := context.Background()

	author, _ := s.authorUseCase.Create(ctx, db, mock.NewAuthor().Create(nil))
	book, _ := s.useCase.Create(ctx, db, mock.NewBookMock().CreateWithAuthor(author.Id, nil))

	result, err := s.useCase.FindAllWithAuthor(ctx, db, s.authorUseCase)

	s.Nil(err)
	s.NotNil(result)
	s.Equal(author.Id, *result[0].Author.Id)
	s.Equal(book.Id, result[0].Book.Id)
}

func (s *BookTestSuite) TestFindAllBookByAuthor() {
	db, _ := s.db.InitDatabaseGetConnection(s.fileInit)
	defer db.Close()
	ctx := context.Background()

	author, _ := s.authorUseCase.Create(ctx, db, mock.NewAuthor().Create(nil))
	book, _ := s.useCase.Create(ctx, db, mock.NewBookMock().CreateWithAuthor(author.Id, nil))

	result, err := s.useCase.FindAllByAuthor(ctx, db, author.Id)

	s.Nil(err)
	s.NotNil(result)
	s.Equal(book.Id, result[0].Id)
	s.Equal(book.Author, result[0].Author)
}

func (s *BookTestSuite) TestFindById() {
	db, err := s.db.InitDatabaseGetConnection(s.fileInit)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx := context.Background()

	description := "test"
	bookCreate := mock.NewBookMock().Create(&description)
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

func (s *BookTestSuite) TestFindByIdWithAuthor() {
	db, _ := s.db.InitDatabaseGetConnection(s.fileInit)
	defer db.Close()
	ctx := context.Background()

	author, _ := s.authorUseCase.Create(ctx, db, mock.NewAuthor().Create(nil))
	book, _ := s.useCase.Create(ctx, db, mock.NewBookMock().CreateWithAuthor(author.Id, nil))

	result, err := s.useCase.FindByIdWithAuthor(ctx, db, book.Id, s.authorUseCase)

	s.Nil(err)
	s.NotNil(result)
	s.Equal(author.Id, *result.Author.Id)
	s.Equal(book.Id, result.Book.Id)
}

func (s *BookTestSuite) TestUpdateBook() {
	db, err := s.db.InitDatabaseGetConnection(s.fileInit)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	ctx := context.Background()

	description := "test"
	bookCreate := mock.NewBookMock().Create(&description)
	book, _ := s.useCase.Create(ctx, db, bookCreate)

	description2 := "test2"
	result, err := s.useCase.Update(ctx, db, book.Id, &dtos.BookUpdate{
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

func (s *BookTestSuite) TestPatchBook() {
	db, err := s.db.InitDatabaseGetConnection(s.fileInit)
	if err != nil {
		panic(err)

	}
	defer db.Close()
	ctx := context.Background()

	description := "test"
	bookCreate := mock.NewBookMock().Create(&description)
	book, _ := s.useCase.Create(ctx, db, bookCreate)

	title := "test2"
	result, err := s.useCase.Patch(ctx, db, book.Id, &dtos.BookPatch{
		Title: &title,
	})

	s.Nil(err)
	s.NotNil(result)

	s.Equal(book.Id, result.Id)
	s.NotEqual(book.Title, result.Title)
	s.Equal(book.Pages, result.Pages)
	s.Equal(*book.Description, *result.Description)
	s.Equal(book.Author, result.Author)
}

func (s *BookTestSuite) TestDeleteBook() {
	db, err := s.db.InitDatabaseGetConnection(s.fileInit)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	description := "test"
	bookCreate := mock.NewBookMock().Create(&description)
	book, _ := s.useCase.Create(ctx, db, bookCreate)

	err = s.useCase.Delete(ctx, db, book.Id)
	s.Nil(err)
	_, err = s.useCase.FindById(ctx, db, book.Id)
	s.NotNil(err)
}

func TestBookSuite(t *testing.T) {
	suite.Run(t, new(BookTestSuite))
}
