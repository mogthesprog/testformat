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

func TestFromJUnit_ExtendedFields(t *testing.T) {
	suite := &junit.TestSuite{
		Name:     "ExtendedTest",
		Tests:    1,
		Failures: 1,
		Hostname: "test-host",
		TestCases: []junit.TestCase{
			{
				Name:       "test_with_all_fields",
				Classname:  "com.example.TestClass",
				Time:       "1.5",
				File:       "src/test/TestClass.java",
				Line:       42,
				Assertions: 5,
				SystemOut:  "Test output line 1\nTest output line 2\n",
				SystemErr:  "Error output\n",
				Properties: &junit.Properties{
					Properties: []junit.Property{
						{Name: "param1", Value: "value1"},
						{Name: "param2", Value: "value2"},
					},
				},
				Failure: &junit.Failure{
					Message: "Test failed",
					Type:    "AssertionError",
					Data:    "Stack trace here",
				},
			},
		},
	}

	report := FromJUnitTestSuite(suite)

	// Verify environment has hostname
	if report.Environment == nil {
		t.Fatal("Expected Environment to be set")
	}
	extra := report.Environment.Extra.(map[string]interface{})
	if extra["hostname"] != "test-host" {
		t.Errorf("Expected hostname in environment, got %v", extra["hostname"])
	}

	// Verify test fields
	test := report.Results.Tests[0]

	if test.FilePath != "src/test/TestClass.java" {
		t.Errorf("Expected filePath to be set, got %s", test.FilePath)
	}

	if test.Line != 42 {
		t.Errorf("Expected line to be 42, got %d", test.Line)
	}

	if test.Type != "AssertionError" {
		t.Errorf("Expected type to be AssertionError, got %s", test.Type)
	}

	if test.RawStatus != "failure" {
		t.Errorf("Expected rawStatus to be failure, got %s", test.RawStatus)
	}

	if len(test.Stdout) != 2 {
		t.Errorf("Expected 2 stdout lines, got %d", len(test.Stdout))
	}

	if len(test.Stderr) != 1 {
		t.Errorf("Expected 1 stderr line, got %d", len(test.Stderr))
	}

	if test.Parameters == nil {
		t.Fatal("Expected parameters to be set")
	}
	params := test.Parameters.(map[string]interface{})
	if params["param1"] != "value1" {
		t.Errorf("Expected param1=value1, got %v", params["param1"])
	}

	if test.Extra == nil {
		t.Fatal("Expected Extra to be set for assertions")
	}
	extraMap := test.Extra.(map[string]interface{})
	if extraMap["assertions"] != 5 {
		t.Errorf("Expected assertions=5, got %v", extraMap["assertions"])
	}
}

func TestFromSurefire_ExtendedFields(t *testing.T) {
	suite := &surefire.TestSuite{
		Name:     "SurefireExtendedTest",
		Tests:    "2",
		Failures: "1",
		Errors:   "0",
		Skipped:  "0",
		Version:  "3.0.0",
		TestCases: []surefire.TestCase{
			{
				Name:      "test_with_group",
				Classname: "com.example.Test",
				Time:      "0.5",
				Group:     "integration",
				SystemOut: "Output from test\n",
			},
			{
				Name:      "flaky_test_with_reruns",
				Classname: "com.example.FlakyTest",
				Time:      "1.0",
				Group:     "flaky",
				FlakyFailures: []surefire.FlakyFailure{
					{
						Message:    "First attempt failed",
						StackTrace: "Stack trace 1",
						SystemOut:  "Attempt 1 output",
					},
				},
				RerunFailures: []surefire.RerunFailure{
					{
						Message:    "Second attempt failed",
						StackTrace: "Stack trace 2",
						SystemErr:  "Attempt 2 error",
					},
				},
				Failures: []surefire.Failure{
					{
						Message: "Final failure",
						Type:    "TestFailure",
						Data:    "Final trace",
					},
					{
						Message: "Additional failure",
						Type:    "AssertionError",
						Data:    "Additional trace",
					},
				},
			},
		},
	}

	report := FromSurefireTestSuite(suite)

	// Verify tool version
	if report.Results.Tool.Version != "3.0.0" {
		t.Errorf("Expected tool version 3.0.0, got %s", report.Results.Tool.Version)
	}

	// Check first test with group
	test1 := report.Results.Tests[0]
	if len(test1.Tags) != 1 || test1.Tags[0] != "integration" {
		t.Errorf("Expected tags=[integration], got %v", test1.Tags)
	}
	if len(test1.Stdout) == 0 {
		t.Error("Expected stdout to be set")
	}

	// Check flaky test
	test2 := report.Results.Tests[1]
	if !test2.Flaky {
		t.Error("Expected flaky to be true")
	}

	if test2.Retries != 2 {
		t.Errorf("Expected 2 retries, got %d", test2.Retries)
	}

	if len(test2.RetryAttempts) != 2 {
		t.Errorf("Expected 2 retry attempts, got %d", len(test2.RetryAttempts))
	}

	if test2.RetryAttempts[0].Message != "First attempt failed" {
		t.Errorf("Expected first retry message, got %s", test2.RetryAttempts[0].Message)
	}

	if test2.Type != "TestFailure" {
		t.Errorf("Expected type TestFailure, got %s", test2.Type)
	}

	// Check multiple failures in Extra
	if test2.Extra == nil {
		t.Fatal("Expected Extra to contain additional failures")
	}
	extraMap := test2.Extra.(map[string]interface{})
	additionalFailures := extraMap["additionalFailures"].([]map[string]string)
	if len(additionalFailures) != 1 {
		t.Errorf("Expected 1 additional failure, got %d", len(additionalFailures))
	}
	if additionalFailures[0]["message"] != "Additional failure" {
		t.Errorf("Expected additional failure message, got %s", additionalFailures[0]["message"])
	}
}

func TestSplitLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"empty string", "", nil},
		{"single line", "line1", []string{"line1"}},
		{"two lines", "line1\nline2", []string{"line1", "line2"}},
		{"with empty lines", "line1\n\nline2\n", []string{"line1", "line2"}},
		{"only whitespace", "   \n  \t  \n", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitLines(tt.input)
			if diff := cmp.Diff(tt.expected, result); diff != "" {
				t.Errorf("splitLines() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
