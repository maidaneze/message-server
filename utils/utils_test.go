package utils

import (
	"testing"
	"time"

	"errors"

	"github.com/stretchr/testify/assert"
)

func TestRetryShouldReturnAtOnceIfItSucceeds(t *testing.T) {
	cases := []struct {
		count int
	}{
		{0},
		{1},
		{2},
		{3},
		{4},
		{5},
	}

	for _, c := range cases {
		t.Run("testRetryShouldReturnAtOnceIfItSucceeds", func(tt *testing.T) {
			var timesCalled = 0
			Retry(func() error {
				timesCalled++
				return nil
			}, c.count, time.Millisecond)
			assert.True(tt, timesCalled <= 1)
		})
	}
}

func TestRetryShouldReturnErrorIfItAlwaysFails(t *testing.T) {
	cases := []struct {
		count int
	}{
		{2},
		{3},
		{4},
		{5},
	}

	for _, c := range cases {
		t.Run("testRetryShouldReturnErrorIfItAlwaysFails", func(tt *testing.T) {
			var timesCalled = 0
			Retry(func() error {
				timesCalled++
				return errors.New("error")
			}, c.count, time.Millisecond)
			assert.Equal(tt, c.count, timesCalled)
		})
	}
}

func TestRetryShouldReturnAfterItSucceeds(t *testing.T) {
	cases := []struct {
		count int
	}{
		{2},
		{3},
		{4},
		{5},
	}

	for _, c := range cases {
		t.Run("testRetryShouldReturnAfterItSucceeds", func(tt *testing.T) {
			var timesCalled = 0
			Retry(func() error {
				timesCalled++
				if timesCalled >= 1 {
					return nil
				}
				return errors.New("error")
			}, c.count, time.Millisecond)
			assert.Equal(tt, 1, timesCalled)
		})
	}
}
