package cryptoutil

import (
	"crypto/tls"
	"crypto/x509"
	"errors"

	"github.com/icub3d/pkg/telemetry"
)

var (
	// ErrInvalidKeyPair is returned when the cert and key pair could
	// not be parsed.
	ErrInvalidKeyPair = errors.New("invalid key pair")

	// ErrNoCertsFound is returned when no certs were added to a cert
	// pool .
	ErrNoCertsFound = errors.New("no certs found")
)

// TLSConfig generates a TLS config from the given certificate
// information and pool. The cert pool will be used for both the
// client and server CA. The config will be setup with a stricter list
// of curves, ciphers, and TLS versions.
func TLSConfig(cert, key, pool []byte) (*tls.Config, error) {
	c, err := tls.X509KeyPair(cert, key)
	if err != nil {
		telemetry.Errorf("cryptoutil.TLSConfig calling X509KeyPair: %v", err)
		return nil, ErrInvalidKeyPair
	}
	p := x509.NewCertPool()
	if !p.AppendCertsFromPEM(pool) {
		return nil, ErrNoCertsFound
	}
	return &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		Certificates: []tls.Certificate{c},
		ClientCAs:    p,
		RootCAs:      p,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}, nil
}
