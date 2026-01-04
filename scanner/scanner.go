package scanner

type PortStatus struct {
	Port        int
	IsAvailable bool
	ProcessName string
	PID         int
	Error       string
	User        string // New: Process owner
	CommandLine string // New: Full command
	StartTime   string // New: Process start time
	MemoryUsage string // New: Memory consumption
}

type PortScanner interface {
	CheckPort(port int) (*PortStatus, error)
}

func NewScanner() PortScanner {
	return &MacScanner{}
}
