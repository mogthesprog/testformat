# JaCoCo XML to JSON Format Documentation

This document describes the JSON format produced by converting JaCoCo XML coverage reports.

## Overview

The converter transforms JaCoCo XML reports (conforming to DTD Report 1.1) into JSON while preserving all structural information and coverage metrics. The conversion is bidirectional - JSON can be converted back to XML without data loss.

## Structure

The JSON structure mirrors the XML hierarchy with the following mappings:

### Report (Root Element)

**XML:**
```xml
<report name="My Coverage Report">
  <!-- content -->
</report>
```

**JSON:**
```json
{
  "name": "My Coverage Report",
  "sessionInfo": [],
  "groups": [],
  "packages": [],
  "counters": []
}
```

### Session Info

Records execution data from test runs.

**XML:**
```xml
<sessioninfo id="test-session-1" start="1705328422000" dump="1705328425000"/>
```

**JSON:**
```json
{
  "id": "test-session-1",
  "start": "1705328422000",
  "dump": "1705328425000"
}
```

**Fields:**
- `id` (string): Session identifier
- `start` (string): Start timestamp (milliseconds since epoch)
- `dump` (string): Dump timestamp (milliseconds since epoch)

### Groups

Logical grouping of packages (e.g., modules, components).

**XML:**
```xml
<group name="Core Services">
  <group name="User Service">
    <!-- nested content -->
  </group>
  <package name="com/example/util">
    <!-- package content -->
  </package>
  <counter type="INSTRUCTION" missed="10" covered="100"/>
</group>
```

**JSON:**
```json
{
  "name": "Core Services",
  "groups": [
    {
      "name": "User Service",
      "groups": [],
      "packages": [],
      "counters": []
    }
  ],
  "packages": [
    {
      "name": "com/example/util",
      "classes": [],
      "sourcefiles": [],
      "counters": []
    }
  ],
  "counters": [
    {
      "type": "INSTRUCTION",
      "missed": 10,
      "covered": 100
    }
  ]
}
```

**Fields:**
- `name` (string): Group name
- `groups` (array): Nested groups (recursive structure)
- `packages` (array): Packages in this group
- `counters` (array): Aggregated coverage metrics

### Packages

Java packages containing classes and source files.

**XML:**
```xml
<package name="com/example/myapp">
  <class name="com/example/myapp/Calculator">
    <!-- class content -->
  </class>
  <sourcefile name="Calculator.java">
    <!-- source file content -->
  </sourcefile>
  <counter type="CLASS" missed="0" covered="1"/>
</package>
```

**JSON:**
```json
{
  "name": "com/example/myapp",
  "classes": [
    {
      "name": "com/example/myapp/Calculator",
      "methods": [],
      "counters": []
    }
  ],
  "sourceFiles": [
    {
      "name": "Calculator.java",
      "lines": [],
      "counters": []
    }
  ],
  "counters": [
    {
      "type": "CLASS",
      "missed": 0,
      "covered": 1
    }
  ]
}
```

**Fields:**
- `name` (string): Package name in VM notation (e.g., `com/example/myapp`)
- `classes` (array): Classes in this package
- `sourceFiles` (array): Source files in this package
- `counters` (array): Package-level coverage metrics

### Classes

Java classes with methods and coverage information.

**XML:**
```xml
<class name="com/example/Calculator" sourcefilename="Calculator.java">
  <method name="add" desc="(II)I" line="10">
    <counter type="METHOD" missed="0" covered="1"/>
  </method>
  <counter type="CLASS" missed="0" covered="1"/>
</class>
```

**JSON:**
```json
{
  "name": "com/example/Calculator",
  "sourceFileName": "Calculator.java",
  "methods": [
    {
      "name": "add",
      "desc": "(II)I",
      "line": 10,
      "counters": [
        {
          "type": "METHOD",
          "missed": 0,
          "covered": 1
        }
      ]
    }
  ],
  "counters": [
    {
      "type": "CLASS",
      "missed": 0,
      "covered": 1
    }
  ]
}
```

**Fields:**
- `name` (string): Fully qualified class name in VM notation
- `sourceFileName` (string, optional): Name of source file
- `methods` (array): Methods in this class
- `counters` (array): Class-level coverage metrics

**Note:** Classes without methods (e.g., interfaces, annotations) are valid and will have an empty `methods` array.

### Methods

Methods within classes.

**XML:**
```xml
<method name="complexMethod" desc="(Ljava/lang/String;I)Z" line="25">
  <counter type="INSTRUCTION" missed="2" covered="15"/>
  <counter type="BRANCH" missed="1" covered="3"/>
</method>
```

**JSON:**
```json
{
  "name": "complexMethod",
  "desc": "(Ljava/lang/String;I)Z",
  "line": 25,
  "counters": [
    {
      "type": "INSTRUCTION",
      "missed": 2,
      "covered": 15
    },
    {
      "type": "BRANCH",
      "missed": 1,
      "covered": 3
    }
  ]
}
```

