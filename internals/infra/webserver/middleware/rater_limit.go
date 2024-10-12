package middleware

import (
	"github.com/go-chi/jwtauth"
	"github.com/guirialli/rater_limit/config"
	"github.com/guirialli/rater_limit/internals/usecases"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

type user struct {
	Trys          int
	Type          string
	AccessTimeout time.Time
	BlockAt       *time.Time
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
	data    map[string]user
	mu      sync.RWMutex
	cfg     config.RaterLimit
	jwtAuth *jwtauth.JWTAuth
}

func NewRaterLimit(cfg config.RaterLimit) (*RaterLimit, error) {
	jwtCfg := config.LoadJwtConfig()
	u, err := usecases.NewUser(jwtCfg.Secret, jwtCfg.ExpireIn, jwtCfg.UnitTime)
	if err != nil {
		return nil, err
	}

	return &RaterLimit{
		data:    make(map[string]user),
		cfg:     cfg,
		jwtAuth: u.NewTokenAuth(),
	}, nil
}

func (rl *RaterLimit) TrackAccess(key string) bool {
	ipRegex := regexp.MustCompile(`^(?:\d{1,3}\.){3}\d{1,3}$`)
	rl.mu.Lock()
	defer rl.mu.Unlock()

	element, exists := rl.data[key]
	if !exists {
		if ipRegex.MatchString(key) {
			rl.data[key] = newUser("ip", rl.cfg.IpTimeout)
		} else {
			rl.data[key] = newUser("jwt", rl.cfg.JwtTimeout)
		}
		return true
	} else if element.BlockAt != nil {
		if element.BlockAt.After(time.Now()) {
			return false
		}
		element.BlockAt = nil
	}

	element.Trys++

	if element.Type == "jwt" && element.Trys >= rl.cfg.JwtTrysMax {
		if element.AccessTimeout.Before(time.Now()) {
			rl.data[key] = newUser(element.Type, rl.cfg.JwtTimeout)
			return true
		}
		dtBlock := time.Now().Add(rl.cfg.BlockTimeout)
		element.BlockAt = &dtBlock
		rl.data[key] = element
		return false

	} else if element.Type == "ip" && element.Trys >= rl.cfg.IpTrysMax {
		if element.AccessTimeout.Before(time.Now()) {
			rl.data[key] = newUser(element.Type, rl.cfg.IpTimeout)
			return true
		}
		dtBlock := time.Now().Add(rl.cfg.BlockTimeout)
		element.BlockAt = &dtBlock
		rl.data[key] = element
		return false
	}

	rl.data[key] = element
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
				http.Error(w, "Too Many Request", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
