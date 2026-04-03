package checker

import (
	"context"
	"fmt"
	"net"
	"netforge/internal/model"
	"netforge/internal/util"
	"time"
)

// TCPChecker implements raw TCP connectivity checks
type TCPChecker struct{}

// NewTCPChecker creates a new TCP checker
func NewTCPChecker() *TCPChecker {
	return &TCPChecker{}
}

// Check attempts to establish a TCP connection to the target
func (c *TCPChecker) Check(ctx context.Context, target *util.TaskTarget) (*model.TCPResult, error) {
	start := time.Now()
	result := &model.TCPResult{
		Status: model.StatusUnknown,
		Host:   target.Host,
		Port:   target.Port,
	}

	dialer := net.Dialer{
		Timeout: 5 * time.Second, // Default internal timeout, can be overridden by ctx
	}

	addr := target.GetAddr()
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		result.Status = model.StatusFailure
		result.Error = err.Error()
		result.Duration = time.Since(start)
		return result, nil
	}
	defer conn.Close()

	result.Status = model.StatusSuccess
	result.Duration = time.Since(start)

	return result, nil
}

// ScanPorts performs a simple TCP connect scan on multiple ports
func (c *TCPChecker) ScanPorts(ctx context.Context, host string, ports []int) (*model.ScanResult, error) {
	result := &model.ScanResult{
		Target: host,
		Ports:  []model.PortResult{},
	}

	for _, port := range ports {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
			addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
			conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
			
			portRes := model.PortResult{
				Port: port,
			}

			if err != nil {
				portRes.Status = "closed"
			} else {
				portRes.Status = "open"
				conn.Close()
			}
			result.Ports = append(result.Ports, portRes)
		}
	}

	return result, nil
}
