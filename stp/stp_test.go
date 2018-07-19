package stp

import (
	"bytes"
	"encoding/binary"
	"io"
	"math/rand"
	"strings"
	"sync"
	"testing"

	"github.com/icub3d/pkg/cryptoutil"
	"github.com/icub3d/pkg/netutil"
)

var ca = []byte(`-----BEGIN CERTIFICATE-----
MIIFjDCCA3SgAwIBAgIJANc2MoM44vJ4MA0GCSqGSIb3DQEBCwUAMFsxCzAJBgNV
BAYTAlVTMQswCQYDVQQIDAJWQTEQMA4GA1UEBwwHTGFuZ2xleTEZMBcGA1UECgwQ
SW50ZXJuYWwgQWZmYWlyczESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTE4MDcxOTE3
MjE0MFoXDTIxMDUwODE3MjE0MFowWzELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAlZB
MRAwDgYDVQQHDAdMYW5nbGV5MRkwFwYDVQQKDBBJbnRlcm5hbCBBZmZhaXJzMRIw
EAYDVQQDDAlsb2NhbGhvc3QwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoIC
AQC47VRlQCAiQ5UmAa2Dl3zn530p89i57lKgm9ZIEb91PIffYUGj5bwOpKuIjgZB
hugsDbMdKDVm8d63rHFVLe3dsEq4Vk1rGDO4sRvSfqQC0FGTd6T7260FgJxvhYst
AB2NzM8QEI4FN3NX4ULFsgAAEHAUwUfw4Ej+DcFWQp34k11Zi1JxLOpCqJyQuv9p
4O78oQBTre2jpeGj2veOaC3skRG8DB/3PLux/8orGa38KplZ/83yVRf7HAiAbqDT
oL3s6+iNctKai9kxKAHRg0KK2Cl//OCCenXUbFyyaNZSfKPIL+eDeKDCGkZN16f0
/gP53fjWgHEeNuTm5WnQol/1C6vMyglSzU9fwTWiJjk5aynYnSRpxgv2JC83htP1
vQCuoyEqb7uXe8osIGRX2feA4X7b3LU0aAoi5kMKeC8ClFeStQA1aMCKBj4eJFve
FQP9zdZnmRKA77DvLxdfnH5mKiG1dq3hTc6OE4abvCTPX2WDSN8w5V6paKlN2YrB
HokGuYsCnl6vA2c5yrfN2wjsh2m/ZhSZoJ3kIF8LigoNNfKPPgoC54jN4G1RHKfv
OBaUvpm1JAR9w034CfNoLDeP8w+5Mgb1JuDJjg2ye+5qZSnbz22gK7Fu3S6jj8Kt
tSz/iSen3ezm3s/57Lq1qono8ozfHwaVa44YT4S92L9A7wIDAQABo1MwUTAdBgNV
HQ4EFgQUBWrn0S0YDb/K8s1v4jTJ3cgmVzkwHwYDVR0jBBgwFoAUBWrn0S0YDb/K
8s1v4jTJ3cgmVzkwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAgEA
JqKA9Cefl84GJMaT/xwaWXgVXeTqr37LUsPcdClLWdB3LeUprMIPRTgu7zYrXspE
FBf9q94RzYl42yGrVGCC0DP/p302cN/+GkFRE3l86AewuBju0iHX1kzlWBgeE6aB
yr6UPjSXmSQ1DxLCDOm9ytC94hVy/0JXoFDgPepvIJO09RoYYNMnx2HZvCJHJUUa
V43w1jstpMAlH6gaoCCdKzHnMkw47qcDjEGcE014vDXapFuzDO8VPiePWugx6UdC
7DqaOJx8/ksB+tGmHEiBAdaoL8O6c18b4yMdzMPH2tjUq7FDUg3gmEQB0HurG75B
AkBkOvawDCQI/9YcDnZcXNfKG+QFLwoJr0hdE6oXGL/LNppGksi1Ac/7Y8b6GJJe
Lt7F0nuxtQqug4crZ3nXt5FC0bJ4CDoClZzKbc0H3Ft8PIz5kabg0jbdQDjeiLen
h0ujNKE+XBqKiRl7Zsnl0nBlhdGGf8WdbZxWLdOXxgw2XqcRc/Q1rUINhhNT+x4T
DureOa27C/DywbZrbYkv845yOmceFhJ1FWQpZNxiD/xLaMW6DAeOg67E8YxpJUhH
3eymkHFNZnt5XKLNkWgUuzBI6ddxMrVu3nwfW3JhwVrN6gjy4BDB+ejCgjmz9Aed
QSU1X0jI+uVf0/Nzu4G1i/WymJAF0witb58iCsesb4A=
-----END CERTIFICATE-----`)

