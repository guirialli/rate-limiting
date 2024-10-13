package usecases

import (
	"context"
	"github.com/go-chi/jwtauth"
	"github.com/guirialli/rater_limit/config"
	"github.com/guirialli/rater_limit/internals/entity"
	"github.com/guirialli/rater_limit/internals/infra/database"
	"regexp"
	"strings"
	"time"
)

type RaterLimit struct {
	rdb     database.IRateLimitDatabase[entity.RaterLimit]
	cfg     config.RaterLimit
	jwtAuth *jwtauth.JWTAuth
}

func NewRaterLimit(userUseCase IUser, cfg config.RaterLimit, rdb database.IRateLimitDatabase[entity.RaterLimit]) (*RaterLimit, error) {
	return &RaterLimit{
		rdb:     rdb,
		cfg:     cfg,
		jwtAuth: userUseCase.NewTokenAuth(),
	}, nil
}

func (rl *RaterLimit) TrackAccess(ctx context.Context, key string) bool {
	ipRegex := regexp.MustCompile(`^(?:\d{1,3}\.){3}\d{1,3}$`)

	element, exists := rl.rdb.Get(ctx, key)
	if !exists {
		if ipRegex.MatchString(key) {
			element = entity.NewRaterLimit("ip", rl.cfg.IpRefresh)
		} else {
			element = entity.NewRaterLimit("jwt", rl.cfg.JwtRefresh)
		}
		err := rl.rdb.Set(ctx, key, element)
		return err == nil
	} else if element.BlockAt != nil && element.BlockAt.After(time.Now()) {
		return false
	}

	element.Trys++

	if element.Type == "jwt" && element.Trys >= rl.cfg.JwtTrysMax {
		if element.AccessTimeout.Before(time.Now()) {
			element = entity.NewRaterLimit(element.Type, rl.cfg.JwtRefresh)
		} else {
			dtBlock := time.Now().Add(rl.cfg.BlockTimeout)
			element.BlockAt = &dtBlock
		}
		if err := rl.rdb.Set(ctx, key, element); err != nil {
			return false
		}

		return element.BlockAt == nil
	}

	if element.Type == "ip" && element.Trys >= rl.cfg.IpTrysMax {
		if element.AccessTimeout.Before(time.Now()) {
			element = entity.NewRaterLimit(element.Type, rl.cfg.IpRefresh)
		} else {
			dtBlock := time.Now().Add(rl.cfg.BlockTimeout)
			element.BlockAt = &dtBlock
		}
		if err := rl.rdb.Set(ctx, key, element); err != nil {
			return false
		}

		return element.BlockAt == nil
	}
	err := rl.rdb.Set(ctx, key, element)
	return err == nil
}

func (rl *RaterLimit) ValidToken(key string) bool {
	if strings.HasPrefix(key, "Bearer") {
		token := key[7:]
		_, err := rl.jwtAuth.Decode(token)
		return err == nil
	}
	return false
}
