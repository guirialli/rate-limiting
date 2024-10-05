package usecases

import (
	"context"
	"database/sql"
	"github.com/guirialli/rater_limit/internals/entity"
	"github.com/guirialli/rater_limit/internals/infra/database"
	"github.com/guirialli/rater_limit/internals/vos"
	"github.com/guirialli/rater_limit/test/mock"
	"github.com/stretchr/testify/suite"
	"math/rand/v2"
	"testing"
	"time"
)

type AuthorTestSuite struct {
	suite.Suite
	useCase  *Author
	db       *database.Sqlite
	fileInit string
}

func (s *AuthorTestSuite) SetupTest() {
	s.db = database.NewSqlite("file::memory:?cache=shared")
	s.useCase = NewAuthor()
	s.fileInit = "../../test/database/init.sql"
}

func (s *AuthorTestSuite) createAuthor(ctx context.Context, db *sql.DB) *entity.Author {
	author, err := s.useCase.Create(ctx, db, mock.NewAuthor().Create())
	if err != nil {
		panic(err)
	}
	return author
}

func (s *AuthorTestSuite) TestCreateAuthor() {
	db, _ := s.db.InitDatabaseGetConnection(s.fileInit)
	defer db.Close()
	ctx := context.Background()
	authorCreate := mock.NewAuthor().Create()

	author, err := s.useCase.Create(ctx, db, authorCreate)

	s.Nil(err)
	s.NotNil(author)
	s.Equal(author.Name, authorCreate.Name)
	s.Equal(author.Description, authorCreate.Description)
	s.Equal(*author.Birthday, *authorCreate.Birthday)

}

func (s *AuthorTestSuite) TestFindAllAuthors() {
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

func TestAuthorTestSuit(t *testing.T) {
	suite.Run(t, new(AuthorTestSuite))
}