var cert = []byte(`-----BEGIN CERTIFICATE-----
MIIEMjCCAhoCCQCOsPRYOZ008jANBgkqhkiG9w0BAQsFADBbMQswCQYDVQQGEwJV
UzELMAkGA1UECAwCVkExEDAOBgNVBAcMB0xhbmdsZXkxGTAXBgNVBAoMEEludGVy
bmFsIEFmZmFpcnMxEjAQBgNVBAMMCWxvY2FsaG9zdDAeFw0xODA3MTkxNzIzMDha
Fw0xOTEyMDExNzIzMDhaMFsxCzAJBgNVBAYTAlVTMQswCQYDVQQIDAJWQTEQMA4G
A1UEBwwHTGFuZ2xleTEZMBcGA1UECgwQSW50ZXJuYWwgQWZmYWlyczESMBAGA1UE
AwwJbG9jYWxob3N0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAu+T1
CpLTRfjD3ayj0Au8/q5nnrlw+UfL/qs4FZn3SvjA5isOYZwqat4PYySLf+5i/l10
Nyl2hZqO/ge9KQ4+PIqUnFRZFX6dE9RgGDoVpkEwlhIDkBMnfXumNsTD5TSnHtBd
rbeLDfNMU9llbQefmeLLdgAg5vZR3hMS5rb6GUkfmHPAKcWRIPWhXHEOlH9Ss3AS
HUjeGgexMowFEUI+Io1YPNxjzTL5a9Iz6hUuCsQ6Fw6DIF/DYKqgUUUnPqZMQdal
YsQOusRAEYHD8zBpmiuBQzxU1MVZKocIgM0TJ2KGtj5ym7EC9u0fpyluryiDLcti
TM+r2TEI55E3s9kylQIDAQABMA0GCSqGSIb3DQEBCwUAA4ICAQAESFREz53bhms4
xjLKYCBxZbLF0jJ4gPIT9iUQzkb89Kv8LhX3CxgMs2fke9SAn356pgPzxLtd5dCX
Hp9uM9EEGOiQM+sYdr7czTFqtmJbYXHnoEpwAUk52Dd00tJTHIzeCYEJFqWDqwx4
qpkfSu7YQSgE2BD8GKDpqq4o06enssMy0vqkNGBhSaUbi7v5wUOtPe/dZ7Xdn2KZ
LLUlGqzy4JosGv78tOUK41NhA1SsIHY0pgd0s8Kunc38aJaFeaJ9vLAZjMRI29LG
FuK7Z6ky47l9atfbXijdP+YldyhkPwElhCMUfq3up2xkgwSP2ZbHMso6YziUThKX
RmU3xQrgNpolwpQ1gHjeMckxsRyUiOeoB3LD9ovkZrllc1w6NZxUIghfqxEIWoPv
yd44uadfnzUPYfz5Mu0jVShmFT7wkS+nq8wTR75unvflbwrZGGlPOHSpwVhZ2Rne
MSdxPJyXXiz66AjderD1+RX1K4FaKUg5ioS2hEbGS0fys4A3uDK1WZTdxNRJeJAA
ULrK2AjdEqmOt/UtbeB9Ow1HDgyNwVokqevfCRPW6+R5Z7tz0sVmUikYoaHAAJrv
oIe+M5iB0RRSGqkmvDW7BCNsZ7qncNK3UUAWupwB3cY1ANBfSUlPKhSrExQsf2LY
VyUtHqWcV5sZUzPz1Var7TQURjAiGA==
-----END CERTIFICATE-----`)

