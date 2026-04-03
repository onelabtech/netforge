package checker

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"netforge/internal/model"
	"netforge/internal/util"
	"time"
)

// TLSChecker implements TLS handshake inspection
type TLSChecker struct{}

// NewTLSChecker creates a new TLS checker
func NewTLSChecker() *TLSChecker {
	return &TLSChecker{}
}

// Check performs a TLS handshake and returns the certificate details
func (c *TLSChecker) Check(ctx context.Context, target *util.TaskTarget) (*model.TLSResult, error) {
	start := time.Now()
	result := &model.TLSResult{
		Status: model.StatusUnknown,
	}

	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	addr := target.GetAddr()
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{
		InsecureSkipVerify: false,
		ServerName:        target.Host,
	})

	if err != nil {
		// Try again with insecure verify if first fails, just to see the details
		// For the doctor, we want to know *why* it failed, even if it's invalid
		connInsecure, err2 := tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{
			InsecureSkipVerify: true,
			ServerName:        target.Host,
		})
		if err2 != nil {
			result.Error = err.Error()
			result.Status = model.StatusFailure
			result.Duration = time.Since(start)
			return result, nil
		}
		defer connInsecure.Close()

		// Certificate exists but is invalid
		state := connInsecure.ConnectionState()
		c.populateTLSResult(result, state)
		result.Error = err.Error()
		result.Status = model.StatusFailure
		result.ChainValidated = false
	} else {
		defer conn.Close()
		state := conn.ConnectionState()
		c.populateTLSResult(result, state)
		result.Status = model.StatusSuccess
		result.ChainValidated = true
	}

	result.Duration = time.Since(start)
	return result, nil
}

func (c *TLSChecker) populateTLSResult(result *model.TLSResult, state tls.ConnectionState) {
	result.Version = c.tlsVersionToString(state.Version)
	result.CipherSuite = tls.CipherSuiteName(state.CipherSuite)

	if len(state.PeerCertificates) > 0 {
		cert := state.PeerCertificates[0]
		result.Issuer = cert.Issuer.CommonName
		result.Subject = cert.Subject.CommonName
		result.Expiry = cert.NotAfter
		result.DaysUntilExpiry = int(time.Until(cert.NotAfter).Hours() / 24)
		result.SANs = cert.DNSNames
	}
}

func (c *TLSChecker) tlsVersionToString(v uint16) string {
	switch v {
	case tls.VersionTLS13:
		return "TLS 1.3"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionSSL30:
		return "SSL 3.0"
	default:
		return fmt.Sprintf("Unknown (%x)", v)
	}
}
