package tracker

import (
	"net"
)

type usageConn struct {
	net.Conn
	user    string
	Tracker UsageTracker
}

type UsageConn interface {
	Read(b []byte) (n int, err error)
}

func NewUsageConn(conn net.Conn, user string, tracker UsageTracker) UsageConn {
	return &usageConn{conn, user, tracker}
}

func (c *usageConn) Read(b []byte) (n int, err error) {
	bytes, err := c.Conn.Read(b)
	if err != nil {
		return bytes, err
	}

	c.Tracker.TrackGlobal(bytes)
	c.Tracker.TrackUser(c.user, bytes)

	return bytes, err
}
