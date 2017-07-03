package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadToken_NoFile(t *testing.T) {
	_, err := loadTokenFromFile("non-existent")
	assert.Error(t, err)
}
