# testformat

A Go library for parsing various test and coverage report formats and converting them to JSON.

## Supported Formats

- **JUnit XML** - Standard test result format used by many testing frameworks (Java, Python, JavaScript, etc.)
- **Maven Surefire XML** - Enhanced JUnit format with additional features like flaky test detection and rerun support
- **JaCoCo XML** - Java code coverage reports with line, branch, method, and class coverage metrics
- **CTRF JSON** - Common Test Report Format - a standardized JSON schema for test reports

## Installation

```bash
go get github.com/mogthesprog/testformat
```

## Use Cases

- Converting XML test reports to JSON for consumption by dashboards and APIs
- Converting JUnit and Surefire reports to the standardized CTRF format
- Building test analytics and reporting tools
- Aggregating test results from multiple test frameworks in CI/CD pipelines
- Processing code coverage data for quality gates and metrics
- Integrating test results into monitoring and observability platforms

## Example

Parse a JUnit XML test report and convert it to JSON:

```go
package main

import (
    "encoding/json"
    "encoding/xml"
    "fmt"
    "log"

    "github.com/mogthesprog/testformat/junit"
)

func main() {
    xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="MyTests" tests="2" failures="1" time="0.5">
    <testcase name="test_success" classname="com.example.Test" time="0.3"/>
    <testcase name="test_failure" classname="com.example.Test" time="0.2">
        <failure message="Expected 5 but got 3" type="AssertionError">
AssertionError: Expected 5 but got 3
    at Test.java:42
        </failure>
    </testcase>
</testsuite>`

    var suite junit.TestSuite
    if err := xml.Unmarshal([]byte(xmlData), &suite); err != nil {
        log.Fatal(err)
    }

    jsonData, err := json.MarshalIndent(suite, "", "  ")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(jsonData))
}
```

This will output:

```json
{
  "name": "MyTests",
  "tests": 2,
  "failures": 1,
  "time": "0.5",
  "testCases": [
    {
      "name": "test_success",
      "className": "com.example.Test",
      "time": "0.3"
    },
    {
      "name": "test_failure",
      "className": "com.example.Test",
      "time": "0.2",
      "failure": {
        "message": "Expected 5 but got 3",
        "type": "AssertionError",
        "data": "\nAssertionError: Expected 5 but got 3\n    at Test.java:42\n        "
      }
    }
  ]
}
```

## Packages

Each format is provided in its own package:

- `github.com/mogthesprog/testformat/junit` - JUnit XML parsing
- `github.com/mogthesprog/testformat/surefire` - Maven Surefire XML parsing
- `github.com/mogthesprog/testformat/jacoco` - JaCoCo coverage XML parsing
- `github.com/mogthesprog/testformat/ctrf` - CTRF format conversion from JUnit and Surefire

All structs include both XML and JSON tags for seamless bidirectional conversion.

## Converting to CTRF Format

The CTRF (Common Test Report Format) package provides functions to convert JUnit and Maven Surefire test reports to the standardized CTRF JSON format:

```go
package main

import (
    "encoding/json"
    "encoding/xml"
    "fmt"
    "log"

    "github.com/mogthesprog/testformat/junit"
    "github.com/mogthesprog/testformat/ctrf"
)

func main() {
    xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="MyTests" tests="2" failures="1" time="0.5">
    <testcase name="test_success" classname="com.example.Test" time="0.3"/>
    <testcase name="test_failure" classname="com.example.Test" time="0.2">
        <failure message="Expected 5 but got 3" type="AssertionError">
AssertionError: Expected 5 but got 3
    at Test.java:42
        </failure>
    </testcase>
</testsuite>`

    var suite junit.TestSuite
    if err := xml.Unmarshal([]byte(xmlData), &suite); err != nil {
        log.Fatal(err)
    }

    // Convert to CTRF format
    report := ctrf.FromJUnitTestSuite(&suite)

    jsonData, err := json.MarshalIndent(report, "", "  ")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(jsonData))
}
```

This will output a CTRF-compliant JSON report:

```json
{
  "reportFormat": "CTRF",
  "specVersion": "0.0.0",
  "timestamp": "2025-10-25T12:00:00Z",
  "generatedBy": "github.com/mogthesprog/testformat",
  "results": {
    "tool": {
      "name": "JUnit"
    },
    "summary": {
      "tests": 2,
      "passed": 1,
      "failed": 1,
      "pending": 0,
      "skipped": 0,
      "other": 0,
      "suites": 1,
      "start": 0,
      "stop": 0
    },
    "tests": [
      {
        "name": "test_success",
        "status": "passed",
        "duration": 300,
        "suite": "com.example.Test"
      },
      {
        "name": "test_failure",
        "status": "failed",
        "duration": 200,
        "suite": "com.example.Test",
        "message": "Expected 5 but got 3",
        "trace": "\nAssertionError: Expected 5 but got 3\n    at Test.java:42\n        "
      }
    ]
  }
}
```

### Available Conversion Functions

- `ctrf.FromJUnitTestSuite(*junit.TestSuite) *ctrf.Report` - Convert a single JUnit test suite
- `ctrf.FromJUnitTestSuites(*junit.TestSuites) *ctrf.Report` - Convert multiple JUnit test suites
- `ctrf.FromSurefireTestSuite(*surefire.TestSuite) *ctrf.Report` - Convert a Maven Surefire test suite

### CTRF Format Details

The CTRF (Common Test Report Format) is a standardized JSON schema for test reports. Learn more at [ctrf.io](https://ctrf.io).

#### Comprehensive Field Mapping

Our CTRF implementation provides extensive field mapping with zero data loss:

**Core Test Fields:**
- `name`, `status`, `duration` - Required CTRF fields
- `suite`, `message`, `trace` - Basic test information
- `type` - Error/failure type (e.g., "AssertionError")
- `filePath`, `line` - Source file location for IDE integration
- `rawStatus` - Original status from test framework (failure/error/skipped)

**Test Output:**
- `stdout` - Array of stdout lines from test execution
- `stderr` - Array of stderr lines from test execution

**Test Metadata:**
- `parameters` - Test properties and configuration (from JUnit Properties)
- `tags` - Test categories and groups (from Surefire Group)

**Flaky Test Support:**
- `flaky` - Boolean flag for flaky tests
- `retries` - Number of retry attempts
- `retryAttempts` - Detailed retry history with status, duration, output

**Extended Data:**
- `extra.assertions` - Number of assertions (JUnit)
- `extra.additionalFailures` - Multiple failures per test (Surefire)
- `environment.hostname` - Test execution host (JUnit)

#### Features
- ✅ Zero data loss - All JUnit and Surefire fields preserved
- ✅ Full CTRF schema compliance
- ✅ Rich metadata for test analytics
- ✅ IDE integration support (file paths and line numbers)
- ✅ Flaky test detection and retry tracking
- ✅ Test output capture (stdout/stderr)
- ✅ Comprehensive test coverage

#### Known Limitations

- **Per-test timing** (`start`/`stop` timestamps) - Not available in JUnit/Surefire XML formats
- **Suite metadata** (package name, suite ID) - Not preserved in CTRF conversion (low priority)

