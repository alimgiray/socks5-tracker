package tracker_test

import (
	"io"
	"testing"

	"github.com/things-go/go-socks5"

	"github.com/alimgiray/socks5-tracker/tracker"
)

type MockAuthenticator struct {
	creds map[string]string
}

func NewMockAuthenticator() *MockAuthenticator {
	return &MockAuthenticator{
		creds: map[string]string{
			"test_user": "pass",
		},
	}
}

func (a *MockAuthenticator) Valid(user, password, _ string) bool {
	pass, ok := a.creds[user]
	return ok && password == pass
}

func (a *MockAuthenticator) Authenticate(io.Reader, io.Writer, string) (*socks5.AuthContext, error) {
	// For the sake of simplicity, let's assume all credentials are valid in this mock
	return &socks5.AuthContext{}, nil
}

func (a *MockAuthenticator) GetCode() uint8 {
	return 0
}

func TestUsageTracker_TrackGlobal(t *testing.T) {
	authenticator := NewMockAuthenticator()
	tr := tracker.NewUsageTracker(100, 200, authenticator)

	tr.TrackGlobal(50)
	if tr.HasGlobalLimitExceeded() {
		t.Errorf("Global limit exceeded too early")
	}

	tr.TrackGlobal(160)
	if !tr.HasGlobalLimitExceeded() {
		t.Errorf("Global limit should be exceeded")
	}
}

func TestUsageTracker_TrackUser(t *testing.T) {
	authenticator := NewMockAuthenticator()
	tr := tracker.NewUsageTracker(100, 200, authenticator)
	user := "test_user"

	tr.TrackUser(user, 50)
	if tr.HasUserLimitExceeded(user) {
		t.Errorf("User limit exceeded too early")
	}

	tr.TrackUser(user, 60)
	if !tr.HasUserLimitExceeded(user) {
		t.Errorf("User limit should be exceeded")
	}
}
