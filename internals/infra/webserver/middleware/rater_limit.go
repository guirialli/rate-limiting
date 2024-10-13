package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/jwtauth"
	"github.com/guirialli/rater_limit/config"
	"github.com/guirialli/rater_limit/internals/infra/database"
	"github.com/guirialli/rater_limit/internals/usecases"
	"github.com/redis/go-redis/v9"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type user struct {
	Trys          int        `json:"trys"`
	Type          string     `json:"type"`
	AccessTimeout time.Time  `json:"access_timeout"`
	BlockAt       *time.Time `json:"block_at"`
}

func newUser(typer string, timeout time.Duration) user {
	return user{
		Trys:          0,
		Type:          typer,
		AccessTimeout: time.Now().Add(timeout),
		BlockAt:       nil,
	}
}

type RaterLimit struct {
	rdb     *redis.Client
	cfg     config.RaterLimit
	jwtAuth *jwtauth.JWTAuth
}

func NewRaterLimit(cfg config.RaterLimit) (*RaterLimit, error) {
	jwtCfg := config.LoadJwtConfig()
	u, err := usecases.NewUser(jwtCfg.Secret, jwtCfg.ExpireIn, jwtCfg.UnitTime)
	if err != nil {
		return nil, err
	}
	rdb := database.NewRedisClient()

	return &RaterLimit{
		rdb:     rdb,
		cfg:     cfg,
		jwtAuth: u.NewTokenAuth(),
	}, nil
}

func (rl *RaterLimit) getUser(key string) (user, bool) {
	ctx := context.Background()
	val, err := rl.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return user{}, false
	} else if err != nil {
		panic(err)
	}

	var u user
	if err = json.Unmarshal([]byte(val), &u); err != nil {
		panic(err)
	}
	return u, true
}

func (rl *RaterLimit) setUser(key string, u user) {
	ctx := context.Background()
	val, _ := json.Marshal(u)
	rl.rdb.Set(ctx, key, val, 0)
}

func (rl *RaterLimit) TrackAccess(key string) bool {
	ipRegex := regexp.MustCompile(`^(?:\d{1,3}\.){3}\d{1,3}$`)

	element, exists := rl.getUser(key)
	if !exists {
		if ipRegex.MatchString(key) {
			element = newUser("ip", rl.cfg.IpRefresh)
		} else {
			element = newUser("jwt", rl.cfg.JwtRefresh)
		}
		rl.setUser(key, element)
		return true
	} else if element.BlockAt != nil && element.BlockAt.After(time.Now()) {
		return false
	}

	element.Trys++

	if element.Type == "jwt" && element.Trys >= rl.cfg.JwtTrysMax {
		if element.AccessTimeout.Before(time.Now()) {
			element = newUser(element.Type, rl.cfg.JwtRefresh)
		} else {
			dtBlock := time.Now().Add(rl.cfg.BlockTimeout)
			element.BlockAt = &dtBlock
		}
		rl.setUser(key, element)
		return element.BlockAt == nil
	}

	if element.Type == "ip" && element.Trys >= rl.cfg.IpTrysMax {
		if element.AccessTimeout.Before(time.Now()) {
			element = newUser(element.Type, rl.cfg.IpRefresh)
		} else {
			dtBlock := time.Now().Add(rl.cfg.BlockTimeout)
			element.BlockAt = &dtBlock
		}
		rl.setUser(key, element)
		return element.BlockAt == nil
	}

	rl.setUser(key, element)
	return true
}

func (rl *RaterLimit) ValidToken(key string) bool {
	if strings.HasPrefix(key, "Bearer") {
		token := key[7:]
		_, err := rl.jwtAuth.Decode(token)
		return err == nil
	}
	return false
}

func (rl *RaterLimit) RateLimitMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("Authorization")
			if !rl.ValidToken(key) {
				key = strings.Split(r.RemoteAddr, ":")[0]
			}

			if !rl.TrackAccess(key) {
				http.Error(w,
					"you have reached the maximum number of requests or actions allowed within a certain time frame",
					http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
