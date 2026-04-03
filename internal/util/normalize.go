package util

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// TaskTarget represents a normalized target for network checks
type TaskTarget struct {
	Original  string
	Host      string
	Port      int
	Scheme    string
	Path      string
	IsIP      bool
	IsIPv6    bool
}

// NormalizeTarget parses a string input into a structured Target
func NormalizeTarget(input string) (*TaskTarget, error) {
	if input == "" {
		return nil, fmt.Errorf("target cannot be empty")
	}

	target := &TaskTarget{
		Original: input,
		Scheme:   "http", // Default scheme
		Port:     80,     // Default port
		Path:     "/",    // Default path
	}

	// Basic check for protocol
	if !strings.Contains(input, "://") && !strings.HasPrefix(input, "localhost") {
		// If it looks like a domain or IP, but doesn't have a protocol, we'll try to guess
		// However, for pure DNS/TCP checks, we might not need a protocol.
		// For consistency, we'll prefix with http if it's not present.
		// input = "http://" + input
	}

	u, err := url.Parse(input)
	if err != nil || u.Host == "" {
		// Try parsing as a hostname:port or just hostname
		host, portStr, err := net.SplitHostPort(input)
		if err != nil {
			// Just a hostname or IP
			target.Host = input
		} else {
			target.Host = host
			fmt.Sscanf(portStr, "%d", &target.Port)
		}
	} else {
		target.Host = u.Hostname()
		target.Scheme = u.Scheme
		target.Path = u.Path
		if u.Port() != "" {
			fmt.Sscanf(u.Port(), "%d", &target.Port)
		} else {
			if u.Scheme == "https" {
				target.Port = 443
			}
		}
	}

	// Check if Host is an IP
	ip := net.ParseIP(target.Host)
	if ip != nil {
		target.IsIP = true
		target.IsIPv6 = ip.To4() == nil
	}

	return target, nil
}

// GetAddr returns the host:port string
func (t *TaskTarget) GetAddr() string {
	return net.JoinHostPort(t.Host, fmt.Sprintf("%d", t.Port))
}

// GetURL returns the full URL string
func (t *TaskTarget) GetURL() string {
	if t.Scheme == "" {
		return fmt.Sprintf("http://%s%s", t.GetAddr(), t.Path)
	}
	return fmt.Sprintf("%s://%s%s", t.Scheme, t.GetAddr(), t.Path)
}
