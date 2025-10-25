package ctrf

import (
	"strconv"

	"github.com/mogthesprog/testformat/surefire"
)

// FromSurefireTestSuite converts a Surefire TestSuite to a CTRF Report
func FromSurefireTestSuite(suite *surefire.TestSuite) *Report {
	report := NewReport("Maven Surefire")

	// Set tool version if available
	if suite.Version != "" {
		report.Results.Tool.Version = suite.Version
	}

	// Parse counts from strings
	tests := parseInt(suite.Tests)
	failures := parseInt(suite.Failures)
	errors := parseInt(suite.Errors)
	skipped := parseInt(suite.Skipped)

	// Process test cases
	for _, tc := range suite.TestCases {
		test := convertSurefireTestCase(&tc, suite.Name)
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

// convertSurefireTestCase converts a Surefire TestCase to a CTRF TestResult
func convertSurefireTestCase(tc *surefire.TestCase, suiteName string) TestResult {
	test := TestResult{
		Name:     tc.Name,
		Status:   determineSurefireStatus(tc),
		Duration: parseTimeString(tc.Time),
	}

	// Set suite name from classname or suite name
	if tc.Classname != "" {
		test.Suite = tc.Classname
	} else if suiteName != "" {
		test.Suite = suiteName
	}

	// Map stdout and stderr
	test.Stdout = splitLines(tc.SystemOut)
	test.Stderr = splitLines(tc.SystemErr)

	// Map group to tags
	if tc.Group != "" {
		test.Tags = []string{tc.Group}
	}

	// Detect flaky tests
	hasFlaky := len(tc.FlakyFailures) > 0 || len(tc.FlakyErrors) > 0
	hasRerun := len(tc.RerunFailures) > 0 || len(tc.RerunErrors) > 0

	if hasFlaky || hasRerun {
		test.Flaky = true
	}

	// Build retry attempts from rerun and flaky data
	var retryAttempts []RetryAttempt

	// Add flaky failures as retry attempts
	for i, ff := range tc.FlakyFailures {
		retryAttempts = append(retryAttempts, RetryAttempt{
			Attempt: i + 1,
			Status:  "failed",
			Message: ff.Message,
			Trace:   ff.StackTrace,
			Stdout:  splitLines(ff.SystemOut),
			Stderr:  splitLines(ff.SystemErr),
		})
	}

	// Add flaky errors as retry attempts
	for _, fe := range tc.FlakyErrors {
		retryAttempts = append(retryAttempts, RetryAttempt{
			Attempt: len(retryAttempts) + 1,
			Status:  "failed",
			Message: fe.Message,
			Trace:   fe.StackTrace,
			Stdout:  splitLines(fe.SystemOut),
			Stderr:  splitLines(fe.SystemErr),
		})
	}

	// Add rerun failures as retry attempts
	for _, rf := range tc.RerunFailures {
		retryAttempts = append(retryAttempts, RetryAttempt{
			Attempt: len(retryAttempts) + 1,
			Status:  "failed",
			Message: rf.Message,
			Trace:   rf.StackTrace,
			Stdout:  splitLines(rf.SystemOut),
			Stderr:  splitLines(rf.SystemErr),
		})
	}

	// Add rerun errors as retry attempts
	for _, re := range tc.RerunErrors {
		retryAttempts = append(retryAttempts, RetryAttempt{
			Attempt: len(retryAttempts) + 1,
			Status:  "failed",
			Message: re.Message,
			Trace:   re.StackTrace,
			Stdout:  splitLines(re.SystemOut),
			Stderr:  splitLines(re.SystemErr),
		})
	}

	if len(retryAttempts) > 0 {
		test.RetryAttempts = retryAttempts
		test.Retries = len(retryAttempts)
	}

	// Add primary failure/error message, trace, and type
	if len(tc.Failures) > 0 {
		test.Message = tc.Failures[0].Message
		test.Trace = tc.Failures[0].Data
		test.Type = tc.Failures[0].Type
		test.RawStatus = "failure"

		// If there are multiple failures, add them to Extra
		if len(tc.Failures) > 1 {
			if test.Extra == nil {
				test.Extra = make(map[string]interface{})
			}
			additionalFailures := []map[string]string{}
			for i := 1; i < len(tc.Failures); i++ {
				additionalFailures = append(additionalFailures, map[string]string{
					"message": tc.Failures[i].Message,
					"type":    tc.Failures[i].Type,
					"data":    tc.Failures[i].Data,
				})
			}
			test.Extra.(map[string]interface{})["additionalFailures"] = additionalFailures
		}
	} else if len(tc.Errors) > 0 {
		test.Message = tc.Errors[0].Message
		test.Trace = tc.Errors[0].Data
		test.Type = tc.Errors[0].Type
		test.RawStatus = "error"

		// If there are multiple errors, add them to Extra
		if len(tc.Errors) > 1 {
			if test.Extra == nil {
				test.Extra = make(map[string]interface{})
			}
			additionalErrors := []map[string]string{}
			for i := 1; i < len(tc.Errors); i++ {
				additionalErrors = append(additionalErrors, map[string]string{
					"message": tc.Errors[i].Message,
					"type":    tc.Errors[i].Type,
					"data":    tc.Errors[i].Data,
				})
			}
			test.Extra.(map[string]interface{})["additionalErrors"] = additionalErrors
		}
	} else if tc.Skipped != nil {
		test.Message = tc.Skipped.Message
		test.RawStatus = "skipped"
	}

	return test
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
