package tracker_test

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alimgiray/socks5-tracker/tracker"
)

// MockConn is a mock implementation of net.Conn
type MockConn struct {
	mock.Mock
}

func (c *MockConn) Read(b []byte) (n int, err error) {
	args := c.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *MockConn) Write(b []byte) (n int, err error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *MockConn) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockConn) LocalAddr() net.Addr {
	args := m.Called()
	return args.Get(0).(net.Addr)
}

func (m *MockConn) RemoteAddr() net.Addr {
	args := m.Called()
	return args.Get(0).(net.Addr)
}

func (m *MockConn) SetDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockConn) SetReadDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockConn) SetWriteDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func TestUsageConn_Read(t *testing.T) {
	data := []byte("test data")
	mockConn := new(MockConn)
	mockConn.On("Read", data).Return(len(data), nil)

	mockTracker := new(MockUsageTracker)
	mockTracker.On("TrackGlobal", len(data))
	mockTracker.On("TrackUser", "testuser", len(data))

	usageConn := tracker.NewUsageConn(mockConn, "testuser", mockTracker)

	n, err := usageConn.Read(data)

	mockConn.AssertExpectations(t)
	mockTracker.AssertExpectations(t)

	assert.NoError(t, err)
	assert.Equal(t, len(data), n)
}