**Fields:**
- `name` (string): Method name (constructors use `<init>`, static initializers use `<clinit>`)
- `desc` (string): Method descriptor in JVM format
- `line` (integer, optional): First source line number
- `counters` (array): Method-level coverage metrics

### Source Files

Source files with line-level coverage.

**XML:**
```xml
<sourcefile name="Calculator.java">
  <line nr="10" mi="0" ci="4" mb="0" cb="0"/>
  <line nr="15" mi="0" ci="3" mb="1" cb="1"/>
  <counter type="LINE" missed="0" covered="2"/>
</sourcefile>
```

**JSON:**
```json
{
  "name": "Calculator.java",
  "lines": [
    {
      "nr": 10,
      "mi": 0,
      "ci": 4,
      "mb": 0,
      "cb": 0
    },
    {
      "nr": 15,
      "mi": 0,
      "ci": 3,
      "mb": 1,
      "cb": 1
    }
  ],
  "counters": [
    {
      "type": "LINE",
      "missed": 0,
      "covered": 2
    }
  ]
}
```

**Fields:**
- `name` (string): Source file name
- `lines` (array): Line-level coverage data
- `counters` (array): File-level coverage metrics

### Lines

Individual source code lines with coverage.

**XML:**
```xml
<line nr="25" mi="0" ci="3" mb="1" cb="1"/>
```

**JSON:**
```json
{
  "nr": 25,
  "mi": 0,
  "ci": 3,
  "mb": 1,
  "cb": 1
}
```

**Fields (all optional except `nr`):**
- `nr` (integer): Line number
- `mi` (integer): Missed instructions
- `ci` (integer): Covered instructions
- `mb` (integer): Missed branches
- `cb` (integer): Covered branches

**Note:** Attributes are optional per DTD. Lines may have only instruction coverage, only branch coverage, or both.

### Counters

Coverage metrics for different measurement types.

**XML:**
```xml
<counter type="INSTRUCTION" missed="10" covered="50"/>
<counter type="BRANCH" missed="5" covered="15"/>
<counter type="LINE" missed="2" covered="20"/>
<counter type="COMPLEXITY" missed="3" covered="12"/>
<counter type="METHOD" missed="1" covered="10"/>
<counter type="CLASS" missed="0" covered="5"/>
```

**JSON:**
```json
[
  {"type": "INSTRUCTION", "missed": 10, "covered": 50},
  {"type": "BRANCH", "missed": 5, "covered": 15},
  {"type": "LINE", "missed": 2, "covered": 20},
  {"type": "COMPLEXITY", "missed": 3, "covered": 12},
  {"type": "METHOD", "missed": 1, "covered": 10},
  {"type": "CLASS", "missed": 0, "covered": 5}
]
```

**Fields:**
- `type` (string): Metric type
- `missed` (integer): Number of missed items
- `covered` (integer): Number of covered items

**Standard Counter Types:**
- `INSTRUCTION`: Bytecode instruction coverage
- `BRANCH`: Branch coverage (if/else, switch, etc.)
- `LINE`: Source line coverage
- `COMPLEXITY`: Cyclomatic complexity
- `METHOD`: Method coverage
- `CLASS`: Class coverage

**Note:** The parser accepts any counter type string, not just the standard types. Custom or future counter types are preserved through round-trip conversion.

## Complete Example

### Input XML

```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<!DOCTYPE report PUBLIC "-//JACOCO//DTD Report 1.1//EN" "report.dtd">
<report name="Example Coverage Report">
  <sessioninfo id="test-run-1" start="1705328422000" dump="1705328425000"/>
  <package name="com/example/calculator">
    <class name="com/example/calculator/Calculator" sourcefilename="Calculator.java">
      <method name="add" desc="(II)I" line="10">
        <counter type="INSTRUCTION" missed="0" covered="4"/>
        <counter type="METHOD" missed="0" covered="1"/>
      </method>
      <method name="multiply" desc="(II)I" line="18">
        <counter type="INSTRUCTION" missed="4" covered="0"/>
        <counter type="METHOD" missed="1" covered="0"/>
      </method>
      <counter type="INSTRUCTION" missed="4" covered="19"/>
      <counter type="METHOD" missed="1" covered="4"/>
      <counter type="CLASS" missed="0" covered="1"/>
    </class>
    <sourcefile name="Calculator.java">
      <line nr="10" mi="0" ci="4" mb="0" cb="0"/>
      <line nr="18" mi="4" ci="0" mb="0" cb="0"/>
      <counter type="INSTRUCTION" missed="5" covered="19"/>
      <counter type="LINE" missed="1" covered="1"/>
    </sourcefile>
    <counter type="INSTRUCTION" missed="5" covered="19"/>
    <counter type="LINE" missed="1" covered="1"/>
    <counter type="METHOD" missed="1" covered="4"/>
    <counter type="CLASS" missed="0" covered="1"/>
  </package>
  <counter type="INSTRUCTION" missed="5" covered="19"/>
  <counter type="LINE" missed="1" covered="1"/>
  <counter type="METHOD" missed="1" covered="4"/>
  <counter type="CLASS" missed="0" covered="1"/>
</report>
```

