package junit

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRoundTripSimpleTestSuite(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="ExampleTests" tests="3" failures="1" errors="0" skipped="1" time="0.245">
	<testcase name="test_success" classname="com.example.TestClass" time="0.123"/>
	<testcase name="test_failure" classname="com.example.TestClass" time="0.100">
		<failure message="Expected 5 but got 3" type="AssertionError">
AssertionError: Expected 5 but got 3
	at TestClass.test_failure(TestClass.java:42)
		</failure>
	</testcase>
	<testcase name="test_skipped" classname="com.example.TestClass" time="0.001">
		<skipped message="Not yet implemented"/>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripTestSuiteWithProperties(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="PropertiesTest" tests="1" failures="0" errors="0" time="0.050" timestamp="2024-01-15T10:30:00" hostname="test-server">
	<properties>
		<property name="java.version" value="11.0.1"/>
		<property name="os.name" value="Linux"/>
	</properties>
	<testcase name="test_example" classname="com.example.Test" time="0.050"/>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripTestSuiteWithSystemOut(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="OutputTest" tests="1" failures="0" errors="0" time="0.100">
	<testcase name="test_with_output" classname="com.example.Test" time="0.100">
		<system-out>This is standard output
Multiple lines of output
</system-out>
		<system-err>This is error output</system-err>
	</testcase>
	<system-out>Suite level output</system-out>
	<system-err>Suite level error</system-err>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripTestSuiteWithError(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="ErrorTest" tests="1" failures="0" errors="1" time="0.010">
	<testcase name="test_with_error" classname="com.example.Test" time="0.010">
		<error message="NullPointerException" type="java.lang.NullPointerException">
java.lang.NullPointerException: Cannot invoke method on null object
	at com.example.Test.test_with_error(Test.java:15)
		</error>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripMultipleTestSuites(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuites name="AllTests" tests="4" failures="1" errors="1" skipped="0" time="0.500">
	<testsuite name="Suite1" tests="2" failures="1" errors="0" time="0.200" package="com.example.suite1">
		<testcase name="test1" classname="com.example.Suite1Test" time="0.100"/>
		<testcase name="test2" classname="com.example.Suite1Test" time="0.100">
			<failure message="Failed" type="AssertionError">Assertion failed</failure>
		</testcase>
	</testsuite>
	<testsuite name="Suite2" tests="2" failures="0" errors="1" time="0.300" package="com.example.suite2">
		<testcase name="test3" classname="com.example.Suite2Test" time="0.200"/>
		<testcase name="test4" classname="com.example.Suite2Test" time="0.100">
			<error message="Error occurred" type="RuntimeError">Runtime error details</error>
		</testcase>
	</testsuite>
</testsuites>`

	testRoundTripTestSuites(t, xmlData)
}

func TestRoundTripEmptyTestSuite(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="EmptyTest" tests="0" failures="0" errors="0" time="0.000"/>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripComplexFailureMessage(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="ComplexTest" tests="1" failures="1" errors="0" time="0.100">
	<testcase name="test_complex" classname="com.example.Test" time="0.100">
		<failure message="Complex failure with special chars: &lt;&gt;&amp;&quot;" type="AssertionError">
<![CDATA[
Expected: <value>
     but: was <other>

Stack trace:
  at Test.java:10
  at Runner.java:20
]]>
		</failure>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripSkippedWithData(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="SkippedTest" tests="1" failures="0" errors="0" skipped="1" time="0.001">
	<testcase name="test_skipped" classname="com.example.Test" time="0.001">
		<skipped message="Skipping test">Additional skip information</skipped>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripAllAttributes(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="FullAttributeTest" tests="1" failures="0" errors="0" skipped="0" time="0.100" timestamp="2024-01-15T10:30:00" hostname="test-host" id="suite-123" package="com.example.full">
	<properties>
		<property name="key1" value="value1"/>
		<property name="key2" value="value2"/>
	</properties>
	<testcase name="test_all_attrs" classname="com.example.FullTest" time="0.100"/>
	<system-out>Suite output</system-out>
	<system-err>Suite errors</system-err>
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

func testRoundTripTestSuites(t *testing.T, xmlData string) {
	t.Helper()

	var original TestSuites
	if err := xml.Unmarshal([]byte(xmlData), &original); err != nil {
		t.Fatalf("Failed to unmarshal original XML: %v", err)
	}

	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}

	var fromJSON TestSuites
	if err := json.Unmarshal(jsonData, &fromJSON); err != nil {
		t.Fatalf("Failed to unmarshal from JSON: %v", err)
	}

	xmlRoundTrip, err := xml.Marshal(fromJSON)
	if err != nil {
		t.Fatalf("Failed to marshal back to XML: %v", err)
	}

	var finalSuites TestSuites
	if err := xml.Unmarshal(xmlRoundTrip, &finalSuites); err != nil {
		t.Fatalf("Failed to unmarshal final XML: %v", err)
	}

	if diff := cmp.Diff(original, finalSuites); diff != "" {
		t.Errorf("Round trip conversion resulted in differences (-original +final):\n%s", diff)
		t.Logf("Original XML:\n%s", xmlData)
		t.Logf("JSON intermediate:\n%s", string(jsonData))
		t.Logf("Final XML:\n%s", string(xmlRoundTrip))
	}
}

func TestRoundTripConventions(t *testing.T) {
	testRoundTripFromFile(t, "testdata/conventions.xml", true)
}

func TestRoundTripJUnitBasic(t *testing.T) {
	testRoundTripFromFile(t, "testdata/junit-basic.xml", true)
}

func TestRoundTripJUnitComplete(t *testing.T) {
	testRoundTripFromFile(t, "testdata/junit-complete.xml", true)
}

func TestRoundTripTestCaseOutput(t *testing.T) {
	testRoundTripFromFile(t, "testdata/testcase-output.xml", true)
}

func TestRoundTripTestCaseProperties(t *testing.T) {
	testRoundTripFromFile(t, "testdata/testcase-properties.xml", true)
}

func testRoundTripFromFile(t *testing.T, filename string, isTestSuites bool) {
	t.Helper()

	xmlData, err := os.ReadFile(filepath.Join(filename))
	if err != nil {
		t.Fatalf("Failed to read test file %s: %v", filename, err)
	}

	if isTestSuites {
		testRoundTripTestSuites(t, string(xmlData))
	} else {
		testRoundTrip(t, string(xmlData))
	}
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
			name:    "invalid xml declaration",
			xmlData: `<?xml version="2.0" encoding="UTF-8"?><testsuite name="test"/>`,
		},
		{
			name:    "no root element",
			xmlData: ``,
		},
		{
			name:    "invalid attribute syntax",
			xmlData: `<testsuite name=test"><testcase name="test1"/></testsuite>`,
		},
		{
			name:    "unclosed CDATA",
			xmlData: `<testsuite name="test"><testcase name="t1"><failure><![CDATA[error</failure></testcase></testsuite>`,
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

func TestErrorHandling_InvalidUTF8(t *testing.T) {
	invalidUTF8 := []byte{0xff, 0xfe, 0xfd}
	xmlData := append([]byte(`<testsuite name="`), invalidUTF8...)
	xmlData = append(xmlData, []byte(`"><testcase name="test1"/></testsuite>`)...)

	var suite TestSuite
	err := xml.Unmarshal(xmlData, &suite)
	if err == nil {
		t.Errorf("Expected error for invalid UTF-8, but got nil")
	}
}

func TestErrorHandling_InvalidJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
	}{
		{
			name:     "invalid json syntax",
			jsonData: `{"name": "test", "testCases": [}`,
		},
		{
			name:     "wrong type for integer field",
			jsonData: `{"name": "test", "tests": "not a number"}`,
		},
		{
			name:     "unclosed string",
			jsonData: `{"name": "test}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var suite TestSuite
			err := json.Unmarshal([]byte(tt.jsonData), &suite)
			if err == nil {
				t.Errorf("Expected error for invalid JSON, but got nil")
			}
		})
	}
}

func TestEdgeCases_MultipleStatuses(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="MultiStatusTest" tests="1" failures="1" errors="1">
	<testcase name="test_both_failure_and_error" classname="com.example.Test" time="0.100">
		<failure message="Assertion failed" type="AssertionError">Failed assertion</failure>
		<error message="Exception occurred" type="RuntimeException">Runtime error</error>
	</testcase>
</testsuite>`

	var suite TestSuite
	if err := xml.Unmarshal([]byte(xmlData), &suite); err != nil {
		t.Fatalf("Failed to unmarshal XML with both failure and error: %v", err)
	}

	if suite.TestCases[0].Failure == nil {
		t.Error("Expected failure to be present")
	}
	if suite.TestCases[0].Error == nil {
		t.Error("Expected error to be present")
	}

	testRoundTrip(t, xmlData)
}

func TestEdgeCases_FailureAndSkipped(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="FailureSkippedTest" tests="1" failures="1" skipped="1">
	<testcase name="test_failure_and_skipped" classname="com.example.Test" time="0.100">
		<failure message="Failed" type="AssertionError">Failed</failure>
		<skipped message="Skipped">Also skipped</skipped>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestEdgeCases_EmptyStatusElements(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="EmptyStatusTest" tests="3" failures="1" errors="1" skipped="1">
	<testcase name="test_empty_failure" classname="com.example.Test" time="0.100">
		<failure/>
	</testcase>
	<testcase name="test_empty_error" classname="com.example.Test" time="0.100">
		<error/>
	</testcase>
	<testcase name="test_empty_skipped" classname="com.example.Test" time="0.100">
		<skipped/>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestEdgeCases_WhitespaceHandling(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="WhitespaceTest" tests="1" failures="0" errors="0">
	<properties>
		<property name="empty" value=""/>
		<property name="spaces" value="   "/>
		<property name="whitespace-only">   </property>
	</properties>
	<testcase name="test1" classname="Test" time="0.1">
		<system-out>   </system-out>
		<system-err></system-err>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestEdgeCases_NegativeAndZeroNumbers(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="NumberTest" tests="0" failures="0" errors="0" time="0">
	<testcase name="test_zero_time" classname="Test" time="0" line="0"/>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestEdgeCases_VeryLargeNumbers(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="LargeNumberTest" tests="999999" failures="999999" errors="999999" time="999999.999999">
	<testcase name="test1" classname="Test" time="999999.999999" assertions="999999" line="999999"/>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestTypeValidation_FloatVsIntegerTime(t *testing.T) {
	tests := []struct {
		name     string
		xmlData  string
		expected string
	}{
		{
			name:     "integer time",
			xmlData:  `<testsuite name="test" time="5"><testcase name="t1" time="3"/></testsuite>`,
			expected: "5",
		},
		{
			name:     "float time",
			xmlData:  `<testsuite name="test" time="5.123"><testcase name="t1" time="3.456"/></testsuite>`,
			expected: "5.123",
		},
		{
			name:     "scientific notation",
			xmlData:  `<testsuite name="test" time="1.5e-3"><testcase name="t1" time="2.5e-4"/></testsuite>`,
			expected: "1.5e-3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var suite TestSuite
			if err := xml.Unmarshal([]byte(tt.xmlData), &suite); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			if suite.Time != tt.expected {
				t.Errorf("Expected time=%q, got %q", tt.expected, suite.Time)
			}
		})
	}
}

func TestSpecialCharacters_XMLEntities(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="Entity&lt;Test&gt;" tests="1">
	<testcase name="test&amp;name" classname="Test&quot;Class&apos;" time="0.1">
		<failure message="Expected &lt;5&gt; but got &lt;3&gt;" type="Error">
Error with &amp; and &lt; and &gt; and &quot; and &apos;
		</failure>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestSpecialCharacters_UnicodeAndEmoji(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="Unicode测试🎉" tests="1">
	<testcase name="test_emoji_🚀" classname="日本語クラス" time="0.1">
		<failure message="Failed with emoji 😢" type="AssertionError">
Stack trace with unicode: 你好世界
And emoji: ✅ ❌ ⚠️
		</failure>
	</testcase>
	<system-out>Output: Καλημέρα κόσμε 🌍</system-out>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestSpecialCharacters_NewlinesAndTabs(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="WhitespaceTest" tests="1">
	<testcase name="test&#10;with&#10;newlines" classname="Test	With	Tabs" time="0.1">
		<failure message="Multi
line
message" type="Error">
Error with
	tabs and
		multiple lines
		</failure>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestSpecialCharacters_ControlCharacters(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="ControlCharTest" tests="1">
	<testcase name="test_control" classname="Test" time="0.1">
		<failure message="Error with control chars" type="Error">
			<![CDATA[Line 1Line 2]]>
		</failure>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestStructuralVariations_EmptyTestSuite(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="EmptyTests" tests="0" failures="0" errors="0"/>`

	var suite TestSuite
	if err := xml.Unmarshal([]byte(xmlData), &suite); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if suite.TestCases != nil && len(suite.TestCases) > 0 {
		t.Errorf("Expected nil or empty testcases, got %v", suite.TestCases)
	}

	testRoundTrip(t, xmlData)
}

func TestStructuralVariations_EmptyProperties(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="Test" tests="1">
	<properties/>
	<testcase name="test1" classname="Test">
		<properties/>
	</testcase>
</testsuite>`

	testRoundTrip(t, xmlData)
}

func TestStructuralVariations_DuplicatePropertyNames(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="Test" tests="1">
	<properties>
		<property name="key" value="value1"/>
		<property name="key" value="value2"/>
		<property name="key" value="value3"/>
	</properties>
	<testcase name="test1" classname="Test"/>
</testsuite>`

	var suite TestSuite
	if err := xml.Unmarshal([]byte(xmlData), &suite); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if suite.Properties == nil || len(suite.Properties.Properties) != 3 {
		t.Errorf("Expected 3 properties, got %v", suite.Properties)
	}

	testRoundTrip(t, xmlData)
}

func TestStructuralVariations_DeeplyNestedSuites(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuites name="Root">
	<testsuite name="Level1">
		<testsuite name="Level2">
			<testsuite name="Level3">
				<testsuite name="Level4">
					<testsuite name="Level5">
						<testcase name="deep_test" classname="Test"/>
					</testsuite>
				</testsuite>
			</testsuite>
		</testsuite>
	</testsuite>
</testsuites>`

	testRoundTripTestSuites(t, xmlData)
}

func TestStructuralVariations_MixedTestCasesAndSuites(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="Mixed" tests="3">
	<testcase name="test1" classname="Test"/>
	<testsuite name="Nested">
		<testcase name="test2" classname="Test"/>
	</testsuite>
	<testcase name="test3" classname="Test"/>
</testsuite>`

	var suite TestSuite
	if err := xml.Unmarshal([]byte(xmlData), &suite); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(suite.TestCases) != 2 {
		t.Errorf("Expected 2 direct testcases, got %d", len(suite.TestCases))
	}

	if len(suite.Suites) != 1 {
		t.Errorf("Expected 1 nested suite, got %d", len(suite.Suites))
	}

	testRoundTrip(t, xmlData)
}

func TestJSONFormat_StructureValidation(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="JSONTest" tests="2" failures="1" time="0.5">
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

	if jsonMap["tests"] != float64(2) {
		t.Errorf("Expected tests=2, got %v", jsonMap["tests"])
	}

	testCases, ok := jsonMap["testCases"].([]interface{})
	if !ok {
		t.Fatalf("Expected testCases to be an array")
	}

	if len(testCases) != 2 {
		t.Errorf("Expected 2 testCases, got %d", len(testCases))
	}

	tc2 := testCases[1].(map[string]interface{})
	if _, hasFailure := tc2["failure"]; !hasFailure {
		t.Error("Expected second testcase to have failure")
	}
}

func TestJSONFormat_OmitEmpty(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="OmitTest">
	<testcase name="test1"/>
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

	if _, hasTests := jsonMap["tests"]; hasTests {
		t.Error("Expected tests field to be omitted when zero")
	}

	if _, hasFailures := jsonMap["failures"]; hasFailures {
		t.Error("Expected failures field to be omitted when zero")
	}

	testCases := jsonMap["testCases"].([]interface{})
	tc := testCases[0].(map[string]interface{})

	if _, hasClassName := tc["className"]; hasClassName {
		t.Error("Expected className to be omitted when empty")
	}

	if _, hasTime := tc["time"]; hasTime {
		t.Error("Expected time to be omitted when empty")
	}
}

func TestJSONFormat_SystemOutErrFieldNames(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="OutputTest">
	<testcase name="test1">
		<system-out>stdout content</system-out>
		<system-err>stderr content</system-err>
	</testcase>
	<system-out>suite stdout</system-out>
	<system-err>suite stderr</system-err>
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

	if jsonMap["systemOut"] != "suite stdout" {
		t.Errorf("Expected systemOut field in JSON, got %v", jsonMap)
	}

	if jsonMap["systemErr"] != "suite stderr" {
		t.Errorf("Expected systemErr field in JSON, got %v", jsonMap)
	}

	testCases := jsonMap["testCases"].([]interface{})
	tc := testCases[0].(map[string]interface{})

	if tc["systemOut"] != "stdout content" {
		t.Errorf("Expected systemOut in testcase, got %v", tc)
	}

	if tc["systemErr"] != "stderr content" {
		t.Errorf("Expected systemErr in testcase, got %v", tc)
	}
}

func TestJSONFormat_PropertiesStructure(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="PropertiesTest">
	<properties>
		<property name="key1" value="value1"/>
		<property name="key2">text content</property>
	</properties>
	<testcase name="test1"/>
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

	props, ok := jsonMap["properties"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected properties to be an object, got %T", jsonMap["properties"])
	}

	propArray, ok := props["property"].([]interface{})
	if !ok {
		t.Fatalf("Expected property to be an array, got %T", props["property"])
	}

	if len(propArray) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(propArray))
	}

	prop1 := propArray[0].(map[string]interface{})
	if prop1["name"] != "key1" || prop1["value"] != "value1" {
		t.Errorf("Property 1 incorrect: %v", prop1)
	}

	prop2 := propArray[1].(map[string]interface{})
	if prop2["name"] != "key2" || prop2["data"] != "text content" {
		t.Errorf("Property 2 incorrect: %v", prop2)
	}
}

func TestJSONFormat_NestedSuitesArray(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="Parent">
	<testsuite name="Child1">
		<testcase name="test1"/>
	</testsuite>
	<testsuite name="Child2">
		<testcase name="test2"/>
	</testsuite>
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

	suites, ok := jsonMap["suites"].([]interface{})
	if !ok {
		t.Fatalf("Expected suites to be an array, got %T", jsonMap["suites"])
	}

	if len(suites) != 2 {
		t.Errorf("Expected 2 nested suites, got %d", len(suites))
	}
}

func TestJSONFormat_TestSuitesWrapper(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<testsuites name="AllTests" tests="2">
	<testsuite name="Suite1">
		<testcase name="test1"/>
	</testsuite>
	<testsuite name="Suite2">
		<testcase name="test2"/>
	</testsuite>
</testsuites>`

	var suites TestSuites
	if err := xml.Unmarshal([]byte(xmlData), &suites); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	jsonData, err := json.Marshal(suites)
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal JSON to map: %v", err)
	}

	if jsonMap["name"] != "AllTests" {
		t.Errorf("Expected name=AllTests, got %v", jsonMap["name"])
	}

	suitesArray, ok := jsonMap["testSuites"].([]interface{})
	if !ok {
		t.Fatalf("Expected testSuites to be an array, got %T", jsonMap["testSuites"])
	}

	if len(suitesArray) != 2 {
		t.Errorf("Expected 2 testSuites, got %d", len(suitesArray))
	}
}

func FuzzTestSuiteXML(f *testing.F) {
	seeds := []string{
		`<testsuite name="test"/>`,
		`<testsuite name="test"><testcase name="t1"/></testsuite>`,
		`<testsuite name="test"><testcase name="t1"><failure message="failed"/></testcase></testsuite>`,
		`<testsuite name="test"><properties><property name="k" value="v"/></properties><testcase name="t1"/></testsuite>`,
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

func FuzzTestSuitesXML(f *testing.F) {
	seeds := []string{
		`<testsuites name="test"/>`,
		`<testsuites><testsuite name="s1"><testcase name="t1"/></testsuite></testsuites>`,
		`<testsuites><testsuite name="s1"><testcase name="t1"><error message="err"/></testcase></testsuite></testsuites>`,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, xmlData string) {
		var suites TestSuites
		err := xml.Unmarshal([]byte(xmlData), &suites)
		if err != nil {
			return
		}

		jsonData, err := json.Marshal(suites)
		if err != nil {
			t.Errorf("Failed to marshal valid XML to JSON: %v", err)
			return
		}

		var fromJSON TestSuites
		err = json.Unmarshal(jsonData, &fromJSON)
		if err != nil {
			t.Errorf("Failed to unmarshal JSON back to struct: %v", err)
		}
	})
}

func FuzzJSON(f *testing.F) {
	seeds := []string{
		`{"name":"test"}`,
		`{"name":"test","testCases":[{"name":"t1"}]}`,
		`{"name":"test","testCases":[{"name":"t1","failure":{"message":"failed"}}]}`,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, jsonData string) {
		var suite TestSuite
		err := json.Unmarshal([]byte(jsonData), &suite)
		if err != nil {
			return
		}

		xmlData, err := xml.Marshal(suite)
		if err != nil {
			t.Errorf("Failed to marshal valid JSON to XML: %v", err)
			return
		}

		var fromXML TestSuite
		err = xml.Unmarshal(xmlData, &fromXML)
		if err != nil {
			t.Errorf("Failed to unmarshal XML back to struct: %v", err)
		}
	})
}
