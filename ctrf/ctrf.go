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
	// Required fields
	Name     string `json:"name"`
	Status   string `json:"status"` // passed, failed, skipped, pending, other
	Duration int64  `json:"duration"` // milliseconds

	// Optional standard CTRF fields
	Suite         string         `json:"suite,omitempty"`
	Message       string         `json:"message,omitempty"`
	Trace         string         `json:"trace,omitempty"`
	FilePath      string         `json:"filePath,omitempty"`
	Line          int            `json:"line,omitempty"`
	RawStatus     string         `json:"rawStatus,omitempty"`
	Tags          []string       `json:"tags,omitempty"`
	Type          string         `json:"type,omitempty"` // Error/failure type
	Retries       int            `json:"retries,omitempty"`
	RetryAttempts []RetryAttempt `json:"retryAttempts,omitempty"`
	Flaky         bool           `json:"flaky,omitempty"`
	Stdout        []string       `json:"stdout,omitempty"`
	Stderr        []string       `json:"stderr,omitempty"`
	Parameters    interface{}    `json:"parameters,omitempty"`
	Extra         interface{}    `json:"extra,omitempty"`
}

// RetryAttempt represents a test retry attempt
type RetryAttempt struct {
	Attempt  int         `json:"attempt"`
	Status   string      `json:"status"`
	Duration int64       `json:"duration,omitempty"`
	Message  string      `json:"message,omitempty"`
	Trace    string      `json:"trace,omitempty"`
	Stdout   []string    `json:"stdout,omitempty"`
	Stderr   []string    `json:"stderr,omitempty"`
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

// splitLines splits a string into lines, filtering out empty lines
func splitLines(s string) []string {
	if s == "" {
		return nil
	}

	lines := []string{}
	for _, line := range splitString(s, '\n') {
		if trimmed := trimSpace(line); trimmed != "" {
			lines = append(lines, trimmed)
		}
	}

	if len(lines) == 0 {
		return nil
	}
	return lines
}

// splitString splits a string by a delimiter
func splitString(s string, delim rune) []string {
	var result []string
	var current string

	for _, r := range s {
		if r == delim {
			result = append(result, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	result = append(result, current)
	return result
}

// trimSpace removes leading and trailing whitespace
func trimSpace(s string) string {
	start := 0
	end := len(s)

	// Trim leading whitespace
	for start < end {
		r := rune(s[start])
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			break
		}
		start++
	}

	// Trim trailing whitespace
	for end > start {
		r := rune(s[end-1])
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			break
		}
		end--
	}

	return s[start:end]
}
