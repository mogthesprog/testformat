package ctrf

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mogthesprog/testformat/junit"
	"github.com/mogthesprog/testformat/surefire"
)

func TestFromJUnitTestSuite_BasicConversion(t *testing.T) {
	suite := &junit.TestSuite{
		Name:     "MyTestSuite",
		Tests:    3,
		Failures: 1,
		Errors:   0,
		Skipped:  1,
		Time:     "1.5",
		TestCases: []junit.TestCase{
			{
				Name:      "test_pass",
				Classname: "com.example.Test",
				Time:      "0.5",
			},
			{
				Name:      "test_fail",
				Classname: "com.example.Test",
				Time:      "0.7",
				Failure: &junit.Failure{
					Message: "Expected 5 but got 3",
					Type:    "AssertionError",
					Data:    "AssertionError: Expected 5 but got 3\n    at Test.java:42",
				},
			},
			{
				Name:      "test_skip",
				Classname: "com.example.Test",
				Time:      "0.3",
				Skipped: &junit.Skipped{
					Message: "Test skipped",
				},
			},
		},
	}

	report := FromJUnitTestSuite(suite)

	// Verify report structure
	if report.ReportFormat != "CTRF" {
		t.Errorf("Expected reportFormat to be CTRF, got %s", report.ReportFormat)
	}

	if report.Results.Tool.Name != "JUnit" {
		t.Errorf("Expected tool name to be JUnit, got %s", report.Results.Tool.Name)
	}

	// Verify summary
	if report.Results.Summary.Tests != 3 {
		t.Errorf("Expected 3 tests, got %d", report.Results.Summary.Tests)
	}
	if report.Results.Summary.Passed != 1 {
		t.Errorf("Expected 1 passed test, got %d", report.Results.Summary.Passed)
	}
	if report.Results.Summary.Failed != 1 {
		t.Errorf("Expected 1 failed test, got %d", report.Results.Summary.Failed)
	}
	if report.Results.Summary.Skipped != 1 {
		t.Errorf("Expected 1 skipped test, got %d", report.Results.Summary.Skipped)
	}

	// Verify tests
	if len(report.Results.Tests) != 3 {
		t.Fatalf("Expected 3 test results, got %d", len(report.Results.Tests))
	}

	// Check passed test
	if report.Results.Tests[0].Status != "passed" {
		t.Errorf("Expected test_pass status to be passed, got %s", report.Results.Tests[0].Status)
	}

	// Check failed test
	if report.Results.Tests[1].Status != "failed" {
		t.Errorf("Expected test_fail status to be failed, got %s", report.Results.Tests[1].Status)
	}
	if report.Results.Tests[1].Message != "Expected 5 but got 3" {
		t.Errorf("Expected failure message, got %s", report.Results.Tests[1].Message)
	}

	// Check skipped test
	if report.Results.Tests[2].Status != "skipped" {
		t.Errorf("Expected test_skip status to be skipped, got %s", report.Results.Tests[2].Status)
	}
}

func TestFromJUnitTestSuites_MultipleeSuites(t *testing.T) {
	suites := &junit.TestSuites{
		Name: "AllTests",
		Suites: []junit.TestSuite{
			{
				Name:     "Suite1",
				Tests:    2,
				Failures: 0,
				TestCases: []junit.TestCase{
					{Name: "test1", Classname: "Suite1", Time: "0.1"},
					{Name: "test2", Classname: "Suite1", Time: "0.2"},
				},
			},
			{
				Name:     "Suite2",
				Tests:    1,
				Failures: 1,
				TestCases: []junit.TestCase{
					{
						Name:      "test3",
						Classname: "Suite2",
						Time:      "0.3",
						Failure: &junit.Failure{
							Message: "Failed",
						},
					},
				},
			},
		},
	}

	report := FromJUnitTestSuites(suites)

	// Verify summary
	if report.Results.Summary.Tests != 3 {
		t.Errorf("Expected 3 total tests, got %d", report.Results.Summary.Tests)
	}
	if report.Results.Summary.Passed != 2 {
		t.Errorf("Expected 2 passed tests, got %d", report.Results.Summary.Passed)
	}
	if report.Results.Summary.Failed != 1 {
		t.Errorf("Expected 1 failed test, got %d", report.Results.Summary.Failed)
	}
	if report.Results.Summary.Suites != 2 {
		t.Errorf("Expected 2 suites, got %d", report.Results.Summary.Suites)
	}

	// Verify all tests are included
	if len(report.Results.Tests) != 3 {
		t.Errorf("Expected 3 test results, got %d", len(report.Results.Tests))
	}
}

