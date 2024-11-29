package utils

import "time"

type IClock interface {
	// Now returns the current local time
	Now() time.Time
}

type Clock struct{}

// Now returns the current local time
func (clock Clock) Now() time.Time {
	return time.Now()
}
