package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_NoFile(t *testing.T) {
	_, err := loadConfigFromJSON("nowhere")
	assert.Error(t, err)
}
