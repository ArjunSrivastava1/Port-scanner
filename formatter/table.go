package formatter

import (
	"fmt"
	"portscanner/scanner"
	"strconv"
	"strings"
)

type TableFormatter struct{}

func NewTableFormatter() *TableFormatter {
	return &TableFormatter{}
}

func (tf *TableFormatter) BriefTable(statuses []*scanner.PortStatus, projectName string) string {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("PORT CONFLICT ANALYSIS: %s\n", projectName))
	sb.WriteString("──────────────────────────────\n\n")

	// Table Header
	sb.WriteString(tf.formatRow("SERVICE", "PORT", "STATUS", "PROCESS", "IMPACT", "UPTIME", "RESOURCES"))
	sb.WriteString(tf.formatRow("───────", "────", "──────", "───────", "──────", "──────", "─────────"))

	// Table Rows
	for _, status := range statuses {
		service := tf.guessService(status.Port)
		statusText := tf.formatStatus(status)
		process := tf.formatProcess(status)
		impact := tf.assessImpact(status)
		uptime := "-"
		resources := tf.assessResources(status)

		sb.WriteString(tf.formatRow(service, strconv.Itoa(status.Port), statusText, process, impact, uptime, resources))
	}

	// Resolution section - ALWAYS show if we have any non-available ports
	conflicts := tf.countConflicts(statuses)
	if conflicts > 0 {
		sb.WriteString("\n")
		sb.WriteString(tf.generateResolutions(statuses))
	}

	//new adds
	if conflicts == 0 {
		sb.WriteString("• All ports are available and ready for use! ✅\n")
		sb.WriteString("• No conflicts detected - development environment is clear 🎉\n")
	}
	return sb.String()
}

func (tf *TableFormatter) formatRow(service, port, status, process, impact, uptime, resources string) string {
	return fmt.Sprintf("%-12s %-6s %-10s %-16s %-8s %-8s %s\n",
		service, port, status, process, impact, uptime, resources)
}

func (tf *TableFormatter) formatStatus(status *scanner.PortStatus) string {
	if status.IsAvailable {
		return "✅ READY"
	}
	// If port is not available, it's a conflict (even if we have error details)
	return "🔴 CONFLICT"
}

func (tf *TableFormatter) formatProcess(status *scanner.PortStatus) string {
	if status.IsAvailable {
		return "-"
	}
	if status.ProcessName != "" && status.PID != 0 {
		return fmt.Sprintf("%s:%d", status.ProcessName, status.PID)
	}
	if status.ProcessName != "" {
		return status.ProcessName
	}
	if status.Error != "" {
		return "unknown"
	}
	return "-"
}

func (tf *TableFormatter) guessService(port int) string {
	serviceMap := map[int]string{
		3000:  "frontend",
		5000:  "backend",
		5173:  "frontend", // Vite
		8000:  "backend",
		8080:  "backend",
		8501:  "streamlit", // Streamlit apps
		5432:  "database",
		6379:  "cache",
		9200:  "search",
		27017: "mongodb",
		3306:  "mysql",
		9000:  "backend",
		4200:  "frontend",
	}

	if service, exists := serviceMap[port]; exists {
		return service
	}
	return "service"
}

func (tf *TableFormatter) assessImpact(status *scanner.PortStatus) string {
	if status.IsAvailable {
		return "-"
	}

	switch status.Port {
	case 5432, 3306, 27017:
		return "HIGH"
	case 6379, 9200:
		return "MEDIUM"
	case 8501:
		return "MEDIUM"
	default:
		return "LOW"
	}
}

func (tf *TableFormatter) assessResources(status *scanner.PortStatus) string {
	if status.IsAvailable {
		return "Available"
	}

	// Port-specific detection
	switch status.Port {
	case 8501:
		return "Streamlit App"
	case 5173:
		return "Vite Dev Server"
	case 8000:
		return "Python/Backend"
	}

	// Process-based detection
	switch status.ProcessName {
	case "postgres", "mysql", "mongod":
		return "Database"
	case "redis":
		return "Cache"
	case "node", "python", "java":
		return "Application"
	default:
		return "System"
	}
}

func (tf *TableFormatter) countConflicts(statuses []*scanner.PortStatus) int {
	count := 0
	for _, status := range statuses {
		if !status.IsAvailable {
			count++
		}
	}
	return count
}

func (tf *TableFormatter) generateResolutions(statuses []*scanner.PortStatus) string {
	var sb strings.Builder
	conflictCount := tf.countConflicts(statuses)

	sb.WriteString(fmt.Sprintf("CONFLICT RESOLUTION (%d conflicts):\n", conflictCount))
	sb.WriteString("\033[32m1. PORT MAPPING\033[0m: Use alternative ports    ✅ RECOMMENDED\n")
	sb.WriteString("\033[33m2. SERVICE RESTART\033[0m: Restart on new ports   ⚠️  LOW RISK\n")
	sb.WriteString("\033[31m3. PROCESS TERMINATION\033[0m: Stop services      🔴 HIGH RISK\n")
	sb.WriteString("\n\033[36mExecute: port-scanner --fix for auto-resolution\033[0m\n")

	return sb.String()
}
