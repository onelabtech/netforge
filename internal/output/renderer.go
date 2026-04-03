package output

import (
	"fmt"
	"netforge/internal/model"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/olekukonko/tablewriter"
)

var (
	// Colors
	successColor = lipgloss.Color("#4BB543")
	warningColor = lipgloss.Color("#FFCC00")
	failureColor = lipgloss.Color("#FF3333")
	dimColor     = lipgloss.Color("#777777")
	titleColor   = lipgloss.Color("#00AAFF")
	highlightColor = lipgloss.Color("#FFFFFF")

	// Styles
	successStyle = lipgloss.NewStyle().Foreground(successColor).Bold(true)
	warningStyle = lipgloss.NewStyle().Foreground(warningColor).Bold(true)
	failureStyle = lipgloss.NewStyle().Foreground(failureColor).Bold(true)
	dimStyle     = lipgloss.NewStyle().Foreground(dimColor)
	titleStyle   = lipgloss.NewStyle().Foreground(titleColor).Bold(true).Underline(true)
	panelStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).Margin(1, 0)
)

// Renderer handles terminal output formatting
type Renderer struct{}

// NewRenderer creates a new renderer
func NewRenderer() *Renderer {
	return &Renderer{}
}

// RenderDoctorReport shows the results of a doctor run
func (r *Renderer) RenderDoctorReport(report *model.SummaryReport) {
	fmt.Println()
	fmt.Println(titleStyle.Render(fmt.Sprintf(" NetForge Doctor Report: %s ", report.Target)))
	fmt.Println(dimStyle.Render(fmt.Sprintf(" Timestamp: %s", report.Timestamp.Format("2006-01-02 15:04:05"))))
	fmt.Println()

	// Overall Health
	healthText := " HEALTHY "
	healthStyle := successStyle
	if report.OverallHealth == model.StatusFailure {
		healthText = " CRITICAL "
		healthStyle = failureStyle
	} else if report.OverallHealth == model.StatusWarning {
		healthText = " WARNING "
		healthStyle = warningStyle
	}
	fmt.Printf("Overall Status: [%s]\n", healthStyle.Render(healthText))
	fmt.Println()

	// Check Summary Table
	table := tablewriter.NewWriter(os.Stdout)
	table.Header("Check", "Status", "Duration", "Details")

	if report.DNS != nil {
		table.Append("DNS Resolution", r.statusText(report.DNS.Status), report.DNS.Duration.String(), strings.Join(report.DNS.IPs, ", "))
	}
	if report.TCP != nil {
		table.Append("TCP Connectivity", r.statusText(report.TCP.Status), report.TCP.Duration.String(), fmt.Sprintf("Port %d", report.TCP.Port))
	}
	if report.TLS != nil {
		table.Append("TLS Handshake", r.statusText(report.TLS.Status), report.TLS.Duration.String(), report.TLS.CipherSuite)
	}
	if report.HTTP != nil {
		table.Append("HTTP Inspection", r.statusText(report.HTTP.Status), report.HTTP.Duration.String(), fmt.Sprintf("Status %d", report.HTTP.StatusCode))
	}

	table.Render()
	fmt.Println()

	// Issues List
	if len(report.Issues) > 0 {
		fmt.Println(titleStyle.Render(" FINDINGS & ACTIONABLE DIAGNOSIS "))
		fmt.Println()
		for i, issue := range report.Issues {
			r.renderIssue(i+1, issue)
		}
	} else {
		fmt.Println(successStyle.Render("✓ No issues detected. Infrastructure looks solid!"))
	}
	fmt.Println()
}

func (r *Renderer) renderIssue(idx int, issue model.DiagnosisIssue) {
	style := warningStyle
	if issue.Severity == model.StatusFailure {
		style = failureStyle
	}

	header := style.Render(fmt.Sprintf("%d. %s [%s]", idx, issue.Title, issue.ID))
	fmt.Println(header)
	fmt.Printf("   Probable Cause: %s\n", issue.ProbableCause)
	fmt.Printf("   Suggested Fix:  %s\n", successStyle.Render(issue.SuggestedFix))
	fmt.Printf("   Evidence:       %s\n", dimStyle.Render(issue.Evidence))
	fmt.Println()
}

func (r *Renderer) statusText(status model.CheckStatus) string {
	switch status {
	case model.StatusSuccess:
		return successStyle.Render("OK")
	case model.StatusFailure:
		return failureStyle.Render("FAIL")
	case model.StatusWarning:
		return warningStyle.Render("WARN")
	default:
		return "UNKNOWN"
	}
}

// RenderJSON outputs the result as JSON
func (r *Renderer) RenderJSON(v interface{}) {
	// Standard JSON output
}
