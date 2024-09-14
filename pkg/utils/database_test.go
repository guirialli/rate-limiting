package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DatabaseTestSuite struct {
	suite.Suite
	Db *sql.DB
}

func (s *DatabaseTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(fmt.Sprintf("Couldn't open database: %v", err))
	}
	s.Db = db

}

func (s *DatabaseTestSuite) TestInitDatabase() {
	err := NewDatabaseUtils().InitDatabase(s.Db, "../../test/database/init.sql")
	s.Nil(err)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}
