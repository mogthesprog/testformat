package junit_test

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"

	"github.com/mogthesprog/testformat/junit"
)

func ExampleTestSuite_unmarshalXML() {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="ExampleTests" tests="2" failures="1" time="0.5">
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

	fmt.Printf("Suite: %s\n", suite.Name)
	fmt.Printf("Tests: %d\n", suite.Tests)
	fmt.Printf("Failures: %d\n", suite.Failures)
	fmt.Printf("First test: %s\n", suite.TestCases[0].Name)
}

func ExampleTestSuite_marshalJSON() {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="ExampleTests" tests="1" time="0.3">
	<testcase name="test_example" classname="com.example.Test" time="0.3"/>
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

func ExampleTestSuites_unmarshalXML() {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuites name="AllTests" tests="2" time="0.5">
	<testsuite name="Suite1" tests="1" time="0.3">
		<testcase name="test1" classname="Test" time="0.3"/>
	</testsuite>
	<testsuite name="Suite2" tests="1" time="0.2">
		<testcase name="test2" classname="Test" time="0.2"/>
	</testsuite>
</testsuites>`

	var suites junit.TestSuites
	if err := xml.Unmarshal([]byte(xmlData), &suites); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total tests: %d\n", suites.Tests)
	fmt.Printf("Number of suites: %d\n", len(suites.Suites))
	fmt.Printf("First suite: %s\n", suites.Suites[0].Name)
}

func ExampleTestCase_withFailure() {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="Tests">
	<testcase name="test_with_failure" classname="Test" time="0.1">
		<failure message="Assertion failed" type="AssertionError">
Expected: 5
Actual: 3
		</failure>
	</testcase>
</testsuite>`

	var suite junit.TestSuite
	if err := xml.Unmarshal([]byte(xmlData), &suite); err != nil {
		log.Fatal(err)
	}

	testCase := suite.TestCases[0]
	if testCase.Failure != nil {
		fmt.Printf("Test failed: %s\n", testCase.Failure.Message)
		fmt.Printf("Type: %s\n", testCase.Failure.Type)
	}
}

func ExampleTestCase_withProperties() {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="Tests">
	<testcase name="test_example" classname="Test">
		<properties>
			<property name="priority" value="high"/>
			<property name="author" value="Alice"/>
			<property name="description">
This test validates the login functionality
			</property>
		</properties>
	</testcase>
</testsuite>`

	var suite junit.TestSuite
	if err := xml.Unmarshal([]byte(xmlData), &suite); err != nil {
		log.Fatal(err)
	}

	if suite.TestCases[0].Properties != nil {
		for _, prop := range suite.TestCases[0].Properties.Properties {
			if prop.Value != "" {
				fmt.Printf("%s: %s\n", prop.Name, prop.Value)
			} else {
				fmt.Printf("%s: %s\n", prop.Name, prop.Data)
			}
		}
	}
}

func ExampleTestSuite_roundTrip() {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="Tests" tests="1">
	<testcase name="test1" classname="Test" time="0.5"/>
</testsuite>`

	var suite junit.TestSuite
	if err := xml.Unmarshal([]byte(xmlData), &suite); err != nil {
		log.Fatal(err)
	}

	jsonData, err := json.Marshal(suite)
	if err != nil {
		log.Fatal(err)
	}

	var fromJSON junit.TestSuite
	if err := json.Unmarshal(jsonData, &fromJSON); err != nil {
		log.Fatal(err)
	}

	xmlOutput, err := xml.MarshalIndent(fromJSON, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(xmlOutput))
}
