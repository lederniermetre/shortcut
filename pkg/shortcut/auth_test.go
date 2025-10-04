package shortcut

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAuth_MissingToken(t *testing.T) {
	// Save original token
	originalToken := os.Getenv("SHORTCUT_API_TOKEN")
	defer func() {
		if originalToken != "" {
			_ = os.Setenv("SHORTCUT_API_TOKEN", originalToken)
		}
	}()

	// Unset token
	_ = os.Unsetenv("SHORTCUT_API_TOKEN")

	// Test that GetAuth returns an error when token is missing
	auth, err := GetAuth()
	assert.Error(t, err)
	assert.Nil(t, auth)
	assert.Contains(t, err.Error(), "SHORTCUT_API_TOKEN")
}

func TestGetAuth_WithToken(t *testing.T) {
	// Save original token
	originalToken := os.Getenv("SHORTCUT_API_TOKEN")
	defer func() {
		if originalToken != "" {
			_ = os.Setenv("SHORTCUT_API_TOKEN", originalToken)
		} else {
			_ = os.Unsetenv("SHORTCUT_API_TOKEN")
		}
	}()

	// Set a test token
	_ = os.Setenv("SHORTCUT_API_TOKEN", "test-token-123")

	// Test that GetAuth returns auth without error
	auth, err := GetAuth()
	assert.NoError(t, err)
	assert.NotNil(t, auth)
}

func TestGetClient(t *testing.T) {
	// Clear client first
	clientSC = nil

	client1 := GetClient()
	assert.NotNil(t, client1)

	client2 := GetClient()
	assert.Equal(t, client1, client2, "GetClient should return the same instance")
}
