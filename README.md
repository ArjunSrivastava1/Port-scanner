<h1>
  <img src="https://raw.githubusercontent.com/ArjunSrivastava1/port-scanner/main/assets/icon.svg" alt="port-scanner" width="100">
</h1>

<h4>A development port conflict detector that simply works like a charm, created to analyse ports and detect conflicts</h4>

<p>
  <a href="https://golang.org"><img src="https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white" alt="Go Version"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-GPL%20v2-blue.svg" alt="License"></a>
  <a href="https://github.com/ArjunSrivastava1/port-scanner/releases"><img src="https://img.shields.io/github/v/release/ArjunSrivastava1/port-scanner" alt="Release"></a>
  <a href="https://goreportcard.com/report/github.com/ArjunSrivastava1/port-scanner"><img src="https://goreportcard.com/badge/github.com/ArjunSrivastava1/port-scanner" alt="Go Report Card"></a>
</p>

<p>
  <a href="#-about">About</a> •
  <a href="#-features">Features</a> •
  <a href="#-quick-start">Quick Start</a> •
  <a href="#-usage">Usage</a> •
  <a href="#-output-examples">Examples</a> •
  <a href="#-architecture">Architecture</a> •
  <a href="#-contributing">Contributing</a>
</p>

<p align="center">
  <img src="https://raw.githubusercontent.com/ArjunSrivastava1/port-scanner/main/assets/demo.gif" alt="Demo" width="600">
</p>

## 🎯 About

**Port Scanner** is a command line interface (CLI) tool that detects and resolves port conflicts in development environments. 

Think of it as a **traffic controller** for your development ports—it shows you exactly what's running where, why there are conflicts, and how to fix them.

No more `Error: listen EADDRINUSE: address already in use` frustration. Get actionable insights with beautiful, informative output.

## ✨ Features

| Category | Features |
|----------|----------|
| **🔍 Smart Detection** | Port availability • Process identification • Service type inference • Impact assessment |
| **🎨 Beautiful Output** | Multiple formats (table/detailed/simple) • Color-coded status • Actionable resolutions • Risk indicators |
| **⚡ Developer UX** | Zero config for basic use • Project-aware scanning • Smart service guessing • Cross-platform ready |
| **🔧 Professional** | Modular architecture • Clean Go code • Comprehensive tests • Full documentation |

## 🚀 Quick Start

### 📦 Installation
```bash
# Install from source
go install github.com/ArjunSrivastava1/port-scanner@latest

# Or clone and build
git clone https://github.com/ArjunSrivastava1/port-scanner
cd port-scanner
go build -o port-scanner main.go
sudo mv port-scanner /usr/local/bin/
```

## 🎮 Usage

### Basic Port Scanning
```bash
# Check specific ports
port-scanner 3000 5432 8080

# Detailed analysis view
port-scanner --format detailed 3000 8501 5173

# Simple output for scripts
port-scanner --format simple 3000 5432
```

### Project-Aware Scanning
```bash
# Scan with project context
port-scanner --project my-app 3000 5432 8080

# Auto-detect development ports (experimental)
port-scanner --auto-detect
```

### Output Formats
```bash
# Brief table view (default)
port-scanner 3000 5432 8080

# Detailed analysis with impact assessment
port-scanner --format detailed 3000 5432

# Simple output for CI/CD pipelines
port-scanner --format simple 3000 5432 | grep -q "CONFLICT" && exit 1
```

## 📊 Output Examples

### Brief Table View
```bash
$ port-scanner 3000 5432 8080

PORT CONFLICT ANALYSIS: project
──────────────────────────────

SERVICE    PORT  STATUS    PROCESS       IMPACT    UPTIME    RESOURCES
frontend   3000  ✅ READY  -             -         -         Available
database   5432  🔴 CONFLICT postgres:8910 HIGH      2h        Production DB
backend    8080  ✅ READY  -             -         -         Available

CONFLICT RESOLUTION (1 conflicts):
1. PORT MAPPING: Use alternative ports    ✅ RECOMMENDED
2. SERVICE RESTART: Restart on new ports   ⚠️  LOW RISK
3. PROCESS TERMINATION: Stop services      🔴 HIGH RISK
```

### Detailed Analysis View
```bash
$ port-scanner --format detailed 3000 5432

DETAILED PORT ANALYSIS: project
──────────────────────────────

SERVICE    PORT  STATUS    PROCESS       PID    USER    MEMORY   UPTIME
database   5432  🔴 CONFLICT postgres    8910   postgres 256MB    2h

IMPACT ANALYSIS:
• database (5432): HIGH - Production database
  - Process: postgres (PID 8910)
  - User: postgres, Memory: 256MB
  - Started: 2 hours ago
  - Risk: Data loss if terminated

DETAILED RESOLUTION PATHS:
1. PORT MAPPING (RECOMMENDED)
   5432 → 5433 (available)
   Impact: Zero downtime, update configuration files

2. SERVICE RESTART (LOW RISK)
   Restart services on alternative ports
   Impact: Brief service interruption

3. PROCESS TERMINATION (HIGH RISK)
   Stop: postgres (PID 8910) - Data loss risk
   Impact: Service disruption, potential data loss
```

## 🏗️ Architecture

Port Scanner follows a clean, modular architecture:

```
port-scanner/
├── scanner/           # Core scanning engine (macOS/Linux/Windows)
├── formatter/         # Output formatting (table/detailed/simple)
├── analyzer/          # Process analysis & service detection
└── main.go           # CLI interface & flag parsing
```

### Key Design Principles:
1. **Platform Abstraction**: Clean interfaces for cross-platform support
2. **Modular Output**: Separate formatting from business logic
3. **Progressive Enhancement**: Simple → Detailed → Advanced features
4. **Developer Experience**: Zero config for 80% use cases

## 🔧 Development

### Building from Source
```bash
# Clone repository
git clone https://github.com/ArjunSrivastava1/port-scanner
cd port-scanner

# Build for current platform
go build -o port-scanner main.go

# Cross-compilation
GOOS=linux GOARCH=amd64 go build -o port-scanner-linux main.go
GOOS=windows GOARCH=amd64 go build -o port-scanner.exe main.go
```

### Running Tests
```bash
# Run unit tests
go test ./...

# With coverage
go test -cover ./...
```

## 🤝 Contributing

We welcome contributions! Here's how to get started:

1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/your-username/port-scanner`
3. **Create a feature branch**: `git checkout -b feature/amazing-feature`
4. **Commit changes** using conventional commits
5. **Push to your branch**: `git push origin feature/amazing-feature`
6. **Open a Pull Request**

### Development Guidelines:
- Follow existing code style and patterns
- Write tests for new functionality
- Update documentation as needed
- Use conventional commit messages

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines. 

## 📄 License
This project is licensed under the GPL v2.0 License - see the [LICENSE](LICENSE) file for details.

---

<p align="center">
  Built with ❤️ by <a href="https://github.com/ArjunSrivastava1">Arjun Srivastava</a>
  <br>
  <sub>Making developer lives easier, one port at a time</sub>
</p>
