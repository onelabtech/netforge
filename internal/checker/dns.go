package checker

import (
	"context"
	"fmt"
	"net"
	"netforge/internal/model"
	"netforge/internal/util"
	"time"
)

// DNSChecker implements DNS resolution checks
type DNSChecker struct{}

// NewDNSChecker creates a new DNS checker
func NewDNSChecker() *DNSChecker {
	return &DNSChecker{}
}

// Check performs DNS resolution for the given target
func (c *DNSChecker) Check(ctx context.Context, target *util.TaskTarget) (*model.DNSResult, error) {
	start := time.Now()
	result := &model.DNSResult{
		Status: model.StatusUnknown,
		Target: target.Host,
	}

	// Use normal net.LookupIP for multi-protocol resolution (A/AAAA)
	ips, err := net.LookupIP(target.Host)
	if err != nil {
		result.Status = model.StatusFailure
		result.Error = err.Error()
		result.Duration = time.Since(start)
		return result, nil
	}

	for _, ip := range ips {
		if ip.To4() != nil {
			result.IPs = append(result.IPs, ip.String())
		} else {
			result.IPv6s = append(result.IPv6s, ip.String())
		}
	}

	// Fetch CNAME
	cname, err := net.LookupCNAME(target.Host)
	if err == nil && cname != "" && cname != target.Host && cname != target.Host+"." {
		result.CNAMEs = append(result.CNAMEs, cname)
	}

	// Fetch MX
	mxs, err := net.LookupMX(target.Host)
	if err == nil {
		for _, mx := range mxs {
			result.MX = append(result.MX, mx.Host)
		}
	}

	// Fetch NS
	nss, err := net.LookupNS(target.Host)
	if err == nil {
		for _, ns := range nss {
			result.NS = append(result.NS, ns.Host)
		}
	}

	result.Status = model.StatusSuccess
	result.Duration = time.Since(start)

	if len(result.IPs) == 0 && len(result.IPv6s) == 0 {
		result.Status = model.StatusWarning
		result.Error = "no IP addresses found for host"
	}

	return result, nil
}

// Trace performs a iterative DNS trace from root servers to authoritative servers
func (c *DNSChecker) Trace(ctx context.Context, domain string) ([]string, error) {
	// Best effort trace with miekg/dns if needed
	// This will be more useful for the 'dns' command specialized flags
	return nil, fmt.Errorf("trace not implemented")
}
