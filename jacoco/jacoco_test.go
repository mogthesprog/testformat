package jacoco

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRoundTripSimpleReport(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<report name="Test Coverage Report">
	<sessioninfo id="session-1" start="1234567890" dump="1234567900"/>
	<package name="com/example/myapp">
		<class name="com/example/myapp/Calculator" sourcefilename="Calculator.java">
			<method name="add" desc="(II)I" line="10">
				<counter type="INSTRUCTION" missed="0" covered="4"/>
				<counter type="BRANCH" missed="0" covered="0"/>
				<counter type="LINE" missed="0" covered="2"/>
				<counter type="METHOD" missed="0" covered="1"/>
			</method>
			<counter type="INSTRUCTION" missed="0" covered="15"/>
			<counter type="BRANCH" missed="2" covered="4"/>
			<counter type="LINE" missed="1" covered="5"/>
			<counter type="METHOD" missed="0" covered="3"/>
			<counter type="CLASS" missed="0" covered="1"/>
		</class>
		<counter type="INSTRUCTION" missed="10" covered="50"/>
		<counter type="BRANCH" missed="5" covered="15"/>
		<counter type="LINE" missed="2" covered="20"/>
		<counter type="METHOD" missed="1" covered="10"/>
		<counter type="CLASS" missed="0" covered="5"/>
	</package>
	<counter type="INSTRUCTION" missed="20" covered="100"/>
	<counter type="BRANCH" missed="10" covered="30"/>
	<counter type="LINE" missed="5" covered="50"/>
	<counter type="METHOD" missed="2" covered="25"/>
	<counter type="CLASS" missed="1" covered="10"/>
</report>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripReportWithSourceFile(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<report name="Coverage Report">
	<sessioninfo id="test-session" start="1234567890" dump="1234567900"/>
	<package name="com/example">
		<sourcefile name="Example.java">
			<line nr="1" mi="0" ci="2" mb="0" cb="0"/>
			<line nr="5" mi="0" ci="3" mb="1" cb="1"/>
			<line nr="10" mi="2" ci="0" mb="0" cb="0"/>
			<counter type="INSTRUCTION" missed="2" covered="5"/>
			<counter type="BRANCH" missed="1" covered="1"/>
			<counter type="LINE" missed="1" covered="2"/>
		</sourcefile>
		<counter type="INSTRUCTION" missed="2" covered="5"/>
		<counter type="BRANCH" missed="1" covered="1"/>
		<counter type="LINE" missed="1" covered="2"/>
	</package>
	<counter type="INSTRUCTION" missed="2" covered="5"/>
	<counter type="BRANCH" missed="1" covered="1"/>
	<counter type="LINE" missed="1" covered="2"/>
</report>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripReportWithGroups(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<report name="Multi-Module Coverage">
	<sessioninfo id="multi-session" start="1234567890" dump="1234567900"/>
	<group name="Module A">
		<package name="com/example/modulea">
			<class name="com/example/modulea/ServiceA">
				<counter type="INSTRUCTION" missed="5" covered="20"/>
				<counter type="CLASS" missed="0" covered="1"/>
			</class>
			<counter type="INSTRUCTION" missed="5" covered="20"/>
			<counter type="CLASS" missed="0" covered="1"/>
		</package>
		<counter type="INSTRUCTION" missed="5" covered="20"/>
		<counter type="CLASS" missed="0" covered="1"/>
	</group>
	<group name="Module B">
		<package name="com/example/moduleb">
			<class name="com/example/moduleb/ServiceB">
				<counter type="INSTRUCTION" missed="10" covered="30"/>
				<counter type="CLASS" missed="0" covered="1"/>
			</class>
			<counter type="INSTRUCTION" missed="10" covered="30"/>
			<counter type="CLASS" missed="0" covered="1"/>
		</package>
		<counter type="INSTRUCTION" missed="10" covered="30"/>
		<counter type="CLASS" missed="0" covered="1"/>
	</group>
	<counter type="INSTRUCTION" missed="15" covered="50"/>
	<counter type="CLASS" missed="0" covered="2"/>
</report>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripNestedGroups(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<report name="Nested Groups Report">
	<sessioninfo id="nested-session" start="1234567890" dump="1234567900"/>
	<group name="Parent Group">
		<group name="Child Group">
			<package name="com/example/child">
				<class name="com/example/child/Child">
					<counter type="INSTRUCTION" missed="0" covered="10"/>
					<counter type="CLASS" missed="0" covered="1"/>
				</class>
				<counter type="INSTRUCTION" missed="0" covered="10"/>
				<counter type="CLASS" missed="0" covered="1"/>
			</package>
			<counter type="INSTRUCTION" missed="0" covered="10"/>
			<counter type="CLASS" missed="0" covered="1"/>
		</group>
		<counter type="INSTRUCTION" missed="0" covered="10"/>
		<counter type="CLASS" missed="0" covered="1"/>
	</group>
	<counter type="INSTRUCTION" missed="0" covered="10"/>
	<counter type="CLASS" missed="0" covered="1"/>
</report>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripEmptyReport(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<report name="Empty Report">
	<counter type="INSTRUCTION" missed="0" covered="0"/>
	<counter type="BRANCH" missed="0" covered="0"/>
	<counter type="LINE" missed="0" covered="0"/>
	<counter type="METHOD" missed="0" covered="0"/>
	<counter type="CLASS" missed="0" covered="0"/>
</report>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripMultipleSessions(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<report name="Multi-Session Report">
	<sessioninfo id="session-1" start="1234567890" dump="1234567900"/>
	<sessioninfo id="session-2" start="1234567910" dump="1234567920"/>
	<sessioninfo id="session-3" start="1234567930" dump="1234567940"/>
	<package name="com/example">
		<class name="com/example/Example">
			<counter type="INSTRUCTION" missed="5" covered="15"/>
			<counter type="CLASS" missed="0" covered="1"/>
		</class>
		<counter type="INSTRUCTION" missed="5" covered="15"/>
		<counter type="CLASS" missed="0" covered="1"/>
	</package>
	<counter type="INSTRUCTION" missed="5" covered="15"/>
	<counter type="CLASS" missed="0" covered="1"/>
</report>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripComplexMethod(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<report name="Method Coverage">
	<sessioninfo id="method-test" start="1234567890" dump="1234567900"/>
	<package name="com/example">
		<class name="com/example/Complex" sourcefilename="Complex.java">
			<method name="&lt;init&gt;" desc="()V" line="5">
				<counter type="INSTRUCTION" missed="0" covered="3"/>
				<counter type="LINE" missed="0" covered="1"/>
				<counter type="METHOD" missed="0" covered="1"/>
			</method>
			<method name="complexMethod" desc="(Ljava/lang/String;I)Z" line="10">
				<counter type="INSTRUCTION" missed="2" covered="15"/>
				<counter type="BRANCH" missed="1" covered="3"/>
				<counter type="LINE" missed="1" covered="5"/>
				<counter type="METHOD" missed="0" covered="1"/>
			</method>
			<counter type="INSTRUCTION" missed="2" covered="18"/>
			<counter type="BRANCH" missed="1" covered="3"/>
			<counter type="LINE" missed="1" covered="6"/>
			<counter type="METHOD" missed="0" covered="2"/>
			<counter type="CLASS" missed="0" covered="1"/>
		</class>
		<counter type="INSTRUCTION" missed="2" covered="18"/>
		<counter type="BRANCH" missed="1" covered="3"/>
		<counter type="LINE" missed="1" covered="6"/>
		<counter type="METHOD" missed="0" covered="2"/>
		<counter type="CLASS" missed="0" covered="1"/>
	</package>
	<counter type="INSTRUCTION" missed="2" covered="18"/>
	<counter type="BRANCH" missed="1" covered="3"/>
	<counter type="LINE" missed="1" covered="6"/>
	<counter type="METHOD" missed="0" covered="2"/>
	<counter type="CLASS" missed="0" covered="1"/>
</report>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripMixedClassesAndSourceFiles(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<report name="Mixed Report">
	<sessioninfo id="mixed-test" start="1234567890" dump="1234567900"/>
	<package name="com/example">
		<class name="com/example/ClassA">
			<counter type="INSTRUCTION" missed="0" covered="10"/>
			<counter type="CLASS" missed="0" covered="1"/>
		</class>
		<sourcefile name="ClassA.java">
			<line nr="1" mi="0" ci="2"/>
			<line nr="2" mi="0" ci="3"/>
			<counter type="INSTRUCTION" missed="0" covered="10"/>
			<counter type="LINE" missed="0" covered="2"/>
		</sourcefile>
		<counter type="INSTRUCTION" missed="0" covered="10"/>
		<counter type="LINE" missed="0" covered="2"/>
		<counter type="CLASS" missed="0" covered="1"/>
	</package>
	<counter type="INSTRUCTION" missed="0" covered="10"/>
	<counter type="LINE" missed="0" covered="2"/>
	<counter type="CLASS" missed="0" covered="1"/>
</report>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripAllCounterTypes(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<report name="All Counter Types">
	<sessioninfo id="counter-test" start="1234567890" dump="1234567900"/>
	<package name="com/example">
		<class name="com/example/Example">
			<counter type="INSTRUCTION" missed="10" covered="40"/>
			<counter type="BRANCH" missed="5" covered="15"/>
			<counter type="LINE" missed="3" covered="12"/>
			<counter type="COMPLEXITY" missed="2" covered="8"/>
			<counter type="METHOD" missed="1" covered="4"/>
			<counter type="CLASS" missed="0" covered="1"/>
		</class>
		<counter type="INSTRUCTION" missed="10" covered="40"/>
		<counter type="BRANCH" missed="5" covered="15"/>
		<counter type="LINE" missed="3" covered="12"/>
		<counter type="COMPLEXITY" missed="2" covered="8"/>
		<counter type="METHOD" missed="1" covered="4"/>
		<counter type="CLASS" missed="0" covered="1"/>
	</package>
	<counter type="INSTRUCTION" missed="10" covered="40"/>
	<counter type="BRANCH" missed="5" covered="15"/>
	<counter type="LINE" missed="3" covered="12"/>
	<counter type="COMPLEXITY" missed="2" covered="8"/>
	<counter type="METHOD" missed="1" covered="4"/>
	<counter type="CLASS" missed="0" covered="1"/>
</report>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripZeroCoverage(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<report name="Zero Coverage">
	<sessioninfo id="zero-test" start="1234567890" dump="1234567900"/>
	<package name="com/example">
		<class name="com/example/Uncovered">
			<method name="neverCalled" desc="()V">
				<counter type="INSTRUCTION" missed="10" covered="0"/>
				<counter type="LINE" missed="3" covered="0"/>
				<counter type="METHOD" missed="1" covered="0"/>
			</method>
			<counter type="INSTRUCTION" missed="20" covered="0"/>
			<counter type="LINE" missed="5" covered="0"/>
			<counter type="METHOD" missed="2" covered="0"/>
			<counter type="CLASS" missed="1" covered="0"/>
		</class>
		<counter type="INSTRUCTION" missed="20" covered="0"/>
		<counter type="LINE" missed="5" covered="0"/>
		<counter type="METHOD" missed="2" covered="0"/>
		<counter type="CLASS" missed="1" covered="0"/>
	</package>
	<counter type="INSTRUCTION" missed="20" covered="0"/>
	<counter type="LINE" missed="5" covered="0"/>
	<counter type="METHOD" missed="2" covered="0"/>
	<counter type="CLASS" missed="1" covered="0"/>
</report>`

	testRoundTrip(t, xmlData)
}

func TestRoundTripFullCoverage(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<report name="Full Coverage">
	<sessioninfo id="full-test" start="1234567890" dump="1234567900"/>
	<package name="com/example">
		<class name="com/example/FullyCovered">
			<method name="tested" desc="()V">
				<counter type="INSTRUCTION" missed="0" covered="10"/>
				<counter type="LINE" missed="0" covered="3"/>
				<counter type="METHOD" missed="0" covered="1"/>
			</method>
			<counter type="INSTRUCTION" missed="0" covered="20"/>
			<counter type="LINE" missed="0" covered="5"/>
			<counter type="METHOD" missed="0" covered="2"/>
			<counter type="CLASS" missed="0" covered="1"/>
		</class>
		<counter type="INSTRUCTION" missed="0" covered="20"/>
		<counter type="LINE" missed="0" covered="5"/>
		<counter type="METHOD" missed="0" covered="2"/>
		<counter type="CLASS" missed="0" covered="1"/>
	</package>
	<counter type="INSTRUCTION" missed="0" covered="20"/>
	<counter type="LINE" missed="0" covered="5"/>
	<counter type="METHOD" missed="0" covered="2"/>
	<counter type="CLASS" missed="0" covered="1"/>
</report>`

	testRoundTrip(t, xmlData)
}

func testRoundTrip(t *testing.T, xmlData string) {
	t.Helper()

	var original Report
	if err := xml.Unmarshal([]byte(xmlData), &original); err != nil {
		t.Fatalf("Failed to unmarshal original XML: %v", err)
	}

	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}

	var fromJSON Report
	if err := json.Unmarshal(jsonData, &fromJSON); err != nil {
		t.Fatalf("Failed to unmarshal from JSON: %v", err)
	}

	xmlRoundTrip, err := xml.Marshal(fromJSON)
	if err != nil {
		t.Fatalf("Failed to marshal back to XML: %v", err)
	}

	var finalReport Report
	if err := xml.Unmarshal(xmlRoundTrip, &finalReport); err != nil {
		t.Fatalf("Failed to unmarshal final XML: %v", err)
	}

	if diff := cmp.Diff(original, finalReport); diff != "" {
		t.Errorf("Round trip conversion resulted in differences (-original +final):\n%s", diff)
		t.Logf("Original XML:\n%s", xmlData)
		t.Logf("JSON intermediate:\n%s", string(jsonData))
		t.Logf("Final XML:\n%s", string(xmlRoundTrip))
	}
}

func TestErrorHandling_MalformedXML(t *testing.T) {
	tests := []struct {
		name    string
		xmlData string
	}{
		{
			name:    "unclosed tag",
			xmlData: `<report name="test"><package name="com/example"</report>`,
		},
		{
			name:    "mismatched tags",
			xmlData: `<report name="test"><package name="com/example"></report></package>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var report Report
			err := xml.Unmarshal([]byte(tt.xmlData), &report)
			if err == nil {
				t.Errorf("Expected error for %s, but got nil", tt.name)
			}
		})
	}
}

func TestJSONFormat_StructureValidation(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<report name="JSON Test">
	<sessioninfo id="json-test" start="1234567890" dump="1234567900"/>
	<package name="com/example">
		<class name="com/example/Test">
			<counter type="INSTRUCTION" missed="5" covered="15"/>
			<counter type="CLASS" missed="0" covered="1"/>
		</class>
		<counter type="INSTRUCTION" missed="5" covered="15"/>
		<counter type="CLASS" missed="0" covered="1"/>
	</package>
	<counter type="INSTRUCTION" missed="5" covered="15"/>
	<counter type="CLASS" missed="0" covered="1"/>
</report>`

	var report Report
	if err := xml.Unmarshal([]byte(xmlData), &report); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	jsonData, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal JSON to map: %v", err)
	}

	if jsonMap["name"] != "JSON Test" {
		t.Errorf("Expected name='JSON Test', got %v", jsonMap["name"])
	}

	sessionInfos, ok := jsonMap["sessionInfo"].([]interface{})
	if !ok {
		t.Fatalf("Expected sessionInfo to be an array")
	}

	if len(sessionInfos) != 1 {
		t.Errorf("Expected 1 sessionInfo, got %d", len(sessionInfos))
	}

	packages, ok := jsonMap["packages"].([]interface{})
	if !ok {
		t.Fatalf("Expected packages to be an array")
	}

	if len(packages) != 1 {
		t.Errorf("Expected 1 package, got %d", len(packages))
	}
}

func FuzzReportXML(f *testing.F) {
	seeds := []string{
		`<report name="test"><counter type="INSTRUCTION" missed="0" covered="0"/></report>`,
		`<report name="test"><package name="com/example"><counter type="CLASS" missed="0" covered="1"/></package><counter type="CLASS" missed="0" covered="1"/></report>`,
		`<report name="test"><sessioninfo id="s1" start="123" dump="456"/><counter type="LINE" missed="0" covered="0"/></report>`,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, xmlData string) {
		var report Report
		err := xml.Unmarshal([]byte(xmlData), &report)
		if err != nil {
			return
		}

		jsonData, err := json.Marshal(report)
		if err != nil {
			t.Errorf("Failed to marshal valid XML to JSON: %v", err)
			return
		}

		var fromJSON Report
		err = json.Unmarshal(jsonData, &fromJSON)
		if err != nil {
			t.Errorf("Failed to unmarshal JSON back to struct: %v", err)
		}
	})
}

func TestRoundTrip_TestDataFiles(t *testing.T) {
	testdataDir := "testdata"
	files, err := filepath.Glob(filepath.Join(testdataDir, "*.xml"))
	if err != nil {
		t.Fatalf("Failed to glob testdata files: %v", err)
	}

	if len(files) == 0 {
		t.Skip("No testdata XML files found")
	}

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			xmlData, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("Failed to read file %s: %v", file, err)
			}

			testRoundTrip(t, string(xmlData))
		})
	}
}

func TestRealWorld_SimpleCalculator(t *testing.T) {
	xmlData, err := os.ReadFile("testdata/simple-single-class.xml")
	if err != nil {
		t.Skip("testdata/simple-single-class.xml not found")
	}

	var report Report
	if err := xml.Unmarshal(xmlData, &report); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if report.Name != "Simple Calculator Coverage" {
		t.Errorf("Expected report name 'Simple Calculator Coverage', got %q", report.Name)
	}

	if len(report.Packages) != 1 {
		t.Fatalf("Expected 1 package, got %d", len(report.Packages))
	}

	pkg := report.Packages[0]
	if pkg.Name != "com/example/calculator" {
		t.Errorf("Expected package name 'com/example/calculator', got %q", pkg.Name)
	}

	if len(pkg.Classes) != 1 {
		t.Fatalf("Expected 1 class, got %d", len(pkg.Classes))
	}

	class := pkg.Classes[0]
	if len(class.Methods) != 5 {
		t.Errorf("Expected 5 methods, got %d", len(class.Methods))
	}

	if len(pkg.SourceFiles) != 1 {
		t.Fatalf("Expected 1 source file, got %d", len(pkg.SourceFiles))
	}

	sourceFile := pkg.SourceFiles[0]
	if sourceFile.Name != "Calculator.java" {
		t.Errorf("Expected source file 'Calculator.java', got %q", sourceFile.Name)
	}

	if len(sourceFile.Lines) != 8 {
		t.Errorf("Expected 8 lines, got %d", len(sourceFile.Lines))
	}
}

func TestRealWorld_MultiModuleApplication(t *testing.T) {
	xmlData, err := os.ReadFile("testdata/complex-multi-module.xml")
	if err != nil {
		t.Skip("testdata/complex-multi-module.xml not found")
	}

	var report Report
	if err := xml.Unmarshal(xmlData, &report); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if report.Name != "Multi-Module E-Commerce Application" {
		t.Errorf("Expected report name 'Multi-Module E-Commerce Application', got %q", report.Name)
	}

	if len(report.SessionInfos) != 2 {
		t.Errorf("Expected 2 session infos (multiple test runs), got %d", len(report.SessionInfos))
	}

	if len(report.Groups) != 2 {
		t.Fatalf("Expected 2 top-level groups, got %d", len(report.Groups))
	}

	coreServices := report.Groups[0]
	if coreServices.Name != "Core Services" {
		t.Errorf("Expected first group name 'Core Services', got %q", coreServices.Name)
	}

	if len(coreServices.Groups) != 2 {
		t.Errorf("Expected 2 nested groups in Core Services, got %d", len(coreServices.Groups))
	}

	userService := coreServices.Groups[0]
	if userService.Name != "User Service" {
		t.Errorf("Expected nested group 'User Service', got %q", userService.Name)
	}

	if len(userService.Packages) != 2 {
		t.Errorf("Expected 2 packages in User Service, got %d", len(userService.Packages))
	}

	infrastructure := report.Groups[1]
	if infrastructure.Name != "Infrastructure" {
		t.Errorf("Expected second group name 'Infrastructure', got %q", infrastructure.Name)
	}

	totalInstructionCounter := findCounter(report.Counters, "INSTRUCTION")
	if totalInstructionCounter == nil {
		t.Fatal("Expected INSTRUCTION counter at report level")
	}

	if totalInstructionCounter.Missed != 50 || totalInstructionCounter.Covered != 211 {
		t.Errorf("Expected INSTRUCTION counter missed=50 covered=211, got missed=%d covered=%d",
			totalInstructionCounter.Missed, totalInstructionCounter.Covered)
	}
}

func TestEdgeCase_UnicodeNames(t *testing.T) {
	xmlData, err := os.ReadFile("testdata/edge-case-unicode-names.xml")
	if err != nil {
		t.Skip("testdata/edge-case-unicode-names.xml not found")
	}

	var report Report
	if err := xml.Unmarshal(xmlData, &report); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if report.Name != "Unicode Test Coverage Résumé" {
		t.Errorf("Expected report name with unicode, got %q", report.Name)
	}

	if len(report.SessionInfos) != 1 {
		t.Fatalf("Expected 1 session info, got %d", len(report.SessionInfos))
	}

	if report.SessionInfos[0].ID != "test-session-中文" {
		t.Errorf("Expected session ID with Chinese characters, got %q", report.SessionInfos[0].ID)
	}

	if len(report.Packages) != 1 {
		t.Fatalf("Expected 1 package, got %d", len(report.Packages))
	}

	pkg := report.Packages[0]
	if pkg.Name != "com/example/测试/パッケージ" {
		t.Errorf("Expected package name with unicode characters, got %q", pkg.Name)
	}

	if len(pkg.Classes) != 1 {
		t.Fatalf("Expected 1 class, got %d", len(pkg.Classes))
	}

	class := pkg.Classes[0]
	if class.Name != "com/example/测试/パッケージ/Übungsklasse" {
		t.Errorf("Expected class name with unicode characters, got %q", class.Name)
	}
}

func TestEdgeCase_LongSourceFile(t *testing.T) {
	xmlData, err := os.ReadFile("testdata/edge-case-long-source-file.xml")
	if err != nil {
		t.Skip("testdata/edge-case-long-source-file.xml not found")
	}

	var report Report
	if err := xml.Unmarshal(xmlData, &report); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(report.Packages) != 1 {
		t.Fatalf("Expected 1 package, got %d", len(report.Packages))
	}

	pkg := report.Packages[0]
	if len(pkg.SourceFiles) != 1 {
		t.Fatalf("Expected 1 source file, got %d", len(pkg.SourceFiles))
	}

	sourceFile := pkg.SourceFiles[0]
	if len(sourceFile.Lines) != 26 {
		t.Errorf("Expected 26 lines in source file, got %d", len(sourceFile.Lines))
	}

	for _, line := range sourceFile.Lines {
		if line.Nr < 1 || line.Nr > 500 {
			t.Errorf("Line number %d out of expected range", line.Nr)
		}
	}

	if sourceFile.Lines[len(sourceFile.Lines)-1].Nr != 500 {
		t.Errorf("Expected last line to be 500, got %d", sourceFile.Lines[len(sourceFile.Lines)-1].Nr)
	}
}

func TestEdgeCase_PartialLineAttributes(t *testing.T) {
	xmlData, err := os.ReadFile("testdata/edge-case-partial-line-attributes.xml")
	if err != nil {
		t.Skip("testdata/edge-case-partial-line-attributes.xml not found")
	}

	var report Report
	if err := xml.Unmarshal(xmlData, &report); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(report.Packages) != 1 {
		t.Fatalf("Expected 1 package, got %d", len(report.Packages))
	}

	pkg := report.Packages[0]
	if len(pkg.SourceFiles) != 1 {
		t.Fatalf("Expected 1 source file, got %d", len(pkg.SourceFiles))
	}

	sourceFile := pkg.SourceFiles[0]
	if len(sourceFile.Lines) != 6 {
		t.Errorf("Expected 6 lines, got %d", len(sourceFile.Lines))
	}

	line1 := sourceFile.Lines[0]
	if line1.MissedInstr != 0 || line1.CoveredInstr != 4 {
		t.Errorf("Line 1: expected mi=0 ci=4, got mi=%d ci=%d", line1.MissedInstr, line1.CoveredInstr)
	}

	if line1.MissedBranches != 0 || line1.CoveredBranches != 0 {
		t.Errorf("Line 1: expected mb=0 cb=0 (omitted attributes), got mb=%d cb=%d",
			line1.MissedBranches, line1.CoveredBranches)
	}

	line3 := sourceFile.Lines[2]
	if line3.MissedBranches != 1 || line3.CoveredBranches != 1 {
		t.Errorf("Line 3: expected mb=1 cb=1, got mb=%d cb=%d",
			line3.MissedBranches, line3.CoveredBranches)
	}
}

func TestEdgeCase_MultiplePackagesNoGroups(t *testing.T) {
	xmlData, err := os.ReadFile("testdata/edge-case-multiple-packages-no-groups.xml")
	if err != nil {
		t.Skip("testdata/edge-case-multiple-packages-no-groups.xml not found")
	}

	var report Report
	if err := xml.Unmarshal(xmlData, &report); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(report.Groups) != 0 {
		t.Errorf("Expected 0 groups, got %d", len(report.Groups))
	}

	if len(report.Packages) != 3 {
		t.Fatalf("Expected 3 packages (no groups), got %d", len(report.Packages))
	}

	expectedPackageNames := []string{
		"com/example/package1",
		"com/example/package2",
		"com/example/package3",
	}

	for i, expectedName := range expectedPackageNames {
		if report.Packages[i].Name != expectedName {
			t.Errorf("Package %d: expected name %q, got %q", i, expectedName, report.Packages[i].Name)
		}
	}
}

func TestEdgeCase_ClassWithoutMethods(t *testing.T) {
	xmlData, err := os.ReadFile("testdata/edge-case-class-without-methods.xml")
	if err != nil {
		t.Skip("testdata/edge-case-class-without-methods.xml not found")
	}

	var report Report
	if err := xml.Unmarshal(xmlData, &report); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(report.Packages) != 1 {
		t.Fatalf("Expected 1 package, got %d", len(report.Packages))
	}

	pkg := report.Packages[0]
	if len(pkg.Classes) != 3 {
		t.Fatalf("Expected 3 classes, got %d", len(pkg.Classes))
	}

	emptyClass := pkg.Classes[0]
	if emptyClass.Name != "com/example/empty/EmptyClass" {
		t.Errorf("Expected class name 'com/example/empty/EmptyClass', got %q", emptyClass.Name)
	}

	if len(emptyClass.Methods) != 0 {
		t.Errorf("Expected 0 methods in EmptyClass, got %d", len(emptyClass.Methods))
	}

	if len(emptyClass.Counters) != 4 {
		t.Errorf("Expected 4 counters in EmptyClass, got %d", len(emptyClass.Counters))
	}

	interfaceClass := pkg.Classes[1]
	if len(interfaceClass.Methods) != 0 {
		t.Errorf("Expected 0 methods in InterfaceClass, got %d", len(interfaceClass.Methods))
	}

	if len(interfaceClass.Counters) != 1 {
		t.Errorf("Expected 1 counter in InterfaceClass, got %d", len(interfaceClass.Counters))
	}

	classCounter := findCounter(interfaceClass.Counters, "CLASS")
	if classCounter == nil {
		t.Fatal("Expected CLASS counter in InterfaceClass")
	}

	if classCounter.Missed != 0 || classCounter.Covered != 1 {
		t.Errorf("Expected CLASS counter missed=0 covered=1, got missed=%d covered=%d",
			classCounter.Missed, classCounter.Covered)
	}
}

func TestInvalidCounterType_StillParses(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<report name="Invalid Counter Type Test">
	<sessioninfo id="invalid-test" start="1234567890" dump="1234567900"/>
	<package name="com/example">
		<class name="com/example/Test">
			<counter type="INSTRUCTION" missed="5" covered="15"/>
			<counter type="INVALID_TYPE" missed="0" covered="0"/>
			<counter type="CLASS" missed="0" covered="1"/>
		</class>
		<counter type="INSTRUCTION" missed="5" covered="15"/>
		<counter type="CUSTOM_METRIC" missed="2" covered="8"/>
		<counter type="CLASS" missed="0" covered="1"/>
	</package>
	<counter type="INSTRUCTION" missed="5" covered="15"/>
	<counter type="CLASS" missed="0" covered="1"/>
</report>`

	var report Report
	if err := xml.Unmarshal([]byte(xmlData), &report); err != nil {
		t.Fatalf("Failed to unmarshal XML with invalid counter types: %v", err)
	}

	if len(report.Packages) != 1 {
		t.Fatalf("Expected 1 package, got %d", len(report.Packages))
	}

	pkg := report.Packages[0]
	if len(pkg.Classes) != 1 {
		t.Fatalf("Expected 1 class, got %d", len(pkg.Classes))
	}

	class := pkg.Classes[0]
	if len(class.Counters) != 3 {
		t.Fatalf("Expected 3 counters (including invalid type), got %d", len(class.Counters))
	}

	invalidCounter := class.Counters[1]
	if invalidCounter.Type != "INVALID_TYPE" {
		t.Errorf("Expected counter type 'INVALID_TYPE', got %q", invalidCounter.Type)
	}

	if len(pkg.Counters) != 3 {
		t.Fatalf("Expected 3 package counters (including custom), got %d", len(pkg.Counters))
	}

	customCounter := pkg.Counters[1]
	if customCounter.Type != "CUSTOM_METRIC" {
		t.Errorf("Expected counter type 'CUSTOM_METRIC', got %q", customCounter.Type)
	}

	jsonData, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("Failed to marshal to JSON with invalid counter types: %v", err)
	}

	var fromJSON Report
	if err := json.Unmarshal(jsonData, &fromJSON); err != nil {
		t.Fatalf("Failed to unmarshal from JSON: %v", err)
	}

	if fromJSON.Packages[0].Classes[0].Counters[1].Type != "INVALID_TYPE" {
		t.Errorf("Invalid counter type not preserved through JSON round-trip")
	}
}

func findCounter(counters []Counter, counterType string) *Counter {
	for i := range counters {
		if counters[i].Type == counterType {
			return &counters[i]
		}
	}
	return nil
}