### Output JSON

```json
{
  "name": "Example Coverage Report",
  "sessionInfo": [
    {
      "id": "test-run-1",
      "start": "1705328422000",
      "dump": "1705328425000"
    }
  ],
  "packages": [
    {
      "name": "com/example/calculator",
      "classes": [
        {
          "name": "com/example/calculator/Calculator",
          "sourceFileName": "Calculator.java",
          "methods": [
            {
              "name": "add",
              "desc": "(II)I",
              "line": 10,
              "counters": [
                {"type": "INSTRUCTION", "missed": 0, "covered": 4},
                {"type": "METHOD", "missed": 0, "covered": 1}
              ]
            },
            {
              "name": "multiply",
              "desc": "(II)I",
              "line": 18,
              "counters": [
                {"type": "INSTRUCTION", "missed": 4, "covered": 0},
                {"type": "METHOD", "missed": 1, "covered": 0}
              ]
            }
          ],
          "counters": [
            {"type": "INSTRUCTION", "missed": 4, "covered": 19},
            {"type": "METHOD", "missed": 1, "covered": 4},
            {"type": "CLASS", "missed": 0, "covered": 1}
          ]
        }
      ],
      "sourceFiles": [
        {
          "name": "Calculator.java",
          "lines": [
            {"nr": 10, "mi": 0, "ci": 4, "mb": 0, "cb": 0},
            {"nr": 18, "mi": 4, "ci": 0, "mb": 0, "cb": 0}
          ],
          "counters": [
            {"type": "INSTRUCTION", "missed": 5, "covered": 19},
            {"type": "LINE", "missed": 1, "covered": 1}
          ]
        }
      ],
      "counters": [
        {"type": "INSTRUCTION", "missed": 5, "covered": 19},
        {"type": "LINE", "missed": 1, "covered": 1},
        {"type": "METHOD", "missed": 1, "covered": 4},
        {"type": "CLASS", "missed": 0, "covered": 1}
      ]
    }
  ],
  "counters": [
    {"type": "INSTRUCTION", "missed": 5, "covered": 19},
    {"type": "LINE", "missed": 1, "covered": 1},
    {"type": "METHOD", "missed": 1, "covered": 4},
    {"type": "CLASS", "missed": 0, "covered": 1}
  ]
}
```

## Usage

### Go

```go
import (
    "encoding/json"
    "encoding/xml"
    "github.com/mogthesprog/testformat/jacoco"
)

// XML to JSON
var report jacoco.Report
xml.Unmarshal(xmlData, &report)
jsonData, _ := json.MarshalIndent(report, "", "  ")

// JSON to XML
var report jacoco.Report
json.Unmarshal(jsonData, &report)
xmlData, _ := xml.MarshalIndent(report, "", "  ")
```

## Implementation Notes

1. **Empty Arrays**: Empty arrays are represented as `[]` in JSON, but omitted in XML per the DTD's `*` and `+` quantifiers.

2. **Optional Fields**: Optional attributes (per DTD `#IMPLIED`) may be omitted from both XML and JSON. Zero values in Go structs are serialized as `0` in JSON but may be omitted in XML if marked as `omitempty`.

3. **Bidirectional Conversion**: The format supports lossless round-trip conversion: XML → JSON → XML produces structurally equivalent output.

4. **Counter Types**: While the DTD specifies six standard counter types, the implementation accepts any string value for the `type` attribute, allowing for custom metrics or future extensions.

5. **Unicode Support**: Full Unicode support in all string fields (names, identifiers, etc.).

6. **Method Descriptors**: Method descriptors use JVM format:
   - `()V` - void method with no parameters
   - `(II)I` - method taking two ints, returning int
   - `(Ljava/lang/String;)Z` - method taking String, returning boolean

## Calculating Coverage Percentages

Coverage percentages can be calculated from counter data:

```javascript
function calculateCoverage(counter) {
  const total = counter.missed + counter.covered;
  if (total === 0) return 0;
  return (counter.covered / total) * 100;
}

// Example: Calculate instruction coverage
const instructionCounter = report.counters.find(c => c.type === "INSTRUCTION");
const coverage = calculateCoverage(instructionCounter);
console.log(`Instruction coverage: ${coverage.toFixed(2)}%`);
```

## References

- [JaCoCo Official Documentation](https://www.jacoco.org/)
- [JaCoCo DTD Report 1.1](../report.dtd)
- [Test Examples](../testdata/)
