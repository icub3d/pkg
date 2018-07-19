package stp

import (
	"crypto/tls"
	"encoding/binary"
	"io"
	"net"

	"github.com/icub3d/pkg/cryptoutil"
)

// Client represents a client connecting to a server. This is used on
// both the client and the server side. Using the client is normally
// done in a for loop where you first call Next() to get details about
// the next message, prepare a properly sized []byte, and then Read()
// the message.
type Client interface {
	// Next gets the information about the next incoming message. The
	// first value is the size of the message. You should pass a
	// []byte of that size to the next Read() call. The second value
	// is the type of message. You define this type yourself and then
	// it lets you know what the next message will contain.
	Next() (int, uint64, error)

	// Read puts the next message into the given []byte. You should
	// call Next() before calling this.
	Read([]byte) error

	// Write the given message of the given type. On the server side,
	// this is sending the message back to the client and on the
	// client side, this is sending the message to the server.
	Write(uint64, []byte) error

	// Close closes the underlying connection.
	Close() error

	// Conn returns the underlying connection. This is mostly to allow
	// you to get connection information or set deadlines. Using
	// Read()/Write may cause issues.
	Conn() net.Conn
}

type client struct {
	next  []byte // temporary binary data from the read in Next()
	write []byte // temporary binary data for the write in Write()
	conn  net.Conn
}

func (c *client) Next() (int, uint64, error) {
	// Read our type and size
	_, err := io.ReadFull(c.conn, c.next)
	if err != nil {
		return 0, 0, err
	}
	return int(binary.LittleEndian.Uint64(c.next[:8])),
		binary.LittleEndian.Uint64(c.next[8:]), nil
}

func (c *client) Read(p []byte) error {
	_, err := io.ReadFull(c.conn, p)
	return err
}

func (c *client) Write(what uint64, p []byte) error {
	// Write out the type.
	binary.LittleEndian.PutUint64(c.write[:8], uint64(len(p)))
	binary.LittleEndian.PutUint64(c.write[8:], what)
	_, err := c.conn.Write(c.write)
	if err != nil {
		return err
	}

	// Write out the data.
	_, err = c.conn.Write(p)
	return err
}

func (c *client) Conn() net.Conn {
	return c.conn
}

func (c *client) Close() error {
	return c.conn.Close()
}

// NewClient returns a new client that connects to the given server
// addr. The cert and key are used for TLS authentication and the pool
// is used to validation the servers certificate.
func NewClient(addr string, cert, key, pool []byte) (Client, error) {
	// Load up our certs into the config.
	cfg, err := cryptoutil.TLSConfig(cert, key, pool)
	if err != nil {
		return nil, err
	}

	// Make the connection
	conn, err := tls.Dial("tcp", addr, cfg)
	if err != nil {
		return nil, err
	}
	return &client{
		next:  make([]byte, 16),
		write: make([]byte, 16),
		conn:  conn,
	}, nil
}

// Server represents a server listening for new connections.
type Server interface {
	// Accept waits for a new client connection and then returns a
	// client suitable for reading messages and sending responses back
	// to the client.
	Accept() (Client, error)

	// Close closes the underlying connection and stops Accept() from
	// waiting.
	Close() error

	// Addr returns the server Addr.
	Addr() net.Addr
}

type server struct {
	listener net.Listener
}

func (s *server) Accept() (Client, error) {
	conn, err := s.listener.Accept()
	return &client{
		next:  make([]byte, 16),
		write: make([]byte, 16),
		conn:  conn,
	}, err
}

func (s *server) Close() error {
	return s.listener.Close()
}

func (s *server) Addr() net.Addr {
	return s.listener.Addr()
}

// NewServer returns a new server listening on the given addr. The
// cert and key are used for TLS and the pool is used for verifying
// client certificates.
func NewServer(addr string, cert, key, pool []byte) (Server, error) {
	// Load up our certs into the config.
	cfg, err := cryptoutil.TLSConfig(cert, key, pool)
	if err != nil {
		return nil, err
	}

	// Start our listener
	l, err := tls.Listen("tcp", addr, cfg)
	if err != nil {
		return nil, err
	}
	return &server{
		listener: l,
	}, nil
}
