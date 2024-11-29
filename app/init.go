package app

import (
	"context"
	"errors"
	"os"

	"github.com/PatricioYegros/uala_challenge/app/repository"
	"github.com/PatricioYegros/uala_challenge/app/service"
	"github.com/PatricioYegros/uala_challenge/app/utils"

	"github.com/redis/go-redis/v9"
)

const (
	CacheURLEnvVar      = "CACHE_URL"
	CachePasswordEnvVar = "CACHE_PASSWORD"
)

var ErrCacheNotConfigured = errors.New("cache env variables not configured")

func NewService() (*service.TwitterService, *redis.Client, error) {
	cacheURL := os.Getenv(CacheURLEnvVar)
	cachePassword := os.Getenv(CachePasswordEnvVar)

	if cacheURL == "" || cachePassword == "" {
		return nil, nil, ErrCacheNotConfigured
	}

	//create redis client
	redis := redis.NewClient(&redis.Options{
		Addr:     cacheURL,
		Password: cachePassword,
	})

	//test connection
	err := redis.Ping(context.Background()).Err()
	if err != nil {
		return nil, nil, err
	}

	//return service
	return &service.TwitterService{
		Repository: repository.Repository{
			Redis: redis,
		},
		Clock: utils.Clock{},
	}, redis, nil
}
