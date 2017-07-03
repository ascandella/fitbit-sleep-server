package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {
	r := randStringRunes(48)
	assert.Equal(t, 48, len(r))
}
