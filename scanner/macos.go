package scanner

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

type MacScanner struct{}

func (ms *MacScanner) CheckPort(port int) (*PortStatus, error) {
	status := &PortStatus{Port: port}

	if ms.isPortAvailable(port) {
		status.IsAvailable = true
		return status, nil
	}

	process, pid, err := ms.findWithLsof(port)
	if err != nil {
		status.Error = err.Error()
		return status, nil
	}

	status.IsAvailable = false
	status.ProcessName = process
	status.PID = pid

	//detailed
	status.User = ms.getProcessUser(pid)
	status.CommandLine = ms.getCommandLine(pid)
	status.MemoryUsage = ms.getMemoryUsage(pid)
	status.StartTime = ms.getStartTime(pid)

	return status, nil
}

func (ms *MacScanner) getProcessUser(pid int) string {
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "user=")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

func (ms *MacScanner) getCommandLine(pid int) string {
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "command=")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

func (ms *MacScanner) getMemoryUsage(pid int) string {
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "rss=")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	rss := strings.TrimSpace(string(output))
	if rssKB, err := strconv.Atoi(rss); err == nil {
		return fmt.Sprintf("%dMB", rssKB/1024)
	}
	return rss + "KB"
}

func (ms *MacScanner) getStartTime(pid int) string {
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "lstart=")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

func (ms *MacScanner) isPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

func (ms *MacScanner) findWithLsof(port int) (string, int, error) {
	// Use full lsof output to get both PID and command
	cmd := exec.Command("lsof", "-i", ":"+strconv.Itoa(port), "-P")
	output, err := cmd.Output()
	if err != nil {
		return "", 0, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) <= 1 { // First line is header
		return "", 0, fmt.Errorf("no lsof output")
	}

	// Parse the first data line (skip header)
	for i, line := range lines {
		if i == 0 { // Skip header
			continue
		}
		if line == "" {
			continue
		}

		// Parse: "Python  28024 user  8u  IPv4 ..."
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			pid, err := strconv.Atoi(fields[1])
			if err != nil {
				continue // Skip if PID not parseable
			}

			processName := fields[0]
			return processName, pid, nil
		}
	}

	return "", 0, fmt.Errorf("no process found in lsof output")
}

func (ms *MacScanner) getProcessName(pid int) (string, error) {
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "comm=")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
