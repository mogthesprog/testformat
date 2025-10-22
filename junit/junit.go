package junit

import "encoding/xml"

// TestSuites represents the root element containing multiple test suites
type TestSuites struct {
	XMLName    xml.Name    `xml:"testsuites" json:"-"`
	Name       string      `xml:"name,attr,omitempty" json:"name,omitempty"`
	Tests      int         `xml:"tests,attr,omitempty" json:"tests,omitempty"`
	Failures   int         `xml:"failures,attr,omitempty" json:"failures,omitempty"`
	Errors     int         `xml:"errors,attr,omitempty" json:"errors,omitempty"`
	Skipped    int         `xml:"skipped,attr,omitempty" json:"skipped,omitempty"`
	Assertions int         `xml:"assertions,attr,omitempty" json:"assertions,omitempty"`
	Time       string      `xml:"time,attr,omitempty" json:"time,omitempty"`
	Timestamp  string      `xml:"timestamp,attr,omitempty" json:"timestamp,omitempty"`
	Suites     []TestSuite `xml:"testsuite" json:"testsuites"`
}

// TestSuite represents a collection of test cases
type TestSuite struct {
	XMLName    xml.Name     `xml:"testsuite" json:"-"`
	Name       string       `xml:"name,attr" json:"name"`
	Tests      int          `xml:"tests,attr,omitempty" json:"tests,omitempty"`
	Failures   int          `xml:"failures,attr,omitempty" json:"failures,omitempty"`
	Errors     int          `xml:"errors,attr,omitempty" json:"errors,omitempty"`
	Skipped    int          `xml:"skipped,attr,omitempty" json:"skipped,omitempty"`
	Assertions int          `xml:"assertions,attr,omitempty" json:"assertions,omitempty"`
	Time       string       `xml:"time,attr,omitempty" json:"time,omitempty"`
	Timestamp  string       `xml:"timestamp,attr,omitempty" json:"timestamp,omitempty"`
	Hostname   string       `xml:"hostname,attr,omitempty" json:"hostname,omitempty"`
	ID         string       `xml:"id,attr,omitempty" json:"id,omitempty"`
	Package    string       `xml:"package,attr,omitempty" json:"package,omitempty"`
	File       string       `xml:"file,attr,omitempty" json:"file,omitempty"`
	Properties *Properties  `xml:"properties,omitempty" json:"properties,omitempty"`
	Suites     []TestSuite  `xml:"testsuite,omitempty" json:"suites,omitempty"`
	TestCases  []TestCase   `xml:"testcase" json:"testcases,omitempty"`
	SystemOut  string       `xml:"system-out,omitempty" json:"system_out,omitempty"`
	SystemErr  string       `xml:"system-err,omitempty" json:"system_err,omitempty"`
}

// TestCase represents an individual test
type TestCase struct {
	XMLName    xml.Name    `xml:"testcase" json:"-"`
	Name       string      `xml:"name,attr" json:"name"`
	Classname  string      `xml:"classname,attr,omitempty" json:"classname,omitempty"`
	Time       string      `xml:"time,attr,omitempty" json:"time,omitempty"`
	Assertions int         `xml:"assertions,attr,omitempty" json:"assertions,omitempty"`
	File       string      `xml:"file,attr,omitempty" json:"file,omitempty"`
	Line       int         `xml:"line,attr,omitempty" json:"line,omitempty"`
	Properties *Properties `xml:"properties,omitempty" json:"properties,omitempty"`
	Skipped    *Skipped    `xml:"skipped,omitempty" json:"skipped,omitempty"`
	Failure    *Failure    `xml:"failure,omitempty" json:"failure,omitempty"`
	Error      *Error      `xml:"error,omitempty" json:"error,omitempty"`
	SystemOut  string      `xml:"system-out,omitempty" json:"system_out,omitempty"`
	SystemErr  string      `xml:"system-err,omitempty" json:"system_err,omitempty"`
}

// Properties represents key-value properties
type Properties struct {
	Properties []Property `xml:"property" json:"property,omitempty"`
}

// Property represents a single key-value property
type Property struct {
	Name  string `xml:"name,attr" json:"name"`
	Value string `xml:"value,attr,omitempty" json:"value,omitempty"`
	Data  string `xml:",chardata" json:"data,omitempty"`
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
