package scanner

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type MacProcessAnalyzer struct{}

func NewMacProcessAnalyzer() *MacProcessAnalyzer {
	return &MacProcessAnalyzer{}
}

func (mpa *MacProcessAnalyzer) AnalyzeProcess(pid int) (*ProcessAnalysis, error) {
	analysis := &ProcessAnalysis{
		PID: pid,
	}

	// Get basic process info
	name, cmd, wd, user, err := mpa.getProcessDetails(pid)
	if err != nil {
		return nil, err
	}

	analysis.Name = name
	analysis.CommandLine = cmd
	analysis.WorkingDir = wd
	analysis.User = user

	// Detect technology
	analysis.Technology = mpa.detectTechnology(analysis)

	// Detect service type
	analysis.ServiceType = mpa.detectServiceType(analysis)

	// Find project root
	projectPath, configFiles := mpa.FindProjectRoot(wd)
	analysis.ProjectPath = projectPath
	analysis.ConfigFiles = configFiles

	// Extract ports
	ports, err := mpa.ExtractPortsFromProcess(pid)
	if err == nil {
		analysis.DetectedPorts = ports
	}

	return analysis, nil
}

func (mpa *MacProcessAnalyzer) getProcessDetails(pid int) (string, string, string, string, error) {
	// Get process name
	nameCmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "comm=")
	nameOutput, err := nameCmd.Output()
	if err != nil {
		return "", "", "", "", err
	}
	processName := strings.TrimSpace(string(nameOutput))

	// Get full command line
	cmdCmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "command=")
	cmdOutput, err := cmdCmd.Output()
	if err != nil {
		return "", "", "", "", err
	}
	commandLine := strings.TrimSpace(string(cmdOutput))

	// Get working directory (for macOS, we need to use lsof or pwdx equivalent)
	wdCmd := exec.Command("lsof", "-p", strconv.Itoa(pid), "-a", "-d", "cwd", "-Fn")
	wdOutput, err := wdCmd.Output()
	workingDir := ""
	if err == nil {
		lines := strings.Split(string(wdOutput), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "n") && len(line) > 1 {
				workingDir = strings.TrimPrefix(line, "n")
				break
			}
		}
	}

	// Get user
	userCmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "user=")
	userOutput, err := userCmd.Output()
	if err != nil {
		return "", "", "", "", err
	}
	user := strings.TrimSpace(string(userOutput))

	return processName, commandLine, workingDir, user, nil
}

func (mpa *MacProcessAnalyzer) detectTechnology(analysis *ProcessAnalysis) string {
	cmd := strings.ToLower(analysis.CommandLine)
	name := strings.ToLower(analysis.Name)

	// Node.js detection
	if strings.Contains(cmd, "node") || strings.Contains(cmd, "npm") ||
		strings.Contains(cmd, "npx") || name == "node" {
		return "node"
	}

	// Python detection
	if strings.Contains(cmd, "python") || strings.Contains(cmd, "python3") ||
		strings.Contains(cmd, "streamlit") || strings.Contains(cmd, "fastapi") ||
		strings.Contains(cmd, "flask") || strings.Contains(cmd, "django") {
		return "python"
	}

	// Postgres detection
	if strings.Contains(name, "postgres") || strings.Contains(cmd, "postgres") {
		return "postgres"
	}

	// Redis detection
	if strings.Contains(name, "redis") || strings.Contains(cmd, "redis") {
		return "redis"
	}

	// Java detection
	if strings.Contains(cmd, "java") || strings.Contains(name, "java") {
		return "java"
	}

	// Go detection
	if strings.Contains(cmd, "go") && !strings.Contains(cmd, "google") {
		return "go"
	}

	// Browser detection
	if strings.Contains(name, "firefox") || strings.Contains(name, "chrome") ||
		strings.Contains(name, "safari") {
		return "browser"
	}

	return "unknown"
}

func (mpa *MacProcessAnalyzer) detectServiceType(analysis *ProcessAnalysis) string {
	switch analysis.Technology {
	case "node", "python", "go", "java":
		// Check if it's a web server
		if strings.Contains(strings.ToLower(analysis.CommandLine), "server") ||
			strings.Contains(strings.ToLower(analysis.CommandLine), "start") ||
			strings.Contains(strings.ToLower(analysis.CommandLine), "run") ||
			strings.Contains(strings.ToLower(analysis.CommandLine), "dev") {
			return "web"
		}
		return "cli"
	case "postgres", "mysql", "mongod":
		return "database"
	case "redis", "memcached":
		return "cache"
	case "browser":
		return "browser"
	default:
		return "system"
	}
}

func (mpa *MacProcessAnalyzer) FindProjectRoot(workingDir string) (string, []string) {
	if workingDir == "" {
		return "", []string{}
	}

	dir := workingDir
	var configFiles []string
	projectMarkers := []string{
		"package.json",     // Node.js
		"go.mod",           // Go
		"requirements.txt", // Python
		"pom.xml",          // Java
		"docker-compose.yml",
		"docker-compose.yaml",
		".git",
		"Cargo.toml",     // Rust
		"Gemfile",        // Ruby
		"pyproject.toml", // Python (modern)
		"composer.json",  // PHP
	}

	for dir != "/" {
		// Check for project markers
		for _, marker := range projectMarkers {
			markerPath := filepath.Join(dir, marker)
			if _, err := os.Stat(markerPath); err == nil {
				configFiles = append(configFiles, markerPath)
			}
		}

		// If we found any config files, this is likely the project root
		if len(configFiles) > 0 {
			return dir, configFiles
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir { // Reached root
			break
		}
		dir = parent
	}

	return workingDir, configFiles
}

func (mpa *MacProcessAnalyzer) ExtractPortsFromProcess(pid int) ([]int, error) {
	var ports []int

	// Use lsof to find network connections for this process
	cmd := exec.Command("lsof", "-p", strconv.Itoa(pid), "-i", "-P", "-n")
	output, err := cmd.Output()
	if err != nil {
		return ports, err
	}

	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		if i == 0 { // Skip header
			continue
		}
		if line == "" {
			continue
		}

		// Parse lsof output to extract port numbers
		// Format: process PID user FD type ... NODE NAME
		fields := strings.Fields(line)
		if len(fields) >= 9 {
			// Look for port in the NODE NAME field (e.g., *:3000 (LISTEN))
			nodeField := fields[8]
			if idx := strings.LastIndex(nodeField, ":"); idx != -1 {
				portStr := nodeField[idx+1:]
				// Remove trailing stuff like (LISTEN)
				if parenIdx := strings.Index(portStr, " "); parenIdx != -1 {
					portStr = portStr[:parenIdx]
				}
				if portStr == "*" {
					continue
				}
				if port, err := strconv.Atoi(portStr); err == nil {
					// Check if port already in list
					found := false
					for _, p := range ports {
						if p == port {
							found = true
							break
						}
					}
					if !found {
						ports = append(ports, port)
					}
				}
			}
		}
	}

	return ports, nil
}
