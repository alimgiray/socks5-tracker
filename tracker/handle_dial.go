package tracker

import (
	"context"
	"net"
)

func (u *usageTracker) Dial(ctx context.Context, network, addr string) (net.Conn, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	return &usageConn{Conn: conn, user: GetUsername(ctx), Tracker: u}, nil
}
