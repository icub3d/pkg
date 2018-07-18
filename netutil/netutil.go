// Package netutil provides helper functions for the net package and
// it's subpackages.
package netutil

import (
	"bytes"
	"net"
	"time"
)

// TestAddr implements the net.Addr interface that returns the values in
// this structure.
type TestAddr struct {
	// N is the network name returned by Network()
	N string

	// S is the string to return when String() is called.
	S string
}

// Network implements the net.Addr interface.
func (a *TestAddr) Network() string { return a.N }

// String implements the net.Addr interface.
func (a *TestAddr) String() string { return a.S }

// NewTestAddr creates an Addr using the the given network and addr.
func NewTestAddr(network, addr string) *TestAddr {
	return &TestAddr{N: network, S: addr}
}

// TestConn implents the net.Conn interface and stores/reads
// information from buffers.
type TestConn struct {
	// R contains the data that Read() will return.
	R *bytes.Buffer

	// W contains all the data Write() received.
	W *bytes.Buffer

	// These are the Addrs that are returned via function calls.
	Local, Remote net.Addr

	// These are the deadlines set by function calls.
	RDeadline, WDeadline time.Time
}

// NewTestConn creates a new TestConn using the givne local and remote
// Addr's. The Read() and Write() buffers are created with an empty
// buffer.
func NewTestConn(local, remote net.Addr) *TestConn {
	return &TestConn{
		R:      &bytes.Buffer{},
		W:      &bytes.Buffer{},
		Local:  local,
		Remote: remote,
	}
}

// Read implements the net.Conn interface.
func (c *TestConn) Read(b []byte) (n int, err error) {
	return c.R.Read(b)
}

// Write implements the net.Conn interface.
func (c *TestConn) Write(b []byte) (n int, err error) {
	return c.W.Write(b)
}

// Close implements the net.Conn interface.
func (c *TestConn) Close() error {
	return nil
}

// LocalAddr implements the net.Conn interface.
func (c *TestConn) LocalAddr() net.Addr {
	return c.Local
}

// RemoteAddr implements the net.Conn interface.
func (c *TestConn) RemoteAddr() net.Addr {
	return c.Remote
}

// SetDeadline implements the net.Conn interface.
func (c *TestConn) SetDeadline(t time.Time) error {
	c.RDeadline = t
	c.WDeadline = t
	return nil
}

// SetReadDeadline implements the net.Conn interface.
func (c *TestConn) SetReadDeadline(t time.Time) error {
	c.RDeadline = t
	return nil
}

// SetWriteDeadline implements the net.Conn interface.
func (c *TestConn) SetWriteDeadline(t time.Time) error {
	c.WDeadline = t
	return nil
}
