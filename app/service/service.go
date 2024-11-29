package service

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"uala_challenge/app/models"
	"uala_challenge/app/repository"
	"uala_challenge/app/utils"

	"github.com/google/uuid"
)

type TwitterService struct {
	Repository repository.IRepository
	Clock      utils.IClock
}

var (
	ErrEqualsIDs              = errors.New("you cant follow yourself")
	ErrFollowing              = errors.New("error adding follow")
	ErrCreatingTweet          = errors.New("error creating tweet")
	ErrorAddingToTimeline     = errors.New("error adding to followers timeline")
	ErrorGettingTimeline      = errors.New("error getting timeline of user")
	ErrorGettingFollowersList = errors.New("error getting followers list")
	ErrorFollowingAlready     = errors.New("error making an already existant follow")
)

const (
	limitTimeLine = 10
)

// Follow makes followerID to follow userID
// Returns ErrEqualsIDs if followerID is the same as userID or
// ErrFollowing if an error occurred
func (service TwitterService) Follow(followerID, userID uint) error {
	if followerID == userID {
		return ErrEqualsIDs
	}

	listOfFollows, err := service.Repository.GetFollowers(userID)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("%w from user %d for check existent follow", ErrFollowing, userID)
	}

	if slices.Contains(listOfFollows, followerID) {
		return ErrorFollowingAlready
	}

	err = service.Repository.AddFollower(userID, followerID)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("%w from user %d to user %d", ErrFollowing, followerID, userID)
	}

	return nil
}

// Tweet creates a Tweet belonging of userID
// Returns ErrTweetTooLong if content len is bigger than 150 characters or
// ErrTweet if an error occurred
func (service TwitterService) Tweet(userID uint, content string) (uuid.UUID, error) {
	tweet, err := models.NewTweet(userID, service.Clock.Now(), content)
	if err != nil {
		return uuid.Nil, err
	}

	tweetID, err := service.Repository.CreateTweet(*tweet)
	if err != nil {
		log.Println(err.Error())
		return uuid.Nil, fmt.Errorf("%w from user %d", ErrCreatingTweet, userID)
	}

	followers, err := service.Repository.GetFollowers(userID)
	if err != nil {
		log.Println(err.Error())
		return uuid.Nil, fmt.Errorf("%w from user %d", ErrorGettingFollowersList, userID)
	}

	for _, follower := range followers {
		err = service.Repository.AddTweetToTimeline(tweetID, follower)
		if err != nil {
			return uuid.Nil, ErrorAddingToTimeline
		}
	}

	return tweetID, nil
}

// GetTimeline returns the list of tweets in user timeline
// Returns ErrTimeline if an error ocurred
func (service TwitterService) GetTimeLine(userID uint) ([]models.Tweet, error) {
	tweetsIDs, err := service.Repository.GetTimeLine(userID)
	if err != nil {
		return nil, ErrorGettingTimeline
	}

	if len(tweetsIDs) <= limitTimeLine {
		return service.Repository.GetTweets(tweetsIDs)
	}

	return service.Repository.GetTweets(tweetsIDs[0:limitTimeLine])
}
