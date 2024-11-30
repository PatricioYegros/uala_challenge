package service_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/PatricioYegros/uala_challenge/app/models"
	"github.com/PatricioYegros/uala_challenge/app/service"
	"github.com/PatricioYegros/uala_challenge/app/utils"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"

	repositoryMocks "github.com/PatricioYegros/uala_challenge/mocks/repository"
	utilsMocks "github.com/PatricioYegros/uala_challenge/mocks/utils"
)

func TestFollowReturnsErrorSameIDs(t *testing.T) {
	followService := service.TwitterService{}
	err := followService.Follow(1, 1)
	assert.Equal(t, err, service.ErrEqualsIDs)
}

func TestFollowReturnsErrorGettingList(t *testing.T) {
	mockRepository := repositoryMocks.NewIRepository(t)

	followService := service.TwitterService{
		Repository: mockRepository,
	}

	mockRepository.On("GetFollowers", uint(1)).Return(nil, errors.New("Error"))

	err := followService.Follow(2, 1)
	assert.Equal(t, err, fmt.Errorf("%w from user %d for check existent follow", service.ErrFollowing, 1))
}

func TestFollowRepositoryReturnsError(t *testing.T) {
	mockRepository := repositoryMocks.NewIRepository(t)
	mockListFollowers := []uint{2, 3}

	followService := service.TwitterService{
		Repository: mockRepository,
	}

	mockRepository.On("GetFollowers", uint(1)).Return(mockListFollowers, nil)
	mockRepository.On("AddFollower", uint(1), uint(4)).Return(service.ErrFollowing)

	err := followService.Follow(4, 1)
	assert.Equal(t, err, fmt.Errorf("%w from user %d to user %d", service.ErrFollowing, uint(4), uint(1)))
}

func TestFollowRetunsErrorAlreadyFollowing(t *testing.T) {
	mockRepository := repositoryMocks.NewIRepository(t)
	mockListFollowers := []uint{2}

	followService := service.TwitterService{
		Repository: mockRepository,
	}

	mockRepository.On("GetFollowers", uint(1)).Return(mockListFollowers, nil)

	err := followService.Follow(2, 1)
	assert.Equal(t, err, service.ErrorFollowingAlready)
}

func TestFollowSuccess(t *testing.T) {
	mockRepository := repositoryMocks.NewIRepository(t)
	mockListFollowers := []uint{2, 3}

	followService := service.TwitterService{
		Repository: mockRepository,
	}

	mockRepository.On("GetFollowers", uint(1)).Return(mockListFollowers, nil)
	mockRepository.On("AddFollower", uint(1), uint(4)).Return(nil)
	err := followService.Follow(4, 1)
	assert.Equal(t, err, nil)
}

func TestTweetErrorMaxLengthExceeded(t *testing.T) {
	mockRepository := repositoryMocks.NewIRepository(t)

	tweetService := service.TwitterService{
		Repository: mockRepository,
		Clock:      utils.Clock{},
	}

	_, err := tweetService.Tweet(0, strings.Repeat("uala_challenge", 15))
	assert.Equal(t, err, models.ErrMaxLengthExceeded)
}

func TestTweetErrorCreatingTweet(t *testing.T) {
	mockRepository := repositoryMocks.NewIRepository(t)
	mockClock := utilsMocks.NewIClock(t)

	tweetService := service.TwitterService{
		Repository: mockRepository,
		Clock:      mockClock,
	}

	now := time.Now()

	tweet := models.Tweet{
		UserID:    1,
		Timestamp: now,
		Body:      "uala_challenge",
	}

	mockClock.On("Now").Return(now)
	mockRepository.On("CreateTweet", tweet).Return(uuid.Nil, errors.New("Error"))

	_, err := tweetService.Tweet(1, "uala_challenge")
	assert.Equal(t, err, fmt.Errorf("%w from user %d", service.ErrCreatingTweet, 1))
}

func TestTweetErrorGettingListOfFollowers(t *testing.T) {
	mockRepository := repositoryMocks.NewIRepository(t)
	mockClock := utilsMocks.NewIClock(t)

	tweetService := service.TwitterService{
		Repository: mockRepository,
		Clock:      mockClock,
	}

	now := time.Now()

	tweet := models.Tweet{
		UserID:    1,
		Timestamp: now,
		Body:      "uala_challenge",
	}

	mockClock.On("Now").Return(now)
	mockRepository.On("CreateTweet", tweet).Return(uuid.New(), nil)
	mockRepository.On("GetFollowers", uint(1)).Return(nil, errors.New("Error"))

	_, err := tweetService.Tweet(1, "uala_challenge")
	assert.Equal(t, err, fmt.Errorf("%w from user %d", service.ErrorGettingFollowersList, 1))
}

