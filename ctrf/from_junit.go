package ctrf

import (
	"github.com/mogthesprog/testformat/junit"
)

// FromJUnitTestSuite converts a JUnit TestSuite to a CTRF Report
func FromJUnitTestSuite(suite *junit.TestSuite) *Report {
	report := NewReport("JUnit")

	// Set environment if hostname is available
	if suite.Hostname != "" {
		report.Environment = &Environment{
			Extra: map[string]interface{}{
				"hostname": suite.Hostname,
			},
		}
	}

	// Process test cases
	for _, tc := range suite.TestCases {
		test := convertJUnitTestCase(&tc, suite.Name)
		report.Results.Tests = append(report.Results.Tests, test)
	}

	// Build summary
	report.Results.Summary = &Summary{
		Tests:   suite.Tests,
		Passed:  suite.Tests - suite.Failures - suite.Errors - suite.Skipped,
		Failed:  suite.Failures + suite.Errors,
		Skipped: suite.Skipped,
		Suites:  1,
		Start:   0, // JUnit doesn't provide start time
		Stop:    0, // JUnit doesn't provide stop time
	}

	return report
}

// FromJUnitTestSuites converts JUnit TestSuites to a CTRF Report
func FromJUnitTestSuites(suites *junit.TestSuites) *Report {
	report := NewReport("JUnit")

	totalTests := 0
	totalPassed := 0
	totalFailed := 0
	totalSkipped := 0

	// Process all test suites
	for _, suite := range suites.Suites {
		for _, tc := range suite.TestCases {
			test := convertJUnitTestCase(&tc, suite.Name)
			report.Results.Tests = append(report.Results.Tests, test)

			// Update totals
			totalTests++
			if test.Status == "passed" {
				totalPassed++
			} else if test.Status == "failed" {
				totalFailed++
			} else if test.Status == "skipped" {
				totalSkipped++
			}
		}
	}

	// Build summary
	report.Results.Summary = &Summary{
		Tests:   totalTests,
		Passed:  totalPassed,
		Failed:  totalFailed,
		Skipped: totalSkipped,
		Suites:  len(suites.Suites),
		Start:   0, // JUnit doesn't provide start time
		Stop:    0, // JUnit doesn't provide stop time
	}

	return report
}

// convertJUnitTestCase converts a JUnit TestCase to a CTRF TestResult
func convertJUnitTestCase(tc *junit.TestCase, suiteName string) TestResult {
	test := TestResult{
		Name:     tc.Name,
		Status:   mapStatus(tc.Failure != nil, tc.Error != nil, tc.Skipped != nil),
		Duration: parseTimeString(tc.Time),
	}

	// Set suite name from classname or suite name
	if tc.Classname != "" {
		test.Suite = tc.Classname
	} else if suiteName != "" {
		test.Suite = suiteName
	}

	// Map file path and line number
	if tc.File != "" {
		test.FilePath = tc.File
	}
	if tc.Line > 0 {
		test.Line = tc.Line
	}

	// Map stdout and stderr
	test.Stdout = splitLines(tc.SystemOut)
	test.Stderr = splitLines(tc.SystemErr)

	// Map properties to parameters
	if tc.Properties != nil && len(tc.Properties.Properties) > 0 {
		params := make(map[string]interface{})
		for _, prop := range tc.Properties.Properties {
			if prop.Value != "" {
				params[prop.Name] = prop.Value
			} else if prop.Data != "" {
				params[prop.Name] = prop.Data
			}
		}
		if len(params) > 0 {
			test.Parameters = params
		}
	}

	// Add failure/error message, trace, type, and rawStatus
	if tc.Failure != nil {
		test.Message = tc.Failure.Message
		test.Trace = tc.Failure.Data
		test.Type = tc.Failure.Type
		test.RawStatus = "failure"
	} else if tc.Error != nil {
		test.Message = tc.Error.Message
		test.Trace = tc.Error.Data
		test.Type = tc.Error.Type
		test.RawStatus = "error"
	} else if tc.Skipped != nil {
		test.Message = tc.Skipped.Message
		test.RawStatus = "skipped"
	}

	// Add assertions count to Extra if present
	if tc.Assertions > 0 {
		if test.Extra == nil {
			test.Extra = make(map[string]interface{})
		}
		test.Extra.(map[string]interface{})["assertions"] = tc.Assertions
	}

	return test
}
