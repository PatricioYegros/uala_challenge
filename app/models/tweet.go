package models

import (
	"encoding/json"
	"errors"
	"time"
)

type Tweet struct {
	UserID    uint      `json:"userId"`
	Timestamp time.Time `json:"timestamp"`
	Body      string    `json:"body"`
}

const MaxLength = 150

var ErrMaxLengthExceeded = errors.New("max Length of 150 exceeded")

// Creates New Tweet
// Returns ErrMaxLengthExceeded if content length is bigger than 150 characters.
func NewTweet(userID uint, timestamp time.Time, content string) (*Tweet, error) {
	if len(content) > MaxLength {
		return nil, ErrMaxLengthExceeded
	}

	return &Tweet{
		UserID:    userID,
		Timestamp: timestamp,
		Body:      content,
	}, nil
}

// Implement encoding.BinaryMarshaler for Redis
func (tweet Tweet) MarshalBinary() (data []byte, err error) {
	return json.Marshal(tweet)
}
