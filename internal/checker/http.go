package checker

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptrace"
	"netforge/internal/model"
	"netforge/internal/util"
	"strings"
	"time"
)

// HTTPChecker implements HTTP response inspection
type HTTPChecker struct {
	client *http.Client
}

// NewHTTPChecker creates a new HTTP checker
func NewHTTPChecker() *HTTPChecker {
	return &HTTPChecker{
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Don't follow redirects automatically if we want to trace them
				// We'll manually handle transparency in some cases
				return nil
			},
			Timeout: 10 * time.Second,
		},
	}
}

// Check performs an HTTP request and measures performance
func (c *HTTPChecker) Check(ctx context.Context, target *util.TaskTarget) (*model.HTTPResult, error) {
	start := time.Now()
	url := target.GetURL()
	result := &model.HTTPResult{
		Status:  model.StatusUnknown,
		URL:     url,
		Headers: make(map[string]string),
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		result.Error = err.Error()
		result.Status = model.StatusFailure
		return result, nil
	}

	// Use httptrace to capture timings
	var ttfb time.Duration
	trace := &httptrace.ClientTrace{
		GotFirstResponseByte: func() {
			ttfb = time.Since(start)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	// Configure transport to be more diagnostic
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false}, // Default to secure
	}
	c.client.Transport = transport

	resp, err := c.client.Do(req)
	if err != nil {
		result.Error = err.Error()
		result.Status = model.StatusFailure
		result.Duration = time.Since(start)
		return result, nil
	}
	defer resp.Body.Close()

	// Fill results
	result.Status = model.StatusSuccess
	result.StatusCode = resp.StatusCode
	result.Proto = resp.Proto
	result.Duration = time.Since(start)
	result.TTFB = ttfb
	result.ContentLength = resp.ContentLength

	for name, values := range resp.Header {
		result.Headers[name] = strings.Join(values, ", ")
	}

	// Detect compression
	if ce := resp.Header.Get("Content-Encoding"); ce != "" {
		result.Compression = ce
	}

	// Follow redirects if needed to build a chain (manually if we want detailed tracking)
	// For the initial version, we'll just show the final status or the single redirect
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		if loc := resp.Header.Get("Location"); loc != "" {
			result.Redirects = append(result.Redirects, loc)
		}
	}

	// Basic body read to ensure we get meaningful data and compression works
	// For diagnostic reasons, we don't want to download large bodies
	io.CopyN(io.Discard, resp.Body, 1024)

	return result, nil
}
