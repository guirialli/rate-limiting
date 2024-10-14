package usecases

import (
	"context"
	"database/sql"
	"github.com/guirialli/rater_limit/config"
	"github.com/guirialli/rater_limit/internals/entity"
	"github.com/guirialli/rater_limit/internals/infra/database"
	databasetest "github.com/guirialli/rater_limit/test/database"
	"github.com/guirialli/rater_limit/test/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

const (
	jwtTry = 10
	ipTry  = 5
)

type RateLimitTestSuite struct {
	suite.Suite
	tableName   string
	fileInit    string
	rateLimit   *RaterLimit
	db          *sql.DB
	userUseCase *User
}

func (s *RateLimitTestSuite) SetupSuite() {
	s.fileInit = "../../test/database/init.sql"
	s.tableName = "rate_limit"
}

func (s *RateLimitTestSuite) SetupTest() {
	db, err := database.NewSqlite("file::memory:?cache=shared").InitDatabaseGetConnection(s.fileInit)
	if err != nil {
		panic(err)
	}
	user, _ := NewUser("a", 10, 'h')
	rateDb := databasetest.NewRaterLimit[entity.RaterLimit](db, s.tableName)
	useCase, _ := NewRaterLimit(user, config.RaterLimit{
		IpRefresh:    time.Duration(10) * time.Second,
		JwtRefresh:   time.Duration(10) * time.Second,
		BlockTimeout: time.Duration(5) * time.Minute,
		JwtTrysMax:   jwtTry,
		IpTrysMax:    ipTry,
	}, rateDb)
	s.rateLimit = useCase
	s.db = db
	s.userUseCase = user
}

func (s *RateLimitTestSuite) TearDownTest() {
	s.db.Close()
	s.rateLimit = nil
	s.userUseCase = nil
}

func (s *RateLimitTestSuite) TestSetup() {
	s.NotNil(s.rateLimit)
	s.NotEmpty(s.tableName)
	s.NotEmpty(s.fileInit)
}

func (s *RateLimitTestSuite) TestValidToken() {
	form := mock.NewUserMock().RegisterForm()
	token, _ := s.userUseCase.Register(context.Background(), s.db, form)
	tokenValid := "Bearer " + token
	tokens := []string{tokenValid, token, "",
		"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"}
	results := []bool{true, false, false, false}
	for i, t := range tokens {
		result := s.rateLimit.ValidToken(t)
		s.Equal(results[i], result)
	}
}

func (s *RateLimitTestSuite) TestRateLimit() {
	form := mock.NewUserMock().RegisterForm()
	token, _ := s.userUseCase.Register(context.Background(), s.db, form)
	keys := []string{"127.0.0.1", token}
	maxAccess := []int{ipTry, jwtTry}
	for i, a := range maxAccess {
		for j := 1; j < a*a; j++ {
			result := s.rateLimit.TrackAccess(context.Background(), keys[i])
			s.Equal(j <= a, result)
		}
	}
}

func TestRateLimitSuite(t *testing.T) {
	suite.Run(t, new(RateLimitTestSuite))
}
