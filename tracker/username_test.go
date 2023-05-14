package tracker_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alimgiray/socks5-tracker/tracker"
)

func TestUsernameContext(t *testing.T) {
	ctx := context.Background()

	// Test when username is put into the context
	username := "testuser"
	ctx = tracker.PutUsername(ctx, username)

	retrievedUsername, ok := tracker.GetUsername(ctx)
	assert.True(t, ok, "The username should be present in the context")
	assert.Equal(t, username, retrievedUsername, "The retrieved username should match the one put into the context")

	// Test when username is not in the context
	ctx = context.Background()
	_, ok = tracker.GetUsername(ctx)
	assert.False(t, ok, "The username should not be present in the context")
}
