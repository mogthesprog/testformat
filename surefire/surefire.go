package surefire

import "encoding/xml"

// TestSuite represents a Surefire test suite
type TestSuite struct {
	XMLName    xml.Name     `xml:"testsuite" json:"-"`
	Version    string       `xml:"version,attr,omitempty" json:"version,omitempty"`
	Name       string       `xml:"name,attr" json:"name"`
	Time       string       `xml:"time,attr,omitempty" json:"time,omitempty"`
	Tests      string       `xml:"tests,attr" json:"tests"`
	Errors     string       `xml:"errors,attr" json:"errors"`
	Skipped    string       `xml:"skipped,attr" json:"skipped"`
	Failures   string       `xml:"failures,attr" json:"failures"`
	Group      string       `xml:"group,attr,omitempty" json:"group,omitempty"`
	Properties *Properties  `xml:"properties,omitempty" json:"properties,omitempty"`
	TestCases  []TestCase   `xml:"testcase" json:"testCases,omitempty"`
	SystemOut  string       `xml:"system-out,omitempty" json:"systemOut,omitempty"`
	SystemErr  string       `xml:"system-err,omitempty" json:"systemErr,omitempty"`
}

// TestCase represents an individual test case
type TestCase struct {
	XMLName        xml.Name        `xml:"testcase" json:"-"`
	Name           string          `xml:"name,attr" json:"name"`
	Classname      string          `xml:"classname,attr,omitempty" json:"className,omitempty"`
	Group          string          `xml:"group,attr,omitempty" json:"group,omitempty"`
	Time           string          `xml:"time,attr" json:"time"`
	Failures       []Failure       `xml:"failure,omitempty" json:"failures,omitempty"`
	RerunFailures  []RerunFailure  `xml:"rerunFailure,omitempty" json:"rerunFailures,omitempty"`
	FlakyFailures  []FlakyFailure  `xml:"flakyFailure,omitempty" json:"flakyFailures,omitempty"`
	Skipped        *Skipped        `xml:"skipped,omitempty" json:"skipped,omitempty"`
	Errors         []Error         `xml:"error,omitempty" json:"errors,omitempty"`
	RerunErrors    []RerunError    `xml:"rerunError,omitempty" json:"rerunErrors,omitempty"`
	FlakyErrors    []FlakyError    `xml:"flakyError,omitempty" json:"flakyErrors,omitempty"`
	SystemOut      string          `xml:"system-out,omitempty" json:"systemOut,omitempty"`
	SystemErr      string          `xml:"system-err,omitempty" json:"systemErr,omitempty"`
}

// Properties represents key-value properties
type Properties struct {
	Properties []Property `xml:"property" json:"property,omitempty"`
}

// Property represents a single key-value property
type Property struct {
	Name  string `xml:"name,attr" json:"name"`
	Value string `xml:"value,attr,omitempty" json:"value,omitempty"`
}

// Skipped represents a skipped test
type Skipped struct {
	Message string `xml:"message,attr,omitempty" json:"message,omitempty"`
	Data    string `xml:",chardata" json:"data,omitempty"`
}

// Failure represents a test failure
type Failure struct {
	Message string `xml:"message,attr,omitempty" json:"message,omitempty"`
	Type    string `xml:"type,attr,omitempty" json:"type,omitempty"`
	Data    string `xml:",chardata" json:"data,omitempty"`
}

// Error represents a test error
type Error struct {
	Message string `xml:"message,attr,omitempty" json:"message,omitempty"`
	Type    string `xml:"type,attr,omitempty" json:"type,omitempty"`
	Data    string `xml:",chardata" json:"data,omitempty"`
}

// RerunFailure represents a failure from a rerun attempt
type RerunFailure struct {
	Message    string `xml:"message,attr,omitempty" json:"message,omitempty"`
	Type       string `xml:"type,attr,omitempty" json:"type,omitempty"`
	StackTrace string `xml:"stackTrace,omitempty" json:"stackTrace,omitempty"`
	SystemOut  string `xml:"system-out,omitempty" json:"systemOut,omitempty"`
	SystemErr  string `xml:"system-err,omitempty" json:"systemErr,omitempty"`
}

// RerunError represents an error from a rerun attempt
type RerunError struct {
	Message    string `xml:"message,attr,omitempty" json:"message,omitempty"`
	Type       string `xml:"type,attr,omitempty" json:"type,omitempty"`
	StackTrace string `xml:"stackTrace,omitempty" json:"stackTrace,omitempty"`
	SystemOut  string `xml:"system-out,omitempty" json:"systemOut,omitempty"`
	SystemErr  string `xml:"system-err,omitempty" json:"systemErr,omitempty"`
}

// FlakyFailure represents a flaky test failure
type FlakyFailure struct {
	Message    string `xml:"message,attr,omitempty" json:"message,omitempty"`
	Type       string `xml:"type,attr,omitempty" json:"type,omitempty"`
	StackTrace string `xml:"stackTrace,omitempty" json:"stackTrace,omitempty"`
	SystemOut  string `xml:"system-out,omitempty" json:"systemOut,omitempty"`
	SystemErr  string `xml:"system-err,omitempty" json:"systemErr,omitempty"`
}

// FlakyError represents a flaky test error
type FlakyError struct {
	Message    string `xml:"message,attr,omitempty" json:"message,omitempty"`
	Type       string `xml:"type,attr,omitempty" json:"type,omitempty"`
	StackTrace string `xml:"stackTrace,omitempty" json:"stackTrace,omitempty"`
	SystemOut  string `xml:"system-out,omitempty" json:"systemOut,omitempty"`
	SystemErr  string `xml:"system-err,omitempty" json:"systemErr,omitempty"`
}
