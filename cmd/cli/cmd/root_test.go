package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCommandStructure(t *testing.T) {
	assert.NotNil(t, rootCmd)
	assert.Equal(t, "shortcut", rootCmd.Use)
	assert.Contains(t, rootCmd.Short, "Supercharge")

	// Test flags
	flags := rootCmd.PersistentFlags()
	flag := flags.Lookup("debug")
	assert.NotNil(t, flag)
	assert.Equal(t, "bool", flag.Value.Type())
}

func TestExecuteFunction(t *testing.T) {
	// Test that Execute function exists and can be called
	// We can't fully test it without mocking, but we can verify it's defined
	assert.NotNil(t, Execute)
}
