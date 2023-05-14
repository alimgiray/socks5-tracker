package tracker

import (
	"context"
	"errors"
	"net"
)

func (u *usageTracker) Dial(ctx context.Context, network, addr string) (net.Conn, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	username, ok := GetUsername(ctx)
	if !ok {
		return nil, errors.New("username is missing")
	}

	return &usageConn{Conn: conn, user: username, Tracker: u}, nil
}