func TestTweetErrorAddToTimeline(t *testing.T) {
	mockRepository := repositoryMocks.NewIRepository(t)
	mockClock := utilsMocks.NewIClock(t)

	tweetService := service.TwitterService{
		Repository: mockRepository,
		Clock:      mockClock,
	}

	now := time.Now()
	tweet := models.Tweet{
		UserID:    1,
		Timestamp: now,
		Body:      "uala_challenge",
	}
	listOfFollowers := []uint{2, 3}

	id := uuid.New()

	mockClock.On("Now").Return(now)
	mockRepository.On("CreateTweet", tweet).Return(id, nil)
	mockRepository.On("GetFollowers", uint(1)).Return(listOfFollowers, nil)
	mockRepository.On("AddTweetToTimeline", id, uint(2)).Return(errors.New("Error"))

	_, err := tweetService.Tweet(1, "uala_challenge")
	assert.Equal(t, err, service.ErrorAddingToTimeline)
}

func TestTweetSuccess(t *testing.T) {
	mockRepository := repositoryMocks.NewIRepository(t)
	mockClock := utilsMocks.NewIClock(t)

	tweetService := service.TwitterService{
		Repository: mockRepository,
		Clock:      mockClock,
	}

	now := time.Now()
	tweet := models.Tweet{
		UserID:    1,
		Timestamp: now,
		Body:      "uala_challenge",
	}
	listOfFollowers := []uint{2}

	id := uuid.New()

	mockClock.On("Now").Return(now)
	mockRepository.On("CreateTweet", tweet).Return(id, nil)
	mockRepository.On("GetFollowers", uint(1)).Return(listOfFollowers, nil)
	mockRepository.On("AddTweetToTimeline", id, uint(2)).Return(nil)

	tweetId, err := tweetService.Tweet(uint(1), "uala_challenge")
	assert.Equal(t, tweetId, id)
	assert.Equal(t, err, nil)
}

func TestGetTimelineErrorRepository(t *testing.T) {
	mockRepository := repositoryMocks.NewIRepository(t)

	timelineService := service.TwitterService{
		Repository: mockRepository,
	}

	mockRepository.On("GetTimeLine", uint(1)).Return(nil, errors.New("Error"))

	_, err := timelineService.GetTimeLine(1)
	assert.Equal(t, err, service.ErrorGettingTimeline)
}

func TestGetTimelineShorterThanLimit(t *testing.T) {
	mockRepository := repositoryMocks.NewIRepository(t)
	mockClock := utilsMocks.NewIClock(t)

	timelineService := service.TwitterService{
		Repository: mockRepository,
		Clock:      mockClock,
	}

	now := time.Now()
	id := uuid.New()
	tweetsTimeLine := []uuid.UUID{id}
	tweet := models.Tweet{
		UserID:    1,
		Timestamp: now,
		Body:      "uala_challenge",
	}
	tweetArray := []models.Tweet{tweet}

	mockRepository.On("GetTimeLine", uint(1)).Return(tweetsTimeLine, nil)
	mockRepository.On("GetTweets", tweetsTimeLine).Return(tweetArray, nil)

	tweets, err := timelineService.GetTimeLine(1)
	assert.Equal(t, tweets, tweetArray)
	assert.Equal(t, err, nil)

}

func TestGetTimelineLargerThanLimit(t *testing.T) {
	mockRepository := repositoryMocks.NewIRepository(t)
	mockClock := utilsMocks.NewIClock(t)

	timelineService := service.TwitterService{
		Repository: mockRepository,
		Clock:      mockClock,
	}

	now := time.Now()
	tweetsTimeLine := []uuid.UUID{uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New()}
	tweet := models.Tweet{
		UserID:    1,
		Timestamp: now,
		Body:      "uala_challenge",
	}
	tweetArray := []models.Tweet{tweet, tweet, tweet, tweet, tweet, tweet, tweet, tweet, tweet, tweet}

	mockRepository.On("GetTimeLine", uint(1)).Return(tweetsTimeLine, nil)
	mockRepository.On("GetTweets", tweetsTimeLine[0:10]).Return(tweetArray, nil)

	tweets, err := timelineService.GetTimeLine(1)
	assert.Equal(t, tweets, tweetArray)
	assert.Equal(t, err, nil)
}
