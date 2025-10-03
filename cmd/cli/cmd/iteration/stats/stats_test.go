package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatsCommandStructure(t *testing.T) {
	cmd := NewCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "stats", cmd.Use)
	assert.Contains(t, cmd.Short, "statistics")

	// Test flags
	flags := cmd.PersistentFlags()

	queryFlag := flags.Lookup("query")
	assert.NotNil(t, queryFlag)
	assert.Equal(t, "string", queryFlag.Value.Type())

	limitFlag := flags.Lookup("limit")
	assert.NotNil(t, limitFlag)
	assert.Equal(t, "int", limitFlag.Value.Type())
}

func TestOwnersCommandStructure(t *testing.T) {
	cmd := newOwnersCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "owners", cmd.Use)
	assert.Contains(t, cmd.Short, "owners")
}

func TestStoriesCommandStructure(t *testing.T) {
	cmd := newStoriesCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "stories", cmd.Use)
	assert.Contains(t, cmd.Short, "stories")
}
