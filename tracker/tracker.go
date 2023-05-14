package tracker

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	"github.com/things-go/go-socks5"
	"github.com/things-go/go-socks5/bufferpool"
)

type usageTracker struct {
	rule          UsageLimitRule
	authenticator socks5.Authenticator
	bufferPool    bufferpool.BufPool

	perUserLimit int
	globalLimit  int

	globalUsage  uint64
	perUserUsage map[string]uint64
}

type UsageTracker interface {
	Dial(ctx context.Context, network, addr string) (net.Conn, error)
	Connect(ctx context.Context, writer io.Writer, req *socks5.Request) error
	Limit() socks5.RuleSet
	TrackGlobal(size int)
	TrackUser(user string, size int)
	BufferPool() bufferpool.BufPool
	LogUsage(interval int)
	HasGlobalLimitExceeded() bool
	HasUserLimitExceeded(user string) bool
	Observe(port string)
}

func NewUsageTracker(perUserLimit, globalLimit int, authenticator socks5.Authenticator) UsageTracker {
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
func (u *usageTracker) LogUsage(interval int) {
	for range time.Tick(time.Duration(interval) * time.Second) {
		l := u.getLogs()
		log.Println("Global usage:", l.GlobalUsage)
		for user, usage := range l.PerUser {
			log.Printf("User %s used %d bytes", user, usage)
		}
	}
}

type logs struct {
	GlobalUsage uint64            `json:"global"`
	PerUser     map[string]uint64 `json:"perUser"`
}

func (u *usageTracker) getLogs() logs {
	return logs{
		GlobalUsage: u.globalUsage,
		PerUser:     u.perUserUsage,
	}
}

func (u *usageTracker) HasGlobalLimitExceeded() bool {
	return uint64(u.globalLimit) < u.globalUsage
}

func (u *usageTracker) HasUserLimitExceeded(username string) bool {
	return uint64(u.perUserLimit) < u.perUserUsage[username]
}
