package formatter

import (
	"fmt"
	"portscanner/scanner"
	"strconv"
	"strings"
)

type DetailedFormatter struct{}

func NewDetailedFormatter() *DetailedFormatter {
	return &DetailedFormatter{}
}

func (df *DetailedFormatter) DetailedTable(statuses []*scanner.PortStatus, projectName string) string {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("DETAILED PORT ANALYSIS: %s\n", projectName))
	sb.WriteString("──────────────────────────────\n\n")

	// Detailed Table Header
	sb.WriteString(df.formatRow("SERVICE", "PORT", "STATUS", "PROCESS", "PID", "USER", "MEMORY", "UPTIME"))
	sb.WriteString(df.formatRow("───────", "────", "──────", "───────", "───", "────", "──────", "──────"))

	// Table Rows
	for _, status := range statuses {
		service := df.guessService(status.Port)
		statusText := df.formatStatus(status)
		process := df.formatProcess(status)
		pid := df.formatPID(status)
		user := df.formatUser(status)
		memory := df.formatMemory(status)
		uptime := df.formatUptime(status)

		sb.WriteString(df.formatRow(service, strconv.Itoa(status.Port), statusText, process, pid, user, memory, uptime))
	}

	// Impact Analysis Section
	sb.WriteString("\n")
	sb.WriteString(df.generateImpactAnalysis(statuses))

	// Resolution Section
	conflicts := df.countConflicts(statuses)
	if conflicts > 0 {
		sb.WriteString("\n")
		sb.WriteString(df.generateDetailedResolutions(statuses))
	}

	// Command Line Details for Conflicts
	sb.WriteString(df.generateProcessDetails(statuses))

	return sb.String()
}

func (df *DetailedFormatter) formatRow(service, port, status, process, pid, user, memory, uptime string) string {
	return fmt.Sprintf("%-12s %-6s %-10s %-16s %-6s %-12s %-8s %s\n",
		service, port, status, process, pid, user, memory, uptime)
}

func (df *DetailedFormatter) formatStatus(status *scanner.PortStatus) string {
	if status.IsAvailable {
		return "✅ READY"
	}
	return "🔴 CONFLICT"
}

func (df *DetailedFormatter) formatProcess(status *scanner.PortStatus) string {
	if status.IsAvailable {
		return "-"
	}
	if status.ProcessName != "" {
		return status.ProcessName
	}
	return "unknown"
}

func (df *DetailedFormatter) formatPID(status *scanner.PortStatus) string {
	if status.IsAvailable || status.PID == 0 {
		return "-"
	}
	return strconv.Itoa(status.PID)
}

func (df *DetailedFormatter) formatUser(status *scanner.PortStatus) string {
	if status.IsAvailable {
		return "-"
	}
	if status.User != "" {
		return status.User
	}
	return "unknown"
}

func (df *DetailedFormatter) formatMemory(status *scanner.PortStatus) string {
	if status.IsAvailable {
		return "-"
	}
	if status.MemoryUsage != "" {
		return status.MemoryUsage
	}
	return "unknown"
}

func (df *DetailedFormatter) formatUptime(status *scanner.PortStatus) string {
	if status.IsAvailable {
		return "-"
	}
	if status.StartTime != "" {
		// For now, just show the start time - we could calculate duration later
		return status.StartTime
	}
	return "unknown"
}

