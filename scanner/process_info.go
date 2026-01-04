package scanner

type ProcessAnalysis struct {
	PID           int
	Name          string
	CommandLine   string
	WorkingDir    string
	User          string
	Technology    string   // "node", "python", "postgres", "redis", "unknown"
	ServiceType   string   // "web", "database", "cache", "cli", "browser"
	DetectedPorts []int    // Ports this process is using
	ProjectPath   string   // Path to project root (if detectable)
	ConfigFiles   []string // package.json, docker-compose.yml, etc.
}

type Dependency struct {
	Type     string // "network", "config", "process"
	Target   string // Port, file path, or process
	Protocol string // HTTP, PostgreSQL, Redis, etc.
}

// ProcessAnalyzer interface
type ProcessAnalyzer interface {
	AnalyzeProcess(pid int) (*ProcessAnalysis, error)
	FindProjectRoot(workingDir string) (string, []string)
	ExtractPortsFromProcess(pid int) ([]int, error)
}
