package ctrf

import (
	"github.com/mogthesprog/testformat/junit"
)

// FromJUnitTestSuite converts a JUnit TestSuite to a CTRF Report
func FromJUnitTestSuite(suite *junit.TestSuite) *Report {
	report := NewReport("JUnit")

	// Process test cases
	for _, tc := range suite.TestCases {
		test := TestResult{
			Name:     tc.Name,
			Status:   mapStatus(tc.Failure != nil, tc.Error != nil, tc.Skipped != nil),
			Duration: parseTimeString(tc.Time),
		}

		// Set suite name from classname or suite name
		if tc.Classname != "" {
			test.Suite = tc.Classname
		} else if suite.Name != "" {
			test.Suite = suite.Name
		}

		// Add failure/error message and trace
		if tc.Failure != nil {
			test.Message = tc.Failure.Message
			test.Trace = tc.Failure.Data
		} else if tc.Error != nil {
			test.Message = tc.Error.Message
			test.Trace = tc.Error.Data
		} else if tc.Skipped != nil {
			test.Message = tc.Skipped.Message
		}

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
			test := TestResult{
				Name:     tc.Name,
				Status:   mapStatus(tc.Failure != nil, tc.Error != nil, tc.Skipped != nil),
				Duration: parseTimeString(tc.Time),
			}

			// Set suite name from classname or suite name
			if tc.Classname != "" {
				test.Suite = tc.Classname
			} else if suite.Name != "" {
				test.Suite = suite.Name
			}

			// Add failure/error message and trace
			if tc.Failure != nil {
				test.Message = tc.Failure.Message
				test.Trace = tc.Failure.Data
			} else if tc.Error != nil {
				test.Message = tc.Error.Message
				test.Trace = tc.Error.Data
			} else if tc.Skipped != nil {
				test.Message = tc.Skipped.Message
			}

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
