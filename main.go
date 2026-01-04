package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"portscanner/formatter"
	"portscanner/scanner"
)

func main() {
	// Define flags
	var (
		format      = flag.String("format", "table", "Output format: table, detailed, or simple")
		project     = flag.String("project", "project", "Project name for analysis")
		showHelp    = flag.Bool("help", false, "Show help message")
		showVersion = flag.Bool("version", false, "Show version")
	)

	// Custom usage function
	flag.Usage = func() {
		printUsage()
	}

	// Parse flags
	flag.Parse()

	// Handle help flag
	if *showHelp {
		printUsage()
		return
	}

	// Handle version flag
	if *showVersion {
		fmt.Println("Port Scanner v1.0.0")
		return
	}

	// Get remaining arguments (ports)
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("❌ No ports provided")
		printUsage()
		return
	}

	// Parse ports from remaining arguments
	ports := parsePorts(args)
	if len(ports) == 0 {
		fmt.Println("❌ No valid ports provided")
		printUsage()
		return
	}

	// Validate format
	validFormats := map[string]bool{"table": true, "detailed": true, "simple": true}
	if !validFormats[*format] {
		fmt.Printf("❌ Invalid format: %s. Use table, detailed, or simple\n", *format)
		printUsage()
		return
	}

	scanPorts(ports, *format, *project)
}

func parsePorts(args []string) []int {
	var ports []int

	for _, arg := range args {
		// Check for port ranges (e.g., 3000-3010)
		if strings.Contains(arg, "-") {
			rangePorts := parsePortRange(arg)
			ports = append(ports, rangePorts...)
			continue
		}

		// Single port
		port, err := strconv.Atoi(arg)
		if err != nil || port < 1 || port > 65535 {
			fmt.Printf("⚠️  Skipping invalid port: %s\n", arg)
			continue
		}
		ports = append(ports, port)
	}
	return ports
}

func parsePortRange(rangeStr string) []int {
	var ports []int
	parts := strings.Split(rangeStr, "-")
	if len(parts) != 2 {
		fmt.Printf("⚠️  Skipping invalid port range: %s\n", rangeStr)
		return ports
	}

	start, err1 := strconv.Atoi(parts[0])
	end, err2 := strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil || start < 1 || end > 65535 || start > end {
		fmt.Printf("⚠️  Skipping invalid port range: %s\n", rangeStr)
		return ports
	}

	for port := start; port <= end; port++ {
		ports = append(ports, port)
	}

	fmt.Printf("🔍 Added port range: %d-%d (%d ports)\n", start, end, end-start+1)
	return ports
}

// The rest of your existing functions remain the same...
func scanPorts(ports []int, format string, projectName string) {
	ps := scanner.NewScanner()

	var statuses []*scanner.PortStatus
	for _, port := range ports {
		status, err := ps.CheckPort(port)
		if err != nil {
			status = &scanner.PortStatus{
				Port:  port,
				Error: err.Error(),
			}
		}
		statuses = append(statuses, status)
	}

	switch format {
	case "simple":
		printSimpleOutput(statuses)
	case "detailed":
		printDetailedOutput(statuses, projectName)
	case "table":
		fallthrough
	default:
		printTableOutput(statuses, projectName)
	}
}

func printDetailedOutput(statuses []*scanner.PortStatus, projectName string) {
	formatter := formatter.NewDetailedFormatter()
	output := formatter.DetailedTable(statuses, projectName)
	fmt.Println(output)
}

func printTableOutput(statuses []*scanner.PortStatus, projectName string) {
	formatter := formatter.NewTableFormatter()
	output := formatter.BriefTable(statuses, projectName)
	fmt.Println(output)
}

func printSimpleOutput(statuses []*scanner.PortStatus) {
	fmt.Printf("🔍 Scanning %d port(s)...\n\n", len(statuses))

	for _, status := range statuses {
		if status.Error != "" {
			fmt.Printf("🚨 Port %d: Error - %s\n", status.Port, status.Error)
		} else if status.IsAvailable {
			fmt.Printf("✅ Port %d: Available\n", status.Port)
		} else {
			fmt.Printf("🚨 Port %d: Occupied by %s (PID %d)\n",
				status.Port, status.ProcessName, status.PID)
		}
	}
}

func printUsage() {
	fmt.Println("Port Scanner - Check if ports are available")
	fmt.Println("")
	fmt.Println("Usage: port-scanner [OPTIONS] <port1> <port2> ...")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  port-scanner 3000 5432 8080")
	fmt.Println("  port-scanner --format detailed 3000 8501 5173")
	fmt.Println("  port-scanner --project my-app 3000 5432")
	fmt.Println("  port-scanner 3000-3010 8080-8085")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  --format string    Output format: table, detailed, or simple (default: table)")
	fmt.Println("  --project string   Project name for analysis (default: project)")
	fmt.Println("  --help             Show this help message")
	fmt.Println("  --version          Show version information")
	fmt.Println("")
	fmt.Println("Ports can be specified as:")
	fmt.Println("  • Single ports: 3000 5432 8080")
	fmt.Println("  • Port ranges: 3000-3010 8080-8085")
}
