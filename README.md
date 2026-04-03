# NetForge: Next-Generation Intelligent Networking Diagnostics

**NetForge** is a production-grade, developer-first networking CLI designed to replace fragmented tools like `ping`, `dig`, `curl`, `traceroute`, and `nmap`. It provides a unified experience for diagnosing DNS, TCP, HTTP, and TLS issues with an intelligent analysis engine.

## 🚀 Key Features

*   **`doctor`**: Flagship intelligent diagnosis command. Orchestrates multiple checks in parallel and provides actionable root-cause analysis.
*   **Intelligent Analysis Engine**: Converts raw networking data into probable causes and suggested fixes (e.g., firewall detections, certificate mismatches, high latency).
*   **Specialized Probes**: Deep-dive into specific layers with `dns`, `http`, `tls`, and `tcp` commands.
*   **Safe Port Scanning**: Diagnostic-friendly TCP connect scans with `scan`.
*   **Endpoint Monitoring**: Periodic health checks with the `monitor` command.
*   **Premium Terminal UI**: Beautiful, color-coded output with structured tables using `lipgloss` and `tablewriter`.
*   **Machine Readable**: Global support for `--json` output for automation and integration.

## 📦 Installation

### Via Curl (Recommended)

To install NetForge instantly:
```bash
curl -sSL https://raw.githubusercontent.com/onelabtech/netforge/main/scripts/install.sh | bash
```
*(Note: As this is a local build, you can use the local install script)*:
```bash
./scripts/install.sh
```

### From Source
```bash
# Requires Go 1.21+
go mod download
go build -o netforge main.go
./netforge --help
```

## 🛠 Usage Examples

### 1. The Flagship: `doctor`
Perfect for when you're thinking "why is this endpoint down?".
```bash
./netforge doctor https://api.example.com
```

### 2. DNS Investigation
Detailed record lookup and TTL analysis.
```bash
./netforge dns example.com
```

### 3. HTTP Performance
Timing breakdown including TTFB and header inspection.
```bash
./netforge http https://example.com/health
```

### 4. TLS/SSL Inspection
Handshake verification, expiry tracking, and cipher suite discovery.
```bash
./netforge tls example.com:443
```

### 5. Safe Port Scan
```bash
./netforge scan example.com
```

### 6. Endpoint Monitoring
```bash
./netforge monitor https://api.example.com --interval 10s
```

## 🧠 Diagnostic Rules

NetForge includes an analysis engine that detects:
- **Firewall blocks**: DNS succeeds but TCP fails.
- **Certificate Chain Errors**: Handshake fails but cert exists (untrusted/missing intermediate).
- **Certificate Expiry**: Warnings if cert expires in < 7 days.
- **Performance Issues**: High TTFB detection.
- **Application Failures**: 5xx status codes with correlation to backend health.
- **Protocol Mismatches**: Weak cipher or protocol version detections.

## 📊 JSON Export
All commands support `--json` for structured reporting.
```bash
./netforge doctor google.com --json > report.json
```

---
Built with ❤️ by NetForge Team.
