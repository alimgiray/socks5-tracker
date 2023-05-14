package tracker

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/things-go/go-socks5"
	"github.com/things-go/go-socks5/bufferpool"
	"github.com/things-go/go-socks5/statute"
)

type trackerCtxKey string

const usernameKey trackerCtxKey = "username"

func (c trackerCtxKey) String() string {
	return string(c)
}

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

func (u *usageTracker) Dial(ctx context.Context, network, addr string) (net.Conn, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	username := ctx.Value(usernameKey).(string)

	return &usageConn{Conn: conn, user: username, Tracker: u}, nil
}

// Connect is copied from the library because there was no other way to modify it
func (u *usageTracker) Connect(ctx context.Context, writer io.Writer, req *socks5.Request) error {

	username := req.AuthContext.Payload[usernameKey.String()]
	ctx = context.WithValue(ctx, usernameKey, username)

	// Attempt to connect
	dial := u.Dial

	target, err := dial(ctx, "tcp", req.DestAddr.String())
	if err != nil {
		msg := err.Error()
		resp := statute.RepHostUnreachable
		if strings.Contains(msg, "refused") {
			resp = statute.RepConnectionRefused
		} else if strings.Contains(msg, "network is unreachable") {
			resp = statute.RepNetworkUnreachable
		}
		if err := socks5.SendReply(writer, resp, nil); err != nil {
			return fmt.Errorf("failed to send reply, %v", err)
		}
		return fmt.Errorf("connect to %v failed, %v", req.RawDestAddr, err)
	}
	defer target.Close()

	// Send success
	if err := socks5.SendReply(writer, statute.RepSuccess, target.LocalAddr()); err != nil {
		return fmt.Errorf("failed to send reply, %v", err)
	}

	// Start proxying
	errCh := make(chan error, 2)
	go func() { errCh <- u.Proxy(target, req.Reader) }()
	go func() { errCh <- u.Proxy(writer, target) }()
	// Wait
	for i := 0; i < 2; i++ {
		e := <-errCh
		if e != nil {
			// return from this function closes target (and conn).
			return e
		}
	}

	return nil
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

func (u *usageTracker) TrackUsage() {
	for range time.Tick(3 * time.Second) {
		log.Println("Global usage:", u.globalUsage)
		for user, usage := range u.perUserUsage {
			log.Printf("User %s used %d bytes", user, usage)
		}
	}
}

type closeWriter interface {
	CloseWrite() error
}

func (u *usageTracker) Proxy(dst io.Writer, src io.Reader) error {
	buf := u.bufferPool.Get()
	defer u.bufferPool.Put(buf)
	_, err := io.CopyBuffer(dst, src, buf[:cap(buf)])
	if tcpConn, ok := dst.(closeWriter); ok {
		tcpConn.CloseWrite() //nolint: errcheck
	}
	return err
}
