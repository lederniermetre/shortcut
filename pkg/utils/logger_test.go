package utils

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetLogger(t *testing.T) {
	// Test with debug = false
	SetLogger(false)
	logger := slog.Default()
	assert.NotNil(t, logger)

	// Test with debug = true
	SetLogger(true)
	logger = slog.Default()
	assert.NotNil(t, logger)

	// Test with SC_DEBUG_SRC
	os.Setenv("SC_DEBUG_SRC", "true")
	SetLogger(true)
	logger = slog.Default()
	assert.NotNil(t, logger)
	os.Unsetenv("SC_DEBUG_SRC")
}
