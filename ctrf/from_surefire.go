package ctrf

import (
	"strconv"

	"github.com/mogthesprog/testformat/surefire"
)

// FromSurefireTestSuite converts a Surefire TestSuite to a CTRF Report
func FromSurefireTestSuite(suite *surefire.TestSuite) *Report {
	report := NewReport("Maven Surefire")

	// Parse counts from strings
	tests := parseInt(suite.Tests)
	failures := parseInt(suite.Failures)
	errors := parseInt(suite.Errors)
	skipped := parseInt(suite.Skipped)

	// Process test cases
	for _, tc := range suite.TestCases {
		test := TestResult{
			Name:     tc.Name,
			Status:   determineSurefireStatus(&tc),
			Duration: parseTimeString(tc.Time),
		}

		// Set suite name from classname or suite name
		if tc.Classname != "" {
			test.Suite = tc.Classname
		} else if suite.Name != "" {
			test.Suite = suite.Name
		}

		// Add failure/error message and trace
		if len(tc.Failures) > 0 {
			test.Message = tc.Failures[0].Message
			test.Trace = tc.Failures[0].Data
		} else if len(tc.Errors) > 0 {
			test.Message = tc.Errors[0].Message
			test.Trace = tc.Errors[0].Data
		} else if tc.Skipped != nil {
			test.Message = tc.Skipped.Message
		}

		report.Results.Tests = append(report.Results.Tests, test)
	}

	// Build summary
	report.Results.Summary = &Summary{
		Tests:   tests,
		Passed:  tests - failures - errors - skipped,
		Failed:  failures + errors,
		Skipped: skipped,
		Suites:  1,
		Start:   0, // Surefire doesn't provide start time
		Stop:    0, // Surefire doesn't provide stop time
	}

	return report
}

// parseInt parses a string to an int, returning 0 on error
func parseInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}

// determineSurefireStatus determines the status of a Surefire test case
func determineSurefireStatus(tc *surefire.TestCase) string {
	hasFailure := len(tc.Failures) > 0 || len(tc.FlakyFailures) > 0 || len(tc.RerunFailures) > 0
	hasError := len(tc.Errors) > 0 || len(tc.FlakyErrors) > 0 || len(tc.RerunErrors) > 0
	hasSkipped := tc.Skipped != nil

	if hasFailure || hasError {
		return "failed"
	}
	if hasSkipped {
		return "skipped"
	}
	return "passed"
}
