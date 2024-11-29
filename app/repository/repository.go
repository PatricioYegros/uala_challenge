package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/PatricioYegros/uala_challenge/app/models"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type IRepository interface {
	// AddFollower adds the newFollowerID to the list of followers of userID
	AddFollower(userID, newFollowerID uint) error
	//GetFollowers returns the list of followers of userID
	GetFollowers(userID uint) ([]uint, error)
	//CreateTweet creates a new tweet and returns the uuid
	CreateTweet(tweet models.Tweet) (uuid.UUID, error)
	//GetTweets returns the list of tweets by ids
	GetTweets(ids []uuid.UUID) ([]models.Tweet, error)
	//AddTweetToTimeline adds a tweetID to the user's timeline. Max 10 tweets.
	AddTweetToTimeline(tweetID uuid.UUID, userID uint) error
	//GetTimeLine returns the list of tweets ids in a user timeline
	GetTimeLine(userID uint) ([]uuid.UUID, error)
}

type Repository struct {
	Redis *redis.Client
}

const (
	TweetTTL            = 24 * time.Hour
	MaxTweetsInTimeline = 10
)

// AddFollower adds the newFollowerID to the list of followers of userID
func (repository Repository) AddFollower(userID, newFollowerID uint) error {
	followersKey := UserFollowersKey(userID)

	return repository.Redis.SAdd(context.Background(), followersKey, newFollowerID).Err()
}

// GetFollowers returns the list of followers of userID
func (repository Repository) GetFollowers(userID uint) ([]uint, error) {
	followersKey := UserFollowersKey(userID)

	idsString, err := repository.Redis.SMembers(context.Background(), followersKey).Result()
	if err != nil {
		return nil, err
	}

	ids := make([]uint, 0, len(idsString))

	for _, idStr := range idsString {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, err
		}
		ids = append(ids, uint(id))
	}

	return ids, nil
}

// CreateTweet creates a new tweet and returns the uuid
func (repository Repository) CreateTweet(tweet models.Tweet) (uuid.UUID, error) {
	tweetID := uuid.New()
	tweetKey := TweetKey(tweetID)

	return tweetID, repository.Redis.Set(context.Background(), tweetKey, tweet, TweetTTL).Err()
}

// GetTweets returns the list of tweets by ids
func (repository Repository) GetTweets(ids []uuid.UUID) ([]models.Tweet, error) {
	tweets := make([]models.Tweet, 0, len(ids))

	for _, id := range ids {
		tweetKey := TweetKey(id)

		tweetString, err := repository.Redis.Get(context.Background(), tweetKey).Result()
		if err == redis.Nil {
			//ttl reached
			break
		} else if err != nil {
			return nil, err
		}

		tweet := models.Tweet{}
		err = json.Unmarshal([]byte(tweetString), &tweet)
		if err != nil {
			return nil, err
		}

		tweets = append(tweets, tweet)
	}

	return tweets, nil
}

// AddTweetToTimeline adds a tweetID to the user's timeline. Max 10 tweets.
func (repository Repository) AddTweetToTimeline(tweetID uuid.UUID, userID uint) error {
	timelineKey := TimelineKey(userID)

	timelineLen, err := repository.Redis.LLen(context.Background(), timelineKey).Result()
	if err != nil {
		return err
	}

	if timelineLen == MaxTweetsInTimeline {
		err := repository.Redis.RPop(context.Background(), timelineKey).Err()
		if err != nil {
			return err
		}
	}
	return repository.Redis.LPush(context.Background(), timelineKey, tweetID.String()).Err()
}

// GetTimeLine returns the list of tweets ids in a user timeline
func (repository Repository) GetTimeLine(userID uint) ([]uuid.UUID, error) {
	timelineKey := TimelineKey(userID)

	idsString, err := repository.Redis.LRange(context.Background(), timelineKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, 0, len(idsString))

	for _, idStr := range idsString {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

// UserFollowersKey returns the key that stores the list of followers of userID in the cache
func UserFollowersKey(userID uint) string {
	return fmt.Sprintf("%d-followers", userID)
}

// TweetKey returns the key that stores a tweet by id
func TweetKey(tweetID uuid.UUID) string {
	return fmt.Sprintf("tweet-%s", tweetID)
}

// TimelineKey returns the key that stores an user's timeline
func TimelineKey(userID uint) string {
	return fmt.Sprintf("tl-%d", userID)
}
