package netutil

import (
	"io"
	"testing"
	"time"
)

func TestNewTestAddr(t *testing.T) {
	a := NewTestAddr("TCP", "IP")
	if a.Network() != "TCP" || a.String() != "IP" {
		t.Errorf("Expected [TCP, IP] but got %v", a)
	}
}

func TestNewTestConn(t *testing.T) {
	l, r := NewTestAddr("TCP", "local"), NewTestAddr("TCP", "remote")
	c := NewTestConn(l, r)
	c.R.WriteString("hello, world")
	b := make([]byte, 12)

	n, err := c.Read(b)
	if n != 12 || err != nil {
		t.Errorf("Read() = %v, %v; expected %v, %v", n, err, 12, nil)
	}
	if string(b) != "hello, world" {
		t.Errorf("didn't read 'hello, world' from the buffer")
	}
	c.ReadErr = io.EOF
	n, err = c.Read(b)
	if n != 0 && err != io.EOF {
		t.Errorf("Read() = %v, %v; expected %v, %v", n, err, 0, io.EOF)
	}

	n, err = c.Write([]byte("g'day, mate"))
	if n != 11 || err != nil {
		t.Errorf("Write() = %v, %v; expected %v, %v", n, err, 11, nil)
	}
	if c.W.String() != "g'day, mate" {
		t.Errorf("didn't read 'g'day mate' from the buffer")
	}
	c.WriteErr = io.EOF
	n, err = c.Write(b)
	if n != 0 && err != io.EOF {
		t.Errorf("Write() = %v, %v; expected %v, %v", n, err, 0, io.EOF)
	}

	if c.LocalAddr() != l {
		t.Errorf("LocalAddr() = %v, expected %v", c.LocalAddr(), l)
	}
	if c.RemoteAddr() != r {
		t.Errorf("RemoteAddr() = %v, expected %v", c.RemoteAddr(), l)
	}
	c.Close()

	now := time.Now()
	c.SetDeadline(now)
	if c.RDeadline != now || c.WDeadline != now {
		t.Errorf("SetDeadline() didn't set read or write deadlines")
	}

	now = time.Now()
	c.SetReadDeadline(now)
	if c.RDeadline != now {
		t.Errorf("SetReadDeadline() didn't set read deadlines")
	}

	now = time.Now()
	c.SetWriteDeadline(now)
	if c.WDeadline != now {
		t.Errorf("SetWriteDeadline() didn't set write deadlines")
	}
}
