package cryptoutil

import (
	"crypto/tls"
	"crypto/x509"
	"reflect"
	"testing"
)

func TestTLSConfig(t *testing.T) {
	key := []byte(`-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQC0elWPpcv7jt0s
hzi5eV6T9RX2L+pMG7xMA+ZXd6yn5PgTAAAvBM/tT0D4wPGKelTaw/0X8ixVTIgp
ji89ZNWfymP4PXRIDtoB4B0byBP/N7/q/xIzWFtR6bbwWtackDTWxxbrkAoklw2I
e04zwshvT/ClDth0oSH2QS8hibQdT7C+JUGoPZrrHDGYfntNn6VSdwcxNQzP0wUo
H+cgf7EvCYDEVQq48Tdk+qYUpNSxqC7QyyEOFMdCqI1SlM+x6PTm5aIL7NTHdRi9
/taNBULYBjRPoSvcMC6/6t8T5AS8y/KHaMpc9WQpBVFQQLsj1ZjgCwDukKcSHkGc
U+npaNoxAgMBAAECggEBAISEYvjD24BNiTcd3tfJN1nalpKa4iWaI+uI3YQR+nOZ
G1IQKRJdLTNpgyJjwbdVVaMAT4Far5S+SiiBH0ysEnNuz3LB5PTX+tlvrs/sXEqE
q+Wn/rw2v27o9guMF5MEC9g8fSbgL6JoS2aQa350ImohP2hi+yq/+cjwWeP9UYRG
95wrhSheTH7Tt6W4wB/BgL9VStdMU2QdirXm7I+mVTIA/gsHRiw03KBfE3k/+WAG
LUWONjRq54Q/xXIsX2ibhATAe4YanntFFBtucnmKZWy3DxamwIEhSBZjxckRSbAa
l98/29LVs1BHREkySyiAeEp49j7cFd+hsMoKPtWSngECgYEA6WaJHJxe/DiaqKEi
998r7XdkYobHX/gAaHQRlTbi2CFxDmaPf2HyQHpdj/27dlkZ9MwUGp9Chu2aQrF+
NgTRKf46bkMCqW+UlFGXOtJNwh3u80o6Bsw0MLszqX2ANNBnXdONogKntdSk+DU/
pBlp1Sax4rOibWKqy/6alWi9NdECgYEAxfP3k5Q43MvLuATEUnAF/J3bNBEBm6Ac
4m3goUhh4fZR4LIxG1OIUgulh6LdoSSTtbwRno7c0+EmLfn4HTIEVZZRqqD3d0vg
q2e3BPEIrgJvE4mM9plPA4l0LAtbz/VkiHpscxwkhnev2e+CxDy5CZtZBLzynZ3F
GMN8DPH/lmECgYAdpfXC2aZDOdmDlDyB6iqFVsY8scw/x5wdqKjXIIq+U7kt0B9o
WgPQ5vpdoSigBv+CFi5zl0l3JFWVJUDYEdKUGARKH18GMPNZx+eDBFX2U6EX7JTk
aiozCoSsghG4I8UNqSLzsoylT2bJVNJJwSSte+Xa22DFrmPt1+DBG/X1EQKBgQCn
zwLQ67Usj+MaE4huLiMVHKjHwabwS9JQrT7g2qCH0q1kYwq4FJ8all7z1cA2K/C/
/jedh5RyVYptLVwFO/Jqr6x5jk1ap0tFYv3GxaJLCSsqj8+c+Sf/YpXGBLcHWwqn
m8i16GSaTXoYsS7UtnlSSIw1NQwjS6zbKlTOEJRP4QKBgQDPv3hZjOadpdf2dfex
xsVUpW/kWuVTf2z6DsIJsJMnwbz+qqFY+xwnGvwiKHEFoWXeMYuXJsmo5/q6ZJLj
+SdDdbXZjXdcGfSl5MU0vFZSEjPC39Pq4icLwe2vfNIff9XETlby4AljiO4SLMaX
pACnfSkMwouZSXK4saepqcbFeQ==
-----END PRIVATE KEY-----`)

	cert := []byte(`-----BEGIN CERTIFICATE-----
MIID/jCCAuagAwIBAgIJAJUU1sB/Ld0DMA0GCSqGSIb3DQEBCwUAMIGTMQswCQYD
VQQGEwJVUzERMA8GA1UECAwIVmlyZ2luaWExEDAOBgNVBAcMB0xhbmdsZXkxFDAS
BgNVBAoMC0FCQyBXaWRnaXRzMRkwFwYDVQQLDBBJbnRlcm5hbCBBZmZhaXJzMREw
DwYDVQQDDAh0bHMudGVzdDEbMBkGCSqGSIb3DQEJARYMdGxzQHRlc3QuY29tMB4X
DTE4MDcxNjE5MDc0M1oXDTE5MDcxNjE5MDc0M1owgZMxCzAJBgNVBAYTAlVTMREw
DwYDVQQIDAhWaXJnaW5pYTEQMA4GA1UEBwwHTGFuZ2xleTEUMBIGA1UECgwLQUJD
IFdpZGdpdHMxGTAXBgNVBAsMEEludGVybmFsIEFmZmFpcnMxETAPBgNVBAMMCHRs
cy50ZXN0MRswGQYJKoZIhvcNAQkBFgx0bHNAdGVzdC5jb20wggEiMA0GCSqGSIb3
DQEBAQUAA4IBDwAwggEKAoIBAQC0elWPpcv7jt0shzi5eV6T9RX2L+pMG7xMA+ZX
d6yn5PgTAAAvBM/tT0D4wPGKelTaw/0X8ixVTIgpji89ZNWfymP4PXRIDtoB4B0b
yBP/N7/q/xIzWFtR6bbwWtackDTWxxbrkAoklw2Ie04zwshvT/ClDth0oSH2QS8h
ibQdT7C+JUGoPZrrHDGYfntNn6VSdwcxNQzP0wUoH+cgf7EvCYDEVQq48Tdk+qYU
pNSxqC7QyyEOFMdCqI1SlM+x6PTm5aIL7NTHdRi9/taNBULYBjRPoSvcMC6/6t8T
5AS8y/KHaMpc9WQpBVFQQLsj1ZjgCwDukKcSHkGcU+npaNoxAgMBAAGjUzBRMB0G
A1UdDgQWBBSni9ENGgGI0OptcpSiMiXZaHokBzAfBgNVHSMEGDAWgBSni9ENGgGI
0OptcpSiMiXZaHokBzAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IB
AQBa6yiQq7r/y34DiO4MX7WNdnWgwTzcwm1kBZAblg8jTxoJn4j4E2xpGQjVvwI0
IJ/R/0REc4ldf1XIwDvIpIeXF9F0cIeP3emy6HmL6sEotFiwU2godoHuUkvE/sD6
ViIns72+vS3OOzMxOh0CErIr9z2llK2zyTgk2bZ+3poQ+bDguth2nTPI9QpJz4SP
SLfK28Hgk5g3+mA6lD32hgha4aD+1euaiQiuPwXJ0waRSrqeyGbFv6WFFId0ah79
ahXCjNx7nODJLQtQr5hiQvrsFjR5G1xcPrA9i3qTZSeQBW15m2K32bCYKF/a7aVi
qh1N6UIFiU8p32KlW/YU2OvZ
-----END CERTIFICATE-----`)

	c, _ := tls.X509KeyPair(cert, key)
	p := x509.NewCertPool()
	p.AppendCertsFromPEM(cert)

	tests := []struct {
		name string
		cert []byte
		key  []byte
		pool []byte
		cfg  *tls.Config
		err  error
	}{
		{
			name: "normal",
			cert: cert,
			key:  key,
			pool: cert,
			cfg: &tls.Config{
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
			},
		},
		{
			name: "bad cert/key",
			pool: cert,
			err:  ErrInvalidKeyPair,
		},
		{
			name: "empty pool",
			cert: cert,
			key:  key,
			err:  ErrNoCertsFound,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg, err := TLSConfig(test.cert, test.key, test.pool)
			if err != test.err {
				t.Fatalf(`expecting err "%v", got "%v"`, test.err, err)
			}
			if !reflect.DeepEqual(cfg, test.cfg) {
				t.Fatalf(`expecting cfg\n%v\ngot\n%v`, test.cfg, cfg)
			}
		})
	}
}
