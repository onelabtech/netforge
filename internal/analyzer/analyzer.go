package analyzer

import (
	"fmt"
	"netforge/internal/model"
	"strings"
	"time"
)

// Analyzer processes raw results through a set of rules to find issues
type Analyzer struct {
	rules []Rule
}

// Rule represents a single diagnostic check
type Rule interface {
	Execute(report *model.SummaryReport) []model.DiagnosisIssue
}

// NewAnalyzer creates an analyzer with default rules
func NewAnalyzer() *Analyzer {
	a := &Analyzer{}
	a.registerDefaultRules()
	return a
}

// Analyze runs all rules against the provided report
func (a *Analyzer) Analyze(report *model.SummaryReport) []model.DiagnosisIssue {
	var issues []model.DiagnosisIssue
	for _, rule := range a.rules {
		issues = append(issues, rule.Execute(report)...)
	}
	return issues
}

func (a *Analyzer) registerDefaultRules() {
	a.rules = append(a.rules,
		&DNSRule{},
		&TCPRule{},
		&TLSRule{},
		&HTTPRule{},
		&ConnectivityRule{},
	)
}

// --- Specific Rule Implementations ---

// DNSRule checks for DNS resolution issues
type DNSRule struct{}

func (r *DNSRule) Execute(report *model.SummaryReport) []model.DiagnosisIssue {
	var issues []model.DiagnosisIssue
	if report.DNS == nil {
		return issues
	}

	if report.DNS.Status == model.StatusFailure {
		issues = append(issues, model.DiagnosisIssue{
			ID:            "DNS_FAILURE",
			Title:         "DNS Resolution Failed",
			Severity:      model.StatusFailure,
			Confidence:    1.0,
			ProbableCause: "The domain name could not be resolved to an IP address.",
			SuggestedFix:  "Check if the domain is correctly registered and has valid A/AAAA records.",
			Evidence:      fmt.Sprintf("Lookup for %s failed: %s", report.Target, report.DNS.Error),
		})
	}

	if len(report.DNS.IPs) > 0 && len(report.DNS.IPv6s) == 0 {
		issues = append(issues, model.DiagnosisIssue{
			ID:            "DNS_IPV6_MISSING",
			Title:         "IPv6 Records Missing (AAAA)",
			Severity:      model.StatusWarning,
			Confidence:    0.9,
			ProbableCause: "The target does not have an IPv6 address configured.",
			SuggestedFix:  "Consider adding AAAA records to support IPv6 clients.",
			Evidence:      "Only IPv4 addresses were returned.",
		})
	}

	return issues
}

// TCPRule checks for target connectivity issues
type TCPRule struct{}

func (r *TCPRule) Execute(report *model.SummaryReport) []model.DiagnosisIssue {
	var issues []model.DiagnosisIssue
	if report.TCP == nil {
		return issues
	}

	if report.TCP.Status == model.StatusFailure {
		cause := "The service is not listening on the port OR a firewall is blocking connections."
		fix := "Check if the service is running and firewall rules allow traffic on this port."
		if strings.Contains(report.TCP.Error, "connection refused") {
			cause = "Service is explicitly not listening on this port."
			fix = "Start the target service or check the configuration."
		} else if strings.Contains(report.TCP.Error, "timeout") {
			cause = "Connection timed out. Likely a firewall drop or severe network congestion."
			fix = "Check NACLs, Security Groups, or local firewalls (iptables/ufw)."
		}

		issues = append(issues, model.DiagnosisIssue{
			ID:            "TCP_CONNECT_FAILURE",
			Title:         "TCP Connection Failed",
			Severity:      model.StatusFailure,
			Confidence:    0.9,
			ProbableCause: cause,
			SuggestedFix:  fix,
			Evidence:      fmt.Sprintf("Failed to connect to %s:%d: %s", report.TCP.Host, report.TCP.Port, report.TCP.Error),
		})
	}

	return issues
}

// ConnectivityRule correlates DNS and TCP/HTTP results
type ConnectivityRule struct{}

