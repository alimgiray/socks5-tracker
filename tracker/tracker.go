package tracker

import (
	"log"
	"time"

	"github.com/things-go/go-socks5"
	"github.com/things-go/go-socks5/bufferpool"
)

type usageTracker struct {
	rule          *usageLimitRule
	authenticator socks5.UserPassAuthenticator
	bufferPool    bufferpool.BufPool

	perUserLimit int
	globalLimit  int

	globalUsage  uint64
	perUserUsage map[string]uint64
}

func NewUsageTracker(perUserLimit, globalLimit int, authenticator socks5.UserPassAuthenticator) *usageTracker {
	usageTracker := &usageTracker{
		authenticator: authenticator,
		bufferPool:    bufferpool.NewPool(10_000_000),
		perUserLimit:  perUserLimit,
		globalLimit:   globalLimit,
		perUserUsage:  make(map[string]uint64),
	}

	usageTracker.rule = &usageLimitRule{tracker: usageTracker}

	return usageTracker
}

func (u *usageTracker) Limit() socks5.RuleSet {
	return u.rule
}

func (u *usageTracker) TrackGlobal(size int) {
	u.globalUsage += uint64(size)
}

func (u *usageTracker) TrackUser(user string, size int) {
	u.perUserUsage[user] += uint64(size)
}

func (u *usageTracker) BufferPool() bufferpool.BufPool {
	return u.bufferPool
}

// TrackUsage periodically prints global usage and per user usage to the console
// interval is time in seconds
func (u *usageTracker) TrackUsage(interval int) {
	for range time.Tick(time.Duration(interval) * time.Second) {
		log.Println("Global usage:", u.globalUsage)
		for user, usage := range u.perUserUsage {
			log.Printf("User %s used %d bytes", user, usage)
		}
	}
}

func (u *usageTracker) HasGlobalLimitExceeded() bool {
	return uint64(u.globalLimit) < u.globalUsage
}

func (u *usageTracker) HasUserLimitExceeded(username string) bool {
	return uint64(u.perUserLimit) < u.perUserUsage[username]
}
