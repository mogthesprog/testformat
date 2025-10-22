# testformat

A Go library for parsing various test and coverage report formats and converting them to JSON.

## Supported Formats

- **JUnit XML** - Standard test result format used by many testing frameworks (Java, Python, JavaScript, etc.)
- **Maven Surefire XML** - Enhanced JUnit format with additional features like flaky test detection and rerun support
- **JaCoCo XML** - Java code coverage reports with line, branch, method, and class coverage metrics

## Installation

```bash
go get github.com/mogthesprog/testformat
```

## Use Cases

- Converting XML test reports to JSON for consumption by dashboards and APIs
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

All structs include both XML and JSON tags for seamless bidirectional conversion.
