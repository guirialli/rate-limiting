package usecases

import (
	"context"
	"github.com/guirialli/rater_limit/internals/infra/database"
	"github.com/guirialli/rater_limit/test/mock"
	"github.com/stretchr/testify/suite"
	"math/rand/v2"
	"testing"
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

func (s *AuthorTestSuite) TestCreateAuthor() {
	db, err := s.db.InitDatabaseGetConnection(s.fileInit)
	if err != nil {
		panic(err)
	}
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
	db, err := s.db.InitDatabaseGetConnection(s.fileInit)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	length := rand.IntN(1000) + 1

	ctx := context.Background()
	authorCreate := mock.NewAuthor().Create()
	for i := 0; i < length; i++ {
		_, err := s.useCase.Create(ctx, db, authorCreate)
		if err != nil {
			panic(err)
		}
	}

	authors, err := s.useCase.FindAll(ctx, db)
	s.Nil(err)
	s.NotNil(authors)
	s.Len(authors, length)
}

func (s *AuthorTestSuite) TestFindById() {
	db, err := s.db.InitDatabaseGetConnection(s.fileInit)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	ctx := context.Background()
	author, err := s.useCase.Create(ctx, db, mock.NewAuthor().Create())
	if err != nil {
		panic(err)
	}

	result, err := s.useCase.FindById(ctx, db, author.Id)

	s.Nil(err)
	s.NotNil(result)
	s.Equal(author.Name, result.Name)
	s.Equal(author.Description, result.Description)
	s.Equal(author.Birthday.UnixMilli(), result.Birthday.UnixMilli())
	s.Equal(author.Id, result.Id)

}

func TestAuthorTestSuit(t *testing.T) {
	suite.Run(t, new(AuthorTestSuite))
}
