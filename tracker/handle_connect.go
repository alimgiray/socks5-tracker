package tracker

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/things-go/go-socks5"
	"github.com/things-go/go-socks5/statute"
)

// Connect is copied from the library because there was no other way to modify it
func (u *usageTracker) Connect(ctx context.Context, writer io.Writer, req *socks5.Request) error {

	username := req.AuthContext.Payload[usernameKey.String()]
	ctx = PutUsername(ctx, username)

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
