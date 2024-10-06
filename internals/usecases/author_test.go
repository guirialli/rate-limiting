package usecases

import (
	"context"
	"database/sql"
	"github.com/guirialli/rater_limit/internals/entity"
	"github.com/guirialli/rater_limit/internals/entity/vos"
	"github.com/guirialli/rater_limit/internals/infra/database"
	"github.com/guirialli/rater_limit/test/mock"
	"github.com/stretchr/testify/suite"
	"math/rand/v2"
	"testing"
	"time"
)

type AuthorTestSuite struct {
	suite.Suite
	useCase     *Author
	bookUseCase *Book
	db          *database.Sqlite
	fileInit    string
}

func (s *AuthorTestSuite) SetupTest() {
	s.db = database.NewSqlite("file::memory:?cache=shared")
	s.useCase = NewAuthor()
	s.bookUseCase = NewBook()
	s.fileInit = "../../test/database/init.sql"
}

func (s *AuthorTestSuite) createAuthor(ctx context.Context, db *sql.DB) *entity.Author {
	d := "teste"
	author, err := s.useCase.Create(ctx, db, mock.NewAuthor().Create(&d))
	if err != nil {
		panic(err)
	}
	return author
}

func (s *AuthorTestSuite) TestCreate() {
	db, _ := s.db.InitDatabaseGetConnection(s.fileInit)
	defer db.Close()
	ctx := context.Background()
	authorCreate := mock.NewAuthor().Create(nil)

	author, err := s.useCase.Create(ctx, db, authorCreate)

	s.Nil(err)
	s.NotNil(author)
	s.Equal(author.Name, authorCreate.Name)
	s.Equal(author.Description, authorCreate.Description)
	s.Equal(*author.Birthday, *authorCreate.Birthday)

}

func (s *AuthorTestSuite) TestFindAll() {
	db, _ := s.db.InitDatabaseGetConnection(s.fileInit)
	defer db.Close()
	length := rand.IntN(1000) + 1

	ctx := context.Background()
	for i := 0; i < length; i++ {
		s.createAuthor(ctx, db)
	}

	authors, err := s.useCase.FindAll(ctx, db)

	s.Nil(err)
	s.NotNil(authors)
	s.Len(authors, length)
}

func (s *AuthorTestSuite) TestFindAllWithBooks() {
	db, _ := s.db.InitDatabaseGetConnection(s.fileInit)
	defer db.Close()
	ctx := context.Background()
	author := s.createAuthor(ctx, db)
	book, _ := s.bookUseCase.Create(ctx, db, mock.NewBookMock().CreateWithAuthor(author.Id, nil))

	result, err := s.useCase.FindAllWithBooks(ctx, db, s.bookUseCase)

	s.Nil(err)
	s.NotNil(result)
	s.Len(result, 1)

	for _, ab := range result {
		s.Equal(ab.Books[0].Id, book.Id)
		s.Equal(ab.Books[0].Author, author.Id)
		s.Equal(ab.Books[0].Description, book.Description)
		s.Equal(ab.Author.Id, author.Id)
	}

}

func (s *AuthorTestSuite) TestFindById() {
	db, _ := s.db.InitDatabaseGetConnection(s.fileInit)
	defer db.Close()
	ctx := context.Background()
	author := s.createAuthor(ctx, db)

	result, err := s.useCase.FindById(ctx, db, author.Id)

	s.Nil(err)
	s.NotNil(result)
	s.Equal(author.Name, result.Name)
	s.Equal(author.Description, result.Description)
	s.Equal(author.Birthday.UnixMilli(), result.Birthday.UnixMilli())
	s.Equal(author.Id, result.Id)

}

func (s *AuthorTestSuite) TestFindByIdNotFound() {
	db, _ := s.db.InitDatabaseGetConnection(s.fileInit)
	defer db.Close()
	result, err := s.useCase.FindById(context.Background(), db, "not-found")

	s.NotNil(err)
	s.Nil(result)
}

func (s *AuthorTestSuite) TestUpdate() {
	db, _ := s.db.InitDatabaseGetConnection(s.fileInit)
	defer db.Close()
	ctx := context.Background()
	author := s.createAuthor(ctx, db)
	birthday := time.Now()
	description := "lorem ipsum"

	result, err := s.useCase.Update(ctx, db, author.Id, &vos.AuthorUpdate{
		Name:        "TestUpdate",
		Birthday:    &birthday,
		Description: &description,
	})

	s.Nil(err)
	s.NotNil(result)
	s.Equal(author.Id, result.Id)
	s.NotEqual(author.Name, result.Name)
	s.NotEqual(author.Birthday.UnixMicro(), result.Birthday.UnixMicro())
	s.NotEqual(author.Description, result.Description)
}

func (s *AuthorTestSuite) TestUpdateIdNotFound() {
	db, _ := s.db.InitDatabaseGetConnection(s.fileInit)
	defer db.Close()

	result, err := s.useCase.Update(context.Background(), db, "test", &vos.AuthorUpdate{
		Name:        "TestUpdateIdNotFound",
		Description: nil,
		Birthday:    nil,
	})

	s.NotNil(err)
	s.Nil(result)
}

func (s *AuthorTestSuite) TestPatchFull() {
	db, _ := s.db.InitDatabaseGetConnection(s.fileInit)
	defer db.Close()
	author := s.createAuthor(context.Background(), db)
	authorPatch := mock.NewAuthor().Patch(nil, nil, nil)

	result, err := s.useCase.Patch(context.Background(), db, author.Id, authorPatch)

	s.Nil(err)
	s.NotNil(result)
	s.NotEqual(author.Name, result.Name)
	s.NotEqual(author.Birthday.UnixMicro(), result.Birthday.UnixMicro())
	s.NotEqual(author.Description, result.Description)
}

func (s *AuthorTestSuite) TestPatchPartial() {
	db, _ := s.db.InitDatabaseGetConnection(s.fileInit)
	defer db.Close()
	author := s.createAuthor(context.Background(), db)
	authorPatch := mock.NewAuthor().Patch(nil, nil, author.Birthday)

	result, err := s.useCase.Patch(context.Background(), db, author.Id, authorPatch)

	s.Nil(err)
	s.NotNil(result)
	s.NotEqual(author.Name, result.Name)
	s.Equal(author.Birthday.UnixMicro(), result.Birthday.UnixMicro())
	s.NotEqual(author.Description, result.Description)
}

func (s *AuthorTestSuite) TestPatchNotFound() {
	db, _ := s.db.InitDatabaseGetConnection(s.fileInit)
	defer db.Close()
	authorPatch := mock.NewAuthor().Patch(nil, nil, nil)

	result, err := s.useCase.Patch(context.Background(), db, "not-found", authorPatch)

	s.NotNil(err)
	s.Nil(result)
}

func (s *AuthorTestSuite) TestDelete() {
	db, _ := s.db.InitDatabaseGetConnection(s.fileInit)
	defer db.Close()
	author := s.createAuthor(context.Background(), db)

	err := s.useCase.Delete(context.Background(), db, author.Id)
	result, resultErr := s.useCase.FindById(context.Background(), db, author.Id)

	s.Nil(err)
	s.NotNil(resultErr)
	s.Nil(result)
}

func (s *AuthorTestSuite) TestDeleteNotFound() {
	db, _ := s.db.InitDatabaseGetConnection(s.fileInit)
	defer db.Close()

	err := s.useCase.Delete(context.Background(), db, "not-found")
	s.NotNil(err)
}

func TestAuthorTestSuit(t *testing.T) {
	suite.Run(t, new(AuthorTestSuite))
}