func (r *ConnectivityRule) Execute(report *model.SummaryReport) []model.DiagnosisIssue {
	var issues []model.DiagnosisIssue
	
	// DNS works, but TCP fails
	if report.DNS != nil && report.DNS.Status == model.StatusSuccess && 
	   report.TCP != nil && report.TCP.Status == model.StatusFailure {
		issues = append(issues, model.DiagnosisIssue{
			ID:            "DNS_OK_TCP_FAIL",
			Title:         "DNS OK but Port Connectivity Fails",
			Severity:      model.StatusFailure,
			Confidence:    0.95,
			ProbableCause: "Name resolution is working, but the specific port is unreachable.",
			SuggestedFix:  "Verify firewall rules (Security Groups, local firewall) for the target port.",
			Evidence:      "DNS resolution succeeded, but TCP handshake timed out or was refused.",
		})
	}

	return issues
}

// TLSRule checks for certificate and protocol issues
type TLSRule struct{}

func (r *TLSRule) Execute(report *model.SummaryReport) []model.DiagnosisIssue {
	var issues []model.DiagnosisIssue
	if report.TLS == nil {
		return issues
	}

	if report.TLS.Status == model.StatusFailure {
		title := "TLS Handshake Failed"
		cause := "The server TLS configuration is invalid or incompatible."
		if !report.TLS.ChainValidated {
			title = "SSL/TLS Certificate Untrusted"
			cause = "The certificate chain could not be verified (Self-signed or missing intermediate)."
		}

		issues = append(issues, model.DiagnosisIssue{
			ID:            "TLS_FAILURE",
			Title:         title,
			Severity:      model.StatusFailure,
			Confidence:    1.0,
			ProbableCause: cause,
			SuggestedFix:  "Ensure the server has a valid certificate from a trusted CA and the full chain is provided.",
			Evidence:      report.TLS.Error,
		})
	}

	if report.TLS.Status == model.StatusSuccess && report.TLS.DaysUntilExpiry < 7 {
		severity := model.StatusWarning
		if report.TLS.DaysUntilExpiry <= 0 {
			severity = model.StatusFailure
		}
		issues = append(issues, model.DiagnosisIssue{
			ID:            "CERT_EXPIRING_SOON",
			Title:         "Certificate Expiring Soon",
			Severity:      severity,
			Confidence:    1.0,
			ProbableCause: fmt.Sprintf("The certificate is scheduled to expire in %d days.", report.TLS.DaysUntilExpiry),
			SuggestedFix:  "Renew the SSL/TLS certificate immediately.",
			Evidence:      fmt.Sprintf("Expiry date: %s", report.TLS.Expiry.Format("2006-01-02")),
		})
	}

	return issues
}

// HTTPRule checks for application-level issues
type HTTPRule struct{}

func (r *HTTPRule) Execute(report *model.SummaryReport) []model.DiagnosisIssue {
	var issues []model.DiagnosisIssue
	if report.HTTP == nil {
		return issues
	}

	if report.HTTP.StatusCode >= 500 {
		issues = append(issues, model.DiagnosisIssue{
			ID:            "HTTP_5XX_SERVER_ERROR",
			Title:         "Server-Side Application Error",
			Severity:      model.StatusFailure,
			Confidence:    0.9,
			ProbableCause: "The origin server or reverse proxy encountered an internal error.",
			SuggestedFix:  "Check application logs on the server and ensure backends are healthy.",
			Evidence:      fmt.Sprintf("Returned HTTP Status: %d", report.HTTP.StatusCode),
		})
	}

	if report.HTTP.TTFB > 1*time.Second {
		issues = append(issues, model.DiagnosisIssue{
			ID:            "HIGH_LATENCY_TTFB",
			Title:         "High TTFB (Time to First Byte)",
			Severity:      model.StatusWarning,
			Confidence:    0.8,
			ProbableCause: "The server or proxy is taking too long to process the request.",
			SuggestedFix:  "Optimize backend queries, check resource utilization, or review proxy configurations.",
			Evidence:      fmt.Sprintf("TTFB recorded: %s", report.HTTP.TTFB),
		})
	}

	return issues
}
