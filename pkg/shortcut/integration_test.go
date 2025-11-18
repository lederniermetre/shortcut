package shortcut

import (
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
)

// TestGetMemberWithoutAuth tests GetMember without authentication
func TestGetMemberWithoutAuth(t *testing.T) {
	// Save and unset token
	t.Setenv("SHORTCUT_API_TOKEN", "")

	uuid := strfmt.UUID("test-uuid")
	member, err := GetMember(uuid)
	assert.Error(t, err)
	assert.Nil(t, member)
	assert.Contains(t, err.Error(), "authentication failed")
}

// TestGetWorkflowWithoutAuth tests GetWorkflow without authentication
func TestGetWorkflowWithoutAuth(t *testing.T) {
	// Save and unset token
	t.Setenv("SHORTCUT_API_TOKEN", "")

	workflow, err := GetWorkflow(123)
	assert.Error(t, err)
	assert.Nil(t, workflow)
	assert.Contains(t, err.Error(), "authentication failed")
}

// TestGetEpicWithoutAuth tests GetEpic without authentication
func TestGetEpicWithoutAuth(t *testing.T) {
	// Save and unset token
	t.Setenv("SHORTCUT_API_TOKEN", "")

	epic, err := GetEpic(456)
	assert.Error(t, err)
	assert.Nil(t, epic)
	assert.Contains(t, err.Error(), "authentication failed")
}

// TestStoriesByIterationWithoutAuth tests StoriesByIteration without authentication
func TestStoriesByIterationWithoutAuth(t *testing.T) {
	// Save and unset token
	t.Setenv("SHORTCUT_API_TOKEN", "")

	stories, err := StoriesByIteration(789)
	assert.Error(t, err)
	assert.Nil(t, stories)
	assert.Contains(t, err.Error(), "authentication failed")
}

// TestRetrieveIterationsWithoutAuth tests RetrieveIterations without authentication
func TestRetrieveIterationsWithoutAuth(t *testing.T) {
	// Save and unset token
	t.Setenv("SHORTCUT_API_TOKEN", "")

	iterations, err := RetrieveIterations("test-query", 10, "")
	assert.Error(t, err)
	assert.Nil(t, iterations)
	assert.Contains(t, err.Error(), "authentication failed")
}

// TestRetrieveIterationsWithInvalidURL tests RetrieveIterations with invalid URL
func TestRetrieveIterationsWithInvalidURL(t *testing.T) {
	// Save and set token
	t.Setenv("SHORTCUT_API_TOKEN", "test-token")

	// Test with invalid URL (contains invalid characters)
	iterations, err := RetrieveIterations("test-query", 10, "://invalid url with spaces")
	assert.Error(t, err)
	assert.Nil(t, iterations)
	assert.Contains(t, err.Error(), "cannot parse URL")
}