var key = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAu+T1CpLTRfjD3ayj0Au8/q5nnrlw+UfL/qs4FZn3SvjA5isO
YZwqat4PYySLf+5i/l10Nyl2hZqO/ge9KQ4+PIqUnFRZFX6dE9RgGDoVpkEwlhID
kBMnfXumNsTD5TSnHtBdrbeLDfNMU9llbQefmeLLdgAg5vZR3hMS5rb6GUkfmHPA
KcWRIPWhXHEOlH9Ss3ASHUjeGgexMowFEUI+Io1YPNxjzTL5a9Iz6hUuCsQ6Fw6D
IF/DYKqgUUUnPqZMQdalYsQOusRAEYHD8zBpmiuBQzxU1MVZKocIgM0TJ2KGtj5y
m7EC9u0fpyluryiDLctiTM+r2TEI55E3s9kylQIDAQABAoIBADOBbr6IQwwCRkbE
5V7JaPRzgIodEV/0E3JrIwOg4g4WJGKz2MxfD62d05/8y5S1F0PdAFeCsS+sj5gJ
LQxNEDhuUJCN+qGOxKZD6FebzIV3h0cMBxf+eCvIsmZ/B6gUB9/FhHhzQNYBJKIz
jL8XgOOhLN7a6CoHyadTmTHPZL17OpIsyLnpZkyo5aFfQY676DnIw+Deu8gzl7qt
nqvQuqCGU/mnDf8+fAkopKz/6M0/v54IAUKT8UQCr0L/QcBIi1tPqjx07qTFieSv
vw8cNvLc8Jl/l5JdUAQL2CjBjh1GAMYO1+LXIf03TR1OHQIyhEAP6OKRMhEuhiSH
MfPGDR0CgYEA6FETHnmZ8ZtkQip/a2LP7RWXDAhj9yoS1xOiTDbHyOaEHksU7CN6
v24chYWkuhcYC67VkTagdzNVVpwmQAKeLJ5ypSe8+6kB9qaa4kWwC8CtX6VviyCH
P3QF5qelvbHJN9U1jPLaBGSy5lYJwdjben7cnFxSs4QtiE6TV/EaagMCgYEAzwyQ
HSAaEG/2UHFM+wQ8SV0GYtbw5CjWOtwTaZMqfsroD0QksgpprUrVu07QVEztImXG
3sC7qoMlnjbQ8grhJ7hSzsRv89ocQFn2vs9pkNz8Sym3oDHUMIWdalflh4zhrs8H
PrKWJTnuRF2+9NZr8ywELWa/R7oxquBdYm7bGYcCgYEA27uLbOCxRb+nRZnyqPxu
FB2+n0f0XVwzM7DDanjJ4HB/+DMl1+/68sNQQM5WLxkWyj4UjGPxxK0OA3xwBU00
lJlNcH60lgeV7thIWMp3J7aXhYbxiNM8eTzhM/fPoSteWevU2br9kPg56fjpA6t8
dFE3ksEUC+1yL6G6ZYyLYkcCgYBU4H+LkQdeSaed3nSCSoA0SvA9QIPL5Hm98035
75xyEdgDuhmY6u/bXFw1xt9cT5S+jx5xXm0QP2cCbtJFBvS6BbPck2aZfoYqUzb1
ja1m91Btt5JRF1w27+9SEupDIRu7P59msmseo5rrZ8rKL8RdEWQ9AGvViZymwgdG
PM+QrQKBgHsHddwtaO7Oxke/ouTPh+zsLGG9WMf+w4bZzrMYHk7xEPmGCMqTqUBG
AlQTsFGJVDB2NtV8CwbCj1bv2TVRcmnHwGne7LK8T+N8YJsEHXQpDxeNbmyDAL/Y
mwbwDW+WZMeKSfJ2jpLFqUrDfj1TANiQ5Yy0X40nat8HTeqERsuu
-----END RSA PRIVATE KEY-----`)

func TestNewServerErrors(t *testing.T) {
	_, err := NewServer("", nil, nil, nil)
	if err != cryptoutil.ErrInvalidKeyPair {
		t.Errorf("sending invalid key pair returned '%v', not '%v'",
			err, cryptoutil.ErrInvalidKeyPair)
	}

	_, err = NewServer(":bad addr", cert, key, ca)
	if err.Error() != "listen tcp: address tcp/bad addr: unknown port" {
		t.Errorf("sending invalid addr returned '%v', not '%v'",
			err, "listen tcp: address tcp/bad addr: unknown port")
	}
}

func TestNewClientErrors(t *testing.T) {
	_, err := NewClient("", nil, nil, nil)
	if err != cryptoutil.ErrInvalidKeyPair {
		t.Errorf("sending invalid key pair returned '%v', not '%v'",
			err, cryptoutil.ErrInvalidKeyPair)
	}

	_, err = NewClient(":bad addr", cert, key, ca)
	if err.Error() != "dial tcp: address tcp/bad addr: unknown port" {
		t.Errorf("sending invalid addr returned '%v', not '%v'",
			err, "dial tcp: address tcp/bad addr: unknown port")
	}
}

func TestClientNext(t *testing.T) {
	tests := []struct {
		name string
		size int
		typ  uint64
		err  error
		buf  func() *bytes.Buffer
	}{
		{
			name: "normal",
			size: 10,
			typ:  8478,
			err:  nil,
			buf: func() *bytes.Buffer {
				buf := &bytes.Buffer{}
				b := make([]byte, 16)
				binary.LittleEndian.PutUint64(b[:8], 10)
				binary.LittleEndian.PutUint64(b[8:], 8478)
				buf.Write(b)
				return buf
			},
		},
		{
			name: "unexpected EOF",
			size: 0,
			typ:  0,
			err:  io.ErrUnexpectedEOF,
			buf: func() *bytes.Buffer {
				buf := &bytes.Buffer{}
				buf.WriteString("asdf")
				return buf
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conn := netutil.NewTestConn(netutil.NewTestAddr("TCP", "local"),
				netutil.NewTestAddr("TCP", "remote"))
			conn.R = test.buf()
			c := &client{
				next:  make([]byte, 16),
				write: make([]byte, 16),
				conn:  conn,
			}
			if c.Conn() != conn {
				t.Errorf("c.Conn() didn't return expected connection")
			}
			size, typ, err := c.Next()
			if err != test.err || size != test.size || typ != test.typ {
				t.Errorf("c.Next() = %v, %v, %v; expected %v, %v, %v",
					size, typ, err, test.size, test.typ, test.err)
			}
		})
	}
}

func TestClientRead(t *testing.T) {
	tests := []struct {
		name string
		p    []byte
		exp  string
		err  error
		buf  func() *bytes.Buffer
	}{
		{
			name: "normal",
			p:    make([]byte, 12),
			exp:  "hello, world",
			err:  nil,
			buf: func() *bytes.Buffer {
				buf := &bytes.Buffer{}
				buf.WriteString("hello, world")
				return buf
			},
		},
		{
			name: "unexpected EOF",
			p:    make([]byte, 20),
			err:  io.ErrUnexpectedEOF,
			buf: func() *bytes.Buffer {
				buf := &bytes.Buffer{}
				buf.WriteString("asdf")
				return buf
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conn := netutil.NewTestConn(netutil.NewTestAddr("TCP", "local"),
				netutil.NewTestAddr("TCP", "remote"))
			conn.R = test.buf()
			c := &client{
				next:  make([]byte, 16),
				write: make([]byte, 16),
				conn:  conn,
			}
			err := c.Read(test.p)
			if err != test.err {
				t.Errorf("c.Read() = %v, %v; expected %v, %v",
					err, string(test.p), test.err, test.exp)
			}
		})
	}
}

func TestClientWrite(t *testing.T) {
	tests := []struct {
		name string
		p    []byte
		what uint64
		err  error
		exp  func() []byte
		conn func(c *netutil.TestConn)
	}{
		{
			name: "normal",
			p:    []byte("hello, world"),
			what: 184871,
			err:  nil,
			exp: func() []byte {
				buf := &bytes.Buffer{}
				b := make([]byte, 16)
				binary.LittleEndian.PutUint64(b[:8], 12)
				binary.LittleEndian.PutUint64(b[8:], 184871)
				buf.Write(b)
				buf.WriteString("hello, world")
				return buf.Bytes()
			},
		},
		{
			name: "err after first write",
			p:    []byte("hello, world"),
			what: 184871,
			err:  io.EOF,
			exp: func() []byte {
				return []byte{}
			},
			conn: func(c *netutil.TestConn) {
				c.WriteErr = io.EOF
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conn := netutil.NewTestConn(netutil.NewTestAddr("TCP", "local"),
				netutil.NewTestAddr("TCP", "remote"))
			c := &client{
				next:  make([]byte, 16),
				write: make([]byte, 16),
				conn:  conn,
			}
			if test.conn != nil {
				test.conn(conn)
			}
			err := c.Write(test.what, test.p)
			if err != test.err || !bytes.Equal(test.exp(), conn.W.Bytes()) {
				t.Errorf("c.Write(%v, %v) = %v, %v; expected %v, %v",
					test.what, test.p, err, conn.W.String(),
					err, conn.W.Bytes())
			}
		})
	}
}

func TestCommunication(t *testing.T) {
	const (
		MessageReverse uint64 = iota
		MessageEcho
	)

	// Start up the server.
	s, err := NewServer("localhost:", cert, key, ca)
	if err != nil {
		t.Fatalf("creating new server: %v", err)
	}
	defer s.Close()
	go func() {
		for {
			c, err := s.Accept()
			if err != nil {
				t.Logf("s.Accept() returned %v", err)
				return
			}
			// Handle the client.
			go func(c Client) {
				for {
					// Get the type and size.
					size, typ, err := c.Next()
					if err == io.EOF {
						return
					} else if err != nil {
						t.Errorf("c.Next() returned %v", err)
						return
					}

					// Get the data.
					b := make([]byte, size)
					err = c.Read(b)
					if err == io.EOF {
						return
					} else if err != nil {
						t.Errorf("c.Read() returned %v", err)
						return
					}

					// Handle the message
					switch typ {
					case MessageEcho:
						err = c.Write(MessageEcho, b)
						if err != nil {
							t.Errorf("writing echo: %v", err)
							return
						}
					case MessageReverse:
						for i := len(b)/2 - 1; i >= 0; i-- {
							opp := len(b) - 1 - i
							b[i], b[opp] = b[opp], b[i]
						}
						err = c.Write(MessageReverse, b)
						if err != nil {
							t.Errorf("writing echo: %v", err)
							return
						}

					}
				}
			}(c)
		}
	}()

	// These are some requests and values we expect to get back
	tests := []struct {
		id  uint64
		msg []byte
		exp []byte
	}{
		{id: MessageEcho, msg: []byte("hello"), exp: []byte("hello")},
		{id: MessageReverse, msg: []byte("hello"), exp: []byte("olleh")},
	}

	// Make some clients
	wg := sync.WaitGroup{}
	for x := 0; x < 100; x++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			addr := s.Addr().String()
			parts := strings.SplitN(addr, ":", 2)
			addr = "localhost:" + parts[1]
			c, err := NewClient(addr, cert, key, ca)
			if err != nil {
				t.Errorf("NewClient() returned %v", err)
				return
			}
			defer c.Close()

			// Send some messages.
			for x := 0; x < 1000; x++ {
				test := tests[rand.Intn(len(tests))]

				// Send the message.
				err := c.Write(test.id, test.msg)
				if err != nil {
					t.Errorf("Write() returned %v", err)
					return
				}

				// Get the id back.
				size, typ, err := c.Next()
				if err != nil || size != len(test.exp) || typ != test.id {
					t.Errorf("Next() returned %v, %v, %v; expected %v, %v, %v",
						size, typ, err, len(test.exp), test.id, nil)
					return
				}

				// Check  the message.
				b := make([]byte, size)
				err = c.Read(b)
				if err != nil {
					t.Errorf("Read() returned %v", err)
					return
				}
				if !bytes.Equal(b, test.exp) {
					t.Errorf("Read() returned %v but expected %v", b, test.exp)
				}
			}
		}()
	}
	wg.Wait()
}