func TestFromSurefireTestSuite_BasicConversion(t *testing.T) {
	suite := &surefire.TestSuite{
		Name:     "MyTestSuite",
		Tests:    "2",
		Failures: "1",
		Errors:   "0",
		Skipped:  "0",
		Time:     "1.0",
		TestCases: []surefire.TestCase{
			{
				Name:      "test_pass",
				Classname: "com.example.Test",
				Time:      "0.5",
			},
			{
				Name:      "test_fail",
				Classname: "com.example.Test",
				Time:      "0.5",
				Failures: []surefire.Failure{
					{
						Message: "Test failed",
						Type:    "AssertionError",
						Data:    "Stack trace here",
					},
				},
			},
		},
	}

	report := FromSurefireTestSuite(suite)

	// Verify report structure
	if report.ReportFormat != "CTRF" {
		t.Errorf("Expected reportFormat to be CTRF, got %s", report.ReportFormat)
	}

	if report.Results.Tool.Name != "Maven Surefire" {
		t.Errorf("Expected tool name to be Maven Surefire, got %s", report.Results.Tool.Name)
	}

	// Verify summary
	if report.Results.Summary.Tests != 2 {
		t.Errorf("Expected 2 tests, got %d", report.Results.Summary.Tests)
	}
	if report.Results.Summary.Passed != 1 {
		t.Errorf("Expected 1 passed test, got %d", report.Results.Summary.Passed)
	}
	if report.Results.Summary.Failed != 1 {
		t.Errorf("Expected 1 failed test, got %d", report.Results.Summary.Failed)
	}

	// Verify tests
	if len(report.Results.Tests) != 2 {
		t.Fatalf("Expected 2 test results, got %d", len(report.Results.Tests))
	}

	// Check failed test has message and trace
	if report.Results.Tests[1].Message != "Test failed" {
		t.Errorf("Expected failure message, got %s", report.Results.Tests[1].Message)
	}
	if report.Results.Tests[1].Trace != "Stack trace here" {
		t.Errorf("Expected stack trace, got %s", report.Results.Tests[1].Trace)
	}
}

func TestFromSurefireTestSuite_FlakyTests(t *testing.T) {
	suite := &surefire.TestSuite{
		Name:     "FlakyTests",
		Tests:    "1",
		Failures: "0",
		Errors:   "0",
		Skipped:  "0",
		TestCases: []surefire.TestCase{
			{
				Name:      "flaky_test",
				Classname: "com.example.FlakyTest",
				Time:      "0.5",
				FlakyFailures: []surefire.FlakyFailure{
					{
						Message: "Flaky failure",
					},
				},
			},
		},
	}

	report := FromSurefireTestSuite(suite)

	// Flaky failures should be treated as failures
	if report.Results.Tests[0].Status != "failed" {
		t.Errorf("Expected flaky test to have failed status, got %s", report.Results.Tests[0].Status)
	}
}

func TestParseTimeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"empty string", "", 0},
		{"0.5 seconds", "0.5", 500},
		{"1.0 seconds", "1.0", 1000},
		{"1.234 seconds", "1.234", 1234},
		{"invalid", "invalid", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseTimeString(tt.input)
			if result != tt.expected {
				t.Errorf("parseTimeString(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMapStatus(t *testing.T) {
	tests := []struct {
		name       string
		hasFailure bool
		hasError   bool
		hasSkipped bool
		expected   string
	}{
		{"passed", false, false, false, "passed"},
		{"failed with failure", true, false, false, "failed"},
		{"failed with error", false, true, false, "failed"},
		{"failed with both", true, true, false, "failed"},
		{"skipped", false, false, true, "skipped"},
		{"failed takes precedence", true, false, true, "failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapStatus(tt.hasFailure, tt.hasError, tt.hasSkipped)
			if result != tt.expected {
				t.Errorf("mapStatus(%v, %v, %v) = %s, want %s",
					tt.hasFailure, tt.hasError, tt.hasSkipped, result, tt.expected)
			}
		})
	}
}

func TestNewReport(t *testing.T) {
	report := NewReport("TestTool")

	if report.ReportFormat != "CTRF" {
		t.Errorf("Expected reportFormat to be CTRF, got %s", report.ReportFormat)
	}

	if report.SpecVersion != "0.0.0" {
		t.Errorf("Expected specVersion to be 0.0.0, got %s", report.SpecVersion)
	}

	if report.Results == nil {
		t.Fatal("Expected results to be initialized")
	}

	if report.Results.Tool.Name != "TestTool" {
		t.Errorf("Expected tool name to be TestTool, got %s", report.Results.Tool.Name)
	}

	if report.Results.Summary == nil {
		t.Fatal("Expected summary to be initialized")
	}

	if report.Results.Tests == nil {
		t.Fatal("Expected tests array to be initialized")
	}

	if report.Timestamp == "" {
		t.Error("Expected timestamp to be set")
	}
}

func TestJSONMarshaling(t *testing.T) {
	suite := &junit.TestSuite{
		Name:     "JSONTest",
		Tests:    1,
		Failures: 0,
		TestCases: []junit.TestCase{
			{
				Name:      "test1",
				Classname: "TestClass",
				Time:      "0.5",
			},
		},
	}

	report := FromJUnitTestSuite(suite)

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal report to JSON: %v", err)
	}

	// Unmarshal back
	var unmarshaled Report
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Compare (excluding timestamp which may vary)
	if unmarshaled.ReportFormat != report.ReportFormat {
		t.Errorf("ReportFormat mismatch after round-trip")
	}

	if unmarshaled.Results.Summary.Tests != report.Results.Summary.Tests {
		t.Errorf("Summary.Tests mismatch after round-trip")
	}

	if len(unmarshaled.Results.Tests) != len(report.Results.Tests) {
		t.Errorf("Tests array length mismatch after round-trip")
	}

	if diff := cmp.Diff(report.Results.Tests[0].Name, unmarshaled.Results.Tests[0].Name); diff != "" {
		t.Errorf("Test name mismatch (-want +got):\n%s", diff)
	}
}
