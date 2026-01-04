// scanner/detectors/
// node_detector.go
//func detectNode(process *ProcessAnalysis) bool {
// Check: command contains "node", "npm", "npx"
// Check: working directory has package.json
// Check: process name is "node"
//	return strings.Contains(process.CommandLine, "node") ||
//		strings.Contains(process.CommandLine, "npm") ||
//		strings.Contains(process.CommandLine, "npx")
//}

// python_detector.go
//func detectPython(process *ProcessAnalysis) bool {
// Check: command contains "python", "pip", "streamlit"
// Check: working directory has requirements.txt, .py files
//return strings.Contains(process.CommandLine, "python") ||
//	strings.Contains(process.CommandLine, "streamlit") ||
//	strings.Contains(process.CommandLine, "fastapi")
//}

// postgres_detector.go
//func detectPostgres(process *ProcessAnalysis) bool {
//	return strings.Contains(process.Name, "postgres") ||
//		strings.Contains(process.CommandLine, "postgres")
//}