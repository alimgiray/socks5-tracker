package tracker

import (
	"net"
)

type usageConn struct {
	net.Conn
	user    string
	Tracker *usageTracker
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