func (df *DetailedFormatter) guessService(port int) string {
	// Same service mapping as brief formatter
	serviceMap := map[int]string{
		3000:  "frontend",
		5000:  "backend",
		5173:  "frontend",
		8000:  "backend",
		8080:  "backend",
		8501:  "streamlit",
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

func (df *DetailedFormatter) countConflicts(statuses []*scanner.PortStatus) int {
	count := 0
	for _, status := range statuses {
		if !status.IsAvailable {
			count++
		}
	}
	return count
}

func (df *DetailedFormatter) generateImpactAnalysis(statuses []*scanner.PortStatus) string {
	var sb strings.Builder
	conflicts := df.countConflicts(statuses)
	sb.WriteString("IMPACT ANALYSIS:\n")

	if conflicts == 0 {
		sb.WriteString("• All ports are available and ready for use! ✅\n")
		sb.WriteString("• No conflicts detected - development environment is clear 🎉\n")
	} else {
		for _, status := range statuses {
			if !status.IsAvailable {
				impact := df.assessImpact(status)
				sb.WriteString(fmt.Sprintf("• \033[31m%s (%d): %s\033[0m\n", df.guessService(status.Port), status.Port, impact))
				sb.WriteString(fmt.Sprintf("  - Process: %s (PID %d)\n", status.ProcessName, status.PID))
				sb.WriteString(fmt.Sprintf("  - User: %s, Memory: %s\n", status.User, status.MemoryUsage))
				sb.WriteString(fmt.Sprintf("  - Started: %s\n", status.StartTime))

				risk := df.assessRisk(status)
				sb.WriteString(fmt.Sprintf("  - \033[33mRisk: %s\033[0m\n", risk))
				sb.WriteString("\n")
			}
		}
	}

	return sb.String()
}

func (df *DetailedFormatter) assessImpact(status *scanner.PortStatus) string {
	switch status.Port {
	case 5432, 3306, 27017:
		return "HIGH - Database service"
	case 6379, 9200:
		return "MEDIUM - Cache/Search service"
	case 8501:
		return "MEDIUM - Streamlit application"
	default:
		return "LOW - Development service"
	}
}

func (df *DetailedFormatter) assessRisk(status *scanner.PortStatus) string {
	switch status.ProcessName {
	case "postgres", "mysql", "mongod":
		return "Data loss if terminated"
	case "redis":
		return "Session data loss"
	case "python", "node", "java":
		return "Service interruption"
	default:
		return "Minimal impact"
	}
}

func (df *DetailedFormatter) generateDetailedResolutions(statuses []*scanner.PortStatus) string {
	var sb strings.Builder
	conflictCount := df.countConflicts(statuses)

	sb.WriteString(fmt.Sprintf("DETAILED RESOLUTION PATHS (%d conflicts):\n", conflictCount))
	sb.WriteString("\n\033[32m1. PORT MAPPING (RECOMMENDED)\033[0m\n")

	// Show specific port mapping suggestions
	for _, status := range statuses {
		if !status.IsAvailable {
			alternative := df.findAlternativePort(status.Port)
			sb.WriteString(fmt.Sprintf("   %d → %d (available)\n", status.Port, alternative))
		}
	}
	sb.WriteString("   Impact: Zero downtime, update configuration files\n")

	sb.WriteString("\n\033[33m2. SERVICE RESTART (LOW RISK)\033[0m\n")
	sb.WriteString("   Restart services on alternative ports\n")
	sb.WriteString("   Impact: Brief service interruption (1-2 minutes)\n")

	sb.WriteString("\n\033[31m3. PROCESS TERMINATION (HIGH RISK)\033[0m\n")
	for _, status := range statuses {
		if !status.IsAvailable {
			sb.WriteString(fmt.Sprintf("   Stop: %s (PID %d) - %s\n", status.ProcessName, status.PID, df.assessRisk(status)))
		}
	}
	sb.WriteString("   Impact: Service disruption, potential data loss\n")

	sb.WriteString("\n\033[36mExecute: port-scanner --fix for auto-resolution\033[0m\n")

	return sb.String()
}

func (df *DetailedFormatter) findAlternativePort(original int) int {
	// Simple alternative port finder
	switch original {
	case 3000:
		return 3001
	case 5432:
		return 5433
	case 6379:
		return 6380
	case 8501:
		return 8502
	case 8080:
		return 8081
	default:
		return original + 1
	}
}

func (df *DetailedFormatter) generateProcessDetails(statuses []*scanner.PortStatus) string {
	var sb strings.Builder
	hasConflicts := false

	for _, status := range statuses {
		if !status.IsAvailable && status.CommandLine != "" {
			if !hasConflicts {
				sb.WriteString("\nPROCESS DETAILS:\n")
				hasConflicts = true
			}
			sb.WriteString(fmt.Sprintf("• %s (PID %d):\n", status.ProcessName, status.PID))
			sb.WriteString(fmt.Sprintf("  Command: %s\n", status.CommandLine))
		}
	}

	return sb.String()
}
