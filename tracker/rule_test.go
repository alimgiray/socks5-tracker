package tracker_test

import (
	"context"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/things-go/go-socks5"
	"github.com/things-go/go-socks5/bufferpool"

	"github.com/alimgiray/socks5-tracker/tracker"
)

var usernameKey = "username"

type MockUsageTracker struct {
	mock.Mock
}

func (m *MockUsageTracker) HasGlobalLimitExceeded() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockUsageTracker) HasUserLimitExceeded(user string) bool {
	args := m.Called(user)
	return args.Bool(0)
}

func (m *MockUsageTracker) Dial(ctx context.Context, network, addr string) (net.Conn, error) {
	args := m.Called(ctx, network, addr)
	return args.Get(0).(net.Conn), args.Error(1)
}

func (m *MockUsageTracker) Connect(ctx context.Context, writer io.Writer, req *socks5.Request) error {
	args := m.Called(ctx, writer, req)
	return args.Error(0)
}

func (m *MockUsageTracker) Limit() socks5.RuleSet {
	args := m.Called()
	return args.Get(0).(socks5.RuleSet)
}

func (m *MockUsageTracker) TrackGlobal(size int) {
	m.Called(size)
}

func (m *MockUsageTracker) TrackUser(user string, size int) {
	m.Called(user, size)
}

func (m *MockUsageTracker) BufferPool() bufferpool.BufPool {
	args := m.Called()
	return args.Get(0).(bufferpool.BufPool)
}

func (m *MockUsageTracker) LogUsage(interval int) {
	m.Called(interval)
}

func TestAllowGlobalLimitExceeded(t *testing.T) {
	mockTracker := new(MockUsageTracker)
	rule := tracker.NewUsageLimitRule(mockTracker)

	mockTracker.On("HasGlobalLimitExceeded").Return(true)
	_, ok := rule.Allow(context.Background(), &socks5.Request{AuthContext: &socks5.AuthContext{Payload: map[string]string{usernameKey: "testuser"}}})
	assert.False(t, ok)
}

func TestAllowUserLimitExceeded(t *testing.T) {
	mockTracker := new(MockUsageTracker)
	rule := tracker.NewUsageLimitRule(mockTracker)

	mockTracker.On("HasGlobalLimitExceeded").Return(false)
	mockTracker.On("HasUserLimitExceeded", "testuser").Return(true)
	_, ok := rule.Allow(context.Background(), &socks5.Request{AuthContext: &socks5.AuthContext{Payload: map[string]string{usernameKey: "testuser"}}})
	assert.False(t, ok)
}

func TestAllowNoLimitExceeded(t *testing.T) {
	mockTracker := new(MockUsageTracker)
	rule := tracker.NewUsageLimitRule(mockTracker)

	mockTracker.On("HasGlobalLimitExceeded").Return(false)
	mockTracker.On("HasUserLimitExceeded", "testuser").Return(false)
	_, ok := rule.Allow(context.Background(), &socks5.Request{AuthContext: &socks5.AuthContext{Payload: map[string]string{usernameKey: "testuser"}}})
	assert.True(t, ok)
}
