package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFriendlyDuration_NoHours(t *testing.T) {
	s := sleep{MinutesAsleep: 24}
	assert.Equal(t, "24 minutes", s.FriendlyDuration())
}

func TestFriendlDuration_SomeHours(t *testing.T) {
	s := sleep{MinutesAsleep: 64}
	d := s.FriendlyDuration()
	assert.Contains(t, d, "1 hour")
	assert.NotContains(t, d, "hours")
	assert.Contains(t, d, "4 minutes")
}

func TestFriendlyDuration_MultiHour(t *testing.T) {
	s := sleep{MinutesAsleep: (2 * 60) + 32}
	d := s.FriendlyDuration()
	assert.Contains(t, d, "2 hours")
	assert.Contains(t, d, "32 minutes")
}

func TestStartTime(t *testing.T) {
	s := sleep{StartTime: "2017-07-01T15:04:05.00"}
	assert.Equal(t, s.Start(), "Saturday, July 1, 15:04 PDT")
}
