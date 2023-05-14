package main

import (
	"github.com/things-go/go-socks5"

	"github.com/alimgiray/socks5-tracker/tracker"
)

func main() {

	perUserLimit := 2 * 1024 * 1024
	globalLimit := 5 * 1024 * 1024
	trackingIntervalInSeconds := 3

	authenticator := socks5.UserPassAuthenticator{
		Credentials: NewStaticCredentials(),
	}

	usageTracker := tracker.NewUsageTracker(perUserLimit, globalLimit, authenticator)

	server := socks5.NewServer(
		socks5.WithDial(usageTracker.Dial),               // for intercepting requests & responses
		socks5.WithConnectHandle(usageTracker.Connect),   // connect handle called before dial, we can get username from here
		socks5.WithRule(usageTracker.Limit()),            // for limiting usage
		socks5.WithBufferPool(usageTracker.BufferPool()), // custom buffer pool needed for custom connect handle
		socks5.WithAuthMethods([]socks5.Authenticator{
			authenticator,
		}),
	)

	go usageTracker.TrackUsage(trackingIntervalInSeconds)

	if err := server.ListenAndServe("tcp", ":8000"); err != nil {
		panic(err)
	}
}

type staticCredentials map[string]string

func NewStaticCredentials() staticCredentials {
	return map[string]string{
		"alim":  "pass",
		"giray": "pass",
	}
}

func (s staticCredentials) Valid(user, password, _ string) bool {
	pass, ok := s[user]
	return ok && password == pass
}
