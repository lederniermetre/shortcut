package iteration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterationCommandStructure(t *testing.T) {
	cmd := NewCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "iteration", cmd.Use)
	assert.Contains(t, cmd.Short, "Work on iteration")
}
