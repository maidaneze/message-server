package utils

import (
	"time"
)

// Executes a function that returns an error the given number of times until it succeeds
// Returns error in case the function call failed all the times
// The function calls are executed after the intervals

func Retry(fn func() error, count int, interval time.Duration) error {
	if count < 0 {
		return nil
	}
	var err error
	for i := 0; i < count; i++ {
		err = fn()
		if err == nil {
			return err
		}
		time.Sleep(interval)
	}
	return err
}

//Returns UTC current time in milliseconds

func UTCTimeMilliseconds() int64 {
	return time.Now().UTC().UnixNano() / int64(time.Millisecond)
}
