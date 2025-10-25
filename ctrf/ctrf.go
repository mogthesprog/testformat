package ctrf

import (
	"strconv"
	"time"
)

// Report represents the root CTRF test report
type Report struct {
	ReportFormat string       `json:"reportFormat"`
	SpecVersion  string       `json:"specVersion"`
	ReportID     string       `json:"reportId,omitempty"`
	Timestamp    string       `json:"timestamp,omitempty"`
	GeneratedBy  string       `json:"generatedBy,omitempty"`
	Results      *Results     `json:"results"`
	Environment  *Environment `json:"environment,omitempty"`
	Extra        interface{}  `json:"extra,omitempty"`
}

// Results contains the test execution data
type Results struct {
	Tool    *Tool         `json:"tool"`
	Summary *Summary      `json:"summary"`
	Tests   []TestResult  `json:"tests"`
	Extra   interface{}   `json:"extra,omitempty"`
}

// Tool identifies the testing framework
type Tool struct {
	Name    string      `json:"name"`
	Version string      `json:"version,omitempty"`
	Extra   interface{} `json:"extra,omitempty"`
}

// Summary contains aggregated test metrics
type Summary struct {
	Tests   int         `json:"tests"`
	Passed  int         `json:"passed"`
	Failed  int         `json:"failed"`
	Pending int         `json:"pending"`
	Skipped int         `json:"skipped"`
	Other   int         `json:"other"`
	Suites  int         `json:"suites,omitempty"`
	Start   int64       `json:"start"`
	Stop    int64       `json:"stop"`
	Extra   interface{} `json:"extra,omitempty"`
}

// TestResult represents an individual test case
type TestResult struct {
	Name     string      `json:"name"`
	Status   string      `json:"status"`
	Duration int64       `json:"duration"`
	Suite    string      `json:"suite,omitempty"`
	Message  string      `json:"message,omitempty"`
	Trace    string      `json:"trace,omitempty"`
	Extra    interface{} `json:"extra,omitempty"`
}

// Environment captures build context and execution details
type Environment struct {
	AppName        string      `json:"appName,omitempty"`
	AppVersion     string      `json:"appVersion,omitempty"`
	OSPlatform     string      `json:"osPlatform,omitempty"`
	OSRelease      string      `json:"osRelease,omitempty"`
	OSVersion      string      `json:"osVersion,omitempty"`
	BuildName      string      `json:"buildName,omitempty"`
	BuildNumber    string      `json:"buildNumber,omitempty"`
	BuildURL       string      `json:"buildUrl,omitempty"`
	RepositoryName string      `json:"repositoryName,omitempty"`
	RepositoryURL  string      `json:"repositoryUrl,omitempty"`
	Branch         string      `json:"branch,omitempty"`
	Commit         string      `json:"commit,omitempty"`
	Extra          interface{} `json:"extra,omitempty"`
}

// NewReport creates a new CTRF report with default values
func NewReport(toolName string) *Report {
	return &Report{
		ReportFormat: "CTRF",
		SpecVersion:  "0.0.0",
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		GeneratedBy:  "github.com/mogthesprog/testformat",
		Results: &Results{
			Tool: &Tool{
				Name: toolName,
			},
			Summary: &Summary{},
			Tests:   []TestResult{},
		},
	}
}

// parseTimeString attempts to parse a time string and convert to milliseconds
func parseTimeString(timeStr string) int64 {
	if timeStr == "" {
		return 0
	}

	// Try parsing as float (seconds)
	if seconds, err := strconv.ParseFloat(timeStr, 64); err == nil {
		return int64(seconds * 1000) // convert to milliseconds
	}

	return 0
}

// mapStatus converts JUnit/Surefire status to CTRF status
func mapStatus(hasFailure, hasError, hasSkipped bool) string {
	if hasFailure || hasError {
		return "failed"
	}
	if hasSkipped {
		return "skipped"
	}
	return "passed"
}
