package service_test

import (
	"errors"
	"testing"

	"github.com/PatricioYegros/uala_challenge/app/service"
	"github.com/go-playground/assert/v2"

	repositoryMocks "github.com/PatricioYegros/uala_challenge/app/mocks/repository"
)

func TestFollowReturnsErrorSameIDs(t *testing.T) {
	followService := service.TwitterService{}
	err := followService.Follow(1, 1)
	assert.Equal(t, err, service.ErrEqualsIDs)
}

func TestFollowRepositoryReturnsError(t *testing.T) {
	mockRepository := repositoryMocks.NewIRepository(t)

	followService := service.TwitterService{
		Repository: mockRepository,
	}

	mockRepository.On("AddFollower", uint(2), uint(1)).Return(errors.New("Error"))

	err := followService.Follow(2, 1)
	assert.Equal(t, err, service.ErrFollowing)
}
