package model

import (
	"time"
)

// CheckStatus represents the status of a specific check
type CheckStatus string

const (
	StatusSuccess CheckStatus = "success"
	StatusWarning CheckStatus = "warning"
	StatusFailure CheckStatus = "failure"
	StatusUnknown CheckStatus = "unknown"
)

// DNSResult holds information about DNS resolution
type DNSResult struct {
	Status    CheckStatus   `json:"status"`
	Target    string        `json:"target"`
	IPs       []string      `json:"ips"`
	IPv6s     []string      `json:"ipv6s"`
	CNAMEs    []string      `json:"cnames"`
	MX        []string      `json:"mx"`
	NS        []string      `json:"ns"`
	TTL       uint32        `json:"ttl"`
	Duration  time.Duration `json:"duration"`
	Error     string        `json:"error,omitempty"`
}

// TCPResult holds information about a TCP connection check
type TCPResult struct {
	Status    CheckStatus   `json:"status"`
	Host      string        `json:"host"`
	Port      int           `json:"port"`
	Duration  time.Duration `json:"duration"`
	Error     string        `json:"error,omitempty"`
}

// HTTPResult holds information about an HTTP request
type HTTPResult struct {
	Status         CheckStatus       `json:"status"`
	URL            string            `json:"url"`
	StatusCode     int               `json:"status_code"`
	Proto          string            `json:"proto"`
	Duration       time.Duration     `json:"duration"`
	TTFB           time.Duration     `json:"ttfb"`
	ContentLength  int64             `json:"content_length"`
	Headers        map[string]string `json:"headers"`
	Redirects      []string          `json:"redirects"`
	Compression    string            `json:"compression"`
	Error          string            `json:"error,omitempty"`
}

// TLSResult holds information about a TLS handshake
type TLSResult struct {
	Status            CheckStatus   `json:"status"`
	Version           string        `json:"version"`
	CipherSuite       string        `json:"cipher_suite"`
	Issuer            string        `json:"issuer"`
	Subject           string        `json:"subject"`
	Expiry            time.Time     `json:"expiry"`
	DaysUntilExpiry   int           `json:"days_until_expiry"`
	SANs              []string      `json:"sans"`
	ChainValidated    bool          `json:"chain_validated"`
	Duration          time.Duration `json:"duration"`
	Error             string        `json:"error,omitempty"`
}

// PathResult holds traceroute-like information
type PathResult struct {
	Status   CheckStatus   `json:"status"`
	Hops     []Hop         `json:"hops"`
	Duration time.Duration `json:"duration"`
	Error    string        `json:"error,omitempty"`
}

// Hop represents a single hop in a path
type Hop struct {
	Number   int           `json:"number"`
	IP       string        `json:"ip"`
	Hostname string        `json:"hostname"`
	Latency  time.Duration `json:"latency"`
}

// ScanResult holds the result of a port scan
type ScanResult struct {
	Target string        `json:"target"`
	Ports  []PortResult  `json:"ports"`
	Error  string        `json:"error,omitempty"`
}

// PortResult represents a single scanned port
type PortResult struct {
	Port    int    `json:"port"`
	Status  string `json:"status"` // open, closed, filtered
	Service string `json:"service"`
	Banner  string `json:"banner,omitempty"`
}

// DiagnosisIssue represents a problem identified by the analyzer
type DiagnosisIssue struct {
	ID             string      `json:"id"`
	Title          string      `json:"title"`
	Severity       CheckStatus `json:"severity"`
	Confidence     float64     `json:"confidence"`
	ProbableCause  string      `json:"probable_cause"`
	SuggestedFix   string      `json:"suggested_fix"`
	Evidence       string      `json:"evidence"`
}

// SummaryReport is the final aggregated report from the doctor
type SummaryReport struct {
	Target       string           `json:"target"`
	Timestamp    time.Time        `json:"timestamp"`
	OverallHealth CheckStatus     `json:"overall_health"`
	DNS          *DNSResult       `json:"dns,omitempty"`
	TCP          *TCPResult       `json:"tcp,omitempty"`
	HTTP         *HTTPResult      `json:"http,omitempty"`
	TLS          *TLSResult       `json:"tls,omitempty"`
	Issues       []DiagnosisIssue `json:"issues"`
}
