package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/jwtauth"
	"github.com/guirialli/rater_limit/config"
	"github.com/guirialli/rater_limit/internals/entity"
	"github.com/redis/go-redis/v9"
	"regexp"
	"strings"
	"time"
)

type RaterLimit struct {
	rdb     *redis.Client
	cfg     config.RaterLimit
	jwtAuth *jwtauth.JWTAuth
}

func NewRaterLimit(userUseCase IUser, cfg config.RaterLimit, rdb *redis.Client) (*RaterLimit, error) {
	return &RaterLimit{
		rdb:     rdb,
		cfg:     cfg,
		jwtAuth: userUseCase.NewTokenAuth(),
	}, nil
}

func (rl *RaterLimit) getUser(key string) (entity.RaterLimit, bool) {
	ctx := context.Background()
	val, err := rl.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return entity.RaterLimit{}, false
	} else if err != nil {
		panic(err)
	}

	var u entity.RaterLimit
	if err = json.Unmarshal([]byte(val), &u); err != nil {
		panic(err)
	}
	return u, true
}

func (rl *RaterLimit) setUser(key string, u entity.RaterLimit) {
	ctx := context.Background()
	val, _ := json.Marshal(u)
	rl.rdb.Set(ctx, key, val, 0)
}

func (rl *RaterLimit) TrackAccess(key string) bool {
	ipRegex := regexp.MustCompile(`^(?:\d{1,3}\.){3}\d{1,3}$`)

	element, exists := rl.getUser(key)
	if !exists {
		if ipRegex.MatchString(key) {
			element = entity.NewRaterLimit("ip", rl.cfg.IpRefresh)
		} else {
			element = entity.NewRaterLimit("jwt", rl.cfg.JwtRefresh)
		}
		rl.setUser(key, element)
		return true
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
		rl.setUser(key, element)
		return element.BlockAt == nil
	}

	if element.Type == "ip" && element.Trys >= rl.cfg.IpTrysMax {
		if element.AccessTimeout.Before(time.Now()) {
			element = entity.NewRaterLimit(element.Type, rl.cfg.IpRefresh)
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
