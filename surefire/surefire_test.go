package surefire

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRoundTripBasicSuccess(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite version="3.0.2" name="com.example.BasicTest" time="0.245" tests="3" errors="0" skipped="0" failures="0">
	<testcase name="testSuccess1" classname="com.example.BasicTest" time="0.082"/>
	<testcase name="testSuccess2" classname="com.example.BasicTest" time="0.081"/>
	<testcase name="testSuccess3" classname="com.example.BasicTest" time="0.082"/>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripWithFailures(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite version="3.0.2" name="com.example.FailureTest" time="0.200" tests="2" errors="0" skipped="0" failures="1">
	<testcase name="testSuccess" classname="com.example.FailureTest" time="0.100"/>
	<testcase name="testFailure" classname="com.example.FailureTest" time="0.100">
		<failure message="Expected 5 but got 3" type="org.junit.ComparisonFailure">
org.junit.ComparisonFailure: Expected 5 but got 3
	at com.example.FailureTest.testFailure(FailureTest.java:42)
		</failure>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripWithErrors(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite version="3.0.2" name="com.example.ErrorTest" time="0.100" tests="1" errors="1" skipped="0" failures="0">
	<testcase name="testError" classname="com.example.ErrorTest" time="0.100">
		<error message="NullPointerException" type="java.lang.NullPointerException">
java.lang.NullPointerException: Cannot invoke method on null object
	at com.example.ErrorTest.testError(ErrorTest.java:15)
		</error>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripWithSkipped(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite version="3.0.2" name="com.example.SkippedTest" time="0.050" tests="2" errors="0" skipped="1" failures="0">
	<testcase name="testSuccess" classname="com.example.SkippedTest" time="0.025"/>
	<testcase name="testSkipped" classname="com.example.SkippedTest" time="0.025">
		<skipped message="Test disabled"/>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripWithProperties(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite version="3.0.2" name="com.example.PropertiesTest" time="0.100" tests="1" errors="0" skipped="0" failures="0">
	<properties>
		<property name="java.version" value="11.0.1"/>
		<property name="os.name" value="Linux"/>
	</properties>
	<testcase name="testSuccess" classname="com.example.PropertiesTest" time="0.100"/>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripWithFlakyFailure(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite version="3.0.2" name="com.example.FlakyTest" time="0.200" tests="1" errors="0" skipped="0" failures="0">
	<testcase name="testFlakyFailure" classname="com.example.FlakyTest" time="0.200">
		<flakyFailure message="Intermittent timeout" type="org.junit.AssertionError">
			<stackTrace>org.junit.AssertionError: Intermittent timeout
	at com.example.FlakyTest.testFlakyFailure(FlakyTest.java:25)</stackTrace>
			<system-out>First attempt stdout</system-out>
			<system-err>First attempt stderr</system-err>
		</flakyFailure>
		<system-out>Final successful attempt stdout</system-out>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripWithFlakyError(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite version="3.0.2" name="com.example.FlakyTest" time="0.150" tests="1" errors="0" skipped="0" failures="0">
	<testcase name="testFlakyError" classname="com.example.FlakyTest" time="0.150">
		<flakyError message="Network connection reset" type="java.net.SocketException">
			<stackTrace>java.net.SocketException: Connection reset
	at com.example.FlakyTest.testFlakyError(FlakyTest.java:35)</stackTrace>
		</flakyError>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripWithRerunFailure(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite version="3.0.2" name="com.example.RerunTest" time="0.350" tests="1" errors="0" skipped="0" failures="0">
	<testcase name="testWithRerunFailure" classname="com.example.RerunTest" time="0.350">
		<rerunFailure message="First run failed" type="org.junit.AssertionError">
			<stackTrace>org.junit.AssertionError: First run failed
	at com.example.RerunTest.testWithRerunFailure(RerunTest.java:20)</stackTrace>
			<system-out>First run output</system-out>
			<system-err>First run error</system-err>
		</rerunFailure>
		<system-out>Final successful run output</system-out>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripWithRerunError(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite version="3.0.2" name="com.example.RerunTest" time="0.250" tests="1" errors="0" skipped="0" failures="0">
	<testcase name="testWithRerunError" classname="com.example.RerunTest" time="0.250">
		<rerunError message="Transient error" type="java.io.IOException">
			<stackTrace>java.io.IOException: Transient error
	at com.example.RerunTest.testWithRerunError(RerunTest.java:30)</stackTrace>
		</rerunError>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripWithGroup(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite version="3.0.2" name="com.example.GroupedTest" time="0.200" tests="2" errors="0" skipped="0" failures="0" group="integration">
	<testcase name="testDatabase" classname="com.example.GroupedTest" group="database" time="0.100"/>
	<testcase name="testAPI" classname="com.example.GroupedTest" group="api" time="0.100"/>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripWithSystemOut(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite version="3.0.2" name="com.example.OutputTest" time="0.100" tests="1" errors="0" skipped="0" failures="0">
	<testcase name="testWithOutput" classname="com.example.OutputTest" time="0.100">
		<system-out>Test output line 1
Test output line 2</system-out>
		<system-err>Error output</system-err>
	</testcase>
	<system-out>Suite level output</system-out>
	<system-err>Suite level error</system-err>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripMultipleFailures(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite version="3.0.2" name="com.example.MultiFailureTest" time="0.100" tests="1" errors="0" skipped="0" failures="2">
	<testcase name="testMultipleFailures" classname="com.example.MultiFailureTest" time="0.100">
		<failure message="First assertion failed" type="org.junit.AssertionError">First failure details</failure>
		<failure message="Second assertion failed" type="org.junit.AssertionError">Second failure details</failure>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripMultipleErrors(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite version="3.0.2" name="com.example.MultiErrorTest" time="0.100" tests="1" errors="2" skipped="0" failures="0">
	<testcase name="testMultipleErrors" classname="com.example.MultiErrorTest" time="0.100">
		<error message="First error" type="java.lang.RuntimeException">First error details</error>
		<error message="Second error" type="java.lang.IllegalStateException">Second error details</error>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripFromFileBasicSuccess(t *testing.T) {
	testRoundTripFromFile(t, "testdata/basic-success.xml")
}

func TestRoundTripFromFileFailuresAndErrors(t *testing.T) {
	testRoundTripFromFile(t, "testdata/failures-and-errors.xml")
}

func TestRoundTripFromFileSkippedTests(t *testing.T) {
	testRoundTripFromFile(t, "testdata/skipped-tests.xml")
}

func TestRoundTripFromFileFlakyTests(t *testing.T) {
	testRoundTripFromFile(t, "testdata/flaky-tests.xml")
}

func TestRoundTripFromFileRerunTests(t *testing.T) {
	testRoundTripFromFile(t, "testdata/rerun-tests.xml")
}

func TestRoundTripFromFileGroupedTests(t *testing.T) {
	testRoundTripFromFile(t, "testdata/grouped-tests.xml")
}

func TestRoundTripFromFileUnicode(t *testing.T) {
	testRoundTripFromFile(t, "testdata/edge-case-unicode.xml")
}

func TestRoundTripFromFileAntUnit(t *testing.T) {
	testRoundTripFromFile(t, "testdata/TEST-AntUnit.xml")
}

func TestRoundTripFromFileNoPackageTest(t *testing.T) {
	testRoundTripFromFile(t, "testdata/TEST-NoPackageTest.xml")
}

func TestRoundTripFromFileNoTimeTestCaseTest(t *testing.T) {
	testRoundTripFromFile(t, "testdata/TEST-NoTimeTestCaseTest.xml")
}

func TestRoundTripFromFileClassWithNoTests(t *testing.T) {
	testRoundTripFromFile(t, "testdata/TEST-classWithNoTests.NoMethodsTestCase.xml")
}

func TestRoundTripFromFileCircleTest(t *testing.T) {
	testRoundTripFromFile(t, "testdata/TEST-com.shape.CircleTest.xml")
}

func TestRoundTripFromFilePointTest(t *testing.T) {
	testRoundTripFromFile(t, "testdata/TEST-com.shape.PointTest.xml")
}

func TestRoundTripFromFileWrapperTestSuite(t *testing.T) {
	testRoundTripFromFile(t, "testdata/TEST-junit.twoTestCaseSuite.WrapperTestSuite.xml")
}

func TestRoundTripWithFailureAndError(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite version="3.0.2" name="com.example.MixedTest" time="0.100" tests="1" errors="1" skipped="0" failures="1">
	<testcase name="testWithBothFailureAndError" classname="com.example.MixedTest" time="0.100">
		<failure message="Assertion failed" type="org.junit.AssertionError">
org.junit.AssertionError: Expected true but was false
	at com.example.MixedTest.testWithBothFailureAndError(MixedTest.java:20)
		</failure>
		<error message="Unexpected exception" type="java.lang.RuntimeException">
java.lang.RuntimeException: Something went wrong during cleanup
	at com.example.MixedTest.testWithBothFailureAndError(MixedTest.java:25)
		</error>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripEmptyTestSuite(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite version="3.0.2" name="com.example.EmptyTest" time="0.000" tests="0" errors="0" skipped="0" failures="0">
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripSkippedWithMessageAndData(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite version="3.0.2" name="com.example.SkippedDetailTest" time="0.050" tests="1" errors="0" skipped="1" failures="0">
	<testcase name="testSkippedWithDetails" classname="com.example.SkippedDetailTest" time="0.050">
		<skipped message="Test disabled for maintenance">Additional details about why this test was skipped.
It requires external service that is currently unavailable.</skipped>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripWithCDATA(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite version="3.0.2" name="com.example.CDATATest" time="0.150" tests="1" errors="0" skipped="0" failures="1">
	<testcase name="testWithCDATA" classname="com.example.CDATATest" time="0.150">
		<failure message="Assertion failed" type="org.junit.AssertionError"><![CDATA[org.junit.AssertionError: Expected <html><body>test</body></html> but got <html><body>wrong</body></html>
	at com.example.CDATATest.testWithCDATA(CDATATest.java:42)
	at java.base/java.lang.reflect.Method.invoke(Method.java:566)
Special characters: & < > " ']]></failure>
		<system-out><![CDATA[Output with special chars: <tag> & "quotes"]]></system-out>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func testRoundTrip(t *testing.T, xmlData string) {
	t.Helper()

	var original TestSuite
	if err := xml.Unmarshal([]byte(xmlData), &original); err != nil {
		t.Fatalf("Failed to unmarshal original XML: %v", err)
	}

	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}

	var fromJSON TestSuite
	if err := json.Unmarshal(jsonData, &fromJSON); err != nil {
		t.Fatalf("Failed to unmarshal from JSON: %v", err)
	}

	xmlRoundTrip, err := xml.Marshal(fromJSON)
	if err != nil {
		t.Fatalf("Failed to marshal back to XML: %v", err)
	}

	var finalSuite TestSuite
	if err := xml.Unmarshal(xmlRoundTrip, &finalSuite); err != nil {
		t.Fatalf("Failed to unmarshal final XML: %v", err)
	}

	if diff := cmp.Diff(original, finalSuite); diff != "" {
		t.Errorf("Round trip conversion resulted in differences (-original +final):\n%s", diff)
		t.Logf("Original XML:\n%s", xmlData)
		t.Logf("JSON intermediate:\n%s", string(jsonData))
		t.Logf("Final XML:\n%s", string(xmlRoundTrip))
	}
}

func testRoundTripFromFile(t *testing.T, filename string) {
	t.Helper()

	xmlData, err := os.ReadFile(filepath.Join(filename))
	if err != nil {
		t.Fatalf("Failed to read test file %s: %v", filename, err)
	}

	testRoundTrip(t, string(xmlData))
}

func TestErrorHandling_MalformedXML(t *testing.T) {
	tests := []struct {
		name    string
		xmlData string
	}{
		{
			name:    "unclosed tag",
			xmlData: `<testsuite name="test"><testcase name="test1"</testsuite>`,
		},
		{
			name:    "mismatched tags",
			xmlData: `<testsuite name="test"><testcase name="test1"></testsuite></testcase>`,
		},
		{
			name:    "no root element",
			xmlData: ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var suite TestSuite
			err := xml.Unmarshal([]byte(tt.xmlData), &suite)
			if err == nil {
				t.Errorf("Expected error for malformed XML, but got nil")
			}
		})
	}
}

func TestErrorHandling_MissingRequiredAttributes(t *testing.T) {
	tests := []struct {
		name    string
		xmlData string
		wantErr bool
	}{
		{
			name:    "missing testsuite name",
			xmlData: `<testsuite tests="0" errors="0" skipped="0" failures="0"/>`,
			wantErr: false,
		},
		{
			name:    "missing tests count",
			xmlData: `<testsuite name="test" errors="0" skipped="0" failures="0"/>`,
			wantErr: false,
		},
		{
			name:    "missing testcase name",
			xmlData: `<testsuite name="test" tests="1" errors="0" skipped="0" failures="0"><testcase time="0.1"/></testsuite>`,
			wantErr: false,
		},
		{
			name:    "empty name attribute",
			xmlData: `<testsuite name="" tests="0" errors="0" skipped="0" failures="0"/>`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var suite TestSuite
			err := xml.Unmarshal([]byte(tt.xmlData), &suite)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error=%v, but got error=%v", tt.wantErr, err)
			}
			if err == nil {
				jsonData, err := json.Marshal(suite)
				if err != nil {
					t.Errorf("Failed to marshal to JSON after parsing XML with missing attributes: %v", err)
				}
				var fromJSON TestSuite
				if err := json.Unmarshal(jsonData, &fromJSON); err != nil {
					t.Errorf("Failed to unmarshal from JSON: %v", err)
				}
			}
		})
	}
}

func TestJSONFormat_StructureValidation(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite version="3.0.2" name="JSONTest" time="0.5" tests="2" errors="0" skipped="0" failures="1">
	<testcase name="test1" classname="Test" time="0.3"/>
	<testcase name="test2" classname="Test" time="0.2">
		<failure message="Failed" type="Error">Stack trace</failure>
	</testcase>
</testsuite>`

	var suite TestSuite
	if err := xml.Unmarshal([]byte(xmlData), &suite); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	jsonData, err := json.Marshal(suite)
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal JSON to map: %v", err)
	}

	if jsonMap["name"] != "JSONTest" {
		t.Errorf("Expected name=JSONTest, got %v", jsonMap["name"])
	}

	if jsonMap["tests"] != "2" {
		t.Errorf("Expected tests=2, got %v", jsonMap["tests"])
	}

	testCases, ok := jsonMap["testCases"].([]interface{})
	if !ok {
		t.Fatalf("Expected testCases to be an array")
	}

	if len(testCases) != 2 {
		t.Errorf("Expected 2 testCases, got %d", len(testCases))
	}
}

func FuzzTestSuiteXML(f *testing.F) {
	seeds := []string{
		`<testsuite name="test" tests="0" errors="0" skipped="0" failures="0"/>`,
		`<testsuite name="test" tests="1" errors="0" skipped="0" failures="0"><testcase name="t1" time="0.1"/></testsuite>`,
		`<testsuite name="test" tests="1" errors="0" skipped="0" failures="1"><testcase name="t1" time="0.1"><failure message="failed"/></testcase></testsuite>`,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, xmlData string) {
		var suite TestSuite
		err := xml.Unmarshal([]byte(xmlData), &suite)
		if err != nil {
			return
		}

		jsonData, err := json.Marshal(suite)
		if err != nil {
			t.Errorf("Failed to marshal valid XML to JSON: %v", err)
			return
		}

		var fromJSON TestSuite
		err = json.Unmarshal(jsonData, &fromJSON)
		if err != nil {
			t.Errorf("Failed to unmarshal JSON back to struct: %v", err)
		}
	})
}
